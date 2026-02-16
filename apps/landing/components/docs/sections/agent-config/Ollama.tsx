"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const OllamaConfigSection = () => (
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
      Ollama Configuration
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Ollama allows you to run AI models locally on your machine. Perfect for privacy, offline use, and avoiding API costs.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Install Ollama</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Download and install Ollama:
    </p>
    <CodeBlock>{`curl https://ollama.ai/install.sh | sh`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Pull a Model</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Download Llama 3.2 (recommended for Mako):
    </p>
    <CodeBlock>{`ollama pull llama3.2`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      Or try other models:
    </p>
    <CodeBlock>{`ollama pull codellama`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Or Mistral:
    </p>
    <CodeBlock>{`ollama pull mistral`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configure Mako</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Set Ollama as your provider:
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

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Environment Variables</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Or configure via .env file:
    </p>
    <CodeBlock>{`LLM_PROVIDER=ollama
LLM_MODEL=llama3.2
LLM_API_BASE=http://localhost:11434`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Verify Setup</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Check Ollama is running:
    </p>
    <CodeBlock>{`ollama list`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      Test in Mako:
    </p>
    <CodeBlock>{`mako health`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Benefits</h2>
    <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li><strong className="text-foreground">Completely free</strong> - no API costs</li>
      <li><strong className="text-foreground">Works offline</strong> - no internet required</li>
      <li><strong className="text-foreground">Privacy first</strong> - data never leaves your machine</li>
      <li><strong className="text-foreground">No rate limits</strong> - use as much as you want</li>
      <li><strong className="text-foreground">Fast responses</strong> - no network latency</li>
    </ul>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-6">Recommended Models</h2>
    <div className="space-y-2">
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-32">llama3.2</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Best balance of performance and resource usage for shell commands
        </p>
      </div>
      
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-32">codellama</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Specialized for code and technical tasks
        </p>
      </div>
      
      <div className="flex items-start gap-4 py-3 hover:bg-muted-foreground/5 transition-colors rounded px-2">
        <code className="font-mono text-sm text-primary shrink-0 w-32">mistral</code>
        <p className="text-muted-foreground text-sm leading-relaxed">
          High quality alternative, requires more resources
        </p>
      </div>
    </div>
  </motion.article>
);
