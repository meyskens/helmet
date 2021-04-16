package template

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pmylund/sortutil"
)

// Chart represents everything in a Helm chart
type Chart struct {
	templateFiles map[string][]byte
	values        map[interface{}]interface{}
	chartInfo     *ChartInfo
	release       Release
}

type tempateInput struct {
	Chart   *ChartInfo
	Values  map[interface{}]interface{}
	Release Release
}

// New generates a new Chart object
func New(chartInfo *ChartInfo, values map[interface{}]interface{}, files map[string][]byte) *Chart {
	return &Chart{
		templateFiles: files,
		values:        values,
		chartInfo:     chartInfo,
	}
}

// MergeValues allows to overwrite values and merge them with the existing ones
func (c *Chart) MergeValues(newValues map[interface{}]interface{}) {
	mergeValues(c.values, newValues)
}

// CreateManifests generates the Kubernetes manifests and NOTES.txt output
func (c *Chart) CreateManifests(release Release) (map[string][]byte, []byte, error) {
	c.release = release

	printTmpl := template.New("Print")
	printTmpl = printTmpl.Funcs(getFuncMap(printTmpl))

	tmpl := template.New("Manifests")
	tmpl = tmpl.Funcs(getFuncMap(tmpl))
	var err error

	var fileNames []string
	outputFiles := []string{}
	for k := range c.templateFiles {
		fileNames = append(fileNames, k)
	}

	sortutil.CiAsc(fileNames)
	for _, name := range fileNames {
		if name == "NOTES.txt" {
			printTmpl, err = printTmpl.Parse(string(c.templateFiles[name]))
		} else if name == "_helpers.tpl" {
			printTmpl, err = printTmpl.Parse(string(c.templateFiles[name]))
			tmpl, err = tmpl.Parse(string(c.templateFiles[name]))
		} else {
			outputFiles = append(outputFiles, name)
		}

		if err != nil {
			return nil, nil, fmt.Errorf("Error parsing %s: %s", name, err)
		}
	}

	manifests := map[string][]byte{}
	for _, file := range outputFiles {
		manifest, err := c.generateOutput(tmpl, c.templateFiles[file])
		if err != nil {
			return nil, nil, fmt.Errorf("error while parsing %s: %s", file, err)
		}
		manifests[file] = manifest
	}
	notes, _ := c.templateToBytes(printTmpl) // can be empty, allowed to fail

	return manifests, notes, nil
}

func (c *Chart) generateOutput(tmpl *template.Template, content []byte) ([]byte, error) {
	t, err := tmpl.Clone()
	if err != nil {
		return nil, err
	}
	t, err = t.Parse(string(content))
	if err != nil {
		return nil, err
	}
	return c.templateToBytes(t)
}

func (c *Chart) templateToBytes(t *template.Template) ([]byte, error) {
	input := tempateInput{
		Chart:   c.chartInfo,
		Values:  c.values,
		Release: c.release,
	}
	buf := new(bytes.Buffer)
	err := t.Execute(buf, input)

	return buf.Bytes(), err
}

// Clone gives a deep copy of the chart
func (c *Chart) Clone() (*Chart, error) {
	info := *c.chartInfo
	return &Chart{
		templateFiles: copyMapB(c.templateFiles),
		values:        copyMapI(c.values),
		chartInfo:     &info,
		release:       c.release,
	}, nil
}

func copyMapB(m map[string][]byte) map[string][]byte {
	cp := make(map[string][]byte)
	for k, v := range m {
		cp[k] = []byte(string(v))
	}

	return cp
}

func copyMapI(m map[interface{}]interface{}) map[interface{}]interface{} {
	cp := make(map[interface{}]interface{})
	for k, v := range m {
		vm, ok := v.(map[interface{}]interface{})
		if ok {
			cp[k] = copyMapI(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}
