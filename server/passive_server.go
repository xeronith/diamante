package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
