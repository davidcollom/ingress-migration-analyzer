package main

import (
	"context"
	"fmt"
	"ingress-migration-analyzer/pkg/analyze"
	"ingress-migration-analyzer/pkg/common"
	"ingress-migration-analyzer/pkg/report"

	"github.com/spf13/cobra"
)

func init() {
	// Scan command flags
	scanCmd.Flags().StringVar(&output, "output", "./reports/", "Output directory for reports")
	scanCmd.Flags().StringVar(&format, "format", "markdown", "Output format (markdown|json)")

	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan cluster for ingress-nginx usage",
	Long: `Scan the Kubernetes cluster for ingress-nginx resources and generate
a migration complexity analysis report.

This command will:
- Connect to your Kubernetes cluster
- Discover all ingress-nginx resources
- Analyze annotation complexity
- Generate a detailed migration report`,
	RunE: runScan,
}

func runScan(cmd *cobra.Command, args []string) error {
	fmt.Printf("🔍 Starting ingress-nginx migration analysis...\n")
	fmt.Printf("📁 Output directory: %s\n", output)
	fmt.Printf("📄 Format: %s\n", format)

	if kubeconfig != "" {
		fmt.Printf("🔧 Kubeconfig: %s\n", kubeconfig)
	}

	if contextName != "" {
		fmt.Printf("🎯 Context: %s\n", contextName)
	}

	if namespace != "" {
		fmt.Printf("📦 Namespace: %s\n", namespace)
	} else {
		fmt.Printf("📦 Scanning all namespaces\n")
	}

	if ingressClassName != "" {
		fmt.Printf("📦 Ingress Class Name: %s\n", ingressClassName)
	} else {
		fmt.Printf("📦 Using default ingress class name: nginx\n")
		ingressClassName = "nginx"
	}

	if controllerNamespace != "" {
		fmt.Printf("📦 Ingress Controller Namespace: %s\n", controllerNamespace)
	} else {
		fmt.Printf("📦 Using default ingress controller namespace: ingress-nginx\n")
		controllerNamespace = "ingress-nginx"
	}

	// Validate flags
	if err := validateFlags(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Create Kubernetes client with validation
	fmt.Println("\n🔌 Testing Kubernetes connection...")
	client, err := common.CreateAnalyzerClient(kubeconfig, contextName)
	if err != nil {
		return err
	}

	// Create analyzer and run analysis
	analyzer := analyze.NewAnalyzer(client, namespace, ingressClassName, controllerNamespace)
	clusterAnalysis, err := analyzer.AnalyzeCluster(context.Background())
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Generate report
	fmt.Println("\n📝 Generating report...")
	var reportPath string

	err = report.GenerateReport(clusterAnalysis, output, format)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	fmt.Printf("✅ Analysis complete! Report saved to: %s\n", reportPath)

	if clusterAnalysis.Summary.HighRiskCount > 0 {
		fmt.Printf("\n⚠️  Warning: Found %d high-risk resources requiring careful migration planning\n",
			clusterAnalysis.Summary.HighRiskCount)
	}

	return nil
}
