//Package pods has been modified from its original state, which can be found here: https://github.com/douglasmakey/admissioncontroller/blob/master/pods/
package pods

import (
	"fmt"

	"github.com/chaospuppy/imageswap/hook"

	_ "crypto/sha256" //Needed in the event the container image contains a digest
	"github.com/docker/distribution/reference"
	"k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

func mutateCreate(hostname string) hook.AdmitFunc {
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

		operations = append(createPatchOperations(pod.Spec.Containers, operations, hostname, "containers"))
		operations = append(createPatchOperations(pod.Spec.InitContainers, operations, hostname, "initContainers"))

		annotations := pod.Annotations
		if annotations == nil {
			annotations = make(map[string]string)
		}

		annotations = createPodAnnotations(allContainers, annotations)
		operations = append(operations, hook.AddPatchOperation("/metadata/annotations", annotations))

		return &hook.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}

func createPodAnnotations(containers []corev1.Container, annotations map[string]string) map[string]string {
	for i, container := range containers {
		annotations[fmt.Sprintf("imageswap.ironbank.dso.mil/%d", i)] = container.Image
	}
	return annotations
}

//createPatchOperations ingests a list of core/v1 Containers and replaces their existing images with images whose hostnames are set to the provided hostname
func createPatchOperations(containers []corev1.Container, operations []hook.PatchOperation, hostname string, containerPath string) []hook.PatchOperation {
	for i, container := range containers {
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
		operations = append(operations, hook.ReplacePatchOperation(fmt.Sprintf("/spec/%s/%d/image", containerPath, i), hostname+"/"+reference.Path(named)+tag+digest))
	}
	return operations
}
