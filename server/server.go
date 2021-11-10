package server

import (
	"imageswap/pods"
	"net/http"
)

func NewHTTPServer(port string) *http.Server {

	// Instances hooks
	podsValidation := pods.NewValidationHook()
	podsMutation := pods.NewMutationHook()

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthz())
	mux.Handle("/validate/pods", ah.Serve(podsValidation))
	mux.Handle("/mutate/pods", ah.Serve(podsMutation))
	mux := http.NewServeMux()
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}
