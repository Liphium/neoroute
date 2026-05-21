package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/Liphium/neoroute"
	"github.com/quic-go/quic-go/http3"
)

//go:generate msgp
type Request struct {
	Field1 string `msg:"field1"`
	Field2 int    `msg:"field2"`
}

//go:generate msgp
type Response struct {
	Field1 string `msg:"field1"`
	Field2 int    `msg:"field2"`
}

func main() {

	// Create HTTP transporter
	hook, t := neoroute.NewHTTPTransporter()

	// Create router and set it for transporter
	router := neoroute.NewNeoRouter(neoroute.Config{
		ErrorHandler: func(err error) string {
			return fmt.Sprintf("error: %v", err)
		},
	})
	t.SetRouter(router)

	// Route: simple.route
	neoroute.Route(router, "simple.route", func(c *neoroute.Ctx[Response], req Request) error {
		return c.Respond(Response{Field1: "simple response", Field2: 68})
	})

	// Create group for group1
	group1 := router.Group("group1")

	// Create subroute for group1
	// Route: group1.route1
	neoroute.Route(group1, "route1", func(c *neoroute.Ctx[Response], req Request) error {
		return c.Respond(Response{
			Field1: "response to " + req.Field1,
			Field2: req.Field2 + 1,
		})
	})

	// Create group2 on top of group1
	group2 := group1.Group("group2")

	// Create subroute for group2
	// Route: group1.group2.route1
	neoroute.Route(group2, "route1", func(c *neoroute.Ctx[Response], req Request) error {
		return c.Respond(Response{
			Field1: "response to " + req.Field1,
			Field2: req.Field2 + 2,
		})
	})

	// Create server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		_, _ = w.Write([]byte("hello over HTTP/3\n"))
	})

	// Hook transporter into /neo route
	mux.HandleFunc("/neo", hook)

	// Load TLS certificate and serve server over HTTP/3
	cert, err := selfSignedCert()
	if err != nil {
		log.Fatal(err)
	}

	server := &http3.Server{
		Addr:    ":6121",
		Handler: mux,
		TLSConfig: http3.ConfigureTLSConfig(&tls.Config{
			Certificates: []tls.Certificate{cert},
		}),
	}

	log.Println("listening on https://localhost:6121 over HTTP/3")
	log.Fatal(server.ListenAndServe())
}

func selfSignedCert() (tls.Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Organization: []string{"neoroute"},
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return tls.X509KeyPair(certPEM, keyPEM)
}
