package api

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
)

type VaultClient interface {
	QuerySecret(path string, field string) (string, error)
	QuerySecretMap(path string) (map[string]interface{}, error)
}

type vaultClient struct {
	apiClient *api.Client
}

func NewVaultClient(vaultEndpoint string, vaultToken string) (VaultClient, error) {
	apiClient, err := api.NewClient(&api.Config{
		Address: vaultEndpoint,
	})

	if err != nil {
		return nil, err
	}

	apiClient.SetToken(strings.TrimSpace(vaultToken))

	vaultClient := &vaultClient{
		apiClient: apiClient,
	}

	return vaultClient, nil
}

func (c *vaultClient) QuerySecretMap(path string) (map[string]interface{}, error) {
	secret, err := c.apiClient.Logical().Read(path)

	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("path '%s' is not found", path)
	}

	return secret.Data, nil
}

func (c *vaultClient) QuerySecret(path string, field string) (string, error) {
	secret, err := c.apiClient.Logical().Read(path)

	if err != nil {
		return "", err
	}

	secretValue, ok := secret.Data[field]

	if !ok {
		return "", fmt.Errorf("secret at path '%s' has no field '%s'", path, field)
	}

	return secretValue.(string), nil
}
