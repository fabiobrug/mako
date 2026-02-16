"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const AskCommand = () => (
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
      mako ask &lt;prompt&gt;
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Generate shell commands from natural language descriptions
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Usage</h3>
          <p className="text-muted-foreground text-sm leading-relaxed mb-3">
            Transform natural language into executable shell commands. Mako analyzes your request and generates the most appropriate command.
          </p>
          <CodeBlock>{`mako ask "find all files larger than 100MB"
# Generates: find . -size +100M -type f

mako ask "compress this directory"
# Generates: tar -czf archive.tar.gz .

mako ask "show me disk usage"
# Generates: df -h`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Advanced Examples</h3>
          <CodeBlock>{`# Complex file operations
mako ask "find and delete all node_modules folders"
# Generates: find . -name "node_modules" -type d -exec rm -rf {} +

# Network operations
mako ask "check which process is using port 3000"
# Generates: lsof -i :3000

# Git operations
mako ask "undo my last git commit but keep the changes"
# Generates: git reset --soft HEAD~1

# System monitoring
mako ask "show top 10 processes by memory usage"
# Generates: ps aux --sort=-%mem | head -n 11`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Context-Aware Requests</h3>
          <p className="text-muted-foreground text-sm leading-relaxed mb-3">
            Mako understands your current directory and can provide context-specific commands.
          </p>
          <CodeBlock>{`# If you're in a Git repository
mako ask "push my changes"
# Generates: git push origin main

# If you're in a Node.js project
mako ask "install dependencies"
# Generates: npm install

# If you're in a Docker project
mako ask "rebuild containers"
# Generates: docker-compose up -d --build`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Options & Flags</h3>
          <CodeBlock>{`# Force execution without confirmation
mako ask "delete all log files" --force

# Get explanation without executing
mako ask "complicated command" --explain

# Generate multiple alternatives
mako ask "backup database" --alternatives 3`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases & Limitations</h3>
          <div className="space-y-3">
            <div className="border-l-2 border-yellow-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Ambiguous Requests</h4>
              <p className="text-muted-foreground text-sm leading-relaxed mb-2">
                If your request is ambiguous, Mako will ask for clarification:
              </p>
              <CodeBlock>{`mako ask "delete files"`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Output: Could you be more specific? Which files should I delete?
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Unsupported Operations</h4>
              <p className="text-muted-foreground text-sm leading-relaxed mb-2">
                Some operations require manual intervention:
              </p>
              <CodeBlock>{`mako ask "edit file and change line 5"`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Error: Interactive editing requires manual intervention. Suggestion: Use <code className="text-primary text-xs">nano filename.txt</code>
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Dangerous Commands</h4>
              <p className="text-muted-foreground text-sm leading-relaxed mb-2">
                Mako will warn about potentially destructive operations:
              </p>
              <CodeBlock>{`mako ask "delete everything in root"`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                WARNING: This command is destructive and affects system files. Do you want to proceed? [y/N]
              </p>
            </div>

            <div className="border-l-2 border-blue-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">No Internet Connection</h4>
              <CodeBlock>{`mako ask "find large files"`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Error: Cannot connect to AI service. Check your internet connection or API key.
              </p>
            </div>
          </div>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Common Failure Cases</h3>
          <div className="space-y-3">
            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <p className="font-mono text-sm font-semibold text-foreground mb-2">Invalid API Key</p>
              <CodeBlock>{`Error: Authentication failed`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Solution: Set your API key via <code className="text-primary text-xs">mako config set api_key your-key-here</code>
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <p className="font-mono text-sm font-semibold text-foreground mb-2">Rate Limit Exceeded</p>
              <CodeBlock>{`Error: API rate limit exceeded`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Solution: Wait a few seconds and try again
              </p>
            </div>

            <div className="border-l-2 border-red-500/50 pl-4 py-2 hover:bg-muted-foreground/5 transition-colors rounded-r">
              <p className="font-mono text-sm font-semibold text-foreground mb-2">Empty Prompt</p>
              <CodeBlock>{`mako ask ""`}</CodeBlock>
              <p className="text-muted-foreground text-sm leading-relaxed mt-2">
                Error: Please provide a description of what you want to do
              </p>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
