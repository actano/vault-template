package main

import (
	"bytes"
	"github.com/Luzifer/rconfig"
	"github.com/hashicorp/vault/api"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	cfg = struct {
		VaultEndpoint string `flag:"vault,v" env:"VAULT_ADDR" default:"https://127.0.0.1:8200" description:"vault API endpoint"`
		VaultToken    string `flag:"vault-token,t" env:"VAULT_TOKEN" description:"The vault token to authenticate"`
		TemplateFile  string `flag:"template" env:"TEMPLATE_FILE" description:"The template file to render"`
		OutputFile    string `flag:"output,o" env:"OUTPUT_FILE" description:"The output file"`
	}{}
)

type vaultPath struct {
	Path string
	Field string
}

func config() {
	rconfig.Parse(&cfg)

	if cfg.VaultToken == "" {
		log.Fatalf("No vault token given")
	}

	if cfg.TemplateFile == "" {
		log.Fatalf("No template file given")
	}

	if cfg.OutputFile == "" {
		log.Fatalf("No output file given")
	}
}

func parsePath(path string) vaultPath {
	split := strings.Split(path, "#")

	if len(split) != 2 {
		log.Fatalf("Unable to parse path %s", path)
	}

	return vaultPath{
		Path: split[0],
		Field: split[1],
	}
}

func querySecret(client *api.Client, queryPath string) string {
	path := parsePath(queryPath)
	secret, err := client.Logical().Read(path.Path)

	if err != nil {
		log.Fatalf("Unable to read secret: %s", err)
	}

	secretValue, ok := secret.Data[path.Field]

	if !ok {
		log.Fatalf("Secrect at path '%s' has no field '%s'", path.Path, path.Field)
	}

	return secretValue.(string)
}

func main() {
	config()

	client, err := api.NewClient(&api.Config{
		Address: cfg.VaultEndpoint,
	})

	if err != nil {
		log.Fatalf("Unable to create client: %s", err)
	}

	client.SetToken(cfg.VaultToken)

	templateContent, err := ioutil.ReadFile(cfg.TemplateFile)

	if err != nil {
		log.Fatalf("Unable to read template file: %s", err)
	}

	query := func(queryPath string) string {
		return querySecret(client, queryPath)
	}

	funcMap := template.FuncMap{
		"vault": query,
	}

	tmpl, err := template.New("template").Funcs(funcMap).Parse(string(templateContent))

	if err != nil {
		log.Fatalf("Unable to create template: %s", err)
	}

	var outputBuffer bytes.Buffer

	if err := tmpl.Execute(&outputBuffer, nil); err != nil {
		log.Fatalf("Unable to execute template: %s", err)
	}

	outputFile, err := os.Create(cfg.OutputFile)

	if err != nil {
		log.Fatalf("Unable to write output file: %s", err)
	}

	defer outputFile.Close()

	outputFile.Write(outputBuffer.Bytes())
}
