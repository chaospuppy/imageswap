package server

import (
	"fmt"
	"imageswap/pods"
	"net/http"
)

func RunHTTPServer(server *http.Server, tlsKey string, tlsCert string) error {
	if err := server.ListenAndServeTLS(tlsCert, tlsKey); err != nil {
		return err
	}
	return nil
}

func NewHTTPServer(port string, ecrHostname string) *http.Server {

	// Instances hooks
	podsValidation := pods.NewValidationHook()
	podsMutation := pods.NewMutationHook(ecrHostname)

	ah := newAdmissionHandler()
	mux := http.NewServeMux()
	mux.Handle("/healthz", healthz())
	mux.Handle("/validate/pods", ah.Serve(podsValidation))
	mux.Handle("/mutate/pods", ah.Serve(podsMutation))
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}
