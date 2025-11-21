package rules

import (
	"testing"

	"ingress-migration-analyzer/internal/models"
)

func TestGetConfigMapRules(t *testing.T) {
	rules := GetConfigMapRules()

	if len(rules) == 0 {
		t.Fatal("Expected ConfigMap rules, got none")
	}

	// Check that we have rules for each risk level
	var autoCount, manualCount, highRiskCount int
	for _, rule := range rules {
		switch rule.RiskLevel {
		case models.RiskAuto:
			autoCount++
		case models.RiskManual:
			manualCount++
		case models.RiskHigh:
			highRiskCount++
		}
	}

	if autoCount == 0 {
		t.Error("Expected at least one AUTO risk rule")
	}
	if manualCount == 0 {
		t.Error("Expected at least one MANUAL risk rule")
	}
	if highRiskCount == 0 {
		t.Error("Expected at least one HIGH_RISK rule")
	}
}

func TestMatchConfigMapSettings(t *testing.T) {
	tests := []struct {
		name       string
		configData map[string]string
		wantCount  int
		wantRisk   models.RiskLevel
	}{
		{
			name: "auto-migratable settings",
			configData: map[string]string{
				"ssl-protocols": "TLSv1.2 TLSv1.3",
				"use-http2":     "true",
				"hsts":          "true",
			},
			wantCount: 3,
			wantRisk:  models.RiskAuto,
		},
		{
			name: "manual review settings",
			configData: map[string]string{
				"proxy-buffer-size":     "8k",
				"proxy-connect-timeout": "60s",
				"ssl-reject-handshake":  "true",
			},
			wantCount: 3,
			wantRisk:  models.RiskManual,
		},
		{
			name: "mixed risk settings",
			configData: map[string]string{
				"ssl-protocols":   "TLSv1.2 TLSv1.3",
				"proxy-body-size": "50m",
				"main-snippet":    "custom nginx config",
			},
			wantCount: 3,
			wantRisk:  models.RiskHigh, // Highest risk wins
		},
		{
			name: "high-risk snippets",
			configData: map[string]string{
				"http-snippet":   "custom http config",
				"server-snippet": "custom server config",
				"stream-snippet": "custom stream config",
			},
			wantCount: 3,
			wantRisk:  models.RiskHigh,
		},
		{
			name: "no matching settings",
			configData: map[string]string{
				"unknown-setting-1": "value1",
				"unknown-setting-2": "value2",
			},
			wantCount: 0,
			wantRisk:  models.RiskAuto,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := MatchConfigMapSettings(tt.configData)

			if len(matches) != tt.wantCount {
				t.Errorf("MatchConfigMapSettings() returned %d matches, want %d", len(matches), tt.wantCount)
			}

			if len(matches) > 0 {
				riskLevel := GetHighestConfigMapRiskLevel(matches)
				if riskLevel != tt.wantRisk {
					t.Errorf("GetHighestConfigMapRiskLevel() = %v, want %v", riskLevel, tt.wantRisk)
				}
			}
		})
	}
}

func TestGetUnknownConfigMapSettings(t *testing.T) {
	configData := map[string]string{
		// Known settings
		"ssl-protocols":   "TLSv1.2 TLSv1.3",
		"proxy-body-size": "50m",

		// Unknown settings
		"custom-unknown-setting":  "value",
		"another-unknown-setting": "value",
	}

	unknown := GetUnknownConfigMapSettings(configData)

	expectedUnknown := []string{
		"custom-unknown-setting",
		"another-unknown-setting",
	}

	if len(unknown) != len(expectedUnknown) {
		t.Errorf("Expected %d unknown settings, got %d", len(expectedUnknown), len(unknown))
	}

	// Convert to map for easier checking
	unknownMap := make(map[string]bool)
	for _, u := range unknown {
		unknownMap[u] = true
	}

	for _, expected := range expectedUnknown {
		if !unknownMap[expected] {
			t.Errorf("Expected unknown setting %s not found", expected)
		}
	}
}

func TestGetConfigMapRuleByKey(t *testing.T) {
	// Test known key
	rule := GetConfigMapRuleByKey("ssl-protocols")
	if rule == nil {
		t.Fatal("Expected to find rule for ssl-protocols")
	}
	if rule.RiskLevel != models.RiskAuto {
		t.Errorf("Expected ssl-protocols to be AUTO risk, got %v", rule.RiskLevel)
	}

	// Test another known key with different risk
	rule = GetConfigMapRuleByKey("proxy-buffer-size")
	if rule == nil {
		t.Fatal("Expected to find rule for proxy-buffer-size")
	}
	if rule.RiskLevel != models.RiskManual {
		t.Errorf("Expected proxy-buffer-size to be MANUAL risk, got %v", rule.RiskLevel)
	}

	// Test high-risk key
	rule = GetConfigMapRuleByKey("main-snippet")
	if rule == nil {
		t.Fatal("Expected to find rule for main-snippet")
	}
	if rule.RiskLevel != models.RiskHigh {
		t.Errorf("Expected main-snippet to be HIGH_RISK, got %v", rule.RiskLevel)
	}

	// Test unknown key
	rule = GetConfigMapRuleByKey("unknown-setting")
	if rule != nil {
		t.Error("Expected nil for unknown setting, got rule")
	}
}

func TestGetHighestConfigMapRiskLevel(t *testing.T) {
	tests := []struct {
		name     string
		rules    []models.ConfigMapRule
		wantRisk models.RiskLevel
	}{
		{
			name:     "empty rules",
			rules:    []models.ConfigMapRule{},
			wantRisk: models.RiskAuto,
		},
		{
			name: "all auto",
			rules: []models.ConfigMapRule{
				{RiskLevel: models.RiskAuto},
				{RiskLevel: models.RiskAuto},
			},
			wantRisk: models.RiskAuto,
		},
		{
			name: "auto and manual",
			rules: []models.ConfigMapRule{
				{RiskLevel: models.RiskAuto},
				{RiskLevel: models.RiskManual},
			},
			wantRisk: models.RiskManual,
		},
		{
			name: "mixed with high",
			rules: []models.ConfigMapRule{
				{RiskLevel: models.RiskAuto},
				{RiskLevel: models.RiskManual},
				{RiskLevel: models.RiskHigh},
			},
			wantRisk: models.RiskHigh,
		},
		{
			name: "high risk first",
			rules: []models.ConfigMapRule{
				{RiskLevel: models.RiskHigh},
				{RiskLevel: models.RiskAuto},
			},
			wantRisk: models.RiskHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHighestConfigMapRiskLevel(tt.rules)
			if got != tt.wantRisk {
				t.Errorf("GetHighestConfigMapRiskLevel() = %v, want %v", got, tt.wantRisk)
			}
		})
	}
}
