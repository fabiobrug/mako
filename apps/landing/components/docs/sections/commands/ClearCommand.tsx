"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const ClearCommand = () => (
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
      mako clear
    </motion.h1>
    <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
      Clear conversation history and context
    </p>
      <div className="space-y-6">
        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Clear Conversation</h3>
          <p className="text-muted-foreground text-sm leading-relaxed mb-3">
            Clears the current AI conversation context. Useful when you want to start fresh or if the AI is confused by previous context.
          </p>
          <CodeBlock>{`# Clear conversation history
mako clear

# Output:
# Conversation history cleared
# The AI will start fresh with no memory of previous interactions

# Note: This does NOT delete your command history
# Command history is managed separately via 'mako history'`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">What Gets Cleared</h3>
          <div className="bg-code border border-border rounded-lg p-4 mb-4">
            <h4 className="font-mono text-sm font-semibold text-foreground mb-3">✓ Cleared:</h4>
            <ul className="list-disc list-inside text-muted-foreground text-sm space-y-1">
              <li>AI conversation context</li>
              <li>Previous prompts and responses</li>
              <li>Multi-turn conversation state</li>
              <li>Temporary context from recent commands</li>
            </ul>
          </div>

          <div className="bg-code border border-border rounded-lg p-4">
            <h4 className="font-mono text-sm font-semibold text-foreground mb-3">✗ NOT Cleared:</h4>
            <ul className="list-disc list-inside text-muted-foreground text-sm space-y-1">
              <li>Command history database</li>
              <li>Configuration settings</li>
              <li>Saved aliases</li>
              <li>Statistics and analytics</li>
            </ul>
          </div>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Clear Options</h3>
          <CodeBlock>{`# Clear with confirmation prompt
mako clear

# Skip confirmation
mako clear --yes

# Clear and show what was cleared
mako clear --verbose

# Output:
# Cleared 15 messages from conversation
# Context window reset
# ✓ Ready for new conversation`}</CodeBlock>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">When to Use</h3>
          <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-4 leading-relaxed">
            <li><strong className="text-foreground">Context confusion:</strong> When AI responses seem unrelated to your current request</li>
            <li><strong className="text-foreground">Topic switching:</strong> Moving from one task to a completely different one</li>
            <li><strong className="text-foreground">Error accumulation:</strong> If multiple failed attempts are affecting suggestions</li>
            <li><strong className="text-foreground">Privacy:</strong> Before switching to a different project or task</li>
          </ul>
        </div>

        <div>
          <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
          <div className="space-y-3">
            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Nothing to Clear</h4>
              <CodeBlock>{`mako clear
# Conversation history is already empty
# Nothing to clear`}</CodeBlock>
            </div>

            <div className="bg-code border border-border rounded-lg p-4">
              <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Accidental Clear</h4>
              <p className="text-muted-foreground text-sm leading-relaxed mb-2">
                Conversation context cannot be recovered after clearing. However, your command history remains intact and can be searched.
              </p>
              <CodeBlock>{`# View recent commands even after clearing
mako history --last 20`}</CodeBlock>
            </div>
          </div>
        </div>
      </div>
    </motion.article>
  );
