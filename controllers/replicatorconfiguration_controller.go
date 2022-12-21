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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"nais/replicator/internal/parser"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	naisiov1 "nais/replicator/api/v1"
)

// ReplicatorConfigurationReconciler reconciles a ReplicatorConfiguration object
type ReplicatorConfigurationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=nais.io,resources=replicatorconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nais.io,resources=replicatorconfigurations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nais.io,resources=replicatorconfigurations/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ReplicatorConfiguration object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *ReplicatorConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	fmt.Println("Reconciling")

	rc := &naisiov1.ReplicatorConfiguration{}
	err := r.Get(ctx, req.NamespacedName, rc)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	namespaces, err := r.listNamespaces(ctx, &rc.Spec.NamespaceSelector)
	if err != nil {
		return ctrl.Result{}, err
	}
	fmt.Printf("Found %d namespaces matching the selector\n", len(namespaces.Items))

	values := &parser.TemplateValues{
		Values: map[string]string{},
	}

	if err := r.loadSecrets(ctx, rc, values); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.loadConfigMaps(ctx, rc, values); err != nil {
		return ctrl.Result{}, err
	}

	for _, ns := range namespaces.Items {
		err := parser.ParseAnnotations(ns.ObjectMeta.Annotations, values)
		if err != nil {
			return ctrl.Result{}, err
		}

		resources, err := parser.Resources(values, rc.Spec.Resources)
		if err != nil {
			return ctrl.Result{}, err
		}

		for _, resource := range resources {
			resource.SetNamespace(ns.Name)
			err = r.createResource(ctx, resource)
			if err != nil {
				fmt.Printf("creating resource: %v\n", err)
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r ReplicatorConfigurationReconciler) createResource(ctx context.Context, resource *unstructured.Unstructured) error {
	err := r.Create(ctx, resource)
	if client.IgnoreAlreadyExists(err) != nil {
		return fmt.Errorf("creating resource: %v", err)
	}
	if errors.IsAlreadyExists(err) {
		err := r.Update(ctx, resource)
		if err != nil {
			return fmt.Errorf("updating resource: %v", err)
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReplicatorConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&naisiov1.ReplicatorConfiguration{}).
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

// TODO clean up printing
func (r *ReplicatorConfigurationReconciler) loadSecrets(ctx context.Context, rc *naisiov1.ReplicatorConfiguration, values *parser.TemplateValues) error {
	for _, s := range rc.Spec.Values.Secrets {

		var secret v1.Secret
		if err := r.Get(ctx, client.ObjectKey{Name: s.Name, Namespace: s.Namespace}, &secret); err != nil {
			return err
		}

		for k, v := range secret.Data {
			values.Values[k] = string(v)
		}
	}
	return nil
}

// TODO: see if we can use the same function for both secrets and configmaps
func (r *ReplicatorConfigurationReconciler) loadConfigMaps(ctx context.Context, rc *naisiov1.ReplicatorConfiguration, values *parser.TemplateValues) error {
	for _, s := range rc.Spec.Values.ConfigMaps {

		var configMap v1.ConfigMap
		if err := r.Get(ctx, client.ObjectKey{Name: s.Name, Namespace: s.Namespace}, &configMap); err != nil {
			return err
		}

		for k, v := range configMap.Data {
			values.Values[k] = v
		}
	}
	return nil
}
