package main

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

	mockVaulClient := api.NewMockVaultClient(mockCtrl)

	mockVaulClient.
		EXPECT().
		QuerySecret("secret/my/test/secret", "field1").
		Return("secret1", nil).
		Times(1)

	template := "The secret is '{{ vault \"secret/my/test/secret\" \"field1\" }}'."

	result, err := renderTemplate(mockVaulClient, template)

	assert.NilError(t, err)
	assert.Equal(t, result.String(), "The secret is 'secret1'.")
}

func TestRenderTemplateQueryError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockVaulClient := api.NewMockVaultClient(mockCtrl)

	mockVaulClient.
		EXPECT().
		QuerySecret("secret/my/test/secret", "field1").
		Return("", errors.New("test error")).
		Times(1)

	template := "The secret is '{{ vault \"secret/my/test/secret\" \"field1\" }}'."

	_, err := renderTemplate(mockVaulClient, template)

	assert.Assert(t, err != nil)
}
