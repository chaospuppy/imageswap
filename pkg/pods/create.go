//Package pods has been modified from its original state, which can be found here: https://github.com/douglasmakey/admissioncontroller/blob/master/pods/
package pods

import (
	"fmt"

	"github.com/chaospuppy/imageswap/pkg/hook"

	_ "crypto/sha256" //Needed in the event the container image contains a digest

	"github.com/docker/distribution/reference"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

func mutateCreate(hostname string, annotation string) hook.AdmitFunc {
	return func(r *v1.AdmissionRequest) (*hook.Result, error) {
		var operations []hook.PatchOperation
		pod, err := parsePod(r.Object.Raw)
		if err != nil {
			return &hook.Result{Msg: err.Error()}, nil
		}

		containersList := [][]corev1.Container{
			pod.Spec.Containers,
			pod.Spec.InitContainers,
		}

		var allContainers []corev1.Container
		for _, s := range containersList {
			allContainers = append(allContainers, s...)
		}

		for i, container := range pod.Spec.Containers {
			operations = append(operations, patchContainerImage(container, hostname, false, i))
		}
		for i, container := range pod.Spec.InitContainers {
			operations = append(operations, patchContainerImage(container, hostname, true, i))
		}

		annotations := pod.Annotations
		if annotations == nil {
			annotations = make(map[string]string)
		}
		for i, container := range allContainers {
			annotations[fmt.Sprintf("%s/%d", annotation, i)] = container.Image
		}
		operations = append(operations, hook.AddPatchOperation("/metadata/annotations", annotations))

		return &hook.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}

func patchContainerImage(container corev1.Container, hostname string, isInit bool, index int) hook.PatchOperation {
	named, err := reference.ParseNormalizedNamed(container.Image)
	if err != nil {
		klog.Fatal(err)
	}
	var digest string
	var tag string
	// Check for valid digest
	if digested, ok := named.(reference.Digested); ok {
		digest = "@" + digested.Digest().String()
	}
	// Check for valid tag
	if tagged, ok := named.(reference.Tagged); ok {
		tag = ":" + tagged.Tag()
	}

	var operation hook.PatchOperation

	if isInit {
		operation = hook.ReplacePatchOperation(fmt.Sprintf("/spec/initContainers/%d/image", index), hostname+"/"+reference.Path(named)+tag+digest)
	} else {
		operation = hook.ReplacePatchOperation(fmt.Sprintf("/spec/containers/%d/image", index), hostname+"/"+reference.Path(named)+tag+digest)
	}

	return operation
}
