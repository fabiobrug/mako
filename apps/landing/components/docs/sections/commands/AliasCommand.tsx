"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const AliasCommand = () => (
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
      mako alias
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Create, manage, and use custom command aliases
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Saving Aliases</h3>
          <CodeBlock>{`# Save a simple alias
mako alias save ll "ls -lah"

# Complex alias with pipes
mako alias save psgrep "ps aux | grep"

# Multi-command alias
mako alias save deploy "git push && ssh prod 'cd /app && ./deploy.sh'"

# Alias with parameters (use $1, $2, etc.)
mako alias save findfile "find . -name '$1' -type f"

# Alias with tags for organization
mako alias save deploy "git push" --tag git
mako alias save dps "docker ps -a" --tag docker`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Listing Aliases</h3>
          <CodeBlock>{`# List all aliases
mako alias list

# List aliases by tag
mako alias list --tag git
mako alias list --tag docker

# Show detailed information
mako alias list --verbose`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Running Aliases</h3>
          <CodeBlock>{`# Run a saved alias
mako alias run ll

# Run alias with parameters
mako alias run findfile "*.js"
mako alias run backup /path/to/dir

# Combine with shell operations
mako alias run ll | grep Documents`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Managing Aliases</h3>
          <CodeBlock>{`# Delete an alias
mako alias delete ll
mako alias delete deploy

# Edit an alias (re-save with same name)
mako alias save ll "ls -lah --color=auto"`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Import & Export</h3>
          <CodeBlock>{`# Export aliases to file
mako alias export aliases.json
mako alias export ~/.mako/my-aliases.json

# Export specific tags
mako alias export git-aliases.json --tag git

# Import aliases from file
mako alias import aliases.json
mako alias import ~/.mako/my-aliases.json

# Import with merge (keep existing)
mako alias import aliases.json --merge`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Common Use Cases</h3>
          <CodeBlock>{`# Git shortcuts
mako alias gst "git status"
mako alias gco "git checkout"
mako alias gp "git push origin $(git branch --show-current)"

# Docker shortcuts
mako alias dps "docker ps -a"
mako alias dcu "docker-compose up -d"
mako alias dcd "docker-compose down"

# Development workflows
mako alias dev "npm install && npm run dev"
mako alias build "npm run build && npm run test"
mako alias deploy "npm run build && npm run deploy"`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Alias Already Exists</h4>
              <CodeBlock>{`mako alias save ll "ls -la"
# Warning: Alias 'll' already exists
# Use --force to overwrite: mako alias save ll "ls -la" --force`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Alias Not Found</h4>
              <CodeBlock>{`mako alias run nonexistent
# Error: Alias 'nonexistent' not found
# Use 'mako alias list' to see available aliases`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Invalid Alias Name</h4>
              <CodeBlock>{`mako alias save "my alias!" "echo hello"
# Error: Alias name contains invalid characters
# Use only letters, numbers, hyphens, and underscores`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Empty Command</h4>
              <CodeBlock>{`mako alias save test ""
# Error: Command cannot be empty`}</CodeBlock>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
