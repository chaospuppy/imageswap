package http

import "net/http"

func newHTTPServer(port string) *http.Server {

	mux := http.NewServeMux()
	mux.Handle("/healthz", healthz())
	mux.Handle("/validate/pods", ah.Serve(podsValidation))
	mux.Handle("/mutate/pods", ah.Serve(podsMutation))
	mux.Handle("/validate/deployments", ah.Serve(deploymentValidation))
	mux := http.NewServeMux()
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}
