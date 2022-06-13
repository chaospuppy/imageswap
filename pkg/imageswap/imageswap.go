package imageswap

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/chaospuppy/imageswap/pkg/pods"
	imageswapserver "github.com/chaospuppy/imageswap/pkg/server"
	"k8s.io/klog/v2"
)

type Imageswap struct {
	TlsKey           string
	TlsCert          string
	HttpPort         string
	Annotation       string
	RegistryHostname string
}

// NewHTTPServer returns a new instance of the *http.Server used to serve the imageswap webhook endpoints
func (c *Imageswap) NewHTTPServer() *http.Server {
	kpr, err := imageswapserver.NewKeypairReloader(c.TlsCert, c.TlsKey)
	if err != nil {
		klog.Fatal(err)
	}

	tlsConfig := &tls.Config{}
	tlsConfig.GetCertificate = kpr.GetCertificateFunc()

	// Instances hooks
	podsMutation := pods.NewMutationHook(c.RegistryHostname, c.Annotation)
	ah := imageswapserver.NewAdmissionHandler()
	mux := http.NewServeMux()

	mux.Handle("/healthz", imageswapserver.Healthz())
	mux.Handle("/mutate", ah.Serve(podsMutation))
	return &http.Server{
		Addr:      fmt.Sprintf(":%s", c.HttpPort),
		Handler:   mux,
		TLSConfig: tlsConfig,
	}
}
