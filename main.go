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
		VaultEndpoint  string `flag:"vault,v" env:"VAULT_ADDR" default:"https://127.0.0.1:8200" description:"Vault API endpoint. Also configurable via VAULT_ADDR."`
		VaultTokenFile string `flag:"vault-token-file,f" env:"VAULT_TOKEN_FILE" description:"The file which contains the vault token. Also configurable via VAULT_TOKEN_FILE."`
		TemplateFile   string `flag:"template,t" env:"TEMPLATE_FILE" description:"The template file to render. Also configurable via TEMPLATE_FILE."`
		OutputFile     string `flag:"output,o" env:"OUTPUT_FILE" description:"The output file. Also configurable via OUTPUT_FILE."`
	}{}
)

func usage(msg string) {
	println(msg)
	rconfig.Usage()
	os.Exit(1)
}

func config() {
	rconfig.Parse(&cfg)

	if cfg.VaultTokenFile == "" {
		usage("No vault token file given")
	}

	if cfg.TemplateFile == "" {
		usage("No template file given")
	}

	if cfg.OutputFile == "" {
		usage("No output file given")
	}
}

func querySecret(client *api.Client, path string, field string) string {
	secret, err := client.Logical().Read(path)

	if err != nil {
		log.Fatalf("Unable to read secret: %s", err)
	}

	secretValue, ok := secret.Data[field]

	if !ok {
		log.Fatalf("Secrect at path '%s' has no field '%s'", path, field)
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

	vaultToken, err := ioutil.ReadFile(cfg.VaultTokenFile)

	if err != nil {
		log.Fatalf("Unable to read vault token file: %s", err)
	}

	client.SetToken(strings.TrimSpace(string(vaultToken)))

	templateContent, err := ioutil.ReadFile(cfg.TemplateFile)

	if err != nil {
		log.Fatalf("Unable to read template file: %s", err)
	}

	query := func(path string, field string) string {
		return querySecret(client, path, field)
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
