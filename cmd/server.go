package cmd

import (
	"fmt"
	"net/http"

	"github.com/mrmonaghan/hook-translator/internal/handlers"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a Stitch server",
	Run: func(cmd *cobra.Command, args []string) {

		// flags
		debug, _ := cmd.Flags().GetBool("debug")
		templateDir, _ := rootCmd.PersistentFlags().GetString("template-dir")
		rulesDir, _ := rootCmd.PersistentFlags().GetString("rules-dir")

		// initialize logger
		var log *zap.Logger
		var err error
		if debug {
			log, err = zap.NewDevelopment()

		} else {
			log, err = zap.NewProduction()
		}
		if err != nil {
			panic(fmt.Errorf("unable to initialize logger: %w", err))
		}
		logger := log.Sugar()
		logger.Debugw("initialized logger", "debug", debug)

		// load rules & templates
		rules, err := stitchCmdInit(templateDir, rulesDir)
		if err != nil {
			logger.Panicw("init error", err)
		}
		logger.Debugw("loaded rules and templates", "rulesDir", rulesDir, "templateDir", templateDir)

		// create handler
		handler := handlers.RuleHandler{
			Rules:  rules,
			Logger: logger,
		}

		// initialize server
		mux := http.NewServeMux()
		mux.HandleFunc("/webhook", handler.HandleWebhooks)
		mux.HandleFunc("/rules", handler.HandleRules)

		port, err := cmd.Flags().GetString("port")
		if err != nil {
			logger.Panicw("error retrieving --port flag value", err)
		}

		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
			panic(fmt.Errorf("error starting HTTP server: %w", err))
		}
	},
}

func init() {
	serverCmd.Flags().String("port", "8888", "specify a port for server to bind to")
	serverCmd.Flags().Bool("debug", false, "enable debug-level logging")
	rootCmd.AddCommand(serverCmd)
}
