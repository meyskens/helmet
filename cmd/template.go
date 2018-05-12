package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/meyskens/helmet/template"
	"github.com/spf13/cobra"
)

type templateCmdOptions struct {
	ExecuteTemplates []string
	Name             string
	TemplateName     string
	Namespace        string
	ShowNotes        bool
	OutputDir        string
	SetArgs          []string
	SetStringArgs    []string
	valueFiles       []string
}

// NewTotalCmd generates the `total` command
func NewTemplateCmd() *cobra.Command {
	o := templateCmdOptions{}
	c := &cobra.Command{
		Use:   "template",
		Short: "Template outputs the generated manifest for a given chart with it's values",
		Long:  `Template outputs the generated manifest for a given chart with it's values`,
		Example: `

		`,
		RunE: o.RunE,
	}

	// TO DO: allow values override

	c.Flags().StringArrayVarP(&o.ExecuteTemplates, "execute", "x", []string{}, "only execute the given templates")
	c.Flags().StringVarP(&o.Name, "name", "n", "RELEASE-NAME", "release name")
	//c.Flags().StringVar(&o.TemplateName, "name-template", "", "specify template used to name the release")
	c.Flags().StringVar(&o.Namespace, "namespace", "", "namespace to install the release into")
	c.Flags().BoolVar(&o.ShowNotes, "notes", false, "show the computed NOTES.txt file as well")
	c.Flags().StringVar(&o.OutputDir, "output-dir", "", "writes the executed templates to files in output-dir instead of stdout")
	//c.Flags().StringArrayVar(&o.SetArgs, "set", []string{}, "set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	//c.Flags().StringArrayVar(&o.SetStringArgs, "set-string", []string{}, "set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
	//c.Flags().StringArrayVarP(&o.valueFiles, "values", "f", []string{}, "specify values in a YAML file (can specify multiple) (default [])")

	return c
}

func (t *templateCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Needs chart as argument")
	}

	err := t.validateArguments()
	if err != nil {
		return err
	}

	// read in values
	chartFile, err := ioutil.ReadFile(path.Join(args[0], "Chart.yaml"))
	if err != nil {
		return err
	}
	valuesFile, err := ioutil.ReadFile(path.Join(args[0], "values.yaml"))
	if err != nil {
		return err
	}

	// parse values
	chartInfo, err := template.ParseChartFile(chartFile)
	if err != nil {
		return err
	}
	values, err := template.ParseValuesFile(valuesFile)
	if err != nil {
		return err
	}

	files := []os.FileInfo{}

	if len(t.ExecuteTemplates) > 0 {
		for _, file := range t.ExecuteTemplates {
			info, err := os.Stat(path.Join(args[0], "templates", file))
			if err != nil {
				return err
			}
			files = append(files, info)
		}
	} else {
		// read templates from fs
		var err error
		files, err = ioutil.ReadDir(path.Join(args[0], "templates"))
		if err != nil {
			return err
		}
	}

	filesMap := map[string][]byte{}
	for _, f := range files {
		content, err := ioutil.ReadFile(path.Join(args[0], "templates", f.Name()))
		if err != nil {
			return err
		}
		filesMap[f.Name()] = content
	}

	chart := template.New(chartInfo, values, filesMap)

	output, notes, err := chart.CreateManifests(template.NewRelease(t.Name, t.Namespace))
	if err != nil {
		return err
	}

	if t.OutputDir != "" {
		for name, fileContent := range output {
			ioutil.WriteFile(path.Join(t.OutputDir, name), fileContent, 0755)
		}
	} else {
		for _, fileContent := range output {
			fmt.Println(string(fileContent))
		}
	}

	if t.ShowNotes {
		fmt.Println(string(notes))
	}

	return nil
}

func (t *templateCmdOptions) validateArguments() error {
	if t.OutputDir != "" && !validatePath(t.OutputDir) {
		return errors.New("output-dir is not a vailid directory")
	}

	for _, file := range t.valueFiles {
		if !validateFile(file) {
			return fmt.Errorf("values file with path \"%s\" is not a vailid file path", file)
		}
	}

	for _, set := range t.SetArgs {
		if !strings.Contains(set, "=") {
			return fmt.Errorf("set option \"%s\" is not a valid syntax", set)
		}
	}

	for _, set := range t.SetStringArgs {
		if !strings.Contains(set, "=") {
			return fmt.Errorf("set option \"%s\" is not a valid syntax", set)
		}
	}

	return nil
}
