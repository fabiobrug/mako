"use client";

import { motion } from "framer-motion";
import { fadeInUp } from "../../animations";

export const AgentConfigOverviewSection = () => (
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
      AI Provider Overview
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Mako supports multiple AI providers, giving you flexibility to choose between local models, free cloud services, or premium AI providers based on your needs.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Supported Providers</h2>
    
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

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-6">Configuration Keys</h2>
    
    <div className="space-y-2">
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-40">llm_provider</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          AI provider name: <code className="text-foreground text-xs">openai</code>, <code className="text-foreground text-xs">anthropic</code>, <code className="text-foreground text-xs">gemini</code>, <code className="text-foreground text-xs">deepseek</code>, <code className="text-foreground text-xs">openrouter</code>, <code className="text-foreground text-xs">ollama</code>
        </p>
      </div>
      
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-40">llm_model</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Model name (provider-specific). Examples: <code className="text-foreground text-xs">gpt-4o-mini</code>, <code className="text-foreground text-xs">claude-3-5-haiku</code>, <code className="text-foreground text-xs">llama3.2</code>
        </p>
      </div>
      
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-40">api_key</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Your API key for the provider (not required for Ollama)
        </p>
      </div>
      
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-40">llm_base_url</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Optional base URL for custom endpoints or self-hosted services
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Choosing a Provider</h2>
    
    <div className="space-y-6">
      <div>
        <h3 className="font-mono text-lg font-semibold text-foreground mb-2">For Privacy & Offline Use</h3>
        <p className="text-muted-foreground leading-relaxed">
          Use <strong className="text-foreground">Ollama</strong> to run models locally. Your data never leaves your machine, works without internet, and is completely free.
        </p>
      </div>

      <div>
        <h3 className="font-mono text-lg font-semibold text-foreground mb-2">For Best Performance</h3>
        <p className="text-muted-foreground leading-relaxed">
          Use <strong className="text-foreground">OpenAI (GPT-4o)</strong> or <strong className="text-foreground">Anthropic (Claude)</strong> for the highest quality command generation and reasoning.
        </p>
      </div>

      <div>
        <h3 className="font-mono text-lg font-semibold text-foreground mb-2">For Free Cloud Service</h3>
        <p className="text-muted-foreground leading-relaxed">
          Use <strong className="text-foreground">Gemini</strong> (default) for a good balance of quality and cost with a generous free tier.
        </p>
      </div>

      <div>
        <h3 className="font-mono text-lg font-semibold text-foreground mb-2">For Model Flexibility</h3>
        <p className="text-muted-foreground leading-relaxed">
          Use <strong className="text-foreground">OpenRouter</strong> to access multiple models from different providers with a single API key.
        </p>
      </div>
    </div>
  </motion.article>
);
