package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (server *defaultServer) startDiagnosticsServer() {
	if server.hudEnabled {
		go server.hud()
	}

	tlsConfiguration := server.Configuration().GetServerConfiguration().GetTLSConfiguration()

	if tlsConfiguration.IsEnabled() {
		certFile := tlsConfiguration.GetCertFile()
		keyFile := tlsConfiguration.GetKeyFile()
		log.Println(http.ListenAndServeTLS(fmt.Sprintf("0.0.0.0:%d", server.diagnosticsPort), certFile, keyFile, nil))
	} else {
		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", server.diagnosticsPort), nil))
	}
}

func (server *defaultServer) hud() {
	for {
		time.Sleep(5000000000)
		fmt.Printf("CON:%010d ", server.ActorsCount())
		fmt.Printf("GRT:%010d\n", runtime.NumGoroutine())
	}
}

func diagnosticsHandler(context echo.Context) error {
	if context.Request().Header.Get("Authorization") != "Bearer a09301fc-911b-4c8a-b0b9-8b69b7409716" {
		return context.String(http.StatusUnauthorized, "")
	}

	clientType := strings.ToLower(context.Param("clientType"))

	file, err := context.FormFile("file")
	if err != nil {
		return err
	}

	if file.Size > 102400 {
		return errors.New("file too big")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer func() { _ = src.Close() }()

	var path string
	switch clientType {
	case "android", "ios":
		path = fmt.Sprintf("./diagnostics/%s/%s", clientType, file.Filename)
	default:
		return errors.New("invalid client")
	}

	_ = os.MkdirAll("./diagnostics/ios/", os.ModePerm)
	_ = os.MkdirAll("./diagnostics/android/", os.ModePerm)

	dst, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() { _ = dst.Close() }()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return context.String(http.StatusOK, "")
}
