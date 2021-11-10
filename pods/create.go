package pods

import (
	"strings"

	"imageswap/hook"
	"imageswap/util"

	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	"regexp"
)

var ecrPattern = regexp.MustCompile(`^(\d{12})\.dkr\.ecr(\-fips)?\.([a-zA-Z0-9][a-zA-Z0-9-_]*)\.(amazonaws\.com(\.cn)?|sc2s\.sgov\.gov|c2s\.ic\.gov)$`)
var registryPattern = regexp.MustCompile(``)

func validateCreate() hook.AdmitFunc {
	return func(r *v1beta1.AdmissionRequest) (*hook.Result, error) {
		pod, err := parsePod(r.Object.Raw)
		if err != nil {
			return &hook.Result{Msg: err.Error()}, nil
		}

		for _, c := range pod.Spec.Containers {
			if strings.HasSuffix(c.Image, ":latest") {
				return &hook.Result{Msg: "You cannot use the tag 'latest' in a container."}, nil
			}
		}

		return &hook.Result{Allowed: true}, nil
	}
}

func mutateCreate() hook.AdmitFunc {
	return func(r *v1beta1.AdmissionRequest) (*hook.Result, error) {
		var operations []hook.PatchOperation
		pod, err := parsePod(r.Object.Raw)
		if err != nil {
			return &hook.Result{Msg: err.Error()}, nil
		}

		for _, c := range pod.Spec.Containers {
			registry := strings.Split(c.Image, "/")[0]
			util.CheckForRegistry(registry)
			operations = append(operations, hook.ReplacePatchOperation("/spec/containers/image", c))
		}
		sideC := v1.Container{
			Name:    "test-sidecar",
			Image:   "busybox:stable",
			Command: []string{"sh", "-c", "while true; do echo 'I am a container injected by mutating webhook'; sleep 2; done"},
		}

		// Add a simple annotation using `AddPatchOperation`
		metadata := map[string]string{"origin": "fromMutation"}
		operations = append(operations, hook.AddPatchOperation("/metadata/annotations", metadata))
		return &hook.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}
