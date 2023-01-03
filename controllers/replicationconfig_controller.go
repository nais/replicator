/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"

	"nais/replicator/internal/replicator"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/record"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	naisiov1 "nais/replicator/api/v1"
)

// ReplicatorConfigurationReconciler reconciles a ReplicationConfig object
type ReplicatorConfigurationReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=nais.io,resources=replicatorconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nais.io,resources=replicatorconfigurations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nais.io,resources=replicatorconfigurations/finalizers,verbs=update
//+kubebuilder:rbac:groups="*",resources=*,verbs=create;update;patch;get;list;watch

func (r *ReplicatorConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	fmt.Println("Reconciling")

	rc := &naisiov1.ReplicationConfig{}
	err := r.Get(ctx, req.NamespacedName, rc)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	namespaces, err := r.listNamespaces(ctx, &rc.Spec.NamespaceSelector)
	if err != nil {
		return ctrl.Result{}, err
	}
	fmt.Printf("Found %d namespaces matching the selector\n", len(namespaces.Items))

	secrets, err := replicator.LoadSecrets(ctx, r.Client, rc)
	if err != nil {
		return ctrl.Result{}, err
	}

	values := merge(rc.Spec.Values, secrets)

	ownerRef := []metav1.OwnerReference{
		{
			APIVersion: rc.APIVersion,
			Kind:       rc.Kind,
			Name:       rc.Name,
			UID:        rc.UID,
		},
	}

	for _, ns := range namespaces.Items {
		nsv := replicator.ExtractValues(ns.Annotations)

		resources, err := replicator.RenderResources(&replicator.TemplateValues{Values: merge(values, nsv)}, rc.Spec.Resources)
		if err != nil {
			r.Recorder.Eventf(rc, "Warning", "RenderResources", "Unable to render resources for namespace %q: %v", ns.Name, err)
			continue
		}

		for _, resource := range resources {
			resource.SetNamespace(ns.Name)
			resource.SetOwnerReferences(ownerRef)
			err = r.createResource(ctx, resource)
			if err != nil {
				r.Recorder.Eventf(rc, "Warning", "CreateResource", "Unable to create resource %v/%v for namespace %q: %v", resource.GetKind(), resource.GetName(), ns.Name, err)
				continue
			}
		}
	}

	rc.Status.LastSynchronized = metav1.Now()
	if err := r.Status().Update(ctx, rc); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func merge(a, b map[string]string) map[string]string {
	if a == nil {
		return b
	}
	for k, v := range b {
		a[k] = v
	}
	return a
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReplicatorConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&naisiov1.ReplicationConfig{}).
		Complete(r)
}

func (r *ReplicatorConfigurationReconciler) listNamespaces(ctx context.Context, ls *metav1.LabelSelector) (v1.NamespaceList, error) {
	selector, err := metav1.LabelSelectorAsSelector(ls)
	if err != nil {
		return v1.NamespaceList{}, err
	}

	var namespaces v1.NamespaceList
	err = r.List(ctx, &namespaces, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		return v1.NamespaceList{}, err
	}
	return namespaces, nil
}

func (r ReplicatorConfigurationReconciler) createResource(ctx context.Context, resource *unstructured.Unstructured) error {
	err := r.Create(ctx, resource)
	if client.IgnoreAlreadyExists(err) != nil {
		return fmt.Errorf("creating resource: %w", err)
	}
	if errors.IsAlreadyExists(err) {
		existing := &unstructured.Unstructured{}
		existing.SetGroupVersionKind(resource.GroupVersionKind())
		err := r.Get(ctx, client.ObjectKeyFromObject(resource), existing)
		if err != nil {
			return fmt.Errorf("getting existing resource: %w", err)
		}
		resource.SetResourceVersion(existing.GetResourceVersion())

		err = r.Update(ctx, resource)
		if err != nil {
			return fmt.Errorf("updating resource: %w", err)
		}
	}
	return nil
}
