"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const VersionCommand = () => (
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
      mako version
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Show Mako version information
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Version</h3>
          <CodeBlock>{`# Show version
mako version

# Output:
# Mako v1.5.0

# Short form
mako --version
mako -v

# Output:
# v1.5.0`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Detailed Version Information</h3>
          <CodeBlock>{`# Show detailed version info
mako version --verbose

# Output:
# Mako v1.5.0
# ─────────────────────────────────────
# Build:        stable
# Commit:       a1b2c3d
# Build Date:   2024-02-12
# Go Version:   go1.24.1
# Platform:     linux/amd64
# 
# Configuration:
# Config Dir:   ~/.mako
# Database:     ~/.mako/mako.db
# LLM Provider: openai (gpt-4o-mini)`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Version Check</h3>
          <CodeBlock>{`# Check if you're running the latest version
mako version --check

# Output (if up to date):
# Current: v1.5.0
# Latest:  v1.5.0
# You're running the latest version

# Output (if update available):
# Current: v1.4.2
# Latest:  v1.5.0
# ⚠ Update available! Run 'mako update install'`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Version Output Formats</h3>
          <CodeBlock>{`# JSON format (useful for scripts)
mako version --json

# Output:
{
  "version": "1.5.0",
  "commit": "a1b2c3d",
  "buildDate": "2024-02-12",
  "goVersion": "go1.24.1",
  "platform": "linux/amd64"
}

# Plain text (parseable)
mako version --short

# Output:
# 1.5.0

# Use in shell scripts
VERSION=$(mako version --short)
echo "Running Mako $VERSION"`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Component Versions</h3>
          <CodeBlock>{`# Show versions of all components
mako version --all

# Output:
# Mako Components
# ─────────────────────────────────────
# mako:         v1.5.0
# mako-menu:    v1.5.0
# database:     v5 (schema)
# config:       v2 (format)
# 
# Dependencies:
# SQLite:       3.45.0
# Go Runtime:   go1.24.1`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Version Compatibility</h3>
          <p className="text-muted-foreground text-sm leading-relaxed mb-3">
            Mako follows semantic versioning (MAJOR.MINOR.PATCH):
          </p>
          <div className="border-l-2 border-primary/30 pl-4 py-3 mb-4">
            <ul className="space-y-2 text-muted-foreground text-sm">
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">•</span>
                <span><strong className="text-foreground">MAJOR (1.x.x):</strong> Breaking changes, may require reconfiguration</span>
              </li>
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">•</span>
                <span><strong className="text-foreground">MINOR (x.5.x):</strong> New features, backward compatible</span>
              </li>
              <li className="flex items-start gap-2">
                <span className="text-primary mt-0.5">•</span>
                <span><strong className="text-foreground">PATCH (x.x.0):</strong> Bug fixes, fully compatible</span>
              </li>
            </ul>
          </div>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="border-l-2 border-blue-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Development Build</h4>
              <CodeBlock>{`mako version`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Output: Mako v1.6.0-dev (Build: development). This is a development build, not for production use.
              </p>
            </div>

            <div className="border-l-2 border-yellow-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Version Mismatch</h4>
              <CodeBlock>{`mako version --all`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Warning: Component version mismatch (mako: v1.5.0, mako-menu: v1.4.2). Recommendation: Run <code className="text-primary text-xs">mako update install</code> to sync versions.
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Cannot Check for Updates</h4>
              <CodeBlock>{`mako version --check`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Error: Cannot connect to update server (Current version: v1.5.0). Unable to check for updates (offline).
              </p>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
