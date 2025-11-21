package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	version     = "0.1.0"
	kubeconfig  string // Path to kubeconfig file
	contextName string // Kubernetes context to use
	namespace   string // Namespace to search

	ingressClassName    string // Ingress class name to search
	controllerNamespace string // Ingress controller namespace

	output string // Output directory for reports
	format string // Output format (markdown|json)
)

var rootCmd = &cobra.Command{
	Use:   "analyzer",
	Short: "Ingress-NGINX Migration Analyzer",
	Long: `Analyze your ingress-nginx usage and plan your migration before the March 2026 EOL.

This tool scans Kubernetes clusters to identify ingress-nginx resources,
classifies migration complexity, and generates actionable reports.`,
	Version: version,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", getDefaultKubeconfig(), "Path to kubeconfig file")
	rootCmd.PersistentFlags().StringVar(&contextName, "context", "", "Kubernetes context to use")
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Specific namespace to scan (default: all namespaces)")

	rootCmd.PersistentFlags().StringVar(&ingressClassName, "ingressclass-name", "nginx", "Specific ingress class name to scan (default: nginx)")
	rootCmd.PersistentFlags().StringVar(&controllerNamespace, "controller-namespace", "ingress-nginx", "Specific ingress controller namespace to scan (default: ingress-nginx)")
}

func getDefaultKubeconfig() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".kube", "config")
	}
	return ""
}

func validateFlags() error {
	// Check if kubeconfig file exists
	if kubeconfig != "" {
		if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
			return fmt.Errorf("kubeconfig file not found: %s", kubeconfig)
		}
	}

	// Validate output format
	if format != "markdown" && format != "json" {
		return fmt.Errorf("invalid format '%s': must be 'markdown' or 'json'", format)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
