/*
Copyright (c) 2016-2017 Bitnami

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
package fake

import (
	v1beta1 "github.com/epimorphics/s3-trigger/pkg/apis/kubeless/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeS3Triggers implements S3TriggerInterface
type FakeS3Triggers struct {
	Fake *FakeKubelessV1beta1
	ns   string
}

var s3triggersResource = schema.GroupVersionResource{Group: "kubeless.io", Version: "v1beta1", Resource: "s3triggers"}

var s3triggersKind = schema.GroupVersionKind{Group: "kubeless.io", Version: "v1beta1", Kind: "S3Trigger"}

// Get takes name of the s3Trigger, and returns the corresponding s3Trigger object, and an error if there is any.
func (c *FakeS3Triggers) Get(name string, options v1.GetOptions) (result *v1beta1.S3Trigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(s3triggersResource, c.ns, name), &v1beta1.S3Trigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.S3Trigger), err
}

// List takes label and field selectors, and returns the list of S3Triggers that match those selectors.
func (c *FakeS3Triggers) List(opts v1.ListOptions) (result *v1beta1.S3TriggerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(s3triggersResource, s3triggersKind, c.ns, opts), &v1beta1.S3TriggerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.S3TriggerList{}
	for _, item := range obj.(*v1beta1.S3TriggerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested s3Triggers.
func (c *FakeS3Triggers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(s3triggersResource, c.ns, opts))

}

// Create takes the representation of a s3Trigger and creates it.  Returns the server's representation of the s3Trigger, and an error, if there is any.
func (c *FakeS3Triggers) Create(s3Trigger *v1beta1.S3Trigger) (result *v1beta1.S3Trigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(s3triggersResource, c.ns, s3Trigger), &v1beta1.S3Trigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.S3Trigger), err
}

// Update takes the representation of a s3Trigger and updates it. Returns the server's representation of the s3Trigger, and an error, if there is any.
func (c *FakeS3Triggers) Update(s3Trigger *v1beta1.S3Trigger) (result *v1beta1.S3Trigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(s3triggersResource, c.ns, s3Trigger), &v1beta1.S3Trigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.S3Trigger), err
}

// Delete takes name of the s3Trigger and deletes it. Returns an error if one occurs.
func (c *FakeS3Triggers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(s3triggersResource, c.ns, name), &v1beta1.S3Trigger{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeS3Triggers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(s3triggersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.S3TriggerList{})
	return err
}

// Patch applies the patch and returns the patched s3Trigger.
func (c *FakeS3Triggers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.S3Trigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(s3triggersResource, c.ns, name, data, subresources...), &v1beta1.S3Trigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.S3Trigger), err
}
