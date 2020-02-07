# `vault-template`

Render templated config files with secrets from [HashiCorp Vault](https://www.vaultproject.io/). Inspired by [vaultenv](https://github.com/channable/vaultenv).

* Define a template for your config file which contains secrets at development time.
* Use `vault-template` to render your config file template by fetching secrets from Vault at runtime.

## Usage

```text
Usage of ./vault-template:
  -o, --output string             The output file.
                                  Also configurable via OUTPUT_FILE.
  -t, --template string           The template file to render.
                                  Also configurable via TEMPLATE_FILE.
  -v, --vault string              Vault API endpoint.
                                  Also configurable via VAULT_ADDR.
                                  (default "http://127.0.0.1:8200")
  -f, --vault-token-file string   The file which contains the vault token.
                                  Also configurable via VAULT_TOKEN_FILE.
```

A [docker image is availabe on Dockerhub.](https://hub.docker.com/r/rplan/vault-template)

## Template

The templates will be rendered using the [Go template](https://golang.org/pkg/text/template/) mechanism. `vault-env` provides a special function for specifying secrets in the template:

```gotemplate
mySecretName = {{ vault "secret/mySecret" "name" }}
mySecretPassword = {{ vault "secret/mySecret" "password" }}
{{ range $name, $secret := vaultMap "secret/mySecret"}}
{{ $name }}: {{ $secret }}
{{- end }}
```

The `vault` function takes two string parameters which specify the path to the secret and the field inside to return.
The `vaultMap` function takes one string parameter which specify the path to the secret to return.

If the secret was created with `vault write secret/mySecret name=john password=secret` the resulting file would be:

```text
mySecretName = john
mySecretPassword = secret
mySecretName: john
mySecretPassword: secret
```

Also it is possible to use environment variables like `{{ .STAGE }}`.
