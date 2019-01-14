#!/bin/bash
mockgen -destination=mocks/api/mock_vault_client.go -package=api github.com/actano/vault-template/api VaultClient
