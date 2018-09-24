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
package v1beta1

import (
	v1beta1 "github.com/epimorphics/s3-trigger/pkg/apis/kubeless/v1beta1"
	scheme "github.com/epimorphics/s3-trigger/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// S3TriggersGetter has a method to return a S3TriggerInterface.
// A group's client should implement this interface.
type S3TriggersGetter interface {
	S3Triggers(namespace string) S3TriggerInterface
}

// S3TriggerInterface has methods to work with S3Trigger resources.
type S3TriggerInterface interface {
	Create(*v1beta1.S3Trigger) (*v1beta1.S3Trigger, error)
	Update(*v1beta1.S3Trigger) (*v1beta1.S3Trigger, error)
	UpdateStatus(*v1beta1.S3Trigger) (*v1beta1.S3Trigger, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta1.S3Trigger, error)
	List(opts v1.ListOptions) (*v1beta1.S3TriggerList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.S3Trigger, err error)
	S3TriggerExpansion
}

// s3Triggers implements S3TriggerInterface
type s3Triggers struct {
	client rest.Interface
	ns     string
}

// newS3Triggers returns a S3Triggers
func newS3Triggers(c *KubelessV1beta1Client, namespace string) *s3Triggers {
	return &s3Triggers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the s3Trigger, and returns the corresponding s3Trigger object, and an error if there is any.
func (c *s3Triggers) Get(name string, options v1.GetOptions) (result *v1beta1.S3Trigger, err error) {
	result = &v1beta1.S3Trigger{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("s3triggers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of S3Triggers that match those selectors.
func (c *s3Triggers) List(opts v1.ListOptions) (result *v1beta1.S3TriggerList, err error) {
	result = &v1beta1.S3TriggerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("s3triggers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested s3Triggers.
func (c *s3Triggers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("s3triggers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a s3Trigger and creates it.  Returns the server's representation of the s3Trigger, and an error, if there is any.
func (c *s3Triggers) Create(s3Trigger *v1beta1.S3Trigger) (result *v1beta1.S3Trigger, err error) {
	result = &v1beta1.S3Trigger{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("s3triggers").
		Body(s3Trigger).
		Do().
		Into(result)
	return
}

// Update takes the representation of a s3Trigger and updates it. Returns the server's representation of the s3Trigger, and an error, if there is any.
func (c *s3Triggers) Update(s3Trigger *v1beta1.S3Trigger) (result *v1beta1.S3Trigger, err error) {
	result = &v1beta1.S3Trigger{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("s3triggers").
		Name(s3Trigger.Name).
		Body(s3Trigger).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *s3Triggers) UpdateStatus(s3Trigger *v1beta1.S3Trigger) (result *v1beta1.S3Trigger, err error) {
	result = &v1beta1.S3Trigger{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("s3triggers").
		Name(s3Trigger.Name).
		SubResource("status").
		Body(s3Trigger).
		Do().
		Into(result)
	return
}

// Delete takes name of the s3Trigger and deletes it. Returns an error if one occurs.
func (c *s3Triggers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("s3triggers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *s3Triggers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("s3triggers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched s3Trigger.
func (c *s3Triggers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.S3Trigger, err error) {
	result = &v1beta1.S3Trigger{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("s3triggers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
