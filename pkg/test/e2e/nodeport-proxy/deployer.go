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

package nodeportproxy

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"k8c.io/kubermatic/v2/pkg/controller/operator/defaults"
	"k8c.io/kubermatic/v2/pkg/controller/operator/seed/resources/nodeportproxy"
	kubermaticv1 "k8c.io/kubermatic/v2/pkg/crd/kubermatic/v1"
	operatorv1alpha1 "k8c.io/kubermatic/v2/pkg/crd/operator/v1alpha1"
	features "k8c.io/kubermatic/v2/pkg/features"
	"k8c.io/kubermatic/v2/pkg/resources/reconciling"
	e2eutils "k8c.io/kubermatic/v2/pkg/test/e2e/utils"
	"k8c.io/kubermatic/v2/pkg/version/kubermatic"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// wait time between poll attempts of a Service vip and/or nodePort.
	// coupled with testTries to produce a net timeout value.
	hitEndpointRetryDelay = 2 * time.Second
	podReadinessTimeout   = 2 * time.Minute
)

// Deployer helps setting up nodeport proxy for testing.
type Deployer struct {
	Log       *zap.SugaredLogger
	Namespace string
	Versions  kubermatic.Versions
	Client    ctrlruntimeclient.Client

	resources []ctrlruntimeclient.Object
}

func (d *Deployer) SetUp() error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.Namespace,
		},
	}
	if d.Namespace == "" {
		ns.ObjectMeta.GenerateName = "nodeport-proxy-"
	}
	d.Log.Debugw("Creating namespace", "service", ns)
	if err := d.Client.Create(context.TODO(), ns); err != nil {
		return errors.Wrap(err, "failed to create namespace")
	}
	d.Namespace = ns.Name
	d.resources = append(d.resources, ns)

	cfg := &operatorv1alpha1.KubermaticConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubermatic",
			Namespace: d.Namespace,
		},
		Spec: operatorv1alpha1.KubermaticConfigurationSpec{
			FeatureGates: sets.NewString(features.TunnelingExposeStrategy),
		},
	}

	recorderFunc := func(create reconciling.ObjectCreator) reconciling.ObjectCreator {
		return func(existing ctrlruntimeclient.Object) (ctrlruntimeclient.Object, error) {
			obj, err := create(existing)
			if err != nil {
				return nil, err
			}

			d.resources = append(d.resources, obj)
			return existing, nil
		}
	}

	seed, err := defaults.DefaultSeed(&kubermaticv1.Seed{}, d.Log)
	if err != nil {
		return errors.Wrap(err, "failed to default seed")
	}

	if err := reconciling.ReconcileServiceAccounts(context.TODO(),
		[]reconciling.NamedServiceAccountCreatorGetter{
			nodeportproxy.ServiceAccountCreator(cfg),
		}, d.Namespace, d.Client, recorderFunc); err != nil {
		return errors.Wrap(err, "failed to reconcile ServiceAcconts")
	}
	if err := reconciling.ReconcileRoles(context.TODO(),
		[]reconciling.NamedRoleCreatorGetter{
			nodeportproxy.RoleCreator(),
		}, d.Namespace, d.Client, recorderFunc); err != nil {
		return errors.Wrap(err, "failed to reconcile Role")
	}
	if err := reconciling.ReconcileRoleBindings(context.TODO(),
		[]reconciling.NamedRoleBindingCreatorGetter{
			nodeportproxy.RoleBindingCreator(cfg),
		}, d.Namespace, d.Client, recorderFunc); err != nil {
		return errors.Wrap(err, "failed to reconcile RoleBinding")
	}
	if err := reconciling.ReconcileClusterRoles(context.TODO(),
		[]reconciling.NamedClusterRoleCreatorGetter{
			nodeportproxy.ClusterRoleCreator(cfg),
		}, "", d.Client, recorderFunc); err != nil {
		return errors.Wrap(err, "failed to reconcile ClusterRole")
	}
	if err := reconciling.ReconcileClusterRoleBindings(context.TODO(),
		[]reconciling.NamedClusterRoleBindingCreatorGetter{
			nodeportproxy.ClusterRoleBindingCreator(cfg),
		}, "", d.Client, recorderFunc); err != nil {
		return errors.Wrap(err, "failed to reconcile ClusterRoleBinding")
	}
	if err := reconciling.ReconcileServices(context.TODO(),
		[]reconciling.NamedServiceCreatorGetter{
			nodeportproxy.ServiceCreator(seed)},
		d.Namespace, d.Client, recorderFunc); err != nil {
		return errors.Wrap(err, "failed to reconcile Services")
	}
	if err := reconciling.ReconcileDeployments(context.TODO(),
		[]reconciling.NamedDeploymentCreatorGetter{
			nodeportproxy.EnvoyDeploymentCreator(cfg, seed, d.Versions),
			nodeportproxy.UpdaterDeploymentCreator(cfg, seed, d.Versions),
		}, d.Namespace, d.Client, recorderFunc); err != nil {
		return errors.Wrap(err, "failed to reconcile Kubermatic Deployments")
	}

	// Wait for pods to be ready
	for _, o := range d.resources {
		if dep, ok := o.(*appsv1.Deployment); ok {
			pods, err := d.waitForPodsCreated(dep)
			if err != nil {
				return errors.Wrap(err, "failed to create pods")
			}
			if err := d.waitForPodsReady(pods...); err != nil {
				return errors.Wrap(err, "failed waiting for pods to be running")
			}
		}
	}
	d.Log.Debugw("deployed nodeport-proxy", "version", d.Versions.Kubermatic)
	return nil
}

// CleanUp deletes the resources.
func (d *Deployer) CleanUp() error {
	for _, o := range d.resources {
		// TODO handle better errors
		_ = d.Client.Delete(context.TODO(), o)
	}
	return nil
}

// GetLbService returns the service used to expose the nodeport proxy pods.
func (d *Deployer) GetLbService() *corev1.Service {
	svc := corev1.Service{}
	if err := d.Client.Get(context.TODO(), types.NamespacedName{Name: nodeportproxy.ServiceName, Namespace: d.Namespace}, &svc); err != nil {
		return nil
	}
	return &svc
}

func (d *Deployer) waitForPodsCreated(dep *appsv1.Deployment) ([]string, error) {
	return e2eutils.WaitForPodsCreated(d.Client, int(*dep.Spec.Replicas), dep.Namespace, dep.Spec.Selector.MatchLabels)
}

func (d *Deployer) waitForPodsReady(pods ...string) error {
	if !e2eutils.CheckPodsRunningReady(d.Client, d.Namespace, pods, podReadinessTimeout) {
		return fmt.Errorf("timeout waiting for %d pods to be ready", len(pods))
	}
	return nil
}
