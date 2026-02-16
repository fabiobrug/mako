"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const ConfigCommand = () => (
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
      mako config
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Manage Mako configuration settings
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">List All Settings</h3>
          <CodeBlock>{`# Show all configuration settings
mako config list

# Output:
# Configuration Settings
# ─────────────────────────────────────
# llm_provider:    openai
# llm_model:       gpt-4o-mini
# llm_base_url:    https://api.openai.com/v1
# api_key:         sk-...abc (hidden)
# max_history:     1000
# auto_suggest:    true
# safety_checks:   true`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Get Configuration Value</h3>
          <CodeBlock>{`# Get a specific setting
mako config get llm_provider
# Output: openai

mako config get llm_model
# Output: gpt-4o-mini

mako config get api_key
# Output: sk-...abc (partially hidden for security)`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Set Configuration Value</h3>
          <CodeBlock>{`# Set AI provider
mako config set llm_provider openai
mako config set llm_model gpt-4o-mini
mako config set api_key sk-your-api-key

# Set base URL for custom endpoints
mako config set llm_base_url https://api.openai.com/v1

# Configure Ollama (local)
mako config set llm_provider ollama
mako config set llm_model llama3.2
mako config set llm_base_url http://localhost:11434

# Set history and behavior options
mako config set max_history 2000
mako config set auto_suggest false
mako config set safety_checks true`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Reset Configuration</h3>
          <CodeBlock>{`# Reset all settings to defaults
mako config reset

# Confirmation prompt:
# Warning: This will reset all configuration to defaults
# Continue? [y/N]

# Force reset without confirmation
mako config reset --force`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Available Configuration Keys</h3>
          <div className="space-y-2">
            <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
              <code className="font-mono text-sm text-primary shrink-0 w-40">llm_provider</code>
              <p className="text-muted-foreground text-sm leading-relaxed">
                AI provider (openai, anthropic, gemini, deepseek, openrouter, ollama)
              </p>
            </div>

            <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
              <code className="font-mono text-sm text-primary shrink-0 w-40">llm_model</code>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Model name for your provider (e.g., gpt-4o-mini, claude-3-5-haiku-20241022)
              </p>
            </div>

            <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
              <code className="font-mono text-sm text-primary shrink-0 w-40">llm_base_url</code>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Base URL for API calls (optional, for custom endpoints or Ollama)
              </p>
            </div>

            <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
              <code className="font-mono text-sm text-primary shrink-0 w-40">api_key</code>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Your API key (not required for Ollama)
              </p>
            </div>

            <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
              <code className="font-mono text-sm text-primary shrink-0 w-40">max_history</code>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Maximum number of commands to store (default: 1000)
              </p>
            </div>

            <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
              <code className="font-mono text-sm text-primary shrink-0 w-40">auto_suggest</code>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Enable automatic command suggestions (default: true)
              </p>
            </div>

            <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
              <code className="font-mono text-sm text-primary shrink-0 w-40">safety_checks</code>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Enable safety validation for dangerous commands (default: true)
              </p>
            </div>
          </div>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Invalid Configuration Key</h4>
              <CodeBlock>{`mako config set invalid_key value`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Error: Unknown configuration key 'invalid_key'. Use <code className="text-primary text-xs">mako config list</code> to see available keys.
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Invalid Value Type</h4>
              <CodeBlock>{`mako config set max_history "not a number"`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Error: Invalid value for 'max_history': expected integer
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Configuration File Issues</h4>
              <CodeBlock>{`Error: Failed to read configuration file`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Solution: Reset configuration with <code className="text-primary text-xs">mako config reset</code>
              </p>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
