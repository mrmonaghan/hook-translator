package rules

import (
	"os"

	"github.com/mrmonaghan/hook-translator/internal/templates"
	"github.com/mrmonaghan/hook-translator/internal/utils"
	"gopkg.in/yaml.v3"
)

type Rule struct {
	Name          string               `yaml:"name"`
	Enabled       bool                 `yaml:"enabled"`
	TemplateNames []string             `yaml:"templates"`
	templates     []templates.Template `yaml:"-"`
}

func (r *Rule) UnmarshalYAML(data []byte) error {
	if err := yaml.Unmarshal(data, &r); err != nil {
		return err
	}
	return nil
}

func (r *Rule) AssociateTemplates(templates []templates.Template) {
	for _, templateName := range r.TemplateNames {
		for _, template := range templates {
			if template.Name == templateName {
				r.templates = append(r.templates, template)
			}
		}
	}
}

func (r *Rule) GetTemplates() []templates.Template {
	return r.templates
}

func LoadRules(dir string, templates []templates.Template) (map[string]Rule, error) {
	files, err := utils.GetYamlFileNamesFromDir(dir)
	if err != nil {
		return map[string]Rule{}, err
	}

	rules := make(map[string]Rule)

	for _, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			return rules, err
		}
		var r Rule

		if err := r.UnmarshalYAML(b); err != nil {
			return rules, err
		}
		r.AssociateTemplates(templates)
		rules[r.Name] = r
	}
	return rules, nil
}
