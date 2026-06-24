# ai-anthropic

Anthropic (Claude) driver for the togo `ai` plugin, using the Messages API.

```bash
togo install togo-framework/ai-anthropic
```

Set `AI_DRIVER=anthropic` and `ANTHROPIC_API_KEY=...`. Default model `claude-3-5-sonnet-latest`. Token usage is reported via `ai.Usage`. (Anthropic has no embeddings API — use `ai-openai` for `Embed`.)

MIT © ToGO
