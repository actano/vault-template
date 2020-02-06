package template

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/actano/vault-template/pkg/api"
	"text/template"
)

type VaultTemplateRenderer struct {
	vaultClient api.VaultClient
}

func NewVaultTemplateRenderer(vaultToken, vaultEndpoint string) (*VaultTemplateRenderer, error) {
	vaultClient, err := api.NewVaultClient(vaultEndpoint, string(vaultToken))

	if err != nil {
		return nil, err
	}

    return &VaultTemplateRenderer{
        vaultClient: vaultClient,
	}, nil
}

func (v *VaultTemplateRenderer) RenderTemplate(templateContent string) (string, error) {
	funcMap := template.FuncMap{
		"vault":    v.vaultClient.QuerySecret,
		"vaultMap": v.vaultClient.QuerySecretMap,
	}

	tmpl, err := template.
		New("template").
		Funcs(sprig.TxtFuncMap()).
		Funcs(funcMap).
		Parse(templateContent)

	if err != nil {
		return "", err
	}

	var outputBuffer bytes.Buffer

	if err := tmpl.Execute(&outputBuffer, nil); err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}
