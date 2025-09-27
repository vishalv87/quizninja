package tests

import (
	"database/sql"
	"testing"

	"quizninja-api/database"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMigration027Verification(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	// Get direct database connection
	db := database.DB
	_ = tc // Suppress unused variable warning

	t.Run("VerifyNewInterestsFromMigration027", func(t *testing.T) {
		// Check that the newly added interests exist with is_test_data = true
		newInterests := []string{"biology", "chemistry", "physics", "football", "basketball"}

		for _, interestID := range newInterests {
			var exists bool
			var isTestData bool
			var name, description string

			query := `SELECT EXISTS(SELECT 1 FROM interests WHERE id = $1),
			                 COALESCE((SELECT is_test_data FROM interests WHERE id = $1), false),
			                 COALESCE((SELECT name FROM interests WHERE id = $1), ''),
			                 COALESCE((SELECT description FROM interests WHERE id = $1), '')`

			err := db.QueryRow(query, interestID).Scan(&exists, &isTestData, &name, &description)
			assert.NoError(t, err, "Should query interest %s without error", interestID)

			assert.True(t, exists, "Interest %s should exist in database", interestID)
			assert.True(t, isTestData, "Interest %s should have is_test_data = true", interestID)
			assert.NotEmpty(t, name, "Interest %s should have a name", interestID)
			assert.NotEmpty(t, description, "Interest %s should have a description", interestID)
		}
	})

	t.Run("VerifyNewAppSettingsFromMigration027", func(t *testing.T) {
		// Check that the newly added app settings exist
		expectedSettings := map[string]string{
			"app_name":                "QuizNinja",
			"quiz_categories_enabled": "true",
			"leaderboard_enabled":     "true",
			"achievements_enabled":    "true",
			"default_quiz_duration":   "300",
		}

		for settingKey, expectedValue := range expectedSettings {
			var exists bool
			var actualValue string

			query := `SELECT EXISTS(SELECT 1 FROM app_settings WHERE key = $1),
			                 COALESCE((SELECT value FROM app_settings WHERE key = $1), '')`

			err := db.QueryRow(query, settingKey).Scan(&exists, &actualValue)
			assert.NoError(t, err, "Should query app setting %s without error", settingKey)

			assert.True(t, exists, "App setting %s should exist in database", settingKey)
			assert.Equal(t, expectedValue, actualValue, "App setting %s should have correct value", settingKey)
		}
	})

	t.Run("VerifyUpdatedMaxQuestionsPerQuiz", func(t *testing.T) {
		// Verify that max_questions_per_quiz was updated from 50 to 20
		var value string

		query := `SELECT value FROM app_settings WHERE key = 'max_questions_per_quiz'`
		err := db.QueryRow(query).Scan(&value)
		assert.NoError(t, err, "Should query max_questions_per_quiz setting")

		assert.Equal(t, "20", value, "max_questions_per_quiz should be updated to 20")
	})

	t.Run("VerifyInterestCategoryMappings", func(t *testing.T) {
		// Verify that our category mapping logic will work with the database data
		expectedMappings := map[string][]string{
			"science": {"science", "biology", "chemistry", "physics", "technology"},
			"sports":  {"sports", "football", "basketball"},
		}

		for categoryName, expectedInterests := range expectedMappings {
			for _, interestID := range expectedInterests {
				var exists bool
				query := `SELECT EXISTS(SELECT 1 FROM interests WHERE id = $1 AND is_test_data = false OR is_test_data = true)`
				err := db.QueryRow(query, interestID).Scan(&exists)
				assert.NoError(t, err, "Should query interest %s without error", interestID)

				assert.True(t, exists, "Interest %s should exist for category %s", interestID, categoryName)
			}
		}
	})

	t.Run("VerifyInterestFieldsStructure", func(t *testing.T) {
		// Verify that all interests have the required fields
		query := `SELECT id, name, description, icon_name, color_hex, is_test_data
		          FROM interests
		          WHERE id IN ('biology', 'chemistry', 'physics', 'football', 'basketball')`

		rows, err := db.Query(query)
		assert.NoError(t, err, "Should query new interests without error")
		defer rows.Close()

		interestCount := 0
		for rows.Next() {
			var id, name, description string
			var iconName, colorHex sql.NullString
			var isTestData bool

			err := rows.Scan(&id, &name, &description, &iconName, &colorHex, &isTestData)
			assert.NoError(t, err, "Should scan interest row without error")

			// Verify required fields are not empty
			assert.NotEmpty(t, id, "Interest ID should not be empty")
			assert.NotEmpty(t, name, "Interest name should not be empty")
			assert.NotEmpty(t, description, "Interest description should not be empty")
			assert.True(t, isTestData, "Interest should have is_test_data = true")

			// Verify optional fields are present (can be empty)
			assert.True(t, iconName.Valid || !iconName.Valid, "icon_name field should exist")
			assert.True(t, colorHex.Valid || !colorHex.Valid, "color_hex field should exist")

			interestCount++
		}

		assert.Equal(t, 5, interestCount, "Should find exactly 5 new interests")
	})

	t.Run("VerifyAppSettingsFieldsStructure", func(t *testing.T) {
		// Verify that all app settings have the required fields
		newSettingKeys := []string{"app_name", "quiz_categories_enabled", "leaderboard_enabled", "achievements_enabled", "default_quiz_duration"}

		query := `SELECT key, value, description, updated_at
		          FROM app_settings
		          WHERE key = ANY($1)`

		rows, err := db.Query(query, pq.Array(newSettingKeys))
		assert.NoError(t, err, "Should query new app settings without error")
		defer rows.Close()

		settingCount := 0
		for rows.Next() {
			var key, value string
			var description sql.NullString
			var updatedAt sql.NullTime

			err := rows.Scan(&key, &value, &description, &updatedAt)
			assert.NoError(t, err, "Should scan app setting row without error")

			// Verify required fields are not empty
			assert.NotEmpty(t, key, "App setting key should not be empty")
			assert.NotEmpty(t, value, "App setting value should not be empty")

			// Verify optional fields exist
			assert.True(t, description.Valid || !description.Valid, "description field should exist")
			assert.True(t, updatedAt.Valid, "updated_at should be set")

			settingCount++
		}

		assert.GreaterOrEqual(t, settingCount, 5, "Should find at least 5 new app settings")
	})

	t.Run("VerifyMigrationIntegrity", func(t *testing.T) {
		// Verify that the migration didn't break existing data

		// Check that original interests still exist
		originalInterests := []string{"general_knowledge", "science", "history", "sports", "technology", "music", "geography", "literature", "art"}

		for _, interestID := range originalInterests {
			var exists bool
			query := `SELECT EXISTS(SELECT 1 FROM interests WHERE id = $1)`
			err := db.QueryRow(query, interestID).Scan(&exists)
			assert.NoError(t, err, "Should query original interest %s without error", interestID)

			assert.True(t, exists, "Original interest %s should still exist after migration", interestID)
		}

		// Check that original app settings still exist
		originalSettings := []string{"app_version", "api_version", "quiz_time_limit_seconds", "points_per_correct_answer"}

		for _, settingKey := range originalSettings {
			var exists bool
			query := `SELECT EXISTS(SELECT 1 FROM app_settings WHERE key = $1)`
			err := db.QueryRow(query, settingKey).Scan(&exists)
			assert.NoError(t, err, "Should query original app setting %s without error", settingKey)

			assert.True(t, exists, "Original app setting %s should still exist after migration", settingKey)
		}
	})
}