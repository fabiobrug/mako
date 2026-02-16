"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../CodeBlock";
import { fadeInUp } from "../animations";

export const AgentConfigSection = () => (
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
      Agent Config
    </motion.h1>
    
    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Supported AI Providers</h2>
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Mako supports multiple AI providers. Configure your preferred provider using environment variables or CLI commands:
    </p>

    <div className="overflow-x-auto mb-8">
      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-b border-border">
            <th className="text-left py-3 px-4 font-mono text-foreground">Provider</th>
            <th className="text-left py-3 px-4 font-mono text-foreground">Type</th>
            <th className="text-left py-3 px-4 font-mono text-foreground">Cost</th>
            <th className="text-left py-3 px-4 font-mono text-foreground">Best For</th>
          </tr>
        </thead>
        <tbody className="text-muted-foreground">
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">Ollama</td>
            <td className="py-3 px-4">Local</td>
            <td className="py-3 px-4 text-success">Free</td>
            <td className="py-3 px-4">Privacy, offline use, no API costs</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">OpenAI</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4">Paid</td>
            <td className="py-3 px-4">Best quality, GPT-4o models</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">Anthropic</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4">Paid</td>
            <td className="py-3 px-4">Claude models, great reasoning</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">Gemini</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4 text-success">Free tier</td>
            <td className="py-3 px-4">Default option, good balance</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">OpenRouter</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4">Paid</td>
            <td className="py-3 px-4">Access to multiple models</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">DeepSeek</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4">Paid</td>
            <td className="py-3 px-4">Cost-effective alternative</td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Example Configurations</h2>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Ollama (Local, Free)</h3>
    <CodeBlock>{`LLM_PROVIDER=ollama
LLM_MODEL=llama3.2
LLM_API_BASE=http://localhost:11434`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">OpenAI</h3>
    <CodeBlock>{`LLM_PROVIDER=openai
LLM_MODEL=gpt-4o-mini
LLM_API_KEY=sk-your-api-key-here`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Anthropic (Claude)</h3>
    <CodeBlock>{`LLM_PROVIDER=anthropic
LLM_MODEL=claude-3-5-haiku-20241022
LLM_API_KEY=sk-ant-your-key`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Google Gemini</h3>
    <CodeBlock>{`LLM_PROVIDER=gemini
LLM_MODEL=gemini-2.5-flash
LLM_API_KEY=your-gemini-key`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Using Ollama (Local Models)</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Ollama allows you to run AI models locally on your machine:
    </p>
    <CodeBlock>{`# Install Ollama
curl https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama3.2

# Configure Mako
mako config set llm_provider ollama
mako config set llm_model llama3.2
mako config set llm_base_url http://localhost:11434

# Start Mako
mako`}</CodeBlock>

    <p className="text-muted-foreground mb-4 leading-relaxed">
      <strong className="text-foreground">Benefits of Ollama:</strong>
    </p>
    <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li>Completely free</li>
      <li>Works offline</li>
      <li>Privacy - data never leaves your machine</li>
      <li>No API rate limits</li>
    </ul>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configuration Methods</h2>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Via CLI Commands</h3>
    <CodeBlock>{`# Set provider and credentials
mako config set llm_provider openai
mako config set llm_model gpt-4o-mini
mako config set api_key sk-your-api-key

# View current configuration
mako config list

# Check provider status
mako health`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Via Configuration File</h3>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Mako stores its configuration in <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">~/.mako/config.json</code>:
    </p>
    <CodeBlock language="json">{`{
  "llm_provider": "openai",
  "llm_model": "gpt-4o-mini",
  "api_key": "your-api-key-here",
  "llm_base_url": "https://api.openai.com/v1",
  "max_history": 1000,
  "auto_suggest": true,
  "safety_checks": true
}`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Getting API Keys</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">OpenAI</h4>
        <a href="https://platform.openai.com/api-keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          platform.openai.com/api-keys
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Anthropic</h4>
        <a href="https://console.anthropic.com/account/keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          console.anthropic.com/account/keys
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Google Gemini</h4>
        <a href="https://ai.google.dev/" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          ai.google.dev
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">OpenRouter</h4>
        <a href="https://openrouter.ai/keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          openrouter.ai/keys
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">DeepSeek</h4>
        <a href="https://platform.deepseek.com/api_keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          platform.deepseek.com/api_keys
        </a>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configuration Options</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">max_history</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Maximum number of commands to store in history (default: 1000)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">auto_suggest</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Enable automatic command suggestions (default: true)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">safety_checks</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Enable safety validation for dangerous commands (default: true)
        </p>
      </div>
    </div>
  </motion.article>
);
