/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/MiniTeks/mks-server/pkg/apis/mkscontroller/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMksTasks implements MksTaskInterface
type FakeMksTasks struct {
	Fake *FakeMkscontrollerV1alpha1
	ns   string
}

var mkstasksResource = schema.GroupVersionResource{Group: "mkscontroller.example.mks", Version: "v1alpha1", Resource: "mkstasks"}

var mkstasksKind = schema.GroupVersionKind{Group: "mkscontroller.example.mks", Version: "v1alpha1", Kind: "MksTask"}

// Get takes name of the mksTask, and returns the corresponding mksTask object, and an error if there is any.
func (c *FakeMksTasks) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.MksTask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mkstasksResource, c.ns, name), &v1alpha1.MksTask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MksTask), err
}

// List takes label and field selectors, and returns the list of MksTasks that match those selectors.
func (c *FakeMksTasks) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.MksTaskList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mkstasksResource, mkstasksKind, c.ns, opts), &v1alpha1.MksTaskList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MksTaskList{ListMeta: obj.(*v1alpha1.MksTaskList).ListMeta}
	for _, item := range obj.(*v1alpha1.MksTaskList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested mksTasks.
func (c *FakeMksTasks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mkstasksResource, c.ns, opts))

}

// Create takes the representation of a mksTask and creates it.  Returns the server's representation of the mksTask, and an error, if there is any.
func (c *FakeMksTasks) Create(ctx context.Context, mksTask *v1alpha1.MksTask, opts v1.CreateOptions) (result *v1alpha1.MksTask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mkstasksResource, c.ns, mksTask), &v1alpha1.MksTask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MksTask), err
}

// Update takes the representation of a mksTask and updates it. Returns the server's representation of the mksTask, and an error, if there is any.
func (c *FakeMksTasks) Update(ctx context.Context, mksTask *v1alpha1.MksTask, opts v1.UpdateOptions) (result *v1alpha1.MksTask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mkstasksResource, c.ns, mksTask), &v1alpha1.MksTask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MksTask), err
}

// Delete takes name of the mksTask and deletes it. Returns an error if one occurs.
func (c *FakeMksTasks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(mkstasksResource, c.ns, name, opts), &v1alpha1.MksTask{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMksTasks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mkstasksResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.MksTaskList{})
	return err
}

// Patch applies the patch and returns the patched mksTask.
func (c *FakeMksTasks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MksTask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mkstasksResource, c.ns, name, pt, data, subresources...), &v1alpha1.MksTask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MksTask), err
}
