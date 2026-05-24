package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	vault "github.com/hashicorp/vault/api"
)

const (
	defaultVaultMount  = "secret"
	defaultVaultPrefix = "cicd"
)

func BootstrapSecrets(ctx context.Context, serviceName string) error {
	addr := strings.TrimSpace(os.Getenv("VAULT_ADDR"))
	if addr == "" {
		return nil
	}

	client, err := newVaultClient(addr)
	if err != nil {
		return err
	}

	mount := EnvOr("VAULT_MOUNT", defaultVaultMount)
	prefix := strings.Trim(strings.TrimSpace(EnvOr("VAULT_KV_PREFIX", defaultVaultPrefix)), "/")

	paths := []string{prefix + "/shared"}
	if serviceName != "" {
		paths = append(paths, prefix+"/"+serviceName)
	}

	for _, path := range paths {
		if err := mergeVaultPath(ctx, client, mount, path); err != nil {
			return fmt.Errorf("vault read %s: %w", path, err)
		}
	}

	return nil
}

func newVaultClient(addr string) (*vault.Client, error) {
	cfg := vault.DefaultConfig()
	cfg.Address = addr
	cfg.Timeout = 10 * time.Second

	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if token := strings.TrimSpace(os.Getenv("VAULT_TOKEN")); token != "" {
		client.SetToken(token)
		return client, nil
	}

	roleID := strings.TrimSpace(os.Getenv("VAULT_ROLE_ID"))
	secretID := strings.TrimSpace(os.Getenv("VAULT_SECRET_ID"))
	if roleID != "" && secretID != "" {
		secret, err := client.Logical().Write("auth/approle/login", map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		})
		if err != nil {
			return nil, fmt.Errorf("vault approle login: %w", err)
		}
		if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
			return nil, fmt.Errorf("vault approle login: empty token")
		}
		client.SetToken(secret.Auth.ClientToken)
		return client, nil
	}

	return nil, fmt.Errorf("vault auth required: set VAULT_TOKEN or VAULT_ROLE_ID+VAULT_SECRET_ID")
}

func mergeVaultPath(ctx context.Context, client *vault.Client, mount, path string) error {
	secret, err := client.KVv2(mount).Get(ctx, path)
	if err != nil {
		if isOptionalPathError(err) {
			return nil
		}
		return err
	}
	if secret == nil || secret.Data == nil {
		return nil
	}

	for key, raw := range secret.Data {
		if os.Getenv(key) != "" {
			continue
		}
		value, ok := raw.(string)
		if !ok || strings.TrimSpace(value) == "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set env %s: %w", key, err)
		}
	}

	return nil
}

func isOptionalPathError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "404") ||
		strings.Contains(msg, "not found") ||
		strings.Contains(msg, "403") ||
		strings.Contains(msg, "permission denied")
}
