package main

import (
	"context"
	"fmt"

	"ingress-migration-analyzer/pkg/analyze"
	"ingress-migration-analyzer/pkg/common"
	"ingress-migration-analyzer/pkg/report"

	"github.com/spf13/cobra"
)

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Generate detailed annotation inventory and usage analysis",
	Long: `Generate a comprehensive inventory of all annotations used in your cluster.

This command provides detailed analysis beyond basic risk classification:
- Complete annotation usage statistics
- Value frequency analysis
- Cross-namespace usage patterns
- Unknown annotation identification
- Migration complexity heat map

This is particularly useful for:
- Understanding the full scope of nginx annotations in use
- Planning annotation-by-annotation migration strategies
- Identifying the most critical annotations to address first
- Creating comprehensive migration documentation`,
	RunE: runInventory,
}

func init() {
	// Add inventory flags
	inventoryCmd.Flags().BoolP("detailed", "d", false, "Include detailed value analysis")
	inventoryCmd.Flags().StringP("sort", "s", "usage", "Sort by: usage, risk, namespace, name")
	inventoryCmd.Flags().IntP("top", "t", 10, "Show top N most used annotations")

	inventoryCmd.Flags().StringVar(&output, "output", "./reports/", "Output directory for reports")
	inventoryCmd.Flags().StringVar(&format, "format", "json", "Output format (json recommended for inventory data)")

	rootCmd.AddCommand(inventoryCmd)
}

func runInventory(cmd *cobra.Command, args []string) error {
	detailed, _ := cmd.Flags().GetBool("detailed")
	sortBy, _ := cmd.Flags().GetString("sort")
	topN, _ := cmd.Flags().GetInt("top")

	fmt.Printf("📋 Generating annotation inventory...\n")
	fmt.Printf("📁 Output directory: %s\n", output)
	fmt.Printf("🔧 Detailed analysis: %v\n", detailed)
	fmt.Printf("📊 Sort by: %s\n", sortBy)
	fmt.Printf("🔝 Top N: %d\n", topN)

	if kubeconfig != "" {
		fmt.Printf("🔧 Kubeconfig: %s\n", kubeconfig)
	}
	if contextName != "" {
		fmt.Printf("🎯 Context: %s\n", contextName)
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

	// Build detailed inventory
	fmt.Println("\n📊 Building annotation inventory...")
	inventory := analyze.BuildAnnotationInventory(clusterAnalysis.Analyses)

	// Print console summary
	printInventorySummary(inventory, topN)

	// Add inventory to cluster analysis
	clusterAnalysis.Inventory = inventory

	// Generate detailed inventory report
	fmt.Println("\n📝 Generating inventory report...")
	var reportPath string

	switch format {
	case "markdown":
		generator := &report.InventoryMarkdownGenerator{
			Detailed:    detailed,
			SortBy:      sortBy,
			TopN:        topN,
			ContextName: contextName,
		}
		reportPath, err = generator.GenerateInventoryReport(inventory, clusterAnalysis, output)
	case "json":
		// JSON includes full inventory data automatically
		generator := report.NewJSONGenerator()
		reportPath, err = generator.GenerateReport(clusterAnalysis, output)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return fmt.Errorf("failed to generate inventory report: %w", err)
	}

	fmt.Printf("✅ Inventory analysis complete! Report saved to: %s\n", reportPath)

	return nil
}

func printInventorySummary(inventory *analyze.AnnotationInventory, topN int) {
	fmt.Println("\n📈 Annotation Inventory Summary:")
	fmt.Printf("   Total Unique Annotations: %d\n", inventory.Summary.TotalUniqueAnnotations)
	fmt.Printf("   Nginx Annotations: %d\n", inventory.Summary.NginxAnnotationsCount)
	fmt.Printf("   Unknown Nginx Annotations: %d\n", inventory.Summary.UnknownAnnotationsCount)

	if inventory.Summary.MostUsedAnnotation != "" {
		fmt.Printf("   Most Used: %s\n", inventory.Summary.MostUsedAnnotation)
	}

	// Show most critical annotations
	critical := inventory.GetMostCriticalAnnotations(topN)
	if len(critical) > 0 {
		fmt.Println("\n🚨 Most Critical Annotations (for migration):")
		for i, annotation := range critical {
			if i >= topN {
				break
			}
			fmt.Printf("   %d. %s (used %d times across %d namespaces)\n",
				i+1, annotation.Key, annotation.UsageCount, len(annotation.Namespaces))
		}
	}

	// Show annotations by risk
	byRisk := inventory.GetAnnotationsByRisk()
	if len(byRisk) > 0 {
		fmt.Println("\n📊 Annotations by Risk Level:")
		for riskLevel, annotations := range byRisk {
			fmt.Printf("   %s: %d annotations\n", riskLevel, len(annotations))
		}
	}
}
