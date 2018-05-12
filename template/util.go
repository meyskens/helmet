package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/Masterminds/sprig"
	toml "github.com/pelletier/go-toml"
	yaml "gopkg.in/yaml.v1"
)

func toYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

func fromYaml(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

func toJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

func fromJson(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := json.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}
func toToml(v interface{}) string {
	b := bytes.NewBuffer(nil)
	e := toml.NewEncoder(b)
	err := e.Encode(v)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

func getFuncMap(t *template.Template) template.FuncMap {
	f := sprig.FuncMap()
	f["toYaml"] = toYaml
	f["fromYaml"] = fromYaml
	f["toToml"] = toToml
	f["toJson"] = toJson
	f["fromJson"] = fromJson

	f["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := t.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	f["required"] = func(warn string, val interface{}) (interface{}, error) {
		if val == nil {
			return val, fmt.Errorf(warn)
		} else if _, ok := val.(string); ok {
			if val == "" {
				return val, fmt.Errorf(warn)
			}
		}
		return val, nil
	}

	return f
}
