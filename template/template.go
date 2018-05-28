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
	values        map[string]interface{}
	chartInfo     *ChartInfo
	release       Release
}

type tempateInput struct {
	Chart   *ChartInfo
	Values  map[string]interface{}
	Release Release
}

// New generates a new Chart object
func New(chartInfo *ChartInfo, values map[string]interface{}, files map[string][]byte) *Chart {
	return &Chart{
		templateFiles: files,
		values:        values,
		chartInfo:     chartInfo,
	}
}

// MergeValues allows to overwrite values and merge them with the existing ones
func (c *Chart) MergeValues(newValues map[string]interface{}) {
	mergeValues(c.values, newValues)
}

// CreateManifests generates the Kubernetes manifests and NOTES.txt output
func (c *Chart) CreateManifests(release Release) (map[string][]byte, []byte, error) {
	c.release = release

	printTmpl := template.New("Manifests")
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
	notes, err := c.templateToBytes(printTmpl)
	if err != nil {
		return nil, nil, err
	}

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
