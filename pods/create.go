//Package pods has been modified from its original state, which can be found here: https://github.com/douglasmakey/admissioncontroller/blob/master/pods/
package pods

import (
	"fmt"
	"strings"

	"imageswap/hook"

	_ "crypto/sha256" //Needed in the event the container image contains a digest
	"github.com/docker/distribution/reference"
	"k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

func validateCreate() hook.AdmitFunc {
	return func(r *v1.AdmissionRequest) (*hook.Result, error) {
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
	return func(r *v1.AdmissionRequest) (*hook.Result, error) {
		var operations []hook.PatchOperation
		pod, err := parsePod(r.Object.Raw)
		if err != nil {
			return &hook.Result{Msg: err.Error()}, nil
		}

		operations = append(createPatchOperations(pod.Spec.Containers, operations, ecrHostname, "containers"))
		operations = append(createPatchOperations(pod.Spec.InitContainers, operations, ecrHostname, "initContainers"))
		return &hook.Result{
			Allowed:  true,
			PatchOps: operations,
		}, nil
	}
}

//createPatchOperations ingests a list of core/v1 Containers and replaces their existing images with images whose hostnames are set to the provided ecrHostname
func createPatchOperations(containers []corev1.Container, operations []hook.PatchOperation, ecrHostname string, containerPath string) []hook.PatchOperation {
	for i, container := range containers {
		named, err := reference.ParseNormalizedNamed(container.Image)
		if err != nil {
			klog.Fatal(err)
		}
		// Ensure latest tag is appended if tag/digest was left blank
		named = reference.TagNameOnly(named)
		// Check for Tag
		if tagged, ok := named.(reference.Tagged); ok {
			operations = append(operations, hook.ReplacePatchOperation(fmt.Sprintf("/spec/%s/%d/image", containerPath, i), ecrHostname+"/"+reference.Path(named)+":"+tagged.Tag()))
			// Check for Digest
		} else if digested, ok := named.(reference.Digested); ok {
			operations = append(operations, hook.ReplacePatchOperation(fmt.Sprintf("/spec/%s/%d/image", containerPath, i), ecrHostname+"/"+reference.Path(named)+"@"+digested.Digest().String()))
			// Fail
		} else {
			klog.Fatalf("Invalid tag/digest format: %s", container.Image)
		}
		// Add annotations indicating original image value
		operations = append(operations, hook.AddPatchOperation(fmt.Sprintf("/metadata/annotations/imageswap.ironbank.dso.mil~1%sOriginalImage", container.Name), container.Image))
	}
	return operations
}
