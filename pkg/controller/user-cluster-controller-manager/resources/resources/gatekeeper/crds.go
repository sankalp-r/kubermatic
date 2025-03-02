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

package gatekeeper

import (
	"k8c.io/kubermatic/v2/pkg/resources"
	"k8c.io/kubermatic/v2/pkg/resources/reconciling"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	configAPIGroup                        = "config.gatekeeper.sh"
	configAPIVersion                      = "v1alpha1"
	constraintTemplateAPIGroup            = "templates.gatekeeper.sh"
	constraintTemplateAPIVersion          = "v1beta1"
	statusAPIGroup                        = "status.gatekeeper.sh"
	constraintPodStatusAPIVersion         = "v1beta1"
	constraintTemplatePodStatusAPIVersion = "v1beta1"
	mutatorPodStatusAPIVersion            = "v1beta1"
	mutatorPodStatusAPIGroup              = "status.gatekeeper.sh"
	assignAPIGroup                        = "mutations.gatekeeper.sh"
	assignAPIVersion                      = "v1alpha1"
	assignMetadataAPIGroup                = "mutations.gatekeeper.sh"
	assignMetadataAPIVersion              = "v1alpha1"
)

// ConfigCRDCreator returns the gatekeeper config CRD definition
func ConfigCRDCreator() reconciling.NamedCustomResourceDefinitionCreatorGetter {
	return func() (string, reconciling.CustomResourceDefinitionCreator) {
		return resources.GatekeeperConfigCRDName, func(crd *apiextensionsv1.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
			crd.Annotations = map[string]string{"controller-gen.kubebuilder.io/version": "v0.5.2"}
			crd.Labels = map[string]string{"gatekeeper.sh/system": "yes"}
			crd.Spec.Group = configAPIGroup
			crd.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    configAPIVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							XPreserveUnknownFields: resources.Bool(true),
							Type:                   "object",
						},
					},
				},
			}
			crd.Spec.Scope = apiextensionsv1.NamespaceScoped
			crd.Spec.Names.Kind = "Config"
			crd.Spec.Names.ListKind = "ConfigList"
			crd.Spec.Names.Plural = "configs"
			crd.Spec.Names.Singular = "config"

			return crd, nil
		}
	}
}

// ConstraintTemplateCRDCreator returns the gatekeeper constraintTemplate CRD definition
func ConstraintTemplateCRDCreator() reconciling.NamedCustomResourceDefinitionCreatorGetter {
	return func() (string, reconciling.CustomResourceDefinitionCreator) {
		return resources.GatekeeperConstraintTemplateCRDName, func(crd *apiextensionsv1.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
			crd.Annotations = map[string]string{"controller-gen.kubebuilder.io/version": "v0.5.2"}
			crd.Labels = map[string]string{"gatekeeper.sh/system": "yes", "controller-tools.k8s.io": "1.0"}
			crd.Spec.Group = constraintTemplateAPIGroup
			crd.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    constraintTemplateAPIVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							XPreserveUnknownFields: resources.Bool(true),
							Type:                   "object",
						},
					},
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
				},
				{
					Name:    "v1alpha1",
					Served:  true,
					Storage: false,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							XPreserveUnknownFields: resources.Bool(true),
							Type:                   "object",
						},
					},
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
				},
			}
			crd.Spec.Scope = apiextensionsv1.ClusterScoped
			crd.Spec.Names.Kind = "ConstraintTemplate"
			crd.Spec.Names.Plural = "constrainttemplates"

			return crd, nil
		}
	}
}

// ConstraintPodStatusCRDCreator returns the gatekeeper ConstraintPodStatus CRD definition
func ConstraintPodStatusCRDCreator() reconciling.NamedCustomResourceDefinitionCreatorGetter {
	return func() (string, reconciling.CustomResourceDefinitionCreator) {
		return resources.GatekeeperConstraintPodStatusCRDName, func(crd *apiextensionsv1.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
			crd.Labels = map[string]string{"gatekeeper.sh/system": "yes"}
			crd.Spec.Group = statusAPIGroup
			crd.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    constraintPodStatusAPIVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							XPreserveUnknownFields: resources.Bool(true),
							Type:                   "object",
						},
					},
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
				},
			}
			crd.Spec.Scope = apiextensionsv1.NamespaceScoped
			crd.Spec.Names.Kind = "ConstraintPodStatus"
			crd.Spec.Names.Plural = "constraintpodstatuses"

			return crd, nil
		}
	}
}

// ConstraintTemplatePodStatusCRDCreator returns the gatekeeper ConstraintTemplatePodStatus CRD definition
func ConstraintTemplatePodStatusCRDCreator() reconciling.NamedCustomResourceDefinitionCreatorGetter {
	return func() (string, reconciling.CustomResourceDefinitionCreator) {
		return resources.GatekeeperConstraintTemplatePodStatusCRDName, func(crd *apiextensionsv1.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
			crd.Labels = map[string]string{"gatekeeper.sh/system": "yes"}
			crd.Spec.Group = statusAPIGroup
			crd.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    constraintTemplatePodStatusAPIVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							XPreserveUnknownFields: resources.Bool(true),
							Type:                   "object",
						},
					},
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
				},
			}
			crd.Spec.Scope = apiextensionsv1.NamespaceScoped
			crd.Spec.Names.Kind = "ConstraintTemplatePodStatus"
			crd.Spec.Names.Plural = "constrainttemplatepodstatuses"

			return crd, nil
		}
	}
}

func MutatorPodStatusCRDCreator() reconciling.NamedCustomResourceDefinitionCreatorGetter {
	return func() (string, reconciling.CustomResourceDefinitionCreator) {
		return resources.GatekeeperMutatorPodStatusCRDName, func(crd *apiextensionsv1.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
			crd.Annotations = map[string]string{"controller-gen.kubebuilder.io/version": "v0.5.2"}
			crd.Labels = map[string]string{"gatekeeper.sh/system": "yes"}
			crd.Spec.Group = mutatorPodStatusAPIGroup
			crd.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    mutatorPodStatusAPIVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							XPreserveUnknownFields: resources.Bool(true),
							Type:                   "object",
						},
					},
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
				},
			}
			crd.Spec.Scope = apiextensionsv1.NamespaceScoped
			crd.Spec.Names.Kind = "MutatorPodStatus"
			crd.Spec.Names.Plural = "mutatorpodstatuses"

			return crd, nil
		}
	}
}

func AssignCRDCreator() reconciling.NamedCustomResourceDefinitionCreatorGetter {
	return func() (string, reconciling.CustomResourceDefinitionCreator) {
		return resources.GatekeeperAssignCRDName, func(crd *apiextensionsv1.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
			crd.Annotations = map[string]string{"controller-gen.kubebuilder.io/version": "v0.5.2"}
			crd.Labels = map[string]string{"gatekeeper.sh/system": "yes"}
			crd.Spec.Group = assignAPIGroup
			crd.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    assignAPIVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							XPreserveUnknownFields: resources.Bool(true),
							Type:                   "object",
						},
					},
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
				},
			}
			crd.Spec.Scope = apiextensionsv1.ClusterScoped
			crd.Spec.Names.Kind = "Assign"
			crd.Spec.Names.Plural = "assign"

			return crd, nil
		}
	}
}

func AssignMetadataCRDCreator() reconciling.NamedCustomResourceDefinitionCreatorGetter {
	return func() (string, reconciling.CustomResourceDefinitionCreator) {
		return resources.GatekeeperAssignMetadataCRDName, func(crd *apiextensionsv1.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
			crd.Annotations = map[string]string{"controller-gen.kubebuilder.io/version": "v0.5.2"}
			crd.Labels = map[string]string{"gatekeeper.sh/system": "yes"}
			crd.Spec.Group = assignMetadataAPIGroup
			crd.Spec.Versions = []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    assignMetadataAPIVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							XPreserveUnknownFields: resources.Bool(true),
							Type:                   "object",
						},
					},
					Subresources: &apiextensionsv1.CustomResourceSubresources{
						Status: &apiextensionsv1.CustomResourceSubresourceStatus{},
					},
				},
			}
			crd.Spec.Scope = apiextensionsv1.ClusterScoped
			crd.Spec.Names.Kind = "AssignMetadata"
			crd.Spec.Names.Plural = "assignmetadata"

			return crd, nil
		}
	}
}
