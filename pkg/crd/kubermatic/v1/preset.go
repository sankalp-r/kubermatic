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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PresetList is the type representing a PresetList
type PresetList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of presets
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md
	Items []Preset `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Preset is the type representing a Preset
type Preset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PresetSpec `json:"spec"`
}

// Presets specifies default presets for supported providers
type PresetSpec struct {
	Digitalocean *Digitalocean `json:"digitalocean,omitempty"`
	Hetzner      *Hetzner      `json:"hetzner,omitempty"`
	Azure        *Azure        `json:"azure,omitempty"`
	VSphere      *VSphere      `json:"vsphere,omitempty"`
	AWS          *AWS          `json:"aws,omitempty"`
	Openstack    *Openstack    `json:"openstack,omitempty"`
	Packet       *Packet       `json:"packet,omitempty"`
	GCP          *GCP          `json:"gcp,omitempty"`
	Kubevirt     *Kubevirt     `json:"kubevirt,omitempty"`
	Alibaba      *Alibaba      `json:"alibaba,omitempty"`
	Anexia       *Anexia       `json:"anexia,omitempty"`
	GKE          *GKE          `json:"gke,omitempty"`
	EKS          *EKS          `json:"eks,omitempty"`
	AKS          *AKS          `json:"aks,omitempty"`

	Fake *Fake `json:"fake,omitempty"`
	// see RequiredEmails
	RequiredEmailDomain string `json:"requiredEmailDomain,omitempty"`
	// RequiredEmails: specify emails and domains
	// RequiredEmailDomain is appended to RequiredEmails for backward compatibility.
	// e.g.:
	//   RequiredEmailDomain: "example.com"
	//   RequiredEmails: ["foo.com", "foo.bar@test.com"]
	// Result:
	//   *@example.com, *@foo.com and foo.bar@test.com can use the Preset
	RequiredEmails []string `json:"requiredEmails,omitempty"`
	Enabled        *bool    `json:"enabled,omitempty"`
}

func (s PresetSpec) IsEnabled() bool {
	if s.Enabled == nil {
		return true
	}

	return *s.Enabled
}

func (s *PresetSpec) SetEnabled(enabled bool) {
	s.Enabled = &enabled
}

type ProviderPreset struct {
	Enabled    *bool  `json:"enabled,omitempty"`
	Datacenter string `json:"datacenter,omitempty"`
}

func (s ProviderPreset) IsEnabled() bool {
	if s.Enabled == nil {
		return true
	}

	return *s.Enabled
}

type Digitalocean struct {
	ProviderPreset `json:",inline"`

	// Token is used to authenticate with the DigitalOcean API.
	Token string `json:"token"`
}

func (s Digitalocean) IsValid() bool {
	return len(s.Token) > 0
}

type Hetzner struct {
	ProviderPreset `json:",inline"`

	// Token is used to authenticate with the Hetzner API.
	Token string `json:"token"`

	// Network is the pre-existing Hetzner network in which the machines are running.
	// While machines can be in multiple networks, a single one must be chosen for the
	// HCloud CCM to work.
	// If this is empty, the network configured on the datacenter will be used.
	Network string `json:"network,omitempty"`
}

func (s Hetzner) IsValid() bool {
	return len(s.Token) > 0
}

type Azure struct {
	ProviderPreset `json:",inline"`

	TenantID       string `json:"tenantId"`
	SubscriptionID string `json:"subscriptionId"`
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`

	ResourceGroup     string `json:"resourceGroup,omitempty"`
	VNetResourceGroup string `json:"vnetResourceGroup,omitempty"`
	VNetName          string `json:"vnet,omitempty"`
	SubnetName        string `json:"subnet,omitempty"`
	RouteTableName    string `json:"routeTable,omitempty"`
	SecurityGroup     string `json:"securityGroup,omitempty"`
	// LoadBalancerSKU sets the LB type that will be used for the Azure cluster, possible values are "basic" and "standard", if empty, "basic" will be used
	LoadBalancerSKU LBSKU `json:"loadBalancerSKU"`
}

func (s Azure) IsValid() bool {
	return len(s.TenantID) > 0 &&
		len(s.SubscriptionID) > 0 &&
		len(s.ClientID) > 0 &&
		len(s.ClientSecret) > 0
}

type VSphere struct {
	ProviderPreset `json:",inline"`

	Username string `json:"username"`
	Password string `json:"password"`

	VMNetName        string `json:"vmNetName,omitempty"`
	Datastore        string `json:"datastore,omitempty"`
	DatastoreCluster string `json:"datastoreCluster,omitempty"`
	ResourcePool     string `json:"resourcePool,omitempty"`
}

func (s VSphere) IsValid() bool {
	return len(s.Username) > 0 && len(s.Password) > 0
}

type AWS struct {
	ProviderPreset `json:",inline"`

	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`

	AssumeRoleARN        string `json:"assumeRoleARN"`
	AssumeRoleExternalID string `json:"assumeRoleExternalID,omitempty"`

	VPCID               string `json:"vpcId,omitempty"`
	RouteTableID        string `json:"routeTableId,omitempty"`
	InstanceProfileName string `json:"instanceProfileName,omitempty"`
	SecurityGroupID     string `json:"securityGroupID,omitempty"`
	ControlPlaneRoleARN string `json:"roleARN,omitempty"`
}

func (s AWS) IsValid() bool {
	return len(s.AccessKeyID) > 0 && len(s.SecretAccessKey) > 0
}

type Openstack struct {
	ProviderPreset `json:",inline"`

	UseToken bool `json:"useToken,omitempty"`

	ApplicationCredentialID     string `json:"applicationCredentialID,omitempty"`
	ApplicationCredentialSecret string `json:"applicationCredentialSecret,omitempty"`

	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`

	Tenant    string `json:"tenant,omitempty"`
	TenantID  string `json:"tenantID,omitempty"`
	Project   string `json:"project,omitempty"`
	ProjectID string `json:"projectID,omitempty"`

	Network        string `json:"network,omitempty"`
	SecurityGroups string `json:"securityGroups,omitempty"`
	FloatingIPPool string `json:"floatingIpPool,omitempty"`
	RouterID       string `json:"routerID,omitempty"`
	SubnetID       string `json:"subnetID,omitempty"`
}

// GetProject returns the the project if defined otherwise fallback to tenant
// Deprecated: the tenant auth var is depreciated in openstack. In pkg/apis/kubermatic/v1/preset.go we will only use Project
func (s Openstack) GetProject() string {
	if len(s.Project) > 0 {
		return s.Project
	} else {
		return s.Tenant
	}
}

// GetProjectId returns the the projectID if defined otherwise fallback to tenantID
// Deprecated: the tenantID auth var is depreciated in openstack. In pkg/apis/kubermatic/v1/preset.go we will only use ProjectID
func (s Openstack) GetProjectId() string {
	if len(s.ProjectID) > 0 {
		return s.ProjectID
	} else {
		return s.TenantID
	}
}

func (s Openstack) IsValid() bool {
	if s.UseToken {
		return true
	}

	if len(s.ApplicationCredentialID) > 0 {
		return len(s.ApplicationCredentialSecret) > 0
	}

	hasProjectOrTenant := len(s.Project) > 0 || len(s.ProjectID) > 0 || len(s.Tenant) > 0 || len(s.TenantID) > 0
	return len(s.Username) > 0 &&
		len(s.Password) > 0 &&
		hasProjectOrTenant &&
		len(s.Domain) > 0
}

type Packet struct {
	ProviderPreset `json:",inline"`

	APIKey    string `json:"apiKey"`
	ProjectID string `json:"projectId"`

	BillingCycle string `json:"billingCycle,omitempty"`
}

func (s Packet) IsValid() bool {
	return len(s.APIKey) > 0 && len(s.ProjectID) > 0
}

type GCP struct {
	ProviderPreset `json:",inline"`

	ServiceAccount string `json:"serviceAccount"`

	Network    string `json:"network,omitempty"`
	Subnetwork string `json:"subnetwork,omitempty"`
}

func (s GCP) IsValid() bool {
	return len(s.ServiceAccount) > 0
}

type Fake struct {
	ProviderPreset `json:",inline"`

	Token string `json:"token"`
}

func (s Fake) IsValid() bool {
	return len(s.Token) > 0
}

type Kubevirt struct {
	ProviderPreset `json:",inline"`

	Kubeconfig string `json:"kubeconfig"`
}

func (s Kubevirt) IsValid() bool {
	return len(s.Kubeconfig) > 0
}

type Alibaba struct {
	ProviderPreset `json:",inline"`

	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

func (s Alibaba) IsValid() bool {
	return len(s.AccessKeyID) > 0 &&
		len(s.AccessKeySecret) > 0
}

type Anexia struct {
	ProviderPreset `json:",inline"`

	// Token is used to authenticate with the Anexia API.
	Token string `json:"token"`
}

func (s Anexia) IsValid() bool {
	return len(s.Token) > 0
}

type GKE struct {
	ProviderPreset `json:",inline"`

	ServiceAccount string `json:"serviceAccount"`
}

func (s GKE) IsValid() bool {
	return len(s.ServiceAccount) > 0
}

type EKS struct {
	ProviderPreset `json:",inline"`

	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

func (s EKS) IsValid() bool {
	return len(s.AccessKeyID) > 0 &&
		len(s.SecretAccessKey) > 0
}

type AKS struct {
	ProviderPreset `json:",inline"`

	TenantID       string `json:"tenantId"`
	SubscriptionID string `json:"subscriptionId"`
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
}

func (s AKS) IsValid() bool {
	return len(s.TenantID) > 0 &&
		len(s.SubscriptionID) > 0 &&
		len(s.ClientID) > 0 &&
		len(s.ClientSecret) > 0
}
