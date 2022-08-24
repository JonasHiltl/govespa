package integration

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jonashiltl/govespa"
	"golang.org/x/net/http2"
)

func createClient() (*http.Client, error) {
	key, err := ioutil.ReadFile("../../../clubo-app/vespa-config/pki/client/client.key")
	if err != nil {
		return nil, err
	}

	crt, err := ioutil.ReadFile("../../../clubo-app/vespa-config/pki/client/client.pem")
	if err != nil {
		return nil, err
	}

	ca, err := ioutil.ReadFile("../../../clubo-app/vespa-config/pki/vespa/ca-vespa.pem")
	if err != nil {
		return nil, err
	}
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(ca)

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil, err
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

	return &http.Client{Transport: trns}, nil

}

func TestPut(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Fatal(err)
	}

	c := govespa.NewClient(govespa.NewClientParams{
		HttpClient: client,
		BaseUrl:    "https://localhost:8090",
	})

	uname := "john.doe"
	fname := "John"
	lname := "Doe"

	err = c.
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
	if err != nil {
		t.Error(err)
	}
}
