"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../CodeBlock";
import { fadeInUp } from "../animations";

export const InstallSection = () => (
  <motion.article 
    initial="hidden"
    animate="visible"
    variants={fadeInUp}
    className="prose prose-invert max-w-none"
  >
    <motion.h1 
      variants={fadeInUp}
      className="font-mono text-3xl font-bold text-foreground mb-6"
    >
      Installation
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Mako works on <strong className="text-foreground">Linux</strong> and <strong className="text-foreground">macOS</strong>. Choose from multiple AI providers including local models (Ollama) or cloud services (OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter).
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Option 1: One-Line Install (Recommended)</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Fast installation with optional environment variable configuration:
    </p>
    <CodeBlock>{`# Basic installation
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash

# Or install with provider configuration
LLM_PROVIDER=openai LLM_MODEL=gpt-4o-mini LLM_API_KEY=sk-your-key \\
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash`}</CodeBlock>
    
    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      <strong className="text-foreground">After installation, configure your AI provider:</strong>
    </p>
    <CodeBlock>{`# Start Mako
mako

# Inside Mako shell, configure your provider:
mako config set llm_provider openai
mako config set llm_model gpt-4o-mini
mako config set api_key sk-your-api-key

# Or for Ollama (local):
mako config set llm_provider ollama
mako config set llm_model llama3.2
mako config set llm_base_url http://localhost:11434

# View all settings:
mako config list`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      <strong className="text-foreground">Supported configuration keys:</strong>
    </p>
    <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li><code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">llm_provider</code> - AI provider (openai, anthropic, gemini, deepseek, openrouter, ollama)</li>
      <li><code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">llm_model</code> - Model name (provider-specific)</li>
      <li><code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">llm_base_url</code> - Base URL (optional, for custom endpoints)</li>
      <li><code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">api_key</code> - Your API key (not required for Ollama)</li>
    </ul>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Option 2: From Source with .env File</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Clone the repository and configure via <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">.env</code> file:
    </p>
    <CodeBlock>{`# Clone the repository
git clone https://github.com/fabiobrug/mako.git
cd mako/apps/cli

# Copy and edit configuration
cp .env.example .env
nano .env  # Edit with your provider settings

# Build
make build

# Install (optional, requires sudo)
make install

# Or run directly
./mako`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      <strong className="text-foreground">Example <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">.env</code> configuration:</strong>
    </p>
    <CodeBlock>{`# OpenAI
LLM_PROVIDER=openai
LLM_MODEL=gpt-4o-mini
LLM_API_KEY=sk-your-key

# Or Ollama (local, free)
LLM_PROVIDER=ollama
LLM_MODEL=llama3.2
LLM_API_BASE=http://localhost:11434`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Verify Installation</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      After installation, verify that Mako is working correctly:
    </p>
    <CodeBlock>{`# Start Mako shell
mako

# Inside Mako shell, try:
mako ask "find files larger than 100MB"
mako history
mako health    # Check configuration status
mako help`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configuration Priority</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Mako checks for configuration in this order:
    </p>
    <ol className="list-decimal list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li><strong className="text-foreground">Environment variables</strong> (<code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">.env</code> file in <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">apps/cli/</code>)</li>
      <li><strong className="text-foreground">Config file</strong> (<code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">~/.mako/config.json</code>) - set via <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">mako config set</code></li>
      <li><strong className="text-foreground">Default values</strong> (Gemini provider)</li>
    </ol>

    <p className="text-muted-foreground mb-4 leading-relaxed">
      For more details on configuring AI providers, see the <button onClick={() => window.dispatchEvent(new CustomEvent('changeSection', { detail: 'agent-config' }))} className="text-primary hover:text-primary/80 underline">Agent Config</button> section.
    </p>
  </motion.article>
);
