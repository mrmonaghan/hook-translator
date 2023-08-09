package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mrmonaghan/hook-translator/internal/stitch"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "execute the specified template locally",
	Long:  `stitch template <template-name> '{"key": "value", "nested": {"key": "nested_value"}}'`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]
		data := args[1]
		templateDir, _ := rootCmd.PersistentFlags().GetString("template-dir")

		templates, err := stitch.LoadTemplates(templateDir)
		if err != nil {
			panic(fmt.Errorf("unable to load templates: %w", err))
		}

		var tmpl stitch.Template
		found := false
		for _, template := range templates {
			if template.Name == templateName {
				tmpl = template
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("no template '%s' found in template directory '%s'\n", templateName, templateDir)
		}

		m := make(map[string]interface{})

		if err := json.Unmarshal([]byte(data), &m); err != nil {
			panic(fmt.Errorf("error processing template data: %w", err))
		}

		rendered, err := tmpl.Render(m)
		if err != nil {
			panic(fmt.Errorf("error rendering template '%s': %w", templateName, err))
		}

		fmt.Println(rendered)

	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
