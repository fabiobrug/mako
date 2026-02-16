"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const UpdateCommand = () => (
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
      mako update
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Check for and install Mako updates
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Check for Updates</h3>
          <CodeBlock>{`# Check if updates are available
mako update check

# Output (if update available):
# Update Available!
# Current version: v1.4.2
# Latest version:  v1.5.0
# 
# Release notes:
# - Added interactive history mode
# - Improved alias management with tags
# - New config command for easier setup
# - Performance improvements
# 
# Run 'mako update install' to update

# Output (if up to date):
# You're running the latest version: v1.4.2`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Install Updates</h3>
          <CodeBlock>{`# Install the latest version
mako update install

# Process:
# 1. Downloading latest version...
# 2. Backing up current installation...
# 3. Installing v1.5.0...
# 4. Verifying installation...
# 
# ✓ Successfully updated to v1.5.0
# Run 'mako --version' to verify

# Skip confirmation prompt
mako update install --yes`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Update Options</h3>
          <CodeBlock>{`# Check for pre-release versions
mako update check --pre-release
mako update check --beta

# Install specific version
mako update install v1.5.0
mako update install v1.5.0-beta.1

# Show what would be updated (dry run)
mako update install --dry-run`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Automatic Update Checks</h3>
          <p className="text-muted-foreground text-sm leading-relaxed mb-3">
            Mako automatically checks for updates once per day. You'll see a notification if an update is available.
          </p>
          <CodeBlock>{`# Disable automatic update checks
mako config set auto_update_check false

# Re-enable automatic checks
mako config set auto_update_check true

# View last update check time
mako config get last_update_check`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">No Internet Connection</h4>
              <CodeBlock>{`mako update check
# Error: Cannot connect to update server
# Check your internet connection and try again`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Permission Denied</h4>
              <CodeBlock>{`mako update install
# Error: Permission denied
# Try: sudo mako update install`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Update Failed</h4>
              <CodeBlock>{`# If update fails, restore from backup
Error: Update failed at step 3
Restoring from backup...
✓ Restored to version v1.4.2

# Manual update via install script
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Downgrade Not Allowed</h4>
              <CodeBlock>{`mako update install v1.3.0
# Error: Downgrading is not supported
# Current version: v1.4.2
# Use manual installation to downgrade`}</CodeBlock>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
