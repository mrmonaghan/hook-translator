package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mrmonaghan/hook-translator/internal/stitch"
)

func stitchCmdInit(templateDir string, rulesDir string) (map[string]stitch.Rule, error) {

	m := make(map[string]stitch.Rule)

	// parse & validate directories
	tDir, err := filepath.Abs(templateDir)
	if err != nil {
		return m, err
	}

	if err := validateDir(tDir); err != nil {
		return m, fmt.Errorf("unable to validate template directory: %w", err)
	}

	rDir, err := filepath.Abs(rulesDir)
	if err != nil {
		return m, err
	}

	if err := validateDir(rDir); err != nil {
		return m, fmt.Errorf("unable to validate rules directory: %w", err)
	}

	templates, err := stitch.LoadTemplates(tDir)
	if err != nil {
		return m, fmt.Errorf("unable to load templates: %w", err)
	}

	rules, err := stitch.LoadRules(rDir, templates)
	if err != nil {
		return m, fmt.Errorf("unable to load rules: %w", err)
	}

	return rules, nil

}

// helper function for determining if a directory exists
func validateDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' does not exist", dir)
	}
	return nil
}
