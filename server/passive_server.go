package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/valyala/fasthttp/reuseport"
	. "github.com/xeronith/diamante/actor"
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/network/http"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/io"
	dispatcher "github.com/xeronith/diamante/network/http"
)

func (server *defaultServer) startPassiveServer() {
	listener, err := reuseport.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", server.passivePort))
	if err != nil {
		log.Fatalf("PASSIVE SERVER LISTENER FATAL ERROR: %s", err)
	}

	server.listeners.Append(listener)

	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{},
		AllowHeaders:     []string{},
		AllowMethods:     []string{"POST"},
		AllowCredentials: false,
		MaxAge:           5600,
	})

	defaultHandler := func(context echo.Context) error {
		var (
			err     error
			message []byte
		)

		message, err = ioutil.ReadAll(context.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "")
		} else {
			var (
				actor  IActor
				result IOperationResult
			)

			writer := CreateHttpWriter(server, context, server.secureCookie)
			actor = CreateActor(writer, false, context.RealIP(), context.Request().UserAgent())
			contentType := context.Request().Header.Get("Content-Type")
			if contentType == "application/json" {
				result = server.OnActorTextData(actor, string(message))
			} else {
				result = server.OnActorBinaryData(actor, message)
			}

			actor.Dispatch(result)
			context.Response().Flush()
			return nil
		}
	}

	passiveServer := echo.New()
	passiveServer.HideBanner = true
	passiveServer.Debug = false
	passiveServer.HidePort = true
	// passiveServer.Logger.SetLevel(labstackLog.OFF)

	// passiveServer.Use(middleware.Logger())
	// passiveServer.Use(middleware.Recover())
	passiveServer.Use(cors)
	passiveServer.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Response().Header().Set(echo.HeaderServer, "Diamante")
			return next(ctx)
		}
	})

	passiveServer.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := err.Error()

		if e, ok := err.(*echo.HTTPError); ok {
			code = e.Code
			message = fmt.Sprintf("%v", e.Message)
		}

		c.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		_ = c.JSON(code,
			struct {
				Message string `json:"message"`
			}{
				Message: message,
			})
	}

	tlsEnabled := server.Configuration().GetServerConfiguration().GetTLSConfiguration().IsEnabled()

	if tlsEnabled {
		passiveServer.Listener = server.createTLSListener(listener)
	} else {
		passiveServer.Listener = listener
	}

	passiveServer.POST("/", defaultHandler)
	passiveServer.POST("/diagnostics/:clientType", diagnosticsHandler)

	const (
		UPLOAD_PATH     = "./media"
		MAX_UPLOAD_SIZE = 10 * 1024 * 1024 // 10MB
	)

	passiveServer.POST("/media", func(ctx echo.Context) error {
		writer := ctx.Response().Writer
		request := ctx.Request()

		request.Body = http.MaxBytesReader(writer, request.Body, MAX_UPLOAD_SIZE)
		if err := request.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "FILE_TOO_BIG")
		}

		file, _, err := request.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE")
		}

		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE")
		}

		// DetectContentType only needs the first 512 bytes
		filetype := http.DetectContentType(fileBytes)
		switch filetype {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
		case "video/x-flv":
		case "video/mp4":
		case "application/x-mpegURL":
		case "video/MP2T":
		case "video/quicktime":
		case "video/3gpp":
		case "video/x-msvideo":
		case "video/x-ms-wmv":
			break
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE_TYPE")
		}

		data := make([]byte, 12)
		rand.Read(data)
		fileName := fmt.Sprintf("%x_%d", data, time.Now().UnixNano())

		fileEndings, err := mime.ExtensionsByType(filetype)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "CANT_READ_FILE_TYPE")
		}

		newPath := filepath.Join(UPLOAD_PATH, fileName+fileEndings[len(fileEndings)-1])
		newFile, err := os.Create(newPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "CANT_WRITE_FILE")
		}

		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "CANT_WRITE_FILE")
		}

		return ctx.JSON(http.StatusOK, struct {
			Url string `json:"url"`
		}{
			Url: fmt.Sprintf("%s://%s/%s",
				server.configuration.GetServerConfiguration().GetProtocol(),
				server.configuration.GetServerConfiguration().GetFQDN(),
				newPath,
			),
		})
	})

	passiveServer.POST("/media/batch", func(ctx echo.Context) error {
		writer := ctx.Response().Writer
		request := ctx.Request()

		files, err := ctx.MultipartForm()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE(S)")
		}

		request.Body = http.MaxBytesReader(writer, request.Body, MAX_UPLOAD_SIZE)
		if err := request.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "FILE_TOO_BIG")
		}

		urls := make([]string, 0)

		for _, fileHeader := range files.File["file"] {
			file, err := fileHeader.Open()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE")
			}

			defer file.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE")
			}

			// DetectContentType only needs the first 512 bytes
			filetype := http.DetectContentType(fileBytes)
			switch filetype {
			case
				"image/jpeg", "image/jpg",
				"image/gif", "image/png",
				"video/x-flv",
				"video/mp4",
				"application/x-mpegURL",
				"video/MP2T",
				"video/quicktime",
				"video/3gpp",
				"video/x-msvideo",
				"video/x-ms-wmv":

			default:
				return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE_TYPE")
			}

			data := make([]byte, 12)
			rand.Read(data)
			fileName := fmt.Sprintf("%x_%d", data, time.Now().UnixNano())

			fileEndings, err := mime.ExtensionsByType(filetype)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "CANT_READ_FILE_TYPE")
			}

			newPath := filepath.Join(UPLOAD_PATH, fileName+fileEndings[len(fileEndings)-1])
			newFile, err := os.Create(newPath)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "CANT_WRITE_FILE")
			}

			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "CANT_WRITE_FILE")
			}

			urls = append(urls, fmt.Sprintf("%s://%s/%s",
				server.configuration.GetServerConfiguration().GetProtocol(),
				server.configuration.GetServerConfiguration().GetFQDN(),
				newPath,
			))
		}

		return ctx.JSON(http.StatusOK, struct {
			Urls []string `json:"urls"`
		}{
			Urls: urls,
		})
	})

	passiveServer.Static("/media", UPLOAD_PATH)

	passiveServer.GET("/health", func(context echo.Context) error {
		return context.String(http.StatusOK, "OK")
	})

	passiveServer.DELETE("/mem", func(context echo.Context) error {
		runtime.GC()
		debug.FreeOSMemory()
		return context.String(200, "Done")
	})

	for path, httpHandler := range server.httpGetHandlers {
		func(path string, method string, handler HttpHandlerFunc) {
			handlerFunc := func(_context echo.Context) error {
				return handler(
					dispatcher.NewDispatcher(
						server,
						_context.Response().Writer,
						_context.Request(),
						_context.QueryParam,
						_context.Param,
						_context.RealIP(),
					),
				)
			}

			passiveServer.GET(path, handlerFunc)
		}(path, httpHandler.Method(), httpHandler.HandlerFunc())
	}

	for path, httpHandler := range server.httpPostHandlers {
		func(path string, method string, handler HttpHandlerFunc) {
			handlerFunc := func(_context echo.Context) error {
				return handler(
					dispatcher.NewDispatcher(
						server,
						_context.Response().Writer,
						_context.Request(),
						_context.QueryParam,
						_context.Param,
						_context.RealIP(),
					),
				)
			}

			passiveServer.POST(path, handlerFunc)
		}(path, httpHandler.Method(), httpHandler.HandlerFunc())
	}

	passiveServer.Server.ReadTimeout = time.Second * 15
	passiveServer.Server.WriteTimeout = time.Second * 15

	if err := passiveServer.Start(""); err != nil {
		// server.logger.Critical(fmt.Sprintf("PASSIVE SERVER FAILURE: %s", err))
	}
}
