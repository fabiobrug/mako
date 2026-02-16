"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../CodeBlock";
import { fadeInUp } from "../animations";
import { Folder, FileCode, Database, Cpu, Layers } from "lucide-react";

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
    
    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-6">Architecture</h2>
    <div className="border-l-2 border-primary/30 pl-6 py-2 bg-primary/5 rounded-r">
      <p className="text-muted-foreground leading-relaxed">
        Mako uses a <strong className="text-foreground">PTY (pseudo-terminal)</strong> architecture to intercept and process terminal I/O. Commands flow through the PTY master to your shell (bash/zsh), with a stream interceptor detecting markers for AI routing.
      </p>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-12 mb-6">File Locations</h2>
    
    <div className="space-y-3">
      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-primary/30">
        <div className="flex items-start gap-3">
          <Folder size={18} className="text-primary mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <code className="font-mono text-sm text-foreground block mb-1.5">~/.mako/</code>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Main configuration directory
            </p>
          </div>
        </div>
      </div>
      
      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-primary/30">
        <div className="flex items-start gap-3">
          <FileCode size={18} className="text-primary mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <code className="font-mono text-sm text-foreground block mb-1.5">~/.mako/config.json</code>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Configuration file with API keys and settings
            </p>
          </div>
        </div>
      </div>
      
      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-primary/30">
        <div className="flex items-start gap-3">
          <Database size={18} className="text-primary mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <code className="font-mono text-sm text-foreground block mb-1.5">~/.mako/mako.db</code>
            <p className="text-muted-foreground text-sm leading-relaxed">
              SQLite database with command history and embeddings
            </p>
          </div>
        </div>
      </div>
      
      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-primary/30">
        <div className="flex items-start gap-3">
          <Cpu size={18} className="text-primary mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <code className="font-mono text-sm text-foreground block mb-1.5">/usr/local/bin/mako</code>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Main Mako binary (shell orchestrator)
            </p>
          </div>
        </div>
      </div>
      
      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-primary/30">
        <div className="flex items-start gap-3">
          <Layers size={18} className="text-primary mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <code className="font-mono text-sm text-foreground block mb-1.5">/usr/local/bin/mako-menu</code>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Interactive menu binary for user choices
            </p>
          </div>
        </div>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-12 mb-6">Tech Stack</h2>
    
    <div className="border border-border/50 rounded-lg overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-border/50 bg-muted-foreground/5">
            <th className="text-left py-3 px-4 font-mono text-sm font-semibold text-foreground">Component</th>
            <th className="text-left py-3 px-4 font-mono text-sm font-semibold text-foreground">Technology</th>
          </tr>
        </thead>
        <tbody className="text-muted-foreground">
          <tr className="border-b border-border/30 hover:bg-muted-foreground/5 transition-colors">
            <td className="py-3 px-4 text-sm">Language</td>
            <td className="py-3 px-4 font-mono text-sm">Go 1.25+</td>
          </tr>
          <tr className="border-b border-border/30 hover:bg-muted-foreground/5 transition-colors">
            <td className="py-3 px-4 text-sm">PTY Handling</td>
            <td className="py-3 px-4 font-mono text-sm">creack/pty</td>
          </tr>
          <tr className="border-b border-border/30 hover:bg-muted-foreground/5 transition-colors">
            <td className="py-3 px-4 text-sm">Database</td>
            <td className="py-3 px-4 font-mono text-sm">SQLite + FTS5</td>
          </tr>
          <tr className="border-b border-border/30 hover:bg-muted-foreground/5 transition-colors">
            <td className="py-3 px-4 text-sm">AI Providers</td>
            <td className="py-3 px-4 font-mono text-sm">OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter, Ollama</td>
          </tr>
          <tr className="border-b border-border/30 hover:bg-muted-foreground/5 transition-colors">
            <td className="py-3 px-4 text-sm">Embeddings</td>
            <td className="py-3 px-4 font-mono text-sm">Provider-specific</td>
          </tr>
          <tr className="hover:bg-muted-foreground/5 transition-colors">
            <td className="py-3 px-4 text-sm">Build System</td>
            <td className="py-3 px-4 font-mono text-sm">Make</td>
          </tr>
        </tbody>
      </table>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-12 mb-6">Environment Variables</h2>
    
    <div className="space-y-2">
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-48">LLM_PROVIDER</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          AI provider: openai, anthropic, gemini, deepseek, openrouter, ollama
        </p>
      </div>

      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-48">LLM_MODEL</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Model name for your provider
        </p>
      </div>

      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-48">LLM_API_KEY</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          API key (not needed for Ollama)
        </p>
      </div>

      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-48">LLM_API_BASE</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Base URL for custom endpoints
        </p>
      </div>
      
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-48">MAKO_CONFIG_DIR</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Override config directory (default: ~/.mako)
        </p>
      </div>
      
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-48">MAKO_DEBUG</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Enable debug logging ("1" or "true")
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-12 mb-6">Development</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      For developers interested in contributing or building from source:
    </p>
    <CodeBlock>{`git clone https://github.com/fabiobrug/mako.git`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Navigate to CLI directory:
    </p>
    <CodeBlock>{`cd mako/apps/cli`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Build the project:
    </p>
    <CodeBlock>{`make build`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Run tests:
    </p>
    <CodeBlock>{`make test`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Install locally:
    </p>
    <CodeBlock>{`make install`}</CodeBlock>
  </motion.article>
);
