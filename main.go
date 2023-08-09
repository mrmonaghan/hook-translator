/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"net/http"

	"github.com/mrmonaghan/hook-translator/pkg/handlers"
	"github.com/mrmonaghan/hook-translator/pkg/templater"
	"go.uber.org/zap"
)

func main() {

	log, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("unable to initialize logger: %w", err))
	}

	logger := log.Sugar()
	logger.Debug("initialized logger")

	tmplDir := "./templates"
	rulesDir := "./rules"

	tmpls, err := templater.LoadTemplates(tmplDir)
	if err != nil {
		logger.Panicw("unable to load templates from directory", "dir", tmplDir, err)
	}

	rules, err := templater.LoadRules(rulesDir, tmpls)
	if err != nil {
		logger.Panicw("unable to load rules from directory", "dir", rulesDir, err)
	}

	handler := handlers.RuleHandler{
		Rules:  rules,
		Logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", handler.ServeHTTP)
	mux.HandleFunc("/rules", handler.ServeHTTP)

	if err := http.ListenAndServe(":3333", mux); err != nil {
		panic(fmt.Errorf("error starting HTTP server: %w", err))
	}
}
