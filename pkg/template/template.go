package template

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/nikita698/vault-template/pkg/api"
	"os"
	"strings"
	"text/template"
)

type VaultTemplateRenderer struct {
	vaultClient api.VaultClient
	openingDelim string
	closingDelim string
}

func NewVaultTemplateRenderer(vaultToken, vaultEndpoint string, openingDelim string, closingDelim string) (*VaultTemplateRenderer, error) {
	vaultClient, err := api.NewVaultClient(vaultEndpoint, string(vaultToken))

	if err != nil {
		return nil, err
	}

	return &VaultTemplateRenderer{
		vaultClient: vaultClient,
		openingDelim: openingDelim,
		closingDelim: closingDelim,
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

	if v.openingDelim != "" && v.closingDelim != "" {
	tmpl, err = template.
                New("template").
                Delims(v.openingDelim, v.closingDelim).
                Funcs(sprig.TxtFuncMap()).
                Funcs(funcMap).
                Parse(templateContent)
	}

	if err != nil {
		return "", err
	}

	var outputBuffer bytes.Buffer

	envMap := envToMap()
	if err := tmpl.Execute(&outputBuffer, envMap); err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}

func envToMap() map[string]string {
	envMap := map[string]string{}

	for _, v := range os.Environ() {
		splitV := strings.Split(v, "=")
		envMap[splitV[0]] = splitV[1]
	}

	return envMap
}
