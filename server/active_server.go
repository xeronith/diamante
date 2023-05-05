package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	. "github.com/xeronith/diamante/actor"
	. "github.com/xeronith/diamante/contracts/actor"
	"github.com/xeronith/diamante/contracts/operation"
	. "github.com/xeronith/diamante/io"
)

func (server *defaultServer) startActiveServer() {
	listener, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", server.activePort))
	if err != nil {
		log.Fatalf("ACTIVE SERVER LISTENER FATAL ERROR: %s", err)
	}

	server.listeners.Append(listener)

	upgrader := websocket.Upgrader{
		// ReadBufferSize: 1024,
		// WriteBufferSize: 1024,
		CheckOrigin: func(request *http.Request) bool {
			return true
		},
	}

	handler := func(context echo.Context) error {
		connection, err := upgrader.Upgrade(context.Response(), context.Request(), nil)
		if err != nil {
			server.logger.Error(fmt.Sprintf("SOCKET UPGRADE: %s", err))
			return err
		}

		var actor IActor
		writer := CreateWebSocketWriter(server, connection, func() {
			server.OnSocketDisconnected(actor)
		})

		actor = CreateActor(writer, true, context.RealIP(), context.Request().UserAgent())

		defer writer.Close()
		server.OnSocketConnected(actor)

		for {
			messageType, message, err := connection.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
					server.logger.Error(fmt.Sprintf("SERVER SOCKET READ ERROR {%s}: %s", actor.Token(), err))
				}

				return nil
			} else {
				switch messageType {
				case websocket.BinaryMessage, websocket.TextMessage:
					var result operation.IOperationResult
					result = server.OnActorBinaryData(actor, message)
					actor.Dispatch(result)
				default:
					server.logger.Error(fmt.Sprintf("UNSUPPORTED SOCKET MESSAGE TYPE: %d", messageType))
				}
			}
		}
	}

	activeServer := echo.New()
	activeServer.HideBanner = true
	activeServer.Debug = false
	activeServer.HidePort = true
	// activeServer.Use(middleware.Logger())
	// activeServer.Logger.SetLevel(labstackLog.OFF)
	// activeServer.Use(middleware.Recover())

	if server.Configuration().GetServerConfiguration().GetTLSConfiguration().IsEnabled() {
		activeServer.Listener = server.createTLSListener(listener)
	} else {
		activeServer.Listener = listener
	}

	activeServer.Server.ReadTimeout = time.Second * 15
	activeServer.Server.WriteTimeout = time.Second * 15

	activeServer.GET("/", handler)

	if err := activeServer.Start(""); err != nil {
		// server.logger.Critical(fmt.Sprintf("ACTIVE SERVER FAILURE: %s", err))
		_ = err
	}
}
