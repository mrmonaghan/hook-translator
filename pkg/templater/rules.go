package templater

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Name          string
	Enabled       bool
	TemplateNames []string   `yaml:"templates"`
	Templates     []Template `yaml:"-"`
}

func NewRule(path string, templates []Template) (Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Rule{}, fmt.Errorf("unable to read file: %w", err)
	}

	var r Rule
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return r, fmt.Errorf("unable to unmarshal rule: %w", err)
	}

	for _, templateName := range r.TemplateNames {
		for _, template := range templates {
			if template.Name == templateName {
				r.Templates = append(r.Templates, template)
			}
		}
	}

	return r, nil
}

func LoadRules(dir string, templates []Template) (map[string]Rule, error) {
	files, err := GetYamlFiles(dir)
	if err != nil {
		return map[string]Rule{}, err
	}

	rules := make(map[string]Rule)

	for _, file := range files {
		r, err := NewRule(file, templates)
		if err != nil {
			return rules, err
		}
		rules[r.Name] = r
	}
	return rules, nil
}
