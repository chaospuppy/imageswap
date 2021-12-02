//Package pods has been modified from its original state, which can be found here: https://github.com/douglasmakey/admissioncontroller/blob/master/pods/
package pods

import (
	"encoding/json"
	"github.com/chaospuppy/imageswap/hook"
	v1 "k8s.io/api/core/v1"
)

// NewMutationHook creates a new instance of pods mutation hook
func NewMutationHook(hostname string) hook.Hook {
	return hook.Hook{
		Create: mutateCreate(hostname),
	}
}

func parsePod(object []byte) (*v1.Pod, error) {
	var pod v1.Pod
	if err := json.Unmarshal(object, &pod); err != nil {
		return nil, err
	}

	return &pod, nil
}
