/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

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

package kubernetes

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"

	k8cuserclusterclient "k8c.io/kubermatic/v2/pkg/cluster/client"
	"k8c.io/kubermatic/v2/pkg/clusterdeletion"
	controllerutil "k8c.io/kubermatic/v2/pkg/controller/util"
	kubermaticv1 "k8c.io/kubermatic/v2/pkg/crd/kubermatic/v1"
	kubermaticv1helper "k8c.io/kubermatic/v2/pkg/crd/kubermatic/v1/helper"
	kuberneteshelper "k8c.io/kubermatic/v2/pkg/kubernetes"
	"k8c.io/kubermatic/v2/pkg/provider"
	"k8c.io/kubermatic/v2/pkg/resources/certificates"
	"k8c.io/kubermatic/v2/pkg/version/kubermatic"

	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	autoscalingv1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	"k8s.io/client-go/tools/record"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	ControllerName = "kubermatic_kubernetes_controller"
)

// userClusterConnectionProvider offers functions to retrieve clients for the given user clusters
type userClusterConnectionProvider interface {
	GetClient(context.Context, *kubermaticv1.Cluster, ...k8cuserclusterclient.ConfigOption) (ctrlruntimeclient.Client, error)
}

type Features struct {
	VPA                          bool
	EtcdDataCorruptionChecks     bool
	KubernetesOIDCAuthentication bool
	EtcdLauncher                 bool
	Konnectivity                 bool
}

// Reconciler is a controller which is responsible for managing clusters
type Reconciler struct {
	ctrlruntimeclient.Client
	log                     *zap.SugaredLogger
	userClusterConnProvider userClusterConnectionProvider
	workerName              string

	externalURL  string
	seedGetter   provider.SeedGetter
	configGetter provider.KubermaticConfigurationGetter

	recorder record.EventRecorder

	overwriteRegistry                string
	nodePortRange                    string
	nodeAccessNetwork                string
	etcdDiskSize                     resource.Quantity
	userClusterMLAEnabled            bool
	dockerPullConfigJSON             []byte
	kubermaticImage                  string
	etcdLauncherImage                string
	dnatControllerImage              string
	machineControllerImageTag        string
	machineControllerImageRepository string
	concurrentClusterUpdates         int
	etcdBackupRestoreController      bool
	backupSchedule                   time.Duration

	oidcIssuerURL      string
	oidcIssuerClientID string

	features Features
	versions kubermatic.Versions

	tunnelingAgentIP string
	caBundle         *certificates.CABundle
}

// NewController creates a cluster controller.
func Add(
	mgr manager.Manager,
	log *zap.SugaredLogger,
	numWorkers int,
	workerName string,
	externalURL string,
	seedGetter provider.SeedGetter,
	configGetter provider.KubermaticConfigurationGetter,
	userClusterConnProvider userClusterConnectionProvider,
	overwriteRegistry string,
	nodePortRange string,
	nodeAccessNetwork string,
	etcdDiskSize resource.Quantity,
	userClusterMLAEnabled bool,
	dockerPullConfigJSON []byte,
	concurrentClusterUpdates int,
	etcdBackupRestoreController bool,
	backupSchedule time.Duration,

	oidcIssuerURL string,
	oidcIssuerClientID string,
	kubermaticImage string,
	etcdLauncherImage string,
	dnatControllerImage string,
	machineControllerImageTag string,
	machineControllerImageRepository string,

	tunnelingAgentIP string,
	caBundle *certificates.CABundle,

	features Features,
	versions kubermatic.Versions) error {

	reconciler := &Reconciler{
		log:                     log.Named(ControllerName),
		Client:                  mgr.GetClient(),
		userClusterConnProvider: userClusterConnProvider,
		workerName:              workerName,

		recorder: mgr.GetEventRecorderFor(ControllerName),

		overwriteRegistry:                overwriteRegistry,
		nodePortRange:                    nodePortRange,
		nodeAccessNetwork:                nodeAccessNetwork,
		etcdDiskSize:                     etcdDiskSize,
		userClusterMLAEnabled:            userClusterMLAEnabled,
		dockerPullConfigJSON:             dockerPullConfigJSON,
		kubermaticImage:                  kubermaticImage,
		etcdLauncherImage:                etcdLauncherImage,
		dnatControllerImage:              dnatControllerImage,
		machineControllerImageTag:        machineControllerImageTag,
		machineControllerImageRepository: machineControllerImageRepository,
		concurrentClusterUpdates:         concurrentClusterUpdates,
		etcdBackupRestoreController:      etcdBackupRestoreController,
		backupSchedule:                   backupSchedule,

		externalURL:  externalURL,
		seedGetter:   seedGetter,
		configGetter: configGetter,

		oidcIssuerURL:      oidcIssuerURL,
		oidcIssuerClientID: oidcIssuerClientID,

		tunnelingAgentIP: tunnelingAgentIP,
		caBundle:         caBundle,

		features: features,
		versions: versions,
	}

	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: reconciler, MaxConcurrentReconciles: numWorkers})
	if err != nil {
		return err
	}

	typesToWatch := []ctrlruntimeclient.Object{
		&corev1.Service{},
		&corev1.ServiceAccount{},
		&corev1.ConfigMap{},
		&corev1.Secret{},
		&corev1.Namespace{},
		&appsv1.StatefulSet{},
		&appsv1.Deployment{},
		&batchv1beta1.CronJob{},
		&policyv1beta1.PodDisruptionBudget{},
		&autoscalingv1beta2.VerticalPodAutoscaler{},
		&rbacv1.Role{},
		&rbacv1.RoleBinding{},
		&networkingv1.NetworkPolicy{},
	}

	for _, t := range typesToWatch {
		if err := c.Watch(&source.Kind{Type: t}, controllerutil.EnqueueClusterForNamespacedObject(mgr.GetClient())); err != nil {
			return fmt.Errorf("failed to create watcher for %T: %v", t, err)
		}
	}

	return c.Watch(&source.Kind{Type: &kubermaticv1.Cluster{}}, &handler.EnqueueRequestForObject{})
}

func (r *Reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.With("request", request)
	log.Debug("Processing")

	cluster := &kubermaticv1.Cluster{}
	if err := r.Get(ctx, request.NamespacedName, cluster); err != nil {
		if kerrors.IsNotFound(err) {
			log.Debug("Could not find cluster")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	log = log.With("cluster", cluster.Name)

	// Add a wrapping here so we can emit an event on error
	result, err := kubermaticv1helper.ClusterReconcileWrapper(
		ctx,
		r.Client,
		r.workerName,
		cluster,
		r.versions,
		kubermaticv1.ClusterConditionClusterControllerReconcilingSuccess,
		func() (*reconcile.Result, error) {
			// only reconcile this cluster if there are not yet too many updates running
			if available, err := controllerutil.ClusterAvailableForReconciling(ctx, r, cluster, r.concurrentClusterUpdates); !available || err != nil {
				log.Infow("Concurrency limit reached, checking again in 10 seconds", "concurrency-limit", r.concurrentClusterUpdates)
				return &reconcile.Result{
					RequeueAfter: 10 * time.Second,
				}, err
			}

			return r.reconcile(ctx, log, cluster)
		},
	)
	if err != nil {
		log.Errorw("Reconciling failed", zap.Error(err))
		r.recorder.Event(cluster, corev1.EventTypeWarning, "ReconcilingError", err.Error())
	}

	if result == nil {
		result = &reconcile.Result{}
	}

	return *result, err
}

func (r *Reconciler) reconcile(ctx context.Context, log *zap.SugaredLogger, cluster *kubermaticv1.Cluster) (*reconcile.Result, error) {
	// synchronize cluster.status.health for Kubernetes clusters
	if err := r.syncHealth(ctx, cluster); err != nil {
		return nil, fmt.Errorf("failed to sync health: %w", err)
	}

	if cluster.DeletionTimestamp != nil {
		log.Debug("Cleaning up cluster")

		// Defer getting the client to make sure we only request it if we actually need it
		userClusterClientGetter := func() (ctrlruntimeclient.Client, error) {
			client, err := r.userClusterConnProvider.GetClient(ctx, cluster)
			if err != nil {
				return nil, fmt.Errorf("failed to get user cluster client: %v", err)
			}
			return client, nil
		}
		// Always requeue a cluster after we executed the cleanup.
		return &reconcile.Result{RequeueAfter: 10 * time.Second}, clusterdeletion.New(r.Client, userClusterClientGetter, r.etcdBackupRestoreController).CleanupCluster(ctx, log, cluster)
	}

	res, err := r.reconcileCluster(ctx, cluster)
	if err != nil {
		updateErr := r.updateClusterError(ctx, cluster, kubermaticv1.ReconcileClusterError, err.Error())
		if updateErr != nil {
			return nil, fmt.Errorf("failed to set the cluster error: %v", updateErr)
		}
		return nil, fmt.Errorf("failed to reconcile cluster: %w", err)
	}

	if err := r.clearClusterError(ctx, cluster); err != nil {
		return nil, fmt.Errorf("failed to clear error on cluster: %v", err)
	}

	return res, nil
}

func (r *Reconciler) updateCluster(ctx context.Context, cluster *kubermaticv1.Cluster, modify func(*kubermaticv1.Cluster), opts ...ctrlruntimeclient.MergeFromOption) error {
	oldCluster := cluster.DeepCopy()
	modify(cluster)
	if reflect.DeepEqual(oldCluster, cluster) {
		return nil
	}

	return r.Patch(ctx, cluster, ctrlruntimeclient.MergeFromWithOptions(oldCluster, opts...))
}

func (r *Reconciler) AddFinalizers(ctx context.Context, cluster *kubermaticv1.Cluster, finalizers ...string) (*reconcile.Result, error) {
	if err := r.updateCluster(ctx, cluster, func(c *kubermaticv1.Cluster) {
		kuberneteshelper.AddFinalizer(c, finalizers...)
	}, ctrlruntimeclient.MergeFromWithOptimisticLock{}); err != nil {
		if !kerrors.IsConflict(err) {
			return nil, fmt.Errorf("failed to add finalizers %v: %w", finalizers, err)
		}
		// In case of conflict we just re-enqueue the item for later
		// processing without returning an error.
		r.log.Infow("failed to add finalizers", "error", err, "finalizers", finalizers)
		return &reconcile.Result{Requeue: true}, nil
	}
	return &reconcile.Result{}, nil
}

func (r *Reconciler) updateClusterError(ctx context.Context, cluster *kubermaticv1.Cluster, reason kubermaticv1.ClusterStatusError, message string) error {
	if cluster.Status.ErrorReason == nil || *cluster.Status.ErrorReason != reason {
		err := r.updateCluster(ctx, cluster, func(c *kubermaticv1.Cluster) {
			c.Status.ErrorMessage = &message
			c.Status.ErrorReason = &reason
		})
		if err != nil {
			return fmt.Errorf("failed to set error status on cluster to: errorReason=%q errorMessage=%q. Could not update cluster: %v", reason, message, err)
		}
	}

	return nil
}

func (r *Reconciler) clearClusterError(ctx context.Context, cluster *kubermaticv1.Cluster) error {
	if cluster.Status.ErrorReason != nil || cluster.Status.ErrorMessage != nil {
		err := r.updateCluster(ctx, cluster, func(c *kubermaticv1.Cluster) {
			c.Status.ErrorMessage = nil
			c.Status.ErrorReason = nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Reconciler) getOwnerRefForCluster(cluster *kubermaticv1.Cluster) metav1.OwnerReference {
	gv := kubermaticv1.SchemeGroupVersion
	return *metav1.NewControllerRef(cluster, gv.WithKind("Cluster"))
}
