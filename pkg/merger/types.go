// The types of the merger package.
// KubeBuilder markers are used to auto-generate OpenAPI schema for the plugin.
//
// +groupName=generators.kustomize.aabouzaid.com
// +versionName=v1alpha1
// +kubebuilder:validation:Required

package merger

import (
	"dario.cat/mergo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Merger manifest configuration.
const (
	ResourceGroup   string = "generators.kustomize.aabouzaid.com"
	ResourceVersion string = "v1alpha1"
	ResourceKind    string = "Merger"
)

// Merger manifest body.
type Merger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              mergerSpec `yaml:"spec" json:"spec"`
}

type mergerSpec struct {
	// +listType=map
	// +listMapKey=name
	Resources []mergerResource `yaml:"resources" json:"resources"`
}

//
// Merger Resource.
//

type mergerResource struct {
	Name   string         `yaml:"name" json:"name"`
	Input  resourceInput  `yaml:"input" json:"input"`
	Merge  resourceMerge  `yaml:"merge" json:"merge"`
	Output resourceOutput `yaml:"output" json:"output"`
}

//
// Merger Resource - Input.
//

type resourceInputFiles struct {
	Sources     []string `yaml:"sources" json:"sources"`
	Destination string   `yaml:"destination" json:"destination"`
}

// +enum
// +kubebuilder:validation:Enum=overlay;patch
type resourceInputMethod string

// Merger resource input method available options.
const (
	Overlay resourceInputMethod = "overlay"
	Patch   resourceInputMethod = "patch"
)

type resourceInput struct {
	Method resourceInputMethod `yaml:"method" json:"method"`
	Files  resourceInputFiles  `yaml:"files" json:"files"`
	items  []resourceInputFiles
}

//
// Merger Resource - Merge.
//

// +enum
// +kubebuilder:validation:Enum=append;combine;replace
type resourceMergeStrategy string

// Merger resource merge strategy available options.
// TODO: Support combine lists by named key.
const (
	Append  resourceMergeStrategy = "append"
	Combine resourceMergeStrategy = "combine"
	Replace resourceMergeStrategy = "replace"
)

type resourceMerge struct {
	Strategy resourceMergeStrategy `yaml:"strategy" json:"strategy"`
	config   func(*mergo.Config)
}

//
// Merger Resource - Output.
//

// +enum
// +kubebuilder:validation:Enum=raw
type resourceOutputFormat string

// TODO: Support ConfigMap and Secret.
const (
	Raw resourceOutputFormat = "raw"
)

type resourceOutput struct {
	Format  resourceOutputFormat `yaml:"format" json:"format"`
	rlItems []*yaml.RNode
}
