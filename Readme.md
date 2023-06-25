# `vault-template`

Render templated config files with secrets from [HashiCorp Vault](https://www.vaultproject.io/). Inspired by [vaultenv](https://github.com/channable/vaultenv).

* Define a template for your config file which contains secrets at development time.
* Use `vault-template` to render your config file template by fetching secrets from Vault at runtime.

This repo has been forked to support vault API v2. 

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

First of all, suppose that the secret was created with `vault write secret/mySecret name=john password=secret`.

The templates will be rendered using the [Go template](https://golang.org/pkg/text/template/) mechanism.

Currently vault-template can render two functions:
- `vault`
- `vaultMap`

Also it is possible to use environment variables like `{{ .STAGE }}`.

The `vault` function takes two string parameters which specify the path to the secret and the field inside to return.

```gotemplate
mySecretName = {{ vault "<kv-engine>/data/mySecretPath" "name" }}
mySecretPassword = {{ vault "<kv-engine>/data/mySecretPath" "password" }}
```

```text
mySecretName = john
mySecretPassword = secret
```

The `vaultMap` function takes one string parameter which specify the path to the secret to return.

```gotemplate
{{ range $name, $secret := vaultMap "secret/mySecret"}}
{{ $name }}: {{ $secret }}
{{- end }}
```

```text
name: john
password: secret
```

More real example:

```gotemplate
---
# Common vars
{{- $customer    := .CUSTOMER }}
{{- $stage       := .STAGE }}
{{- $project     := .PROJECT }}
{{- $postgres    := print "kv/data/" $customer "/" $stage "/" $project "/postgres" }}
{{- $postgresMap := vaultMap $postgres }}

postgresql:
  postgresqlUsername: {{ $postgresMap.data.user }}
  postgresqlPassword: {{ $postgresMap.data.password }}
  postgresqlDatabase: {{ $postgresMap.data.db }}

app:
  postgres:
{{ range $name, $secret := $postgresMap }}
    {{ $name }}: {{ $secret }}
{{- end }}
```

And command that use this template in kubernetes:
```
CUSTOMER=internal STAGE=test PROJECT=myprj vault-template -o values.yaml -t values.tmpl -v "http://vault.default.svc.cluster.local:8200" -f token
```
