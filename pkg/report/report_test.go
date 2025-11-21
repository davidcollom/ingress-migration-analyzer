package report

import (
	"ingress-migration-analyzer/internal/models"
	"os"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestGenerateReport(t *testing.T) {
	analysis := &models.ClusterAnalysis{}
	outputDir, err := os.MkdirTemp(t.TempDir(), "report")
	require.NoError(t, err)

	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{
			name:    "markdown format",
			format:  "markdown",
			wantErr: false,
		},
		{
			name:    "json format",
			format:  "json",
			wantErr: false,
		},
		{
			name:    "unsupported format",
			format:  "xml",
			wantErr: true,
		},
		{
			name:    "empty format",
			format:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateReport(analysis, outputDir, tt.format)

			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, "unsupported report format")
				if tt.format != "" {
					assert.ErrorContains(t, err, tt.format)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
