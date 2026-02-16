"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const AdvancedConfigSection = () => (
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
      Advanced Configuration
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Fine-tune Mako's behavior with advanced configuration options for power users.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configuration File</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Mako stores configuration in <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">~/.mako/config.json</code>:
    </p>
    <CodeBlock language="json">{`{
  "llm_provider": "openai",
  "llm_model": "gpt-4o-mini",
  "api_key": "your-api-key",
  "llm_base_url": "https://api.openai.com/v1",
  "max_history": 1000,
  "auto_suggest": true,
  "safety_checks": true
}`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-6">Available Options</h2>
    
    <div className="space-y-2">
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-40">max_history</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Maximum number of commands to store in history. Default: <code className="text-foreground text-xs">1000</code>
        </p>
      </div>
      
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-40">auto_suggest</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Enable automatic command suggestions based on context. Default: <code className="text-foreground text-xs">true</code>
        </p>
      </div>
      
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-40">safety_checks</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Enable validation for potentially dangerous commands. Default: <code className="text-foreground text-xs">true</code>
        </p>
      </div>

      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-40">llm_base_url</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Custom API endpoint for self-hosted or proxy services. Optional.
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Custom Base URLs</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Use a custom API endpoint:
    </p>
    <CodeBlock>{`mako config set llm_base_url https://your-proxy.com/v1`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      Useful for:
    </p>
    <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li>Self-hosted AI services</li>
      <li>Corporate proxy endpoints</li>
      <li>Custom load balancers</li>
      <li>Development/testing environments</li>
    </ul>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">View Current Settings</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      List all configuration:
    </p>
    <CodeBlock>{`mako config list`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Reset Configuration</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Remove a specific setting:
    </p>
    <CodeBlock>{`mako config set auto_suggest false`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      Or manually edit the config file:
    </p>
    <CodeBlock>{`nano ~/.mako/config.json`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-6">Configuration Priority</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Settings are loaded in this order (highest to lowest priority):
    </p>
    <div className="border-l-2 border-primary/30 pl-6 py-2">
      <ol className="space-y-3 text-muted-foreground leading-relaxed">
        <li className="flex items-start gap-3">
          <span className="text-primary font-bold shrink-0">1.</span>
          <div>
            <strong className="text-foreground">Environment variables</strong> - from <code className="font-mono text-primary text-xs">.env</code>
          </div>
        </li>
        <li className="flex items-start gap-3">
          <span className="text-primary font-bold shrink-0">2.</span>
          <div>
            <strong className="text-foreground">Config file</strong> - from <code className="font-mono text-primary text-xs">~/.mako/config.json</code>
          </div>
        </li>
        <li className="flex items-start gap-3">
          <span className="text-primary font-bold shrink-0">3.</span>
          <div>
            <strong className="text-foreground">Defaults</strong> - built-in values
          </div>
        </li>
      </ol>
    </div>

    <div className="mt-6 border-l-2 border-primary/30 pl-4 py-2 bg-primary/5 rounded-r">
      <p className="text-muted-foreground text-sm leading-relaxed">
        This allows you to override config file settings with environment variables for testing without modifying the config.
      </p>
    </div>
  </motion.article>
);
