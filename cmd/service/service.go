package service

import (
	"es-content-export/collecters"
	"es-content-export/settings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	configFile     string
	StartServerCmd = &cobra.Command{
		Use:   "service",
		Short: "start service",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartServerCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yml", "config file path")
}

func setup() {
	settings.SetupConfig(configFile)
}

func run() {
	reg := prometheus.NewRegistry()
	collecters.RegistryEsCollect(reg)
	go collecters.TickerTask()
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	go http.ListenAndServe(":9999", nil)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
