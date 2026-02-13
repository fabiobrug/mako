# Mako Setup Guide

This guide will help you configure Mako with your preferred AI provider.

## Quick Start

1. **Navigate to CLI directory**
   ```bash
   cd apps/cli
   ```

2. **Copy the configuration template**
   ```bash
   cp .env.example .env
   ```

3. **Choose your AI provider** (see options below)

4. **Edit the `.env` file** with your provider settings

5. **Start Mako**
   ```bash
   ./mako
   ```

## Provider Setup Instructions

### Option 1: Ollama (Local, Free) ⭐ Recommended for Privacy

**Perfect for:** Privacy-conscious users, offline work, no API costs

```bash
# 1. Install Ollama
curl https://ollama.ai/install.sh | sh

# 2. Pull a model (choose one)
ollama pull llama3.2        # Fast, good quality (2GB)
ollama pull llama3.2:1b     # Smallest, fastest (1GB)
ollama pull mistral         # Alternative option (4GB)

# 3. Navigate to CLI directory and configure Mako
cd apps/cli
cat > .env << EOF
LLM_PROVIDER=ollama
LLM_MODEL=llama3.2
LLM_API_BASE=http://localhost:11434
EOF

# 4. Test it
./mako
# Inside Mako: mako ask "list files"
```

**Pros:**
- ✅ Completely free
- ✅ Works offline
- ✅ Data never leaves your machine
- ✅ No rate limits
- ✅ No API key needed

**Cons:**
- ⚠️ Requires ~2-8GB disk space per model
- ⚠️ Slower on older hardware
- ⚠️ Quality may vary by model

---

### Option 2: OpenAI (Cloud, Paid)

**Perfect for:** Best quality, production use, GPT-4 access

```bash
# 1. Get API key from https://platform.openai.com/api-keys

# 2. Navigate to CLI directory and configure Mako
cd apps/cli
cat > .env << EOF
LLM_PROVIDER=openai
LLM_MODEL=gpt-4o-mini
LLM_API_KEY=sk-your-api-key-here
EOF

# 3. Start Mako
./mako
```

**Model Options:**
- `gpt-4o-mini` - Fast, cost-effective (recommended)
- `gpt-4o` - Best quality, higher cost
- `gpt-3.5-turbo` - Cheapest, good for simple tasks

**Pricing:** https://openai.com/pricing

---

### Option 3: Anthropic Claude (Cloud, Paid)

**Perfect for:** Great reasoning, detailed explanations, safety-focused

```bash
# 1. Get API key from https://console.anthropic.com/

# 2. Navigate to CLI directory and configure Mako
cd apps/cli
cat > .env << EOF
LLM_PROVIDER=anthropic
LLM_MODEL=claude-3-5-haiku-20241022
LLM_API_KEY=sk-ant-your-key-here
EOF

# 3. Start Mako
./mako
```

**Model Options:**
- `claude-3-5-haiku-20241022` - Fast, cost-effective (recommended)
- `claude-3-5-sonnet-20241022` - Best quality, balanced cost
- `claude-3-opus-20240229` - Highest quality, premium pricing

**Pricing:** https://www.anthropic.com/pricing

---

### Option 4: Google Gemini (Cloud, Free Tier Available)

**Perfect for:** Trying Mako without cost, good balance of quality and speed

```bash
# 1. Get API key from https://ai.google.dev/

# 2. Navigate to CLI directory and configure Mako
cd apps/cli
cat > .env << EOF
LLM_PROVIDER=gemini
LLM_MODEL=gemini-2.5-flash
LLM_API_KEY=your-gemini-key-here
EOF

# 3. Start Mako
./mako
```

**Model Options:**
- `gemini-2.5-flash` - Fast, efficient (recommended)
- `gemini-2.0-flash-exp` - Experimental features
- `gemini-1.5-pro` - Highest quality

**Pricing:** https://ai.google.dev/pricing

---

### Option 5: OpenRouter (Cloud, Paid)

**Perfect for:** Access to multiple models through one API, cost optimization

```bash
# 1. Get API key from https://openrouter.ai/keys

# 2. Navigate to CLI directory and configure Mako
cd apps/cli
cat > .env << EOF
LLM_PROVIDER=openrouter
LLM_MODEL=deepseek/deepseek-chat
LLM_API_KEY=sk-or-v1-your-key-here
EOF

# 3. Start Mako
./mako
```

**Popular Models:**
- `deepseek/deepseek-chat` - Very cost-effective
- `anthropic/claude-3.5-sonnet` - Access Claude via OpenRouter
- `openai/gpt-4o` - Access GPT-4 via OpenRouter

**Pricing:** https://openrouter.ai/docs/pricing

---

### Option 6: DeepSeek (Cloud, Paid)

**Perfect for:** Cost-effective alternative to other cloud providers

```bash
# 1. Get API key from https://platform.deepseek.com/

# 2. Navigate to CLI directory and configure Mako
cd apps/cli
cat > .env << EOF
LLM_PROVIDER=deepseek
LLM_MODEL=deepseek-chat
LLM_API_KEY=your-deepseek-key-here
EOF

# 3. Start Mako
./mako
```

**Model Options:**
- `deepseek-chat` - Main chat model
- `deepseek-coder` - Optimized for code

**Pricing:** https://platform.deepseek.com/pricing

---

## Advanced Configuration

### Using Different Providers for Embeddings

You can use one provider for command generation and another for semantic search:

```bash
# Navigate to CLI directory
cd apps/cli

# Edit .env file
# Main provider for commands
LLM_PROVIDER=anthropic
LLM_MODEL=claude-3-5-haiku-20241022
LLM_API_KEY=sk-ant-your-key

# Separate provider for embeddings (e.g., local Ollama)
EMBEDDING_PROVIDER=ollama
EMBEDDING_MODEL=nomic-embed-text
EMBEDDING_API_BASE=http://localhost:11434
```

This is useful for:
- Using free local embeddings with paid LLM
- Optimizing costs (embeddings are cheaper with some providers)
- Privacy (keep search history local)

### Docker + Ollama Configuration

If you're running Mako in Docker and Ollama on your host:

```bash
# On macOS/Windows
LLM_API_BASE=http://host.docker.internal:11434

# On Linux (use your host IP)
LLM_API_BASE=http://192.168.1.100:11434
# Or use --network=host when running Docker
```

---

## Verification

Test your configuration:

```bash
# Navigate to CLI directory
cd apps/cli

# Start Mako
./mako

# Try a simple command
mako ask "list files"

# Check configuration
mako config list
```

---

## Troubleshooting

### "API key not found"
- Check your `.env` file is in the Mako directory
- Verify the API key is correct (no extra spaces)
- Try setting the environment variable directly: `export LLM_API_KEY=your-key`

### "Ollama not reachable"
- Make sure Ollama is running: `ollama serve`
- Check if the model is installed: `ollama list`
- Verify the URL: `curl http://localhost:11434/api/tags`

### "API error (status 401)"
- Your API key is invalid or expired
- For OpenAI/Anthropic: Check you have credits in your account
- For Gemini: Ensure the API is enabled in your project

### "Model not found"
- Check the model name is correct (case-sensitive)
- For Ollama: Pull the model first with `ollama pull <model-name>`
- For cloud providers: Check their documentation for available models

---

## Cost Comparison

Approximate costs for generating 1000 commands:

| Provider | Model | Cost (USD) |
|----------|-------|------------|
| Ollama | Any | $0.00 (Free) |
| Gemini | gemini-2.5-flash | $0.01 (free tier) |
| OpenAI | gpt-4o-mini | $0.10 |
| Anthropic | claude-3-5-haiku | $0.15 |
| DeepSeek | deepseek-chat | $0.05 |
| OpenRouter | deepseek via OR | $0.05 |

*Estimates based on ~50 tokens per command generation*

---

## Switching Providers

You can easily switch between providers:

```bash
# Navigate to CLI directory
cd apps/cli

# Edit your .env file
nano .env

# Change the provider
LLM_PROVIDER=anthropic  # Change to desired provider
LLM_MODEL=claude-3-5-haiku-20241022
LLM_API_KEY=your-new-key

# Restart Mako
exit  # Exit current session
./mako  # Start new session
```

No data is lost when switching providers. Your command history and preferences are preserved.

---

## Getting Help

- **Documentation:** [README.md](../README.md)
- **Issues:** [GitHub Issues](https://github.com/fabiobrug/mako/issues)
- **Discussions:** [GitHub Discussions](https://github.com/fabiobrug/mako/discussions)
