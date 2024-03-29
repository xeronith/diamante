package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/codingsince1985/checksum"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	. "github.com/xeronith/diamante/actor"
	. "github.com/xeronith/diamante/contracts/actor"
	. "github.com/xeronith/diamante/contracts/network/http"
	. "github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/io"
	dispatcher "github.com/xeronith/diamante/network/http"
	"github.com/xeronith/diamante/utility"
)

type uploadedMedia struct {
	ContentType  string `json:"contentType"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	Url          string `json:"url"`
	Checksum     string `json:"checksum"`
}

func (server *defaultServer) startPassiveServer() {
	listener, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", server.passivePort))
	if err != nil {
		log.Fatalf("PASSIVE SERVER LISTENER FATAL ERROR: %s", err)
	}

	server.listeners.Append(listener)

	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: server.Configuration().GetAllowedOrigins(),
		AllowHeaders: []string{},
		ExposeHeaders: []string{
			"X-Request-Timestamp",
			"X-Response-Signature",
			"X-Turbo",
			"X-Note",
			"X-Quote",
		},
		AllowMethods:     []string{"POST"},
		AllowCredentials: true,
		MaxAge:           5600,
	})

	defaultHandler := func(context echo.Context) error {
		var (
			err     error
			message []byte
		)

		message, err = io.ReadAll(context.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "")
		} else {
			var (
				actor  IActor
				result IOperationResult
			)

			writer := CreateHttpWriter(
				server,
				context,
				server.secureCookie,
			)

			actor = CreateActor(
				writer,
				false,
				context.Request().Header.Get("X-Request-Signature"),
				context.RealIP(),
				context.Request().UserAgent(),
			)

			if writer.ContentType() == echo.MIMEApplicationJSON {
				var body map[string]interface{}
				if err := json.Unmarshal(message, &body); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, "")
				}

				if _, exists := body["payload"]; exists {
					payload, _ := json.Marshal(body["payload"])
					body["payload"] = payload
					message, _ = json.Marshal(body)
				}
			}

			result = server.OnData(actor, message)

			actor.Dispatch(result)
			context.Response().Flush()
			return nil
		}
	}

	passiveServer := echo.New()
	passiveServer.HideBanner = true
	passiveServer.Debug = false
	passiveServer.HidePort = true

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
		PLAYGROUND_PATH = "./playground"
		MAX_UPLOAD_SIZE = 1024 * 1024 * 1024 // 1GB
	)

	passiveServer.POST("/media", func(ctx echo.Context) error {
		writer := ctx.Response().Writer
		request := ctx.Request()

		request.Body = http.MaxBytesReader(writer, request.Body, MAX_UPLOAD_SIZE)
		if err := request.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "FILE_TOO_BIG")
		}

		file, header, err := request.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE")
		}

		defer file.Close()
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE")
		}

		// DetectContentType only needs the first 512 bytes
		contentType := http.DetectContentType(fileBytes)
		/* switch fileType {
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
		} */

		data := make([]byte, 12)
		rand.Read(data)
		tempFileName := fmt.Sprintf("%x_%d", data, time.Now().UnixNano())
		fileExtension := strings.ToLower(filepath.Ext(header.Filename))
		if fileExtension == ".error" {
			return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE_TYPE")
		}

		tempPath := filepath.Join(UPLOAD_PATH, fmt.Sprintf("%s%s", tempFileName, fileExtension))
		tempFile, err := os.Create(tempPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "CANT_WRITE_FILE")
		}

		defer tempFile.Close()
		if _, err := tempFile.Write(fileBytes); err != nil || tempFile.Close() != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "CANT_WRITE_FILE")
		}

		sha256Checksum, err := checksum.SHA256sum(tempPath)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "CANT_CALCULATE_CHECKSUM")
		}

		config := server.configuration.GetServerConfiguration()
		protocol := config.GetProtocol()
		fqdn := config.GetFQDN()

		uploadedFilePath := filepath.Join(UPLOAD_PATH, fmt.Sprintf("%s%s", sha256Checksum, fileExtension))
		uploadedFileUrl := fmt.Sprintf("%s://%s/%s", protocol, fqdn, uploadedFilePath)

		thumbnailPath := filepath.Join(UPLOAD_PATH, fmt.Sprintf("%s_thumbnail%s", sha256Checksum, fileExtension))
		thumbnailUrl := fmt.Sprintf("%s://%s/%s", protocol, fqdn, thumbnailPath)

		if fileStat, err := os.Stat(uploadedFilePath); os.IsNotExist(err) {
			if err := os.Rename(tempPath, uploadedFilePath); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "CANT_MOVE_FILE")
			}

			if err := utility.CreateThumbnail(uploadedFilePath, thumbnailPath, 400, 400); err != nil {
				thumbnailUrl = ""
			}
		} else {
			_ = os.Remove(tempPath)

			if fileStat.Size() != int64(len(fileBytes)) {
				return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE_CONTENT")
			}
		}

		return ctx.JSON(http.StatusOK, uploadedMedia{
			ContentType:  contentType,
			ThumbnailUrl: thumbnailUrl,
			Url:          uploadedFileUrl,
			Checksum:     sha256Checksum,
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

		uploadedFiles := make([]uploadedMedia, 0)

		for _, fileHeader := range files.File["file"] {
			file, err := fileHeader.Open()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE")
			}

			defer file.Close()
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "INVALID_FILE")
			}

			// DetectContentType only needs the first 512 bytes
			fileType := http.DetectContentType(fileBytes)
			/* switch fileType {
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
			} */

			data := make([]byte, 12)
			rand.Read(data)
			fileName := fmt.Sprintf("%x_%d", data, time.Now().UnixNano())

			fileEndings, err := mime.ExtensionsByType(fileType)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "CANT_READ_FILE_TYPE")
			}

			newPath := filepath.Join(UPLOAD_PATH, fileName+fileEndings[len(fileEndings)-1])
			thumbnailPath := filepath.Join(UPLOAD_PATH, fileName+"_thumbnail.jpg")
			newFile, err := os.Create(newPath)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "CANT_WRITE_FILE")
			}

			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "CANT_WRITE_FILE")
			}

			thumbnailUrl := ""
			if err := utility.CreateThumbnail(newPath, thumbnailPath, 400, 400); err == nil {
				thumbnailUrl = fmt.Sprintf("%s://%s/%s",
					server.configuration.GetServerConfiguration().GetProtocol(),
					server.configuration.GetServerConfiguration().GetFQDN(),
					thumbnailPath,
				)
			}

			uploadedFiles = append(uploadedFiles, uploadedMedia{
				ContentType:  fileType,
				ThumbnailUrl: thumbnailUrl,
				Url: fmt.Sprintf("%s://%s/%s",
					server.configuration.GetServerConfiguration().GetProtocol(),
					server.configuration.GetServerConfiguration().GetFQDN(),
					newPath,
				),
			})
		}

		return ctx.JSON(http.StatusOK, struct {
			Files []uploadedMedia `json:"files"`
		}{
			Files: uploadedFiles,
		})
	})

	passiveServer.Static("/media", UPLOAD_PATH)

	if !server.configuration.IsProductionEnvironment() {
		passiveServer.Static("/playground", PLAYGROUND_PATH)
	}

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
			handlerFunc := func(ctx echo.Context) error {
				writer := CreateHttpWriter(server, ctx, server.secureCookie)
				return handler(
					dispatcher.NewDispatcher(
						server,
						writer,
						ctx.Response().Writer,
						ctx.Request(),
						ctx.QueryParam,
						ctx.Param,
						ctx.RealIP(),
					),
				)
			}

			passiveServer.GET(path, handlerFunc)
		}(path, httpHandler.Method(), httpHandler.HandlerFunc())
	}

	for path, httpHandler := range server.httpPostHandlers {
		func(path string, method string, handler HttpHandlerFunc) {
			handlerFunc := func(ctx echo.Context) error {
				writer := CreateHttpWriter(server, ctx, server.secureCookie)
				return handler(
					dispatcher.NewDispatcher(
						server,
						writer,
						ctx.Response().Writer,
						ctx.Request(),
						ctx.QueryParam,
						ctx.Param,
						ctx.RealIP(),
					),
				)
			}

			passiveServer.POST(path, handlerFunc)
		}(path, httpHandler.Method(), httpHandler.HandlerFunc())
	}

	passiveServer.Server.ReadTimeout = time.Second * 15
	passiveServer.Server.WriteTimeout = time.Second * 15

	server.logger.SysComp(fmt.Sprintf("┄ Listening on port %d", server.passivePort))
	if err := passiveServer.Start(""); err != nil {
		// server.logger.Critical(fmt.Sprintf("PASSIVE SERVER FAILURE: %s", err))
		_ = err
	}
}
