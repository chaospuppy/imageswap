package pods

import (
	"strings"

	"imageswap/hook"
	"imageswap/util"

	"k8s.io/api/admission/v1beta1"
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

func mutateCreate(ecrHostname string) hook.AdmitFunc {
	return func(r *v1beta1.AdmissionRequest) (*hook.Result, error) {
		var operations []hook.PatchOperation
		pod, err := parsePod(r.Object.Raw)
		if err != nil {
			return &hook.Result{Msg: err.Error()}, nil
		}

		for _, c := range pod.Spec.Containers {
			splitImage := strings.Split(c.Image, "/")
			registry := splitImage[0]
			imagePath := splitImage[1:]
			if util.CheckForRegistry(registry) {
				// Replace existing registry with ecr registry
				operations = append(operations, hook.ReplacePatchOperation("/spec/containers/image", ecrHostname+strings.Join(imagePath[:], ",")))
			} else {
				operations = append(operations, hook.ReplacePatchOperation("/spec/containers/image", ecrHostname+c.Image))
			}
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
