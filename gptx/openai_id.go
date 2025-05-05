package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const monitoringWarning string = `
User identification is required to monitor for abuse.
Monitoring is handled by OpenAI, For more information see:
https://platform.openai.com/docs/guides/safety-best-practices#end-user-ids
The user id is stored in the user's config directory.
OpenAI user id file: %s
`

var userID string = func() string {
	userID := filepath.Join(AppConfigDir, "user_id")
	warning := fmt.Sprintf(monitoringWarning, userID)
	if AppConfigDir == "" {
		warn(fmt.Errorf(Theme.Warning+"%s", warning))
		exit(fmt.Errorf("openai_id: user config dir not set"))
	}

	// create new user id file if it doesn't exist
	if _, err := os.Stat(userID); os.IsNotExist(err) {
		warn(fmt.Errorf(Theme.Warning+"%s", warning))
		if err := os.MkdirAll(AppConfigDir, 0755); err != nil {
			exit(fmt.Errorf("openai_id: %w", err))
		}
		id := uuid.New().String()
		if err := os.WriteFile(userID, []byte(id), 0644); err != nil {
			exit(fmt.Errorf("openai_id: %w", err))
		}
	}

	// read the user id from the file
	data, err := os.ReadFile(userID)
	if err != nil {
		warn(fmt.Errorf(Theme.Warning+"%s", warning))
		exit(fmt.Errorf("openai_id: %w", err))
	}

	// validate the user id
	id, err := uuid.Parse(string(data))
	if err != nil {
		warn(fmt.Errorf(Theme.Warning+"%s", warning))
		exit(fmt.Errorf("openai_id: %w", err))
	}
	return id.String()
}()
