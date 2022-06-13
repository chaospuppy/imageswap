//Package server has been renamed and modified from its original state, which can be found here: https://github.com/douglasmakey/admissioncontroller/blob/master/http/
package server

import (
	"errors"
	"net/http"

	"k8s.io/klog/v2"
)

// RunHTTPServer runs a new instance of an http.Server over TLS using the provided paths to the TLS cert and key files
func RunHTTPServer(server *http.Server) {
	if err := server.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
		klog.Fatalf("Failed to listen and serve: %v", err)
	}
}
