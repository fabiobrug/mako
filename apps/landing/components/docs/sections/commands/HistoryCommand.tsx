"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const HistoryCommand = () => (
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
      mako history
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      View, search, and manage your command history
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Usage</h3>
          <CodeBlock>{`# Show recent command history
mako history

# Search by keyword
mako history docker
mako history "git push"

# Show history with timestamps
mako history --timestamps`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Keyword Search</h3>
          <p className="text-muted-foreground text-sm leading-relaxed mb-3">
            Quick search for commands containing specific keywords (case-insensitive).
          </p>
          <CodeBlock>{`# Search by single keyword
mako history docker

# Search by phrase
mako history "npm install"

# Search with multiple terms
mako history git push origin`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Semantic Search</h3>
          <p className="text-muted-foreground text-sm leading-relaxed mb-3">
            Find commands by meaning, not just exact text matches. Uses AI embeddings for intelligent search.
          </p>
          <CodeBlock>{`# Find database-related commands
mako history semantic "database backup"
# Finds: pg_dump, mysqldump, mongodump, etc.

# Find deployment commands
mako history semantic "deploy to production"
# Finds: git push, kubectl apply, docker push, etc.

# Find file compression commands
mako history semantic "compress files"
# Finds: tar, gzip, zip, 7z, etc.`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Filtering by Status</h3>
          <CodeBlock>{`# Show only failed commands
mako history --failed

# Show only successful commands
mako history --success

# Combine with search
mako history docker --failed
mako history semantic "database backup" --success`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Interactive Mode</h3>
          <p className="text-muted-foreground text-sm leading-relaxed mb-3">
            Browse your command history interactively with arrow keys and search.
          </p>
          <CodeBlock>{`# Launch interactive history browser
mako history --interactive

# Controls:
# ↑/↓ - Navigate through commands
# Enter - Copy selected command to clipboard
# / - Search within results
# Esc - Exit interactive mode`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Additional Options</h3>
          <CodeBlock>{`# Filter by date range
mako history --since "2024-01-01"
mako history --since "yesterday"
mako history --since "1 week ago"

# Filter by directory
mako history --dir "/home/user/projects"

# Limit results
mako history --limit 20

# Combine filters
mako history docker --success --since "1 week ago"`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Output Formats</h3>
          <CodeBlock>{`# JSON output
mako history --format json

# CSV output
mako history --format csv

# Detailed output with metadata
mako history --verbose

# Copy command to clipboard
mako history --limit 1 --copy`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">History Management</h3>
          <CodeBlock>{`# Clear all history
mako history --clear

# Clear last N entries
mako history --clear --limit 10

# Remove specific command by ID
mako history --remove 1234

# Clear failed commands
mako history --clear --failed`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="border-l-2 border-blue-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Empty History</h4>
              <CodeBlock>{`mako history`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Output: No command history found. Run some commands to start building your history.
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Corrupted Database</h4>
              <CodeBlock>{`Error: Database error: unable to read history`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Solution: Run <code className="text-primary text-xs">mako health</code> to diagnose and repair
              </p>
            </div>

            <div className="border-l-2 border-yellow-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Semantic Search Without Internet</h4>
              <CodeBlock>{`mako history semantic "backup"`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Falls back to text search when offline. Using text search instead (AI features unavailable)
              </p>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
