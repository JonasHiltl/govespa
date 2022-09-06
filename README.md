## govespa 

Govespa is a client implementation for the [vespa-engine](https://github.com/vespa-engine). It uses the Document Api and can be used over http or http/2, since you can supply your own `http.Client`. It's goal is to support all functionalities exposed by the Document Api.

#### Features
- `Put`: writes a Document by ID and allows binding a `struct`/`map` to the document fields. 
- `Get`: returns a Document by ID and allows mapping the Response to a struct.
- `Update`: updates fields of a Document by ID.
- `Remove`: removes a Document by ID.
- `Query`: executes yql and allows binding the Response to a struct.

#### Getting Started
```
go get -u github.com/jonashiltl/govespa
```
Create an `http.Client`, for http/2 add your TLS Certificates.
```go
key, err := ioutil.ReadFile("client.key")
...

crt, err := ioutil.ReadFile("client.pem")
...

ca, err := ioutil.ReadFile("ca-vespa.pem")
...

rootCAs := x509.NewCertPool()
rootCAs.AppendCertsFromPEM(ca)

cert, err := tls.X509KeyPair(crt, key)
... 

tls := &tls.Config{
  Certificates: []tls.Certificate{cert},
  RootCAs:      rootCAs,
  ServerName:   "localhost",
}

trns := &http2.Transport{
  TLSClientConfig: tls,
  AllowHTTP:       false,
}

client := &http.Client{Transport: trns}
```
Create the VespaClient with a `baseUrl` and the `http.Client`.
```go
c := govespa.NewClient(govespa.NewClientParams{
  HttpClient: client,
  BaseUrl:    "https://localhost:8090",
})
```

#### Examples
Have a look at the [integration](integration/) folder to see examples.

#### TODO
- [ ] Extend the DocumentId with the key/value pair section
- [ ] Decide where and how to use concurrency 
- [ ] Improve the `ParseDocId` function, current benchmark can be seen in [documentid_test.go](documentid_test.go)
- [ ] Implement an Iterator which can be used for "update where", "visit", "delete where". It can either be a single struct used by Update/Get/Delete/Query but that may be harder to implement than just defining the logic for every instance separately.

#### Disclaimer 
> This client implementation is work in progress and not an official client created by the Team [@vespa](https://github.com/vespa-engine).
