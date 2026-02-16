"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const OneLineInstallSection = () => (
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
      Quick Install
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      The fastest way to get started with Mako. Install with a single command and optionally configure your AI provider inline.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Basic Installation</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Download and install Mako with default configuration (Gemini):
    </p>
    <CodeBlock>{`curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash`}</CodeBlock>
    
    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Install with Provider Configuration</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Configure your preferred AI provider during installation:
    </p>
    <CodeBlock>{`LLM_PROVIDER=openai LLM_MODEL=gpt-4o-mini LLM_API_KEY=sk-your-key curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">What Gets Installed</h2>
    <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li>Mako binary installed to <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">/usr/local/bin/mako</code></li>
      <li>Shell integration functions for bash/zsh</li>
      <li>Configuration directory at <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">~/.mako/</code></li>
      <li>Automatic shell profile update (with backup)</li>
    </ul>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Post-Installation</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Restart your shell or source your profile:
    </p>
    <CodeBlock>{`source ~/.bashrc`}</CodeBlock>
    
    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Or for zsh:
    </p>
    <CodeBlock>{`source ~/.zshrc`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      Start Mako and test your installation:
    </p>
    <CodeBlock>{`mako`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Inside Mako, verify everything is working:
    </p>
    <CodeBlock>{`mako health`}</CodeBlock>
  </motion.article>
);
