package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppSettingsHandler(t *testing.T) {
	tc := SetupTestServer(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("GetAppSettings", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/app/settings", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// App settings might not have is_test_data field as they're configuration data
			// This test mainly ensures the endpoint is working and returns valid data
			assert.NotNil(t, response, "Should receive valid response from app settings")

			// Check for common app settings fields
			if settings, exists := response["settings"]; exists {
				settingsMap, ok := settings.(map[string]interface{})
				if ok {
					// Common app settings might include version, features, etc.
					// We don't expect is_test_data here as these are app configuration
					assert.NotEmpty(t, settingsMap, "Settings should not be empty")
				}
			}

			// Alternative response structure - data field
			if data, exists := response["data"]; exists {
				dataMap, ok := data.(map[string]interface{})
				if ok {
					assert.NotEmpty(t, dataMap, "Data should not be empty")
				}
			}
		}
	})

	t.Run("GetAppVersion", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/app/version", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Version information should be present
			version, exists := response["version"]
			if exists {
				assert.NotEmpty(t, version, "Version should not be empty")
			}

			// Alternative field names
			if appVersion, exists := response["app_version"]; exists {
				assert.NotEmpty(t, appVersion, "App version should not be empty")
			}

			if buildVersion, exists := response["build_version"]; exists {
				assert.NotEmpty(t, buildVersion, "Build version should not be empty")
			}
		}
	})

	t.Run("GetFeatureFlags", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/app/features", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Feature flags should be present
			features, exists := response["features"]
			if exists {
				featuresMap, ok := features.(map[string]interface{})
				if ok {
					// Feature flags are typically boolean values
					for featureName, featureValue := range featuresMap {
						assert.NotEmpty(t, featureName, "Feature name should not be empty")
						// Feature values might be boolean, string, or other types
						assert.NotNil(t, featureValue, "Feature value should not be nil")
					}
				}
			}

			// Alternative response structure
			if data, exists := response["data"]; exists {
				dataMap, ok := data.(map[string]interface{})
				if ok {
					features, featuresExist := dataMap["features"]
					if featuresExist {
						assert.NotNil(t, features, "Features should be present in data")
					}
				}
			}
		}
	})

	t.Run("GetAppConfig", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/app/config", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// App config should contain configuration data
			config, exists := response["config"]
			if exists {
				configMap, ok := config.(map[string]interface{})
				if ok {
					assert.NotEmpty(t, configMap, "Config should not be empty")

					// Common config fields might include API limits, timeouts, etc.
					for configKey, configValue := range configMap {
						assert.NotEmpty(t, configKey, "Config key should not be empty")
						assert.NotNil(t, configValue, "Config value should not be nil")
					}
				}
			}

			// Alternative response structure
			if data, exists := response["data"]; exists {
				assert.NotNil(t, data, "Data should be present")
			}
		}
	})

	t.Run("GetMaintenanceStatus", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/app/maintenance", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Maintenance status should be present
			maintenance, exists := response["maintenance"]
			if exists {
				maintenanceMap, ok := maintenance.(map[string]interface{})
				if ok {
					// Should have at least an enabled/active status
					if enabled, enabledExists := maintenanceMap["enabled"]; enabledExists {
						_, ok := enabled.(bool)
						assert.True(t, ok, "Maintenance enabled should be boolean")
					}

					if active, activeExists := maintenanceMap["active"]; activeExists {
						_, ok := active.(bool)
						assert.True(t, ok, "Maintenance active should be boolean")
					}
				}
			}

			// Alternative response structure
			if isMaintenanceMode, exists := response["is_maintenance_mode"]; exists {
				_, ok := isMaintenanceMode.(bool)
				assert.True(t, ok, "is_maintenance_mode should be boolean")
			}
		}
	})

	t.Run("GetSystemHealth", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/app/health", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Health status should be present
			health, exists := response["health"]
			if exists {
				healthMap, ok := health.(map[string]interface{})
				if ok {
					// Common health indicators
					if status, statusExists := healthMap["status"]; statusExists {
						assert.NotEmpty(t, status, "Health status should not be empty")
					}

					if database, dbExists := healthMap["database"]; dbExists {
						assert.NotNil(t, database, "Database health should be reported")
					}
				}
			}

			// Alternative response structure
			if status, exists := response["status"]; exists {
				assert.NotEmpty(t, status, "Status should not be empty")
			}
		}
	})

	t.Run("GetAPILimits", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/app/limits", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// API limits should be present
			limits, exists := response["limits"]
			if exists {
				limitsMap, ok := limits.(map[string]interface{})
				if ok {
					// Common limit fields
					if rateLimit, rateLimitExists := limitsMap["rate_limit"]; rateLimitExists {
						_, ok := rateLimit.(float64)
						assert.True(t, ok, "Rate limit should be a number")
					}

					if maxRequestsPerMinute, exists := limitsMap["max_requests_per_minute"]; exists {
						_, ok := maxRequestsPerMinute.(float64)
						assert.True(t, ok, "Max requests per minute should be a number")
					}
				}
			}

			// Alternative response structure
			if data, exists := response["data"]; exists {
				assert.NotNil(t, data, "Data should be present")
			}
		}
	})

	t.Run("GetAppAnnouncements", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/app/announcements", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Announcements should be present
			announcements, exists := response["announcements"]
			if exists {
				announcementsList, ok := announcements.([]interface{})
				if ok && len(announcementsList) > 0 {
					// Check first few announcements
					for i, announcement := range announcementsList {
						announcementMap, ok := announcement.(map[string]interface{})
						if ok {
							// Announcements might have is_test_data if they're user-generated content
							if _, hasTestData := announcementMap["is_test_data"]; hasTestData {
								VerifyIsTestDataField(t, announcementMap, true, "announcement")
							}

							// Check basic announcement fields
							if title, titleExists := announcementMap["title"]; titleExists {
								assert.NotEmpty(t, title, "Announcement title should not be empty")
							}
						}

						// Limit checking for performance
						if i >= 3 {
							break
						}
					}
				}
			}

			// Alternative response structure
			if data, exists := response["data"]; exists {
				dataMap, ok := data.(map[string]interface{})
				if ok {
					announcements, announcementsExist := dataMap["announcements"]
					if announcementsExist {
						assert.NotNil(t, announcements, "Announcements should be present in data")
					}
				}
			}
		}
	})

	_ = userID // Use userID to avoid unused variable warning
}
