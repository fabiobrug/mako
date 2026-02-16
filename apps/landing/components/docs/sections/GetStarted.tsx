"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../CodeBlock";
import { fadeInUp } from "../animations";

export const GetStartedSection = () => (
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
      Get Started
    </motion.h1>
    
    <motion.h2 
      variants={fadeInUp}
      className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4"
    >
      What is Mako?
    </motion.h2>
    <motion.p 
      variants={fadeInUp}
      className="text-muted-foreground mb-4 leading-relaxed"
    >
      <strong className="text-foreground">Mako</strong> is an AI-native shell orchestrator that wraps around your existing shell (bash/zsh) to provide intelligent command assistance. Generate commands from natural language, search your history semantically, and work faster with an AI that understands context.
    </motion.p>
    <motion.p 
      variants={fadeInUp}
      className="text-muted-foreground mb-6 leading-relaxed"
    >
      Unlike traditional command-line tools, Mako intercepts terminal I/O through a PTY (pseudo-terminal) and routes commands to AI for natural language processing, making your command-line experience more intuitive and productive.
    </motion.p>

    <motion.h2 
      variants={fadeInUp}
      className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4"
    >
      How It Works
    </motion.h2>
    <ol className="list-decimal list-inside text-muted-foreground space-y-3 mb-6 leading-relaxed">
      <li><strong className="text-foreground">Start Mako</strong> - Wraps around your bash/zsh shell</li>
      <li><strong className="text-foreground">Natural Language</strong> - Type <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">mako ask "compress this video"</code></li>
      <li><strong className="text-foreground">AI Generation</strong> - Your configured AI provider generates the appropriate shell command</li>
      <li><strong className="text-foreground">Review & Execute</strong> - Review the command before running it</li>
      <li><strong className="text-foreground">Learn & Improve</strong> - Mako learns your preferences over time</li>
    </ol>

    <motion.h2 
      variants={fadeInUp}
      className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4"
    >
      Key Features
    </motion.h2>
    
    <motion.h3 
      variants={fadeInUp}
      className="font-mono text-lg font-semibold text-foreground mt-6 mb-3"
    >
      Natural Language Commands
    </motion.h3>
    <motion.p 
      variants={fadeInUp}
      className="text-muted-foreground mb-4 leading-relaxed"
    >
      Generate shell commands from plain English. No need to remember complex syntax.
    </motion.p>
    <CodeBlock>{`mako ask "find all PDF files modified in the last week"
# Generates: find . -name "*.pdf" -mtime -7`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Semantic History Search</h3>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Search your command history by meaning, not exact text. Find that command you ran months ago.
    </p>
    <CodeBlock>{`mako history semantic "backup database"
# Finds commands like: pg_dump -U postgres mydb > backup.sql`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Context-Aware AI</h3>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Mako understands your current directory, recent output, and command patterns to provide better suggestions.
    </p>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Safety First</h3>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Detects potentially dangerous commands before execution. Prevents accidental data loss.
    </p>
    <CodeBlock>{`# Mako will warn you about:
rm -rf /
sudo dd if=/dev/zero of=/dev/sda`}</CodeBlock>
  </motion.article>
);
