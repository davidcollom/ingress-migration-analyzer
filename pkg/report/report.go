package report

import (
	"fmt"
	"ingress-migration-analyzer/internal/models"
)

var reportFormats = []string{"markdown", "json"}

func GetReportFormats() []string {
	return reportFormats
}

func GenerateReport(clusterAnalysis *models.ClusterAnalysis, outputDir, format string) error {
	var generator ReportInterface

	switch format {
	case "markdown":
		generator = NewMarkdownGenerator()
	case "json":
		generator = NewJSONGenerator()
	default:
		return fmt.Errorf("unsupported report format: %s", format)
	}

	_, err := generator.GenerateReport(clusterAnalysis, outputDir)
	return err
}

type ReportInterface interface {
	GenerateReport(analysis *models.ClusterAnalysis, outputDir string) (string, error)
}
