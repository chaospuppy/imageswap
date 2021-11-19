//Package server has been renamed and modified from its original state, which can be found here: https://github.com/douglasmakey/admissioncontroller/blob/master/http/
package server

import (
	"fmt"
	"imageswap/pods"
	"net/http"
)

// RunHTTPServer runs a new instance of an http.Server over TLS using the provided paths to the TLS cert and key files
func RunHTTPServer(server *http.Server, tlsKey string, tlsCert string) error {
	if err := server.ListenAndServeTLS(tlsCert, tlsKey); err != nil {
		return err
	}
	return nil
}

// NewHTTPServer returns a new instance of the *http.Server used to serve the imageswap webhook endpoints
func NewHTTPServer(port string, ecrHostname string) *http.Server {

	// Instances hooks
	podsValidation := pods.NewValidationHook()
	podsMutation := pods.NewMutationHook(ecrHostname)

	ah := newAdmissionHandler()
	mux := http.NewServeMux()
	mux.Handle("/healthz", healthz())
	mux.Handle("/validate/pods", ah.Serve(podsValidation))
	mux.Handle("/mutate", ah.Serve(podsMutation))
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}
