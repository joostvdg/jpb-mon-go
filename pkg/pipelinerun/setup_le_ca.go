package pipelinerun

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	localCertFile = "/usr/src/letsencryptauthorityx3.pem"
)

func GetTransportWithLetsEncryptRootCA() *http.Transport {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := ioutil.ReadFile(localCertFile)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", localCertFile, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		RootCAs: rootCAs,
	}
	return &http.Transport{TLSClientConfig: config}
}
