package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"nais/replicator/internal/content"

	"github.com/davecgh/go-spew/spew"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"nais/replicator/internal/replicator"

	log "github.com/sirupsen/logrus"

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

	log.Debugf("reconciling %s%q to %d namespaces\n", rc.Kind, rc.Name, len(namespaces.Items))

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
			return ctrl.Result{}, err
		}

		log.Debugf("rendered %d resources for namespace %q", len(renderResources), ns.Name)

		for _, resource := range renderResources {
			log.Debugf("reconciling resource %s%q", resource.GetKind(), resource.GetName())
			if os.Getenv("DEBUG") == "true" {
				spew.Dump(resource)
			}

			resource.SetNamespace(ns.Name)
			resource.SetOwnerReferences(ownerRef)
			err = r.createUpdateResource(ctx, resource)
			if err != nil {
				if apierrors.HasStatusCause(err, v1.NamespaceTerminatingCause) {
					log.Infof("namespace %q is terminating, skipping resource %v/%v", ns.Name, resource.GetKind(), resource.GetName())
					continue
				}
				r.Recorder.Eventf(rc, "Warning", "createUpdateResource", "Unable to create/update resource %v/%v for namespace %q: %v", resource.GetKind(), resource.GetName(), ns.Name, err)
				return ctrl.Result{}, err
			}
		}
	}

	// Get the latest version of the ReplicationConfig before updating status.
	rc = &naisiov1.ReplicationConfig{}
	err = r.Get(ctx, req.NamespacedName, rc)
	if err != nil {
		return ctrl.Result{}, err
	}

	rc.Status.SynchronizationTimestamp = metav1.Now()
	rc.Status.SynchronizationHash = hash
	if err := r.Status().Update(ctx, rc); err != nil {
		r.Recorder.Eventf(rc, "Warning", "UpdateStatus", "Unable to update status for %q: %v", rc.Name, err)
		return ctrl.Result{}, err
	}

	log.Infof("finished reconcile %s%q to %d namespaces\n", rc.Kind, rc.Name, len(namespaces.Items))

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
	existing := &unstructured.Unstructured{}
	existing.SetGroupVersionKind(resource.GroupVersionKind())
	err := r.Get(ctx, client.ObjectKeyFromObject(resource), existing)
	if client.IgnoreNotFound(err) != nil {
		return err
	}

	if apierrors.IsNotFound(err) {
		err := r.Create(ctx, resource)
		if client.IgnoreAlreadyExists(err) != nil {
			return err
		}
		log.Infof("created resource %v/%v for namespace %q", resource.GetKind(), resource.GetName(), resource.GetNamespace())
		return nil
	}

	err = r.updateResource(ctx, resource, existing)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReplicationConfigReconciler) updateResource(ctx context.Context, resource, existing *unstructured.Unstructured) error {
	resourceContent, err := content.Get(resource)
	if err != nil {
		log.Warnf("unable to set resource content type: %v", err)
	}

	existingContent, err := content.Get(existing)
	if err != nil {
		log.Warnf("unable to set existing content type: %v", err)
	}

	if resourceContent.Equals(existingContent) {
		log.Debugf("unchanged resource %s%q for namespace %q", resource.GetKind(), resource.GetName(), resource.GetNamespace())
		return nil
	}

	resource.SetResourceVersion(existing.GetResourceVersion())
	err = r.Update(ctx, resource)
	if err != nil {
		return fmt.Errorf("updating resource: %w", err)
	}
	log.Infof("updated resource %s%q to namespace %q", resource.GetKind(), resource.GetName(), resource.GetNamespace())
	return nil
}

func (r *ReplicationConfigReconciler) needsSync(timestamp time.Time) bool {
	window := time.Now().Add(-r.SyncInterval)
	return timestamp.Before(window)
}
