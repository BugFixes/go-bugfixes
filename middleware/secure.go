package middleware

import (
	"fmt"
	"net/http"
)

func (s *System) SetSecure(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SecureEnabled = enabled
}

func (s *System) SetSecureConfig(domain string, config SecureConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.SecureConfigs == nil {
		s.SecureConfigs = make(map[string]SecureConfig)
	}
	s.SecureConfigs[domain] = config
}

func (s *System) AddSecureConfig(domain string, config SecureConfig) {
	s.SetSecureConfig(domain, config)
}

func (s *System) isSecureEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.SecureEnabled
}

func (s *System) getSecureConfig(domain string) SecureConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if cfg, ok := s.SecureConfigs[domain]; ok {
		return cfg
	}

	return DefaultSecureConfig
}

func (s *System) Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.isSecureEnabled() {
			next.ServeHTTP(w, r)
			return
		}

		domain := r.Host
		config := s.getSecureConfig(domain)

		if config.XFrameOptions != "" {
			w.Header().Set("X-Frame-Options", config.XFrameOptions)
		}

		if config.XContentTypeOptions {
			w.Header().Set("X-Content-Type-Options", "nosniff")
		}

		if config.XXSSProtection != "" {
			w.Header().Set("X-XSS-Protection", config.XXSSProtection)
		}

		if config.HSTSEnabled {
			hsts := fmt.Sprintf("max-age=%d", int(config.HSTSMaxAge.Seconds()))
			if config.HSTSIncludeSubdomains {
				hsts += "; includeSubDomains"
			}
			if config.HSTSPreload {
				hsts += "; preload"
			}
			w.Header().Set("Strict-Transport-Security", hsts)
		}

		if config.CSP != "" {
			w.Header().Set("Content-Security-Policy", config.CSP)
		}

		if config.ReferrerPolicy != "" {
			w.Header().Set("Referrer-Policy", config.ReferrerPolicy)
		}

		if config.PermissionsPolicy != "" {
			w.Header().Set("Permissions-Policy", config.PermissionsPolicy)
		}

		next.ServeHTTP(w, r)
	})
}
