"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const ExportCommand = () => (
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
      mako export
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Export command history to JSON format
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Export</h3>
          <CodeBlock>{`# Export all history to JSON
mako export > history.json

# Export last N commands
mako export --last 100 > recent.json
mako export --last 50 > last-50.json

# Export to specific file
mako export > ~/backups/mako-backup-$(date +%Y%m%d).json`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Filtering Exports</h3>
          <CodeBlock>{`# Export only successful commands
mako export --success > success.json

# Export only failed commands
mako export --failed > failed.json

# Export commands from specific date range
mako export --since "2024-01-01" > jan-2024.json
mako export --since "1 week ago" > last-week.json

# Combine filters
mako export --last 100 --success > recent-success.json`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Export Options</h3>
          <CodeBlock>{`# Pretty print JSON (default)
mako export > history.json

# Compact JSON (no formatting)
mako export --compact > history.json

# Include metadata
mako export --with-metadata > history.json

# Output:
# {
#   "version": "1.5.0",
#   "exported_at": "2024-02-12T10:30:00Z",
#   "total_commands": 150,
#   "commands": [...]
# }`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Backup & Restore Workflow</h3>
          <CodeBlock>{`# Create timestamped backup
mako export > ~/backups/mako-$(date +%Y%m%d-%H%M%S).json

# Automated daily backup (add to crontab)
0 2 * * * mako export > ~/backups/mako-$(date +%Y%m%d).json

# Restore from backup
mako import ~/backups/mako-20240212.json

# Verify export
mako export --last 10 | jq '.commands | length'
# Output: 10`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Empty History</h4>
              <CodeBlock>{`mako export > history.json
# Output: No commands to export
# Exit code: 1`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Permission Denied</h4>
              <CodeBlock>{`mako export > /root/history.json
# bash: /root/history.json: Permission denied
# Solution: Export to a directory you have write access to`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Large Export</h4>
              <CodeBlock>{`# Exporting 50,000 commands
mako export > large-history.json
# This may take a few seconds...
# ✓ Exported 50,000 commands (15.2 MB)

# Use filters to reduce size
mako export --last 1000 > recent.json`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Corrupt Export</h4>
              <CodeBlock>{`# Verify JSON is valid
mako export > history.json
jq . history.json > /dev/null && echo "Valid JSON" || echo "Invalid JSON"

# If corrupted, re-export
rm history.json
mako export > history.json`}</CodeBlock>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
