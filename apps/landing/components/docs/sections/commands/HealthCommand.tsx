"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const HealthCommand = () => (
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
      mako health
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Check system health, configuration, and diagnose issues
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Health Check</h3>
          <CodeBlock>{`mako health

# Output:
# Mako Health Check
# ─────────────────────────────────────
# [OK] Configuration loaded
# [OK] API Key configured
# [OK] API connection successful
# [OK] Database initialized
# [OK] Database integrity check passed
# [OK] Shell integration active
# [OK] PTY working correctly
# [OK] All binaries present
# 
# System Status: Healthy`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Detailed Diagnostics</h3>
          <CodeBlock>{`# Run detailed diagnostics
mako health --verbose

# Output includes:
# - Configuration file location and contents
# - API key status and validity
# - Database file size and entry count
# - Shell type and version
# - PTY capabilities
# - Network connectivity
# - Disk space available
# - Memory usage`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Specific Checks</h3>
          <CodeBlock>{`# Check only API connection
mako health --api

# Check only database
mako health --database

# Check only configuration
mako health --config

# Check for updates
mako health --check-updates`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Repair & Fix</h3>
          <CodeBlock>{`# Auto-fix common issues
mako health --fix

# Repair database
mako health --repair-database

# Reset configuration to defaults
mako health --reset-config

# Clear cache
mako health --clear-cache

# Reinstall shell integration
mako health --reinstall-shell`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Common Issues</h3>
          <div className="space-y-3">
            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">API Key Not Found</h4>
              <CodeBlock>{`[FAIL] API Key not configured`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Solution: Configure your API key with <code className="text-primary text-xs">mako config set api_key your-key-here</code> or set environment variable for your provider
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Database Corruption</h4>
              <CodeBlock>{`[FAIL] Database integrity check failed`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Solution: Repair database with <code className="text-primary text-xs">mako health --repair-database</code>. Or backup and recreate by removing <code className="text-primary text-xs">~/.mako/mako.db</code> (will recreate on next start)
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Shell Integration Inactive</h4>
              <CodeBlock>{`[FAIL] Shell integration not detected`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Solution: Reinstall shell hooks with <code className="text-primary text-xs">mako health --reinstall-shell</code>, then restart your shell with <code className="text-primary text-xs">exec $SHELL</code>
              </p>
            </div>

            <div className="border-l-2 border-yellow-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Low Disk Space</h4>
              <CodeBlock>{`[WARNING] Low disk space detected (< 1GB available)`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Database size: 850MB. Recommendation: Clean up old history with <code className="text-primary text-xs">mako history --clear --before "6 months ago"</code>
              </p>
            </div>

            <div className="border-l-2 border-blue-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Outdated Version</h4>
              <CodeBlock>{`[INFO] Update available: v1.5.0 (current: v1.4.2)`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Run: <code className="text-primary text-xs">curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash</code>
              </p>
            </div>
          </div>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Performance Metrics</h3>
          <CodeBlock>{`mako health --performance

# Output:
# Performance Metrics
# ─────────────────────────────────────
# Average command generation time: 0.3s
# Average semantic search time: 0.5s
# Database query time: 0.02s
# Memory usage: 45MB
# Disk usage: 120MB
# API response time: 0.25s`}</CodeBlock>
        </div>
      </div>
    </motion.article>
  );
