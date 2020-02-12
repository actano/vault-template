package template

import (
	"github.com/actano/vault-template/mocks/api"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"gotest.tools/assert"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockVaultClient := api.NewMockVaultClient(mockCtrl)

	mockVaultClient.
		EXPECT().
		QuerySecret("secret/my/test/secret", "field1").
		Return("secret1", nil).
		Times(1)

	mockVaultClient.
		EXPECT().
		QuerySecretMap("secret/my/test/secret").
		Return(map[string]interface{}{"field1": "secret1"}, nil).
		Times(1)

	template := "The secret is '{{ vault \"secret/my/test/secret\" \"field1\" }}'."

	renderer := VaultTemplateRenderer{
		vaultClient: mockVaultClient,
	}

	result, err := renderer.RenderTemplate(template)

	assert.NilError(t, err)
	assert.Equal(t, result, "The secret is 'secret1'.")
}

func TestRenderTemplateQueryError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockVaultClient := api.NewMockVaultClient(mockCtrl)

	mockVaultClient.
		EXPECT().
		QuerySecret("secret/my/test/secret", "field1").
		Return("", errors.New("test error")).
		Times(1)

	mockVaultClient.
		EXPECT().
		QuerySecretMap("secret/my/test/secret").
		Return(nil, errors.New("test error")).
		Times(1)

	template := "The secret is '{{ vault \"secret/my/test/secret\" \"field1\" }}'."

	renderer := VaultTemplateRenderer{
		vaultClient: mockVaultClient,
	}

	_, err := renderer.RenderTemplate(template)

	assert.Assert(t, err != nil)
}
