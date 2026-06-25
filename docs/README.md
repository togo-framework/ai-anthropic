# ai-anthropic — documentation

Anthropic Claude driver for togo ai

## Overview

Package anthropic is an Anthropic (Claude) driver for togo ai using the Messages API.
Blank-import it and set AI_DRIVER=anthropic + ANTHROPIC_API_KEY.

## Install

```bash
togo install togo-framework/ai-anthropic
```

Set `AI_DRIVER=anthropic`.

## Configuration

Environment variables read by this plugin (extracted from the source — see the gateway/provider docs for each value):

| Env var |
|---|
| `ANTHROPIC_API_KEY"` |

## Usage

```go
provider := ai.FromKernel(k)
resp, err := provider.Chat(ctx, []ai.Message{{Role: "user", Content: "Hello"}}, ai.Options{})
// streaming + provider.Embed(ctx, texts) for vectors; resp.Usage carries token counts
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai-anthropic
- Full README: ../README.md
