"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const CompletionCommand = () => (
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
      mako completion
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Generate shell completion scripts for tab-completion
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Generate Completion Scripts</h3>
          <CodeBlock>{`# Generate completion for bash
mako completion bash > ~/.mako-completion.bash

# Generate completion for zsh
mako completion zsh > ~/.mako-completion.zsh

# Generate completion for fish
mako completion fish > ~/.config/fish/completions/mako.fish`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Bash Setup</h3>
          <CodeBlock>{`# Generate completion script
mako completion bash > ~/.mako-completion.bash

# Add to ~/.bashrc
echo "source ~/.mako-completion.bash" >> ~/.bashrc

# Reload shell
source ~/.bashrc

# Test tab completion
mako <TAB>        # Shows all available commands
mako h<TAB>       # Completes to 'mako history'
mako alias <TAB>  # Shows alias subcommands`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Zsh Setup</h3>
          <CodeBlock>{`# Generate completion script
mako completion zsh > ~/.mako-completion.zsh

# Add to ~/.zshrc
echo "source ~/.mako-completion.zsh" >> ~/.zshrc

# Reload shell
source ~/.zshrc

# Enable completion if not already enabled
# Add to ~/.zshrc before sourcing completion:
autoload -Uz compinit && compinit`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Fish Setup</h3>
          <CodeBlock>{`# Generate completion (fish auto-loads from this directory)
mako completion fish > ~/.config/fish/completions/mako.fish

# Reload completions
fish_update_completions

# Test
mako <TAB>  # Shows available commands with descriptions`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Completion Features</h3>
          <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-4 leading-relaxed">
            <li><strong className="text-foreground">Command completion:</strong> Tab-complete all mako commands</li>
            <li><strong className="text-foreground">Subcommand completion:</strong> Complete subcommands like <code className="font-mono text-primary bg-code px-1 py-0.5 rounded text-xs">alias save</code></li>
            <li><strong className="text-foreground">Flag completion:</strong> Complete command flags like <code className="font-mono text-primary bg-code px-1 py-0.5 rounded text-xs">--failed</code></li>
            <li><strong className="text-foreground">Dynamic completion:</strong> Complete saved alias names, config keys, etc.</li>
            <li><strong className="text-foreground">Contextual help:</strong> Shows command descriptions during completion</li>
          </ul>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Completion Examples</h3>
          <CodeBlock>{`# Command completion
$ mako <TAB>
ask        -- Generate command from natural language
history    -- Show and search command history
alias      -- Manage command aliases
config     -- Manage configuration
stats      -- Show usage statistics
health     -- Check system health

# Subcommand completion
$ mako history <TAB>
semantic      -- Search by meaning
--failed      -- Show only failed commands
--success     -- Show only successful
--interactive -- Browse interactively

# Alias name completion
$ mako alias run <TAB>
deploy    ll    findfile    psgrep

# Config key completion
$ mako config get <TAB>
llm_provider  llm_model  api_key  max_history`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Unsupported Shell</h4>
              <CodeBlock>{`mako completion powershell
# Error: Unsupported shell: powershell
# Supported shells: bash, zsh, fish`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Completion Not Working</h4>
              <CodeBlock>{`# Verify completion is sourced
$ which mako
/usr/local/bin/mako

# Check if completion file exists
$ ls -la ~/.mako-completion.bash

# Regenerate and reload
mako completion bash > ~/.mako-completion.bash
source ~/.bashrc`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Conflicting Completions</h4>
              <CodeBlock>{`# If you have multiple completion files
# Remove old ones and keep only one

# For bash
rm ~/.mako-completion.bash.old
mako completion bash > ~/.mako-completion.bash
source ~/.bashrc`}</CodeBlock>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
