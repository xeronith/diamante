package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func (server *defaultServer) createTLSListener(listener net.Listener) net.Listener {
	tlsConfiguration := server.Configuration().GetServerConfiguration().GetTLSConfiguration()

	if !tlsConfiguration.IsEnabled() {
		server.Logger().Warning(fmt.Sprintf("WARNING: SSL certificate or key not provided. TLS disabled for: %s", listener.Addr()))
		return listener
	}

	config := new(tls.Config)
	config.MinVersion = tls.VersionTLS10
	// config.SessionTicketsDisabled = true
	config.PreferServerCipherSuites = true
	config.NextProtos = append(config.NextProtos, "h2")
	/*config.CipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	}*/

	certFile := tlsConfiguration.GetCertFile()
	keyFile := tlsConfiguration.GetKeyFile()

	config.Certificates = make([]tls.Certificate, 1)
	var err error
	if config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		log.Fatal(err)
	}

	config.BuildNameToCertificate()
	return tls.NewListener(listener, config)
}
