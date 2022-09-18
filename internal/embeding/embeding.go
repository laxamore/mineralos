package embeding

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

//go:embed cert/*
var f embed.FS

func LoadClientTLSCert(embedCert bool, certPath string) (cred credentials.TransportCredentials, err error) {
	var b []byte
	if embedCert {
		b, err = f.ReadFile("cert/config-public-cert.pem")
		if err != nil {
			return
		}
	} else {
		b, err = ioutil.ReadFile(certPath)
		if err != nil {
			return
		}
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		err = fmt.Errorf("credentials: failed to append certificates")
		return
	}

	cred, err = credentials.NewTLS(&tls.Config{
		RootCAs: cp,
	}), nil
	return
}

func LoadServerTLSCert() (credentials.TransportCredentials, error) {
	certPEMBlock, err := f.ReadFile("cert/config-public-cert.pem")
	if err != nil {
		return nil, err
	}
	keyPEMBlock, err := f.ReadFile("cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	}), nil
}
