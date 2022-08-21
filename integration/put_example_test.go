package integration

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/clubo-app/govespa"
	"golang.org/x/net/http2"
)

func createClient() *http.Client {
	key, err := ioutil.ReadFile("../../vespa-config/pki/client/client.key")
	if err != nil {
		log.Println(err)
		return nil
	}

	crt, err := ioutil.ReadFile("../../vespa-config/pki/client/client.pem")
	if err != nil {
		log.Println(err)
		return nil
	}

	ca, err := ioutil.ReadFile("../../vespa-config/pki/vespa/ca-vespa.pem")
	if err != nil {
		log.Println(err)
		return nil
	}
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(ca)

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Println(err)
		return nil
	}
	tls := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCAs,
		ServerName:   "localhost",
	}

	trns := &http2.Transport{
		TLSClientConfig: tls,
		AllowHTTP:       false,
	}

	return &http.Client{Transport: trns}

}

func TestPut(t *testing.T) {
	client := createClient()
	if client == nil {
		t.Fatal("Error creating Http Client")
	}

	c := govespa.NewClient(govespa.NewClientParams{
		HttpClient: client,
		BaseUrl:    "https://localhost:8090",
	})

	uname := "john.doe"
	fname := "John"
	lname := "Doe"

	apiErr := c.
		Put(govespa.DocumentId{
			Namespace:    "default",
			DocType:      "user",
			UserSpecific: "awgaw-1w234a-dw14ag-w1414a",
		}).
		WithContext(context.Background()).
		BindMap(map[string]any{
			"username":  uname,
			"firstname": fname,
			"lastname":  lname,
		}).
		Exec()
	if apiErr != nil {
		t.Error(apiErr)
	}
}
