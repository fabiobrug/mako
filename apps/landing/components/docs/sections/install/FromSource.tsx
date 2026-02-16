"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const FromSourceSection = () => (
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
      Build From Source
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Build Mako from source for development, customization, or to use the latest features. Requires Go 1.25 or higher.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Prerequisites</h2>
    <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li>Go 1.25+ (<a href="https://go.dev/dl/" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 underline">download here</a>)</li>
      <li>Git</li>
      <li>Make (usually pre-installed on Linux/macOS)</li>
    </ul>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Clone Repository</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Clone the Mako repository from GitHub:
    </p>
    <CodeBlock>{`git clone https://github.com/fabiobrug/mako.git`}</CodeBlock>
    
    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Navigate to the CLI directory:
    </p>
    <CodeBlock>{`cd mako/apps/cli`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configure Environment</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Copy the example environment file:
    </p>
    <CodeBlock>{`cp .env.example .env`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Edit the <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">.env</code> file with your provider settings:
    </p>
    <CodeBlock>{`nano .env`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      <strong className="text-foreground">Example .env configuration:</strong>
    </p>
    <CodeBlock>{`LLM_PROVIDER=openai
LLM_MODEL=gpt-4o-mini
LLM_API_KEY=sk-your-key`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Build</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Build Mako using the Makefile:
    </p>
    <CodeBlock>{`make build`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      This creates the <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">mako</code> and <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">mako-menu</code> binaries in the current directory.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Install (Optional)</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Install to <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">/usr/local/bin</code> (requires sudo):
    </p>
    <CodeBlock>{`make install`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      Or run directly from the build directory:
    </p>
    <CodeBlock>{`./mako`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Manual Build</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Build without Make:
    </p>
    <CodeBlock>{`go build -o mako ./cmd/mako`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Build the menu binary:
    </p>
    <CodeBlock>{`go build -o mako-menu ./cmd/mako-menu`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Development Mode</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      For development, use the .env file for quick configuration changes without rebuilding. Environment variables take precedence over config file settings.
    </p>
  </motion.article>
);
