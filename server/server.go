//Package server has been renamed and modified from its original state, which can be found here: https://github.com/douglasmakey/admissioncontroller/blob/master/http/
package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/chaospuppy/imageswap/pods"
	"k8s.io/klog/v2"
	"net/http"
	"sync"
	"time"
)

type keypairReloader struct {
	certMu   sync.RWMutex
	cert     *tls.Certificate
	certPath string
	keyPath  string
}

func newKeypairReloader(certPath, keyPath string) (*keypairReloader, error) {
	result := &keypairReloader{
		certPath: certPath,
		keyPath:  keyPath,
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	result.cert = &cert
	go func() {
		// TODO: Consider replacing with syscall.SIGHUP channel - send upon updates to TLS secret
		ticker := time.NewTicker(60 * time.Second)
		for range ticker.C {
			klog.Infof("Reloading TLS certificate and key from %s and %s", certPath, keyPath)
			if err := result.maybeReload(); err != nil {
				klog.Infof("Keeping old TLS certificate because the new one could not be loaded: %v", err)
			}
		}
	}()
	return result, nil
}

func (kpr *keypairReloader) maybeReload() error {
	newCert, err := tls.LoadX509KeyPair(kpr.certPath, kpr.keyPath)
	if err != nil {
		return err
	}
	kpr.certMu.Lock()
	defer kpr.certMu.Unlock()
	kpr.cert = &newCert
	return nil
}

func (kpr *keypairReloader) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		klog.Info("Received ClientHello!")
		kpr.certMu.RLock()
		defer kpr.certMu.RUnlock()
		return kpr.cert, nil
	}
}

// RunHTTPServer runs a new instance of an http.Server over TLS using the provided paths to the TLS cert and key files
func RunHTTPServer(server *http.Server) {
	if err := server.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
		klog.Fatalf("Failed to listen and serve: %v", err)
	}
}

// NewHTTPServer returns a new instance of the *http.Server used to serve the imageswap webhook endpoints
func NewHTTPServer(port, registryHostname, tlsCert, tlsKey string) *http.Server {
	kpr, err := newKeypairReloader(tlsCert, tlsKey)
	if err != nil {
		klog.Fatal(err)
	}

	tlsConfig := &tls.Config{}
	tlsConfig.GetCertificate = kpr.GetCertificateFunc()

	// Instances hooks
	podsMutation := pods.NewMutationHook(registryHostname)
	ah := newAdmissionHandler()
	mux := http.NewServeMux()

	mux.Handle("/healthz", healthz())
	mux.Handle("/mutate", ah.Serve(podsMutation))
	return &http.Server{
		Addr:      fmt.Sprintf(":%s", port),
		Handler:   mux,
		TLSConfig: tlsConfig,
	}
}
