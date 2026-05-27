# discord-bot

Bot do servidor Discord da Nullspace Lab, escrito em Go.

## Funcionalidades

- Embed de regras com reação para auto-aprovação de membros
- Atribuição automática de role ao reagir com ✅ no canal de regras
- Promoção diária de novos membros para membro oficial
- Comando `/ping`

## Variáveis de ambiente

| Variável | Descrição |
|---|---|
| `DISCORD_TOKEN` | Token do bot |
| `DISCORD_GUILD_ID` | ID do servidor |
| `DISCORD_RULES_CHANNEL_ID` | ID do canal de regras |
| `DISCORD_NEW_MEMBER_ROLE_ID` | ID da role de novo membro |
| `DISCORD_OFICIAL_MEMBER_ROLE_ID` | ID da role de membro oficial |

Em desenvolvimento, crie um arquivo `.env` na raiz com essas variáveis.

## Rodando localmente

```sh
go run main.go
```

## Deploy

O bot roda em Kubernetes (k3s) gerenciado pelo ArgoCD. Os segredos são provisionados pelo [External Secrets Operator](https://external-secrets.io) a partir do HashiCorp Vault.

### Estrutura dos manifests

```
k8s/
├── base/          # recursos base (Deployment, SecretStore, ExternalSecret)
└── overlays/
    └── prod/      # overlay de produção
```

O ArgoCD deve apontar para `k8s/overlays/prod`.

### Configuração do Vault

Habilite o Kubernetes auth e crie a policy e role necessárias:

```sh
vault auth enable kubernetes

vault write auth/kubernetes/config \
  kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443"

vault policy write discord-bot - <<EOF
path "secret/data/discord-bot" { capabilities = ["read"] }
EOF

vault write auth/kubernetes/role/discord-bot \
  bound_service_account_names=discord-bot \
  bound_service_account_namespaces=discord-bot \
  policies=discord-bot \
  ttl=1h
```

Escreva os segredos no Vault:

```sh
vault kv put secret/discord-bot \
  DISCORD_TOKEN=... \
  DISCORD_GUILD_ID=... \
  DISCORD_RULES_CHANNEL_ID=... \
  DISCORD_NEW_MEMBER_ROLE_ID=... \
  DISCORD_OFICIAL_MEMBER_ROLE_ID=...
```

## CI/CD

A cada push na `main`, o GitHub Actions builda e publica a imagem no GHCR:

```
ghcr.io/nullspace-lab/discord-bot:latest
ghcr.io/nullspace-lab/discord-bot:<sha>
```

## Licença

[MIT](LICENSE)
