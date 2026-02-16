"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../CodeBlock";
import { fadeInUp } from "../animations";

export const HelpSection = () => (
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
      Help
    </motion.h1>
    
    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Troubleshooting</h2>
    
    <div className="space-y-6">
      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-base font-semibold text-foreground mb-3">Mako won't start</h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          Check that:
        </p>
        <ul className="list-disc list-inside text-muted-foreground text-sm space-y-1 mb-3">
          <li>Mako is installed: <code className="font-mono text-primary bg-background px-1 py-0.5 rounded">which mako</code></li>
          <li>Binary has execute permissions: <code className="font-mono text-primary bg-background px-1 py-0.5 rounded">chmod +x /usr/local/bin/mako</code></li>
          <li>Your shell is supported (bash or zsh)</li>
        </ul>
      </div>

      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-base font-semibold text-foreground mb-3">API key not working</h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          Verify your API key:
        </p>
        <CodeBlock>{`# Check current configuration
mako config list

# Or check config file
cat ~/.mako/config.json

# Test API connection
mako health`}</CodeBlock>
      </div>

      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-base font-semibold text-foreground mb-3">Commands not being intercepted</h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          Ensure shell integration is active:
        </p>
        <CodeBlock>{`# Check if running inside Mako
echo $MAKO_ACTIVE

# Restart Mako
exit
mako`}</CodeBlock>
      </div>

      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-base font-semibold text-foreground mb-3">Slow response times</h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          If AI responses are slow:
        </p>
        <ul className="list-disc list-inside text-muted-foreground text-sm space-y-1">
          <li>Check your internet connection</li>
          <li>Verify API quota with your provider</li>
          <li>Consider using a faster model or switching to local Ollama</li>
        </ul>
      </div>

      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-base font-semibold text-foreground mb-3">Database corruption</h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          If history search fails:
        </p>
        <CodeBlock>{`# Backup existing database
cp ~/.mako/mako.db ~/.mako/mako.db.backup

# Rebuild database
rm ~/.mako/mako.db
mako

# Database will be recreated on next start`}</CodeBlock>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Frequently Asked Questions</h2>
    
    <div className="space-y-6">
      <div>
        <h3 className="font-mono text-lg font-semibold text-foreground mb-2">Is my command history private?</h3>
        <p className="text-muted-foreground leading-relaxed">
          Your command history is stored locally in <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">~/.mako/mako.db</code>. Only command prompts are sent to your chosen AI provider for processing. Output and sensitive data stay on your machine. For maximum privacy, use Ollama (local models).
        </p>
      </div>

      <div>
        <h3 className="font-mono text-lg font-semibold text-foreground mb-2">Does Mako work offline?</h3>
        <p className="text-muted-foreground leading-relaxed">
          Mako requires an internet connection for AI features (command generation and semantic search). Your shell functions normally, but AI-powered features will be unavailable offline.
        </p>
      </div>

      <div>
        <h3 className="font-mono text-lg font-semibold text-foreground mb-2">Can I use my own AI model?</h3>
        <p className="text-muted-foreground leading-relaxed">
          Yes! Mako supports multiple AI providers including OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter, and local models via Ollama. You can configure your preferred provider using <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">mako config set</code>.
        </p>
      </div>

      <div>
        <h3 className="font-mono text-lg font-semibold text-foreground mb-2">How much does it cost?</h3>
        <p className="text-muted-foreground leading-relaxed">
          Mako itself is free and open source. Costs depend on your AI provider: Ollama is completely free (local), Gemini has a generous free tier, while OpenAI, Anthropic, and others are pay-per-use. Most users with cloud providers stay within free tiers.
        </p>
      </div>

      <div>
        <h3 className="font-mono text-lg font-semibold text-foreground mb-2">Can I sync history across machines?</h3>
        <p className="text-muted-foreground leading-relaxed">
          Not yet. History synchronization is planned for a future release. Currently, each machine maintains its own local history.
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Getting Help</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-base font-semibold text-foreground mb-3">GitHub Issues</h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          Report bugs and request features on our GitHub repository:
        </p>
        <a 
          href="https://github.com/fabiobrug/mako/issues" 
          target="_blank" 
          rel="noopener noreferrer"
          className="text-primary hover:text-primary/80 text-sm"
        >
          github.com/fabiobrug/mako/issues
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-base font-semibold text-foreground mb-3">GitHub Discussions</h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          Ask questions and share ideas with the community:
        </p>
        <a 
          href="https://github.com/fabiobrug/mako/discussions" 
          target="_blank" 
          rel="noopener noreferrer"
          className="text-primary hover:text-primary/80 text-sm"
        >
          github.com/fabiobrug/mako/discussions
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-base font-semibold text-foreground mb-3">Contributing</h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          We welcome contributions! Check out our contributing guide:
        </p>
        <a 
          href="https://github.com/fabiobrug/mako/blob/dev/docs/CONTRIBUTING.md" 
          target="_blank" 
          rel="noopener noreferrer"
          className="text-primary hover:text-primary/80 text-sm"
        >
          Contributing Guide
        </a>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Uninstalling</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      If you need to uninstall Mako:
    </p>
    <CodeBlock>{`# Remove binaries
sudo rm /usr/local/bin/mako /usr/local/bin/mako-menu

# Remove configuration and data
rm -rf ~/.mako

# Remove shell integration (remove from ~/.bashrc or ~/.zshrc)
# Look for lines containing "mako" and remove them

# Reload shell
source ~/.bashrc  # or source ~/.zshrc`}</CodeBlock>
  </motion.article>
);
