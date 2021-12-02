//Package server has been renamed and modified from its original state, which can be found here: https://github.com/douglasmakey/admissioncontroller/blob/master/http/
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chaospuppy/imageswap/hook"

	admission "k8s.io/api/admission/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog/v2"
)

var patchType admission.PatchType = admission.PatchTypeJSONPatch

// admissionHandler represents the HTTP handler for an admission webhook
type admissionHandler struct {
	decoder runtime.Decoder
}

// newAdmissionHandler returns an instance of AdmissionHandler
func newAdmissionHandler() *admissionHandler {
	return &admissionHandler{
		decoder: serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer(),
	}
}

// Serve returns a http.HandlerFunc for an admission webhook
func (h *admissionHandler) Serve(hook hook.Hook) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprint("invalid method only POST requests are allowed"), http.StatusMethodNotAllowed)
			return
		}

		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			http.Error(w, fmt.Sprint("only content type 'application/json' is supported"), http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not read request body: %v", err), http.StatusBadRequest)
			return
		}

		var review admission.AdmissionReview
		if _, _, err := h.decoder.Decode(body, nil, &review); err != nil {
			http.Error(w, fmt.Sprintf("could not deserialize request: %v", err), http.StatusBadRequest)
			return
		}

		if review.Request == nil {
			http.Error(w, "malformed admission review: request is nil", http.StatusBadRequest)
			return
		}

		result, err := hook.Execute(review.Request)
		if err != nil {
			klog.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		admissionReview := admission.AdmissionReview{
			TypeMeta: meta.TypeMeta{
				Kind:       "AdmissionReview",
				APIVersion: "admission.k8s.io/v1",
			},
			Response: &admission.AdmissionResponse{
				UID:       review.Request.UID,
				Allowed:   result.Allowed,
				Result:    &meta.Status{Message: result.Msg},
				PatchType: &patchType,
			},
		}

		// set the patch operations for mutating admission
		if len(result.PatchOps) > 0 {
			patchBytes, err := json.Marshal(result.PatchOps)
			if err != nil {
				klog.Error(err)
				http.Error(w, fmt.Sprintf("could not marshal JSON patch: %v", err), http.StatusInternalServerError)
			}
			admissionReview.Response.Patch = patchBytes
		}

		res, err := json.Marshal(admissionReview)
		if err != nil {
			klog.Error(err)
			http.Error(w, fmt.Sprintf("could not marshal response: %v", err), http.StatusInternalServerError)
			return
		}

		klog.Infof("res: %+v", admissionReview)
		klog.Infof("Webhook [%s - %s] - Allowed: %t - PatchOps: %s", r.URL.Path, review.Request.Operation, result.Allowed, result.PatchOps)
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func healthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
