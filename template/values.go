package template

import yaml "gopkg.in/yaml.v1"

// ParseValuesFile parses the content of values.yaml
func ParseValuesFile(in []byte) (map[string]interface{}, error) {
	values := map[string]interface{}{}
	err := yaml.Unmarshal(in, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}
