// Methods inherited by fn/framework.
package merger

import (
	_ "embed"

	"k8s.io/kube-openapi/pkg/validation/spec"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

//go:embed "schema/generators.kustomize.aabouzaid.com_mergers.yaml"
var mergerSchemaDefinition string

// Schema returns the OpenAPI schema definition for Merger.
func (m *Merger) Schema() (*spec.Schema, error) {
	schema, err := framework.SchemaFromFunctionDefinition(
		resid.NewGvk(ResourceGroup, ResourceVersion, ResourceKind),
		mergerSchemaDefinition)
	return schema, errors.WrapPrefixf(err, "failed to parse Merger schema")
}

// Default sets default values for Merger resources.
func (m *Merger) Default() error {
	for index := range m.Spec.Resources {
		// Defaults input files.
		if err := m.Spec.Resources[index].setInputFiles(); err != nil {
			return err
		}
		// Defaults merge strategy.
		m.Spec.Resources[index].setMergeStrategy()
	}
	return nil
}

// Validate checks in Merger resource against its OpenAPI schema.
func (m *Merger) Validate() error {
	return nil
}

// Filter performs the merging of configuration files for Merger resources.
func (m *Merger) Filter(rlItems []*yaml.RNode) ([]*yaml.RNode, error) {
	for _, resource := range m.Spec.Resources {
		resource.merge(resource.Input.items)
		rlItems = append(rlItems, resource.Output.rlItems...)
	}
	return rlItems, nil
}
