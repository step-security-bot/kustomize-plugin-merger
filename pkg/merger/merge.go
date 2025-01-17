// The main logic for Merger plugin.
package merger

import (
	"log"

	"dario.cat/mergo"
	"github.com/knadh/koanf"
	koanfYaml "github.com/knadh/koanf/parsers/yaml"
	koanfFile "github.com/knadh/koanf/providers/file"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// setMergeStrategy sets "mergo" merging strategy, mainly for merging arrays.
func (r *mergerResource) setMergeStrategy() {
	switch r.Merge.Strategy {
	case Replace:
		r.Merge.config = mergo.WithOverride
	case Append:
		r.Merge.config = mergo.WithAppendSlice
	case Combine:
		r.Merge.config = mergo.WithSliceDeepCopy
	default:
		r.Merge.config = func(*mergo.Config) {}
	}
}

func (r *mergerResource) setInputFilesOverlay() {
	for _, inputFileSource := range r.Input.Files.Sources {
		r.Input.items = append(r.Input.items,
			resourceInputFiles{
				Sources:     []string{inputFileSource},
				Destination: r.Input.Files.Destination,
			})
	}
}

func (r *mergerResource) setInputFilesPatch() {
	r.Input.items = append(r.Input.items, r.Input.Files)
}

// setInputFiles determines the input file sources based on the input method (Overlay or Patch).
func (r *mergerResource) setInputFiles() error {
	switch r.Input.Method {
	case Overlay:
		r.setInputFilesOverlay()
	case Patch:
		r.setInputFilesPatch()
	}
	return nil
}

// merge performs the actual merging of configuration files from resourceInputFiles sources.
func (r *mergerResource) merge(rfs []resourceInputFiles) {
	for _, rf := range rfs {
		k := koanf.New(".")

		if err := k.Load(koanfFile.Provider(rf.Destination), koanfYaml.Parser()); err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		for _, srcFile := range rf.Sources {
			err := k.Load(koanfFile.Provider(srcFile), koanfYaml.Parser(),
				koanf.WithMergeFunc(func(src, dst map[string]interface{}) error {
					if err := mergo.Merge(&dst, src, mergo.WithOverride, r.Merge.config); err != nil {
						log.Fatalf("Error merging config: %v", err)
					}
					return nil
				}))
			if err != nil {
				log.Fatalf("Error loading config: %v", err)
			}
		}

		mergedData, err := k.Marshal(koanfYaml.Parser())
		if err != nil {
			log.Fatalf("Error marshaling yaml: %v", err)
		}
		rNode := yaml.MustParse(string(mergedData))
		r.Output.rlItems = append(r.Output.rlItems, rNode)
	}
}
