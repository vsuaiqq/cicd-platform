package config

import (
	"fmt"
	"strings"
)

func RequireNonEmpty(fields map[string]string) error {
	var missing []string
	for name, value := range fields {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("required configuration missing (set in .env or Vault): %s", strings.Join(missing, ", "))
	}
	return nil
}
