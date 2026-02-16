"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const SyncCommand = () => (
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
      mako sync
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Sync bash/zsh history to Mako database
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Sync</h3>
          <CodeBlock>{`# Sync shell history to Mako
mako sync

# Output:
# Syncing history from ~/.bash_history...
# Found 523 commands
# New commands: 45
# Syncing... [████████████████] 100%
# ✓ Successfully synced 45 new commands

# Mako automatically detects your shell (bash/zsh)`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Sync Options</h3>
          <CodeBlock>{`# Sync from specific history file
mako sync --file ~/.bash_history
mako sync --file ~/.zsh_history

# Sync only recent commands (last N)
mako sync --last 100

# Sync commands from specific time range
mako sync --since "2024-01-01"
mako sync --since "1 week ago"

# Dry run (preview what would be synced)
mako sync --dry-run`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Automatic Sync</h3>
          <p className="text-muted-foreground text-sm leading-relaxed mb-3">
            Configure Mako to automatically sync your shell history in the background.
          </p>
          <CodeBlock>{`# Enable automatic sync (every 5 minutes)
mako config set auto_sync true
mako config set auto_sync_interval 300

# Disable automatic sync
mako config set auto_sync false

# View sync status
mako sync status

# Output:
# Auto-sync: enabled
# Last sync: 2 minutes ago
# Commands synced: 1,234
# Next sync: in 3 minutes`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Sync Strategy</h3>
          <CodeBlock>{`# Bidirectional sync (Mako → Shell)
# Adds Mako-generated commands to shell history
mako sync --bidirectional

# One-way sync from shell only (default)
mako sync --one-way

# Skip duplicates
mako sync --skip-duplicates

# Include command timestamps
mako sync --preserve-timestamps`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">History File Not Found</h4>
              <CodeBlock>{`mako sync
# Error: Shell history file not found
# Expected: ~/.bash_history or ~/.zsh_history
# Solution: Ensure your shell writes history to disk`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">No New Commands</h4>
              <CodeBlock>{`mako sync
# No new commands to sync
# All commands are already in Mako database`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Permission Issues</h4>
              <CodeBlock>{`mako sync
# Error: Permission denied reading ~/.bash_history
# Solution: Check file permissions
chmod 600 ~/.bash_history`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Large History File</h4>
              <CodeBlock>{`mako sync
# Found 50,000 commands in history
# This may take a few minutes. Continue? [y/N]

# Process in background
mako sync --background`}</CodeBlock>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
