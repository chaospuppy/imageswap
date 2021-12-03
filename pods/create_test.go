package pods

import (
	"github.com/chaospuppy/imageswap/hook"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

const expectedRegistry = "expected.registry"

type TestCasePatchOperations struct {
	ExpectedContainerResult interface{}
	Container               v1.Container
	InitContainer           v1.Container
}

func TestCreatePatchOperations(t *testing.T) {
	assert := assert.New(t)
	testCasesPatchOperations := []TestCasePatchOperations{
		{
			ExpectedContainerResult: hook.PatchOperation{Op: "replace", Path: "/spec/containers/0/image", From: "", Value: "expected.registry/foo"},
			Container: v1.Container{
				Name:  "testcontainer0",
				Image: "registry.example.com/foo",
			},
			InitContainer: v1.Container{
				Name:  "testcontainer0",
				Image: "registry.example.com/foo",
			},
		},
		{
			ExpectedContainerResult: hook.PatchOperation{Op: "replace", Path: "/spec/containers/1/image", From: "", Value: "expected.registry/foo/bar"},
			Container: v1.Container{
				Name:  "testcontainer1",
				Image: "registry.example.com/foo/bar",
			},
			InitContainer: v1.Container{
				Name:  "testcontainer1",
				Image: "registry.example.com/foo/bar",
			},
		},
		{
			ExpectedContainerResult: hook.PatchOperation{Op: "replace", Path: "/spec/containers/2/image", From: "", Value: "expected.registry/foo/bar:1.2.3"},
			Container: v1.Container{
				Name:  "testcontainer2",
				Image: "registry.example.com/foo/bar:1.2.3",
			},
			InitContainer: v1.Container{
				Name:  "testcontainer2",
				Image: "registry.example.com/foo/bar:1.2.3",
			},
		},
		{
			ExpectedContainerResult: hook.PatchOperation{Op: "replace", Path: "/spec/containers/3/image", From: "", Value: "expected.registry/library/foo"},
			Container: v1.Container{
				Name:  "testcontainer3",
				Image: "foo",
			},
			InitContainer: v1.Container{
				Name:  "testcontainer3",
				Image: "foo",
			},
		},
		{
			ExpectedContainerResult: hook.PatchOperation{Op: "replace", Path: "/spec/containers/4/image", From: "", Value: "expected.registry/foo/bar"},
			Container: v1.Container{
				Name:  "testcontainer4",
				Image: "foo/bar",
			},
			InitContainer: v1.Container{
				Name:  "testcontainer4",
				Image: "foo/bar",
			},
		},
		{
			ExpectedContainerResult: hook.PatchOperation{Op: "replace", Path: "/spec/containers/5/image", From: "", Value: "expected.registry/foo"},
			Container: v1.Container{
				Name:  "testcontainer5",
				Image: "localhost/foo",
			},
			InitContainer: v1.Container{
				Name:  "testcontainer5",
				Image: "localhost/foo",
			},
		},
		{
			ExpectedContainerResult: hook.PatchOperation{Op: "replace", Path: "/spec/containers/6/image", From: "", Value: "expected.registry/foo"},
			Container: v1.Container{
				Name:  "testcontainer6",
				Image: "localhost:5000/foo",
			},
			InitContainer: v1.Container{
				Name:  "testcontainer6",
				Image: "localhost:5000/foo",
			},
		},
		{
			ExpectedContainerResult: hook.PatchOperation{Op: "replace", Path: "/spec/containers/7/image", From: "", Value: "expected.registry/foo"},
			Container: v1.Container{
				Name:  "testcontainer7",
				Image: "172.0.0.1/foo",
			},
			InitContainer: v1.Container{
				Name:  "testcontainer7",
				Image: "172.0.0.1/foo",
			},
		},
		{
			ExpectedContainerResult: hook.PatchOperation{Op: "replace", Path: "/spec/containers/8/image", From: "", Value: "expected.registry/foo/bar@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2"},
			Container: v1.Container{
				Name:  "testcontainer8",
				Image: "registry.example.com/foo/bar@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2",
			},
			InitContainer: v1.Container{
				Name:  "testcontainer8",
				Image: "registry.example.com/foo/bar@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2",
			},
		},
	}
	for i, testCase := range testCasesPatchOperations {
		assert.Equal(testCase.ExpectedContainerResult, patchContainerImage(testCase.Container, expectedRegistry, false, i))
	}
}
