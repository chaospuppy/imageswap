package pods

import (
	"strings"

	"imageswap"

	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
)

func validateCreate() imageswap.AdmitFunc {
	return func(r *v1beta1.AdmissionRequest) (*imageswap.Result, error) {
		pod, err := parsePod(r.Object.Raw)
		if err != nil {
			return &imageswap.Result{Msg: err.Error()}, nil
		}

		for _, c := range pod.Spec.Containers {
			if strings.HasSuffix(c.Image, ":latest") {
				return &imageswap.Result{Msg: "You cannot use the tag 'latest' in a container."}, nil
			}
		}

		return &imageswap.Result{Allowed: true}, nil
	}
}

func mutateCreate() imageswap.AdmitFunc {
	return func(r *v1beta1.AdmissionRequest) (*imageswap.Result, error) {
		var operations []imageswap.PatchOperation
		pod, err := parsePod(r.Object.Raw)
		if err != nil {
			return &imageswap.Result{Msg: err.Error()}, nil
		}

		// Very simple logic to inject a new "sidecar" container.
		if ironbank := pod.Labels["ironbank"]; ironbank == "imageswap" {
			var containers []v1.Container
			containers = append(containers, pod.Spec.Containers...)
			sideC := v1.Container{
				Name:    "test-sidecar",
				Image:   "busybox:stable",
				Command: []string{"sh", "-c", "while true; do echo 'I am a container injected by mutating webhook'; sleep 2; done"},
			}
			containers = append(containers, sideC)
			operations = append(operations, imageswap.ReplacePatchOperation("/spec/containers", containers))
		}

		// Add a simple annotation using `AddPatchOperation`
		metadata := map[string]string{"origin": "fromMutation"}
		operations = append(operations, imageswap.AddPatchOperation("/metadata/annotations", metadata))
		return &imageswap.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}
