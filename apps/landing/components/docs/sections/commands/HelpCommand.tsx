"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const HelpCommand = () => (
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
      mako help
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Display help information and command documentation
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">General Help</h3>
          <CodeBlock>{`# Show all available commands
mako help

# Show version information
mako --version
mako -v

# Show quick start guide
mako help quickstart`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Command-Specific Help</h3>
          <CodeBlock>{`# Get help for specific command
mako help ask
mako help history
mako help stats
mako help export
mako help alias
mako help health

# Alternative syntax
mako ask --help
mako history --help`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Examples & Tutorials</h3>
          <CodeBlock>{`# Show common examples
mako help examples

# Show interactive tutorial
mako tutorial

# Show best practices
mako help best-practices`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Documentation Links</h3>
          <CodeBlock>{`# Open online documentation
mako help --docs

# Open GitHub repository
mako help --github

# Report an issue
mako help --report-bug

# Request a feature
mako help --request-feature`}</CodeBlock>
        </div>
      </div>
    </motion.article>
  );
