# discord-bot

Bot do servidor Discord da Nullspace Lab, escrito em Go.

## Funcionalidades

- Embed de regras com reação para auto-aprovação de membros
- Atribuição automática de role ao reagir com ✅ no canal de regras
- Promoção diária de novos membros para membro oficial
- Comando `/ping`

## Variáveis de ambiente

| Variável | Descrição | Fonte em produção |
|---|---|---|
| `DISCORD_TOKEN` | Token do bot | Infisical (projeto `nullspace-bot`, env `prod`), via `ExternalSecret` |
| `DISCORD_GUILD_ID` | ID do servidor | `ConfigMap` (`k8s/base/configmap.yaml`) |
| `DISCORD_RULES_CHANNEL_ID` | ID do canal de regras | `ConfigMap` (`k8s/base/configmap.yaml`) |
| `DISCORD_NEW_MEMBER_ROLE_ID` | ID da role de novo membro | `ConfigMap` (`k8s/base/configmap.yaml`) |
| `DISCORD_OFICIAL_MEMBER_ROLE_ID` | ID da role de membro oficial | `ConfigMap` (`k8s/base/configmap.yaml`) |

Só `DISCORD_TOKEN` é segredo de verdade; as demais são IDs públicos do servidor Discord.

## Rodando localmente

```sh
cp .env.example .env   # preencha as variáveis
go run main.go
```

## Deploy

O bot roda em Kubernetes (k3s), gerenciado pelo ArgoCD via GitOps. A `Application`
(`kubernetes/applications/discord-bot.yaml`) vive no repo
[`vps-infrastructure`](https://github.com/eduardoschelive/vps-infrastructure) e aponta
para `k8s/overlays/prod` deste repo.

### Estrutura dos manifests

```
k8s/
├── base/
│   ├── namespace.yaml       # namespace discord-bot
│   ├── secret-store.yaml    # ClusterSecretStore infisical-discord-bot
│   ├── external-secret.yaml # DISCORD_TOKEN, via Infisical
│   ├── configmap.yaml       # demais variáveis (não-sensíveis)
│   └── deployment.yaml
└── overlays/
    └── prod/                # overlay de produção (pin da tag de imagem)
```

### Configuração do Infisical

`DISCORD_TOKEN` fica no projeto `nullspace-bot` (ambiente `prod`). O `ClusterSecretStore`
`infisical-discord-bot` autentica com a machine identity já compartilhada no cluster —
ela só precisa de acesso de leitura concedido a esse projeto.

`reloader.stakater.com/auto: "true"` no Deployment reinicia o pod automaticamente quando
o Secret ou o ConfigMap mudam.

## CI/CD

A cada push na `main`, o GitHub Actions (`.github/workflows/build.yaml`):

1. builda e publica a imagem no GHCR (`latest` e `<sha>`);
2. atualiza a tag em `k8s/overlays/prod/kustomization.yaml` e commita de volta
   (`ci: bump image tag to <sha> [skip ci]`) — o ArgoCD sincroniza a partir daí.

```
ghcr.io/nullspace-lab/discord-bot:latest
ghcr.io/nullspace-lab/discord-bot:<sha>
```

## Licença

[MIT](LICENSE)
