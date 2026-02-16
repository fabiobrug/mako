"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const PostInstallSection = () => (
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
      Post-Installation Setup
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      After installing Mako, configure your AI provider and verify everything is working correctly.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Start Mako</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Launch Mako shell:
    </p>
    <CodeBlock>{`mako`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configure Provider</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Set your AI provider (OpenAI example):
    </p>
    <CodeBlock>{`mako config set llm_provider openai`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Set your model:
    </p>
    <CodeBlock>{`mako config set llm_model gpt-4o-mini`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Set your API key:
    </p>
    <CodeBlock>{`mako config set api_key sk-your-api-key`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">For Ollama (Local)</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Configure for local Ollama:
    </p>
    <CodeBlock>{`mako config set llm_provider ollama`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Set the model:
    </p>
    <CodeBlock>{`mako config set llm_model llama3.2`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Set the base URL:
    </p>
    <CodeBlock>{`mako config set llm_base_url http://localhost:11434`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Verify Configuration</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      View your current settings:
    </p>
    <CodeBlock>{`mako config list`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      Check provider health:
    </p>
    <CodeBlock>{`mako health`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Test Your Setup</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Try your first AI command:
    </p>
    <CodeBlock>{`mako ask "list files in current directory"`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      View command history:
    </p>
    <CodeBlock>{`mako history`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      Get help:
    </p>
    <CodeBlock>{`mako help`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-6">Configuration Priority</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Mako checks for configuration in this order (highest to lowest priority):
    </p>
    <div className="border-l-2 border-primary/30 pl-6 py-2">
      <ol className="space-y-3 text-muted-foreground leading-relaxed">
        <li className="flex items-start gap-3">
          <span className="text-primary font-bold shrink-0">1.</span>
          <div>
            <strong className="text-foreground">Environment variables</strong> - from <code className="font-mono text-primary text-xs">.env</code> file in <code className="font-mono text-primary text-xs">apps/cli/</code>
          </div>
        </li>
        <li className="flex items-start gap-3">
          <span className="text-primary font-bold shrink-0">2.</span>
          <div>
            <strong className="text-foreground">Config file</strong> - <code className="font-mono text-primary text-xs">~/.mako/config.json</code> (set via <code className="font-mono text-primary text-xs">mako config set</code>)
          </div>
        </li>
        <li className="flex items-start gap-3">
          <span className="text-primary font-bold shrink-0">3.</span>
          <div>
            <strong className="text-foreground">Default values</strong> - Gemini provider with free tier
          </div>
        </li>
      </ol>
    </div>
  </motion.article>
);
