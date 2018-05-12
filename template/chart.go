package template

import yaml "gopkg.in/yaml.v1"

// ChartInfo describes the content of chart.yaml
type ChartInfo struct {
	APIVersion  string `yaml:"apiVersion"`
	APPVersion  string `yaml:"appVersion"`
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
}

// ParseChartFile parses the content of Chart.yaml
func ParseChartFile(in []byte) (*ChartInfo, error) {
	values := ChartInfo{}
	err := yaml.Unmarshal(in, &values)
	if err != nil {
		return nil, err
	}

	return &values, nil
}
