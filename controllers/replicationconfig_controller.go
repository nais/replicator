package controllers

import (
	"context"
	"fmt"
	"nais/replicator/internal/resources"
	"time"

	"nais/replicator/internal/replicator"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/record"

	naisiov1 "nais/replicator/api/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReplicationConfigReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	Recorder     record.EventRecorder
	SyncInterval time.Duration
}

// +kubebuilder:rbac:groups=nais.io,resources=replicationconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nais.io,resources=replicationconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nais.io,resources=replicationconfigs/finalizers,verbs=update
// +kubebuilder:rbac:groups="*",resources=*,verbs=create;update;patch;get;list;watch
func (r *ReplicationConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	rc := &naisiov1.ReplicationConfig{}
	err := r.Get(ctx, req.NamespacedName, rc)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	hash, err := replicator.Hash(&rc.Spec)
	if err != nil {
		return ctrl.Result{}, err
	}

	// skip reconciliation if hash is unchanged and timestamp is within sync interval
	// reconciliation is triggered when status subresource is updated, so we need this check to avoid infinite loop
	if rc.Status.SynchronizationHash == hash && !r.needsSync(rc.Status.SynchronizationTimestamp.Time) {
		log.Debugf("skipping reconciliation of %q, hash %q is unchanged and changed within syncInterval window", rc.Name, hash)
		return ctrl.Result{}, nil
	} else {
		log.Debugf("reconciling: hash changed: %v, outside syncInterval window: %v", rc.Status.SynchronizationHash != hash, r.needsSync(rc.Status.SynchronizationTimestamp.Time))
	}

	namespaces, err := r.listNamespaces(ctx, &rc.Spec.NamespaceSelector)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Debugf("reconciling %q to %d namespaces\n", rc.Name, len(namespaces.Items))

	secrets, err := replicator.LoadSecrets(ctx, r.Client, rc)
	if err != nil {
		return ctrl.Result{}, err
	}

	values := replicator.Merge(rc.Spec.TemplateValues.Values, secrets)

	ownerRef := []metav1.OwnerReference{
		{
			APIVersion: rc.APIVersion,
			Kind:       rc.Kind,
			Name:       rc.Name,
			UID:        rc.UID,
		},
	}

	for _, ns := range namespaces.Items {
		nsv := replicator.ExtractValues(ns, rc.Spec.TemplateValues.Namespace)

		renderResources, err := replicator.RenderResources(&replicator.TemplateValues{Values: replicator.Merge(values, nsv)}, rc.Spec.Resources)
		if err != nil {
			r.Recorder.Eventf(rc, "Warning", "RenderResources", "Unable to render resources for namespace %q: %v", ns.Name, err)
			continue
		}
		log.Debugf("rendered %d resources for namespace %q", len(renderResources), ns.Name)

		for _, resource := range renderResources {
			log.Debugf("resource: %s %s", resource.GetKind(), resource.GetName())
			spew.Dump(resource)

			resource.SetNamespace(ns.Name)
			resource.SetOwnerReferences(ownerRef)
			err = r.createUpdateResource(ctx, resource)
			if err != nil {
				r.Recorder.Eventf(rc, "Warning", "CreateResource", "Unable to create resource %v/%v for namespace %q: %v", resource.GetKind(), resource.GetName(), ns.Name, err)
				continue
			}
		}
	}

	rc.Status.SynchronizationTimestamp = metav1.Now()
	rc.Status.SynchronizationHash = hash
	if err := r.Status().Update(ctx, rc); err != nil {
		return ctrl.Result{}, err
	}

	log.Infof("finished reconciling %q to %d namespaces\n", rc.Name, len(namespaces.Items))

	return ctrl.Result{}, nil
}

func (r *ReplicationConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&naisiov1.ReplicationConfig{}).
		Complete(r)
}

func (r *ReplicationConfigReconciler) listNamespaces(ctx context.Context, ls *metav1.LabelSelector) (v1.NamespaceList, error) {
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

func (r *ReplicationConfigReconciler) createUpdateResource(ctx context.Context, resource *unstructured.Unstructured) error {
	err := r.Create(ctx, resource)
	if client.IgnoreAlreadyExists(err) != nil {
		return fmt.Errorf("creating resources: %w", err)
	}

	if errors.IsAlreadyExists(err) {
		if err := r.updateResource(ctx, resource); err != nil {
			return fmt.Errorf("updating resources: %w", err)
		}
		return nil
	}

	log.Infof("created resource %v/%v for namespace %q", resource.GetKind(), resource.GetName(), resource.GetNamespace())
	return nil
}

func (r *ReplicationConfigReconciler) updateResource(ctx context.Context, resource *unstructured.Unstructured) error {
	existing := &unstructured.Unstructured{}
	existing.SetGroupVersionKind(resource.GroupVersionKind())
	err := r.Get(ctx, client.ObjectKeyFromObject(resource), existing)
	if err != nil {
		return fmt.Errorf("getting existing resource: %w", err)
	}

	changed, err := resources.HasChanged(existing, resource)
	log.Debugf("resource %v/%v has changed: %v", existing.GetKind(), existing.GetName(), changed)
	if err != nil {
		return fmt.Errorf("comparing resources: %w", err)
	}

	if changed {
		resource.SetResourceVersion(existing.GetResourceVersion())
		err = r.Update(ctx, resource)
		if err != nil {
			return fmt.Errorf("updating resource: %w", err)
		}
		log.Infof("updated resource %v/%v for namespace %q", resource.GetKind(), resource.GetName(), resource.GetNamespace())
		return nil
	}

	log.Infof("resource %v/%v for namespace %q is unchanged", resource.GetKind(), resource.GetName(), resource.GetNamespace())
	return nil
}

func (r *ReplicationConfigReconciler) needsSync(timestamp time.Time) bool {
	window := time.Now().Add(-r.SyncInterval)
	return timestamp.Before(window)
}
