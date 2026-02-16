"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../CodeBlock";
import { fadeInUp } from "../animations";

export const ReferenceSection = () => (
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
      Reference
    </motion.h1>
    
    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Architecture</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Mako uses a PTY (pseudo-terminal) architecture to intercept and process terminal I/O. Commands flow through the PTY master to your shell (bash/zsh), with a stream interceptor detecting markers for AI routing.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">File Locations</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">~/.mako/</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Main configuration directory for Mako
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">~/.mako/config.json</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Configuration file with API keys and settings
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">~/.mako/mako.db</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          SQLite database containing command history and embeddings
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">/usr/local/bin/mako</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Main Mako binary (shell orchestrator)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">/usr/local/bin/mako-menu</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Interactive menu binary for user choices
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Tech Stack</h2>
    
    <div className="overflow-x-auto">
      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-b border-border">
            <th className="text-left py-3 px-4 font-mono text-foreground">Component</th>
            <th className="text-left py-3 px-4 font-mono text-foreground">Technology</th>
          </tr>
        </thead>
        <tbody className="text-muted-foreground">
          <tr className="border-b border-border">
            <td className="py-3 px-4">Language</td>
            <td className="py-3 px-4 font-mono">Go 1.24+</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">PTY Handling</td>
            <td className="py-3 px-4 font-mono">creack/pty</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">Database</td>
            <td className="py-3 px-4 font-mono">SQLite with FTS5 (modernc.org/sqlite)</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">AI Providers</td>
            <td className="py-3 px-4 font-mono">OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter, Ollama</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">Embeddings</td>
            <td className="py-3 px-4 font-mono">Provider-specific (Gemini, OpenAI, Ollama)</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">Build System</td>
            <td className="py-3 px-4 font-mono">Make</td>
          </tr>
        </tbody>
      </table>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Environment Variables</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">LLM_PROVIDER</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          AI provider to use (openai, anthropic, gemini, deepseek, openrouter, ollama)
        </p>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">LLM_MODEL</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Model name for your chosen provider
        </p>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">LLM_API_KEY</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Your API key for the chosen provider (not needed for Ollama)
        </p>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">LLM_API_BASE</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Base URL for API calls (optional, for custom endpoints or Ollama)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">MAKO_CONFIG_DIR</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Override default config directory (default: ~/.mako)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">MAKO_DEBUG</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Enable debug logging (set to "1" or "true")
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Development</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      For developers interested in contributing or building from source:
    </p>
    <CodeBlock>{`# Clone repository
git clone https://github.com/fabiobrug/mako.git
cd mako/apps/cli

# Build
make build

# Run tests
make test

# Install locally
make install`}</CodeBlock>
  </motion.article>
);
