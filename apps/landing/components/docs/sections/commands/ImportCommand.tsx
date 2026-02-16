"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const ImportCommand = () => (
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
      mako import
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Import command history from JSON file
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Import</h3>
          <CodeBlock>{`# Import commands from JSON file
mako import history.json
mako import ~/backups/mako-backup-20240212.json

# Output:
# Importing commands from history.json...
# Found 150 commands
# Importing... [████████████████] 100%
# Successfully imported 150 commands`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Import Options</h3>
          <CodeBlock>{`# Merge with existing history (skip duplicates)
mako import history.json --merge

# Replace all existing history
mako import history.json --replace

# Preview what would be imported
mako import history.json --dry-run

# Import only successful commands
mako import history.json --success-only`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Import from Different Sources</h3>
          <CodeBlock>{`# Import from Mako export
mako import mako-export.json

# Import from bash history
mako import ~/.bash_history --format bash

# Import from zsh history
mako import ~/.zsh_history --format zsh

# Import from fish history
mako import ~/.local/share/fish/fish_history --format fish`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">File Not Found</h4>
              <CodeBlock>{`mako import nonexistent.json`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Error: File not found: nonexistent.json
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Invalid JSON Format</h4>
              <CodeBlock>{`mako import invalid.json`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Error: Invalid JSON format. Expected array of command objects with 'command' and 'timestamp' fields.
              </p>
            </div>

            <div className="border-l-2 border-yellow-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Duplicate Commands</h4>
              <CodeBlock>{`mako import history.json`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Found 100 commands, skipped 25 duplicates. Successfully imported 75 new commands. Use --replace to import all including duplicates.
              </p>
            </div>

            <div className="border-l-2 border-yellow-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Large Import</h4>
              <CodeBlock>{`mako import large-history.json`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Warning: This file contains 10,000 commands. This may take a few minutes. Use <code className="text-primary text-xs">--yes</code> to skip confirmation.
              </p>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
