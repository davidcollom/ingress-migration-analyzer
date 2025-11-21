package rules

import (
	"ingress-migration-analyzer/internal/models"
)

// GetConfigMapRules returns the complete set of ConfigMap setting classification rules
func GetConfigMapRules() []models.ConfigMapRule {
	return []models.ConfigMapRule{
		// Tier A - AUTO (settings with established Gateway API equivalents)
		{
			Name:        "SSL Protocols",
			Key:         "ssl-protocols",
			RiskLevel:   models.RiskAuto,
			Description: "Defines the SSL/TLS protocol versions to use",
			MigrationNote: "Gateway API supports TLS configuration via Gateway listeners. " +
				"Most Gateway implementations allow TLS version configuration through Gateway or policy resources.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#ssl-protocols",
		},
		{
			Name:        "SSL Ciphers",
			Key:         "ssl-ciphers",
			RiskLevel:   models.RiskAuto,
			Description: "Specifies the enabled SSL/TLS ciphers",
			MigrationNote: "Gateway API TLS configuration can specify cipher suites. " +
				"Check your Gateway implementation's TLS policy support.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#ssl-ciphers",
		},
		{
			Name:        "HTTP2",
			Key:         "use-http2",
			RiskLevel:   models.RiskAuto,
			Description: "Enable HTTP/2 support",
			MigrationNote: "Gateway API supports HTTP/2 via Gateway listener protocol configuration. " +
				"Standard feature across most Gateway implementations.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#use-http2",
		},
		{
			Name:        "HSTS",
			Key:         "hsts",
			RiskLevel:   models.RiskAuto,
			Description: "Enable HTTP Strict Transport Security",
			MigrationNote: "HSTS can be configured via Gateway API response header modifiers or " +
				"implementation-specific security policies.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#hsts",
		},
		{
			Name:        "HSTS Max Age",
			Key:         "hsts-max-age",
			RiskLevel:   models.RiskAuto,
			Description: "Sets the max-age directive for HSTS",
			MigrationNote: "Configure via Gateway API response header modifiers. " +
				"Standard HTTPRoute filter functionality.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#hsts-max-age",
		},
		{
			Name:        "HSTS Include Subdomains",
			Key:         "hsts-include-subdomains",
			RiskLevel:   models.RiskAuto,
			Description: "Add includeSubDomains directive to HSTS header",
			MigrationNote: "Configure via Gateway API response header modifiers. " +
				"Can be set through HTTPRoute filters.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#hsts-include-subdomains",
		},
		{
			Name:        "HSTS Preload",
			Key:         "hsts-preload",
			RiskLevel:   models.RiskAuto,
			Description: "Add preload directive to HSTS header",
			MigrationNote: "Configure via Gateway API response header modifiers. " +
				"Standard HTTPRoute filter functionality.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#hsts-preload",
		},

		// Tier B - MANUAL (medium complexity, requires review)
		{
			Name:        "Proxy Buffer Size",
			Key:         "proxy-buffer-size",
			RiskLevel:   models.RiskManual,
			Description: "Sets the buffer size for reading the first part of the response",
			MigrationNote: "No standardized Gateway API equivalent. Gateway implementations may support " +
				"buffer configurations via vendor-specific policies. Check your Gateway documentation.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#proxy-buffer-size",
		},
		{
			Name:        "Proxy Buffers",
			Key:         "proxy-buffers-number",
			RiskLevel:   models.RiskManual,
			Description: "Sets the number of buffers used for reading response from backend",
			MigrationNote: "Implementation-specific setting. Review if your application requires " +
				"specific buffering behavior and configure at Gateway or service mesh level.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#proxy-buffers-number",
		},
		{
			Name:        "Proxy Body Size",
			Key:         "proxy-body-size",
			RiskLevel:   models.RiskManual,
			Description: "Maximum size allowed for client request body",
			MigrationNote: "Gateway implementations may support request size limits via vendor-specific policies. " +
				"Check your Gateway documentation or implement at application level.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#proxy-body-size",
		},
		{
			Name:        "Proxy Connect Timeout",
			Key:         "proxy-connect-timeout",
			RiskLevel:   models.RiskManual,
			Description: "Timeout for establishing connection to backend",
			MigrationNote: "Check Gateway implementation support for connection timeouts. " +
				"May require vendor-specific BackendPolicy or service mesh configuration.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#proxy-connect-timeout",
		},
		{
			Name:        "Proxy Read Timeout",
			Key:         "proxy-read-timeout",
			RiskLevel:   models.RiskManual,
			Description: "Timeout for reading response from backend",
			MigrationNote: "Gateway API may support timeouts via implementation-specific policies. " +
				"Consider BackendPolicy extensions or service mesh timeout configuration.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#proxy-read-timeout",
		},
		{
			Name:        "Proxy Send Timeout",
			Key:         "proxy-send-timeout",
			RiskLevel:   models.RiskManual,
			Description: "Timeout for transmitting request to backend",
			MigrationNote: "Similar to read timeout - check Gateway implementation policy support. " +
				"May need vendor-specific policies or service mesh configuration.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#proxy-send-timeout",
		},
		{
			Name:        "Client Body Buffer Size",
			Key:         "client-body-buffer-size",
			RiskLevel:   models.RiskManual,
			Description: "Buffer size for reading client request body",
			MigrationNote: "Implementation-specific setting. Review if your application requires " +
				"specific buffering behavior and implement accordingly.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#client-body-buffer-size",
		},
		{
			Name:        "Client Header Buffer Size",
			Key:         "client-header-buffer-size",
			RiskLevel:   models.RiskManual,
			Description: "Buffer size for reading client request header",
			MigrationNote: "Gateway implementations may have different header size limits. " +
				"Check your Gateway documentation for header size configuration.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#client-header-buffer-size",
		},
		{
			Name:        "Large Client Header Buffers",
			Key:         "large-client-header-buffers",
			RiskLevel:   models.RiskManual,
			Description: "Maximum number and size of buffers for large client headers",
			MigrationNote: "Review Gateway implementation limits for header sizes. " +
				"May require vendor-specific configuration or policy.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#large-client-header-buffers",
		},
		{
			Name:        "Keep Alive",
			Key:         "keep-alive",
			RiskLevel:   models.RiskManual,
			Description: "Timeout for keep-alive connections with clients",
			MigrationNote: "Gateway implementations handle keep-alive differently. " +
				"Check your Gateway's connection management configuration.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#keep-alive",
		},
		{
			Name:        "Keep Alive Requests",
			Key:         "keep-alive-requests",
			RiskLevel:   models.RiskManual,
			Description: "Maximum number of requests through one keep-alive connection",
			MigrationNote: "Gateway-specific configuration. Review your Gateway implementation's " +
				"connection pooling and keep-alive settings.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#keep-alive-requests",
		},
		{
			Name:        "Upstream Keep Alive",
			Key:         "upstream-keepalive-connections",
			RiskLevel:   models.RiskManual,
			Description: "Maximum number of idle keepalive connections to upstream servers",
			MigrationNote: "Backend connection pooling is implementation-specific. " +
				"Check your Gateway's backend connection management settings.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#upstream-keepalive-connections",
		},
		{
			Name:        "Worker Processes",
			Key:         "worker-processes",
			RiskLevel:   models.RiskManual,
			Description: "Number of NGINX worker processes",
			MigrationNote: "Gateway implementations handle worker/thread configuration differently. " +
				"Review your Gateway's scaling and performance tuning documentation.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#worker-processes",
		},
		{
			Name:        "Worker Connections",
			Key:         "max-worker-connections",
			RiskLevel:   models.RiskManual,
			Description: "Maximum number of simultaneous connections per worker",
			MigrationNote: "Gateway-specific capacity planning. Review your Gateway implementation's " +
				"concurrency and resource limit settings.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#max-worker-connections",
		},
		{
			Name:        "SSL Reject Handshake",
			Key:         "ssl-reject-handshake",
			RiskLevel:   models.RiskManual,
			Description: "Reject SSL handshake for requests without valid certificates",
			MigrationNote: "TLS handshake behavior is Gateway-specific. Check your Gateway's " +
				"TLS policy and certificate validation configuration options.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#ssl-reject-handshake",
		},
		{
			Name:        "Enable OCSP",
			Key:         "enable-ocsp",
			RiskLevel:   models.RiskManual,
			Description: "Enable OCSP stapling for SSL certificates",
			MigrationNote: "OCSP support varies by Gateway implementation. " +
				"Check your Gateway's TLS certificate validation features.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#enable-ocsp",
		},
		{
			Name:        "SSL Session Cache",
			Key:         "ssl-session-cache",
			RiskLevel:   models.RiskManual,
			Description: "Enable shared SSL session cache",
			MigrationNote: "TLS session management is implementation-specific. " +
				"Review your Gateway's TLS performance optimization settings.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#ssl-session-cache",
		},
		{
			Name:        "SSL Session Timeout",
			Key:         "ssl-session-timeout",
			RiskLevel:   models.RiskManual,
			Description: "Timeout for reusing SSL session parameters",
			MigrationNote: "TLS session timeout configuration is Gateway-specific. " +
				"Check your Gateway's TLS session management settings.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#ssl-session-timeout",
		},
		{
			Name:        "Enable Brotli",
			Key:         "enable-brotli",
			RiskLevel:   models.RiskManual,
			Description: "Enable Brotli compression",
			MigrationNote: "Compression support varies by Gateway implementation. " +
				"Some Gateways support compression policies, or implement at application level.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#enable-brotli",
		},
		{
			Name:        "Gzip Level",
			Key:         "gzip-level",
			RiskLevel:   models.RiskManual,
			Description: "Gzip compression level",
			MigrationNote: "Compression is typically handled at application or Gateway level. " +
				"Check if your Gateway supports compression policies.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#gzip-level",
		},
		{
			Name:        "Gzip Types",
			Key:         "gzip-types",
			RiskLevel:   models.RiskManual,
			Description: "MIME types to compress with gzip",
			MigrationNote: "Review Gateway compression support or implement compression " +
				"at application level for better control.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#gzip-types",
		},
		{
			Name:        "Log Format",
			Key:         "log-format-upstream",
			RiskLevel:   models.RiskManual,
			Description: "Custom log format for upstream traffic",
			MigrationNote: "Gateway implementations have different logging mechanisms. " +
				"Review your Gateway's logging configuration and format options.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#log-format-upstream",
		},
		{
			Name:        "Access Log Path",
			Key:         "access-log-path",
			RiskLevel:   models.RiskManual,
			Description: "Path for access logs",
			MigrationNote: "Gateway logging destinations are implementation-specific. " +
				"Configure according to your Gateway's logging documentation.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#access-log-path",
		},
		{
			Name:        "Error Log Path",
			Key:         "error-log-path",
			RiskLevel:   models.RiskManual,
			Description: "Path for error logs",
			MigrationNote: "Gateway error logging configuration varies. " +
				"Review your Gateway's error logging and monitoring options.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#error-log-path",
		},

		// Tier C - HIGH_RISK (complex configurations needing careful planning)
		{
			Name:        "Main Snippet",
			Key:         "main-snippet",
			RiskLevel:   models.RiskHigh,
			Description: "Custom NGINX main context configuration",
			MigrationNote: "Main snippets contain custom NGINX configuration that has no Gateway API equivalent. " +
				"Review the configuration and implement equivalent functionality using Gateway-level policies, " +
				"infrastructure changes, or consider staying with NGINX Inc commercial controller.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#main-snippet",
		},
		{
			Name:        "HTTP Snippet",
			Key:         "http-snippet",
			RiskLevel:   models.RiskHigh,
			Description: "Custom NGINX http context configuration",
			MigrationNote: "HTTP snippets affect global behavior and have no direct Gateway API mapping. " +
				"Requires careful analysis and potential migration to Gateway-level policies or infrastructure changes.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#http-snippet",
		},
		{
			Name:        "Server Snippet",
			Key:         "server-snippet",
			RiskLevel:   models.RiskHigh,
			Description: "Custom NGINX server context configuration",
			MigrationNote: "Server snippets contain custom NGINX server block configuration. " +
				"Review and map to Gateway API policies, service mesh, or application-level changes.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#server-snippet",
		},
		{
			Name:        "Stream Snippet",
			Key:         "stream-snippet",
			RiskLevel:   models.RiskHigh,
			Description: "Custom NGINX stream context configuration for TCP/UDP",
			MigrationNote: "Stream snippets are for Layer 4 routing. Gateway API supports TCP/UDP via " +
				"TCPRoute/UDPRoute, but custom stream logic requires complete reimplementation.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#stream-snippet",
		},
		{
			Name:        "Plugins",
			Key:         "plugins",
			RiskLevel:   models.RiskHigh,
			Description: "Enable custom Lua plugins",
			MigrationNote: "Custom plugins require complete reimplementation. Consider Gateway API " +
				"extension mechanisms, service mesh capabilities, or Gateway-specific plugin systems.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#plugins",
		},
		{
			Name:        "Lua Shared Dicts",
			Key:         "lua-shared-dicts",
			RiskLevel:   models.RiskHigh,
			Description: "Define Lua shared memory dictionaries",
			MigrationNote: "Lua-specific functionality requires complete redesign. " +
				"Evaluate Gateway extension points or service mesh for equivalent functionality.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#lua-shared-dicts",
		},
		{
			Name:        "Custom HTTP Errors",
			Key:         "custom-http-errors",
			RiskLevel:   models.RiskHigh,
			Description: "Custom error page backend service",
			MigrationNote: "Custom error handling varies significantly across Gateway implementations. " +
				"May require custom filters, service mesh error handling, or application-level changes.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#custom-http-errors",
		},
		{
			Name:        "Default Backend",
			Key:         "default-backend-service",
			RiskLevel:   models.RiskHigh,
			Description: "Override default backend for 404 responses",
			MigrationNote: "Gateway API doesn't have a direct equivalent for default backends. " +
				"Requires implementing catch-all routes or Gateway-specific fallback mechanisms.",
			SourceURL: "https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#default-backend-service",
		},
	}
}

// GetConfigMapRuleByKey returns the rule that matches a ConfigMap key
func GetConfigMapRuleByKey(key string) *models.ConfigMapRule {
	rules := GetConfigMapRules()

	for _, rule := range rules {
		if rule.Key == key {
			return &rule
		}
	}

	return nil
}

// MatchConfigMapSettings finds all rules that match the given ConfigMap data
func MatchConfigMapSettings(configData map[string]string) []models.ConfigMapRule {
	var matchedRules []models.ConfigMapRule
	rules := GetConfigMapRules()

	for key := range configData {
		for _, rule := range rules {
			if rule.Key == key {
				matchedRules = append(matchedRules, rule)
				break
			}
		}
	}

	return matchedRules
}

// GetUnknownConfigMapSettings identifies ConfigMap settings not in our rules
func GetUnknownConfigMapSettings(configData map[string]string) []string {
	var unknown []string
	rules := GetConfigMapRules()

	// Create a set of known keys
	knownKeys := make(map[string]bool)
	for _, rule := range rules {
		knownKeys[rule.Key] = true
	}

	for key := range configData {
		if !knownKeys[key] {
			unknown = append(unknown, key)
		}
	}

	return unknown
}

// GetHighestConfigMapRiskLevel determines the highest risk level from a set of ConfigMap rules
func GetHighestConfigMapRiskLevel(rules []models.ConfigMapRule) models.RiskLevel {
	if len(rules) == 0 {
		return models.RiskAuto
	}

	highestRisk := models.RiskAuto

	for _, rule := range rules {
		switch rule.RiskLevel {
		case models.RiskHigh:
			return models.RiskHigh // Highest possible, return immediately
		case models.RiskManual:
			highestRisk = models.RiskManual
		}
	}

	return highestRisk
}
