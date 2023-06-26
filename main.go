package main

import (
	"github.com/Luzifer/rconfig"
	"github.com/nikita698/vault-template/pkg/template"
	"io/ioutil"
	"log"
	"os"
)

var (
	cfg = struct {
		VaultEndpoint  string `flag:"vault,v" env:"VAULT_ADDR" default:"https://127.0.0.1:8200" description:"Vault API endpoint. Also configurable via VAULT_ADDR."`
		VaultTokenFile string `flag:"vault-token-file,f" env:"VAULT_TOKEN_FILE" description:"The file which contains the vault token. Also configurable via VAULT_TOKEN_FILE."`
		TemplateFile   string `flag:"template,t" env:"TEMPLATE_FILE" description:"The template file to render. Also configurable via TEMPLATE_FILE."`
		OutputFile     string `flag:"output,o" env:"OUTPUT_FILE" description:"The output file. Also configurable via OUTPUT_FILE."`
                OpeningDelim   string `flag:"opening-delim,s" env:"OPENING_DELIM" description:"Optional overwrite of the go template opening delimiter. The default is {{."`
                ClosingDelim   string `flag:"closing-delim,e" env:"CLOSING_DELIM" description:"Optional overwrite of the go template closing delimiter. The default is }}."`
		
	}{}
)

func usage(msg string) {
	println(msg)
	rconfig.Usage()
	os.Exit(1)
}

func config() {
	err := rconfig.Parse(&cfg)

	if err != nil {
		log.Fatalf("Error while parsing the command line arguments: %s", err)
	}

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

func main() {
	config()

	vaultToken, err := ioutil.ReadFile(cfg.VaultTokenFile)

	if err != nil {
		log.Fatalf("Unable to read vault token file: %s", err)
	}

	renderer, err := template.NewVaultTemplateRenderer(string(vaultToken), cfg.VaultEndpoint, cfg.OpeningDelim, cfg.ClosingDelim)

	if err != nil {
		log.Fatalf("Unable to create renderer: %s", err)
	}

	templateContent, err := ioutil.ReadFile(cfg.TemplateFile)

	if err != nil {
		log.Fatalf("Unable to read template file: %s", err)
	}

	renderedContent, err := renderer.RenderTemplate(string(templateContent))

	if err != nil {
		log.Fatalf("Unable to render template: %s", err)
	}

	outputFile, err := os.Create(cfg.OutputFile)

	if err != nil {
		log.Fatalf("Unable to write output file: %s", err)
	}

	defer func() {
		err := outputFile.Close()
		if err != nil {
			log.Fatalf("Error while closing the output file: %s", err)
		}
	}()

	_, err = outputFile.Write([]byte(renderedContent))

	if err != nil {
		log.Fatalf("Error while writing the output file: %s", err)
	}
}
