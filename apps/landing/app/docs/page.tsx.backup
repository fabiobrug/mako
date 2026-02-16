"use client";

import { useState, useRef } from "react";
import Link from "next/link";
import { motion, useInView } from "framer-motion";
import { Github, Menu, X, Terminal, Settings, Wrench, Monitor, BookOpen, HelpCircle, Rocket, ChevronDown, Copy, Check } from "lucide-react";

const sections = [
  { id: "get-started", label: "Get Started", icon: Rocket },
  { id: "install", label: "Install", icon: Terminal },
  { id: "agent-config", label: "Agent Config", icon: Settings },
  { id: "commands", label: "Commands", icon: Wrench },
  { id: "platforms", label: "Platforms", icon: Monitor },
  { id: "reference", label: "Reference", icon: BookOpen },
  { id: "help", label: "Help", icon: HelpCircle },
];

export default function Docs() {
  const [activeSection, setActiveSection] = useState("get-started");
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [commandsExpanded, setCommandsExpanded] = useState(false);
  const [activeCommand, setActiveCommand] = useState<string | null>(null);

  return (
    <div className="min-h-screen flex flex-col">
      {/* Header */}
      <motion.header 
        initial={{ y: -100 }}
        animate={{ y: 0 }}
        transition={{ type: "spring", stiffness: 100, damping: 20 }}
        className="border-b border-border sticky top-0 z-20 bg-background/95 backdrop-blur"
      >
        <div className="max-w-6xl mx-auto px-6 h-12 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <motion.button
              className="md:hidden text-muted-foreground hover:text-foreground"
              onClick={() => setSidebarOpen(!sidebarOpen)}
              aria-label="Toggle menu"
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.9 }}
            >
              {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
            </motion.button>
            <Link href="/" className="font-mono font-bold text-lg tracking-tight text-foreground hover:text-primary transition-colors">
              Mako
            </Link>
          </div>
          <motion.a
            href="https://github.com/fabiobrug/mako"
            target="_blank"
            rel="noopener noreferrer"
            className="text-primary hover:text-primary/80 transition-colors"
            aria-label="GitHub"
            whileHover={{ scale: 1.1, rotate: 5 }}
            whileTap={{ scale: 0.95 }}
          >
            <Github size={20} />
          </motion.a>
        </div>
      </motion.header>

      <div className="flex flex-1 max-w-6xl mx-auto w-full">
        {/* Sidebar */}
        <aside
          className={`
            fixed inset-y-0 left-0 z-10 w-60 bg-background border-r border-border pt-14 px-4
            transform transition-transform md:relative md:translate-x-0 md:pt-0 md:mt-0
            ${sidebarOpen ? "translate-x-0" : "-translate-x-full"}
          `}
        >
          <nav className="py-6 space-y-1">
            {sections.map((s) => (
              <div key={s.id}>
                {s.id === "commands" ? (
                  <>
                    <button
                      onClick={() => {
                        if (commandsExpanded) {
                          setCommandsExpanded(false);
                        } else {
                          setCommandsExpanded(true);
                          setActiveSection("commands");
                          setActiveCommand("ask");
                        }
                      }}
                      className={`w-full text-left px-3 py-2 rounded text-sm font-medium transition-all flex items-center justify-between cursor-pointer ${
                        activeSection === "commands"
                          ? "text-primary bg-secondary"
                          : "text-muted-foreground hover:text-foreground hover:bg-secondary/50"
                      }`}
                    >
                      <div className="flex items-center gap-2">
                        <s.icon size={16} />
                        {s.label}
                      </div>
                      <motion.div
                        animate={{ rotate: commandsExpanded ? 180 : 0 }}
                        transition={{ duration: 0.3 }}
                      >
                        <ChevronDown size={14} />
                      </motion.div>
                    </button>
                    <motion.div
                      initial={false}
                      animate={{ height: commandsExpanded ? "auto" : 0, opacity: commandsExpanded ? 1 : 0 }}
                      transition={{ duration: 0.3, ease: "easeInOut" }}
                      className="overflow-hidden"
                    >
                      <div className="ml-6 mt-1 space-y-1 border-l border-border pl-3">
                        {[
                          { id: "ask", label: "mako ask" },
                          { id: "history", label: "mako history" },
                          { id: "stats", label: "mako stats" },
                          { id: "export", label: "mako export" },
                          { id: "alias", label: "mako alias" },
                          { id: "health", label: "mako health" },
                          { id: "help", label: "mako help" },
                        ].map((cmd) => (
                          <button
                            key={cmd.id}
                            onClick={() => {
                              setActiveCommand(cmd.id);
                              setActiveSection("commands");
                              setSidebarOpen(false);
                            }}
                            className={`w-full text-left px-3 py-1.5 rounded text-xs font-mono transition-all cursor-pointer ${
                              activeCommand === cmd.id
                                ? "text-primary bg-secondary/70"
                                : "text-muted-foreground hover:text-foreground hover:bg-secondary/30"
                            }`}
                          >
                            {cmd.label}
                          </button>
                        ))}
                      </div>
                    </motion.div>
                  </>
                ) : (
                  <button
                    onClick={() => { setActiveSection(s.id); setSidebarOpen(false); }}
                    className={`w-full text-left px-3 py-2 rounded text-sm font-medium transition-all flex items-center gap-2 cursor-pointer ${
                      activeSection === s.id
                        ? "text-primary bg-secondary"
                        : "text-muted-foreground hover:text-foreground hover:bg-secondary/50"
                    }`}
                  >
                    <s.icon size={16} />
                    {s.label}
                  </button>
                )}
              </div>
            ))}
          </nav>
        </aside>

        {/* Overlay for mobile */}
        {sidebarOpen && (
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-[5] bg-background/60 md:hidden" 
            onClick={() => setSidebarOpen(false)} 
          />
        )}

        {/* Content */}
        <motion.main 
          key={activeSection}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4 }}
          className="flex-1 px-6 md:px-12 py-10 min-w-0"
        >
          {activeSection === "get-started" && <GetStartedSection />}
          {activeSection === "install" && <InstallSection />}
          {activeSection === "agent-config" && <AgentConfigSection />}
          {activeSection === "commands" && <CommandsSection activeCommand={activeCommand} />}
          {activeSection === "platforms" && <PlatformsSection />}
          {activeSection === "reference" && <ReferenceSection />}
          {activeSection === "help" && <HelpSection />}
        </motion.main>
      </div>
    </div>
  );
}

const CodeBlock = ({ children, language = "bash" }: { children: string; language?: string }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(children);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  // Parse the code to add syntax highlighting
  const highlightBashSyntax = (code: string) => {
    const lines = code.split('\n');
    return lines.map((line, index) => {
      // Comment lines
      if (line.trim().startsWith('#')) {
        return (
          <div key={index} className="text-muted-foreground/60">
            {line}
          </div>
        );
      }
      
      // Lines with command prompt
      if (line.startsWith('$') || line.startsWith('mako')) {
        const parts = line.split(' ');
        return (
          <div key={index}>
            {parts.map((part, i) => {
              if (part === '$') {
                return <span key={i} className="text-success mr-2">{part}</span>;
              }
              if (part === 'mako' || part.startsWith('mako')) {
                return <span key={i} className="text-primary mr-1">{part}</span>;
              }
              if (part.startsWith('"') || part.startsWith("'")) {
                return <span key={i} className="text-[#f59e0b] mr-1">{part}</span>;
              }
              if (part.startsWith('--') || part.startsWith('-')) {
                return <span key={i} className="text-[#8b5cf6] mr-1">{part}</span>;
              }
              return <span key={i} className="text-foreground mr-1">{part}</span>;
            })}
          </div>
        );
      }
      
      // Output lines
      return (
        <div key={index} className="text-muted-foreground/80">
          {line}
        </div>
      );
    });
  };

  return (
    <motion.div 
      initial={{ opacity: 0, y: 10 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ duration: 0.3 }}
      className="my-4 rounded-lg overflow-hidden shadow-lg group relative"
    >
      {/* Mac terminal header */}
      <div className="bg-[#2d2d2d] px-4 py-2 flex items-center gap-2 border-b border-[#1a1a1a] justify-between">
        <div className="flex gap-1.5">
          <motion.div 
            className="w-3 h-3 rounded-full bg-[#ff5f57]"
            whileHover={{ scale: 1.2 }}
            transition={{ type: "spring", stiffness: 400 }}
          />
          <motion.div 
            className="w-3 h-3 rounded-full bg-[#febc2e]"
            whileHover={{ scale: 1.2 }}
            transition={{ type: "spring", stiffness: 400 }}
          />
          <motion.div 
            className="w-3 h-3 rounded-full bg-[#28c840]"
            whileHover={{ scale: 1.2 }}
            transition={{ type: "spring", stiffness: 400 }}
          />
        </div>
        <motion.button
          onClick={handleCopy}
          className="opacity-0 group-hover:opacity-100 transition-opacity p-1.5 rounded hover:bg-[#1a1a1a]"
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.9 }}
        >
          {copied ? (
            <Check size={14} className="text-success" />
          ) : (
            <Copy size={14} className="text-muted-foreground" />
          )}
        </motion.button>
      </div>
      {/* Terminal content */}
      <pre className="bg-[#1e1e1e] p-4 overflow-x-auto">
        <code className="font-mono text-sm leading-relaxed">
          {language === "bash" ? highlightBashSyntax(children) : children}
        </code>
      </pre>
    </motion.div>
  );
};

// Animation variants for sections
const fadeInUp = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.5 } }
};

const staggerContainer = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.1 }
  }
};

const GetStartedSection = () => (
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

const InstallSection = () => (
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
      Installation
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Mako works on <strong className="text-foreground">Linux</strong> and <strong className="text-foreground">macOS</strong>. Choose from multiple AI providers including local models (Ollama) or cloud services (OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter).
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Option 1: One-Line Install (Recommended)</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Fast installation with optional environment variable configuration:
    </p>
    <CodeBlock>{`# Basic installation
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash

# Or install with provider configuration
LLM_PROVIDER=openai LLM_MODEL=gpt-4o-mini LLM_API_KEY=sk-your-key \\
curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash`}</CodeBlock>
    
    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      <strong className="text-foreground">After installation, configure your AI provider:</strong>
    </p>
    <CodeBlock>{`# Start Mako
mako

# Inside Mako shell, configure your provider:
mako config set llm_provider openai
mako config set llm_model gpt-4o-mini
mako config set api_key sk-your-api-key

# Or for Ollama (local):
mako config set llm_provider ollama
mako config set llm_model llama3.2
mako config set llm_base_url http://localhost:11434

# View all settings:
mako config list`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      <strong className="text-foreground">Supported configuration keys:</strong>
    </p>
    <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li><code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">llm_provider</code> - AI provider (openai, anthropic, gemini, deepseek, openrouter, ollama)</li>
      <li><code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">llm_model</code> - Model name (provider-specific)</li>
      <li><code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">llm_base_url</code> - Base URL (optional, for custom endpoints)</li>
      <li><code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">api_key</code> - Your API key (not required for Ollama)</li>
    </ul>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Option 2: From Source with .env File</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Clone the repository and configure via <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">.env</code> file:
    </p>
    <CodeBlock>{`# Clone the repository
git clone https://github.com/fabiobrug/mako.git
cd mako/apps/cli

# Copy and edit configuration
cp .env.example .env
nano .env  # Edit with your provider settings

# Build
make build

# Install (optional, requires sudo)
make install

# Or run directly
./mako`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-6 leading-relaxed">
      <strong className="text-foreground">Example <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">.env</code> configuration:</strong>
    </p>
    <CodeBlock>{`# OpenAI
LLM_PROVIDER=openai
LLM_MODEL=gpt-4o-mini
LLM_API_KEY=sk-your-key

# Or Ollama (local, free)
LLM_PROVIDER=ollama
LLM_MODEL=llama3.2
LLM_API_BASE=http://localhost:11434`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Verify Installation</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      After installation, verify that Mako is working correctly:
    </p>
    <CodeBlock>{`# Start Mako shell
mako

# Inside Mako shell, try:
mako ask "find files larger than 100MB"
mako history
mako health    # Check configuration status
mako help`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configuration Priority</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Mako checks for configuration in this order:
    </p>
    <ol className="list-decimal list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li><strong className="text-foreground">Environment variables</strong> (<code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">.env</code> file in <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">apps/cli/</code>)</li>
      <li><strong className="text-foreground">Config file</strong> (<code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">~/.mako/config.json</code>) - set via <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">mako config set</code></li>
      <li><strong className="text-foreground">Default values</strong> (Gemini provider)</li>
    </ol>

    <p className="text-muted-foreground mb-4 leading-relaxed">
      For more details on configuring AI providers, see the <button onClick={() => window.dispatchEvent(new CustomEvent('changeSection', { detail: 'agent-config' }))} className="text-primary hover:text-primary/80 underline">Agent Config</button> section.
    </p>
  </motion.article>
);

const AgentConfigSection = () => (
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
      Agent Config
    </motion.h1>
    
    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Supported AI Providers</h2>
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Mako supports multiple AI providers. Configure your preferred provider using environment variables or CLI commands:
    </p>

    <div className="overflow-x-auto mb-8">
      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-b border-border">
            <th className="text-left py-3 px-4 font-mono text-foreground">Provider</th>
            <th className="text-left py-3 px-4 font-mono text-foreground">Type</th>
            <th className="text-left py-3 px-4 font-mono text-foreground">Cost</th>
            <th className="text-left py-3 px-4 font-mono text-foreground">Best For</th>
          </tr>
        </thead>
        <tbody className="text-muted-foreground">
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">Ollama</td>
            <td className="py-3 px-4">Local</td>
            <td className="py-3 px-4 text-success">Free</td>
            <td className="py-3 px-4">Privacy, offline use, no API costs</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">OpenAI</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4">Paid</td>
            <td className="py-3 px-4">Best quality, GPT-4o models</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">Anthropic</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4">Paid</td>
            <td className="py-3 px-4">Claude models, great reasoning</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">Gemini</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4 text-success">Free tier</td>
            <td className="py-3 px-4">Default option, good balance</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">OpenRouter</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4">Paid</td>
            <td className="py-3 px-4">Access to multiple models</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4 font-mono">DeepSeek</td>
            <td className="py-3 px-4">Cloud</td>
            <td className="py-3 px-4">Paid</td>
            <td className="py-3 px-4">Cost-effective alternative</td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Example Configurations</h2>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Ollama (Local, Free)</h3>
    <CodeBlock>{`LLM_PROVIDER=ollama
LLM_MODEL=llama3.2
LLM_API_BASE=http://localhost:11434`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">OpenAI</h3>
    <CodeBlock>{`LLM_PROVIDER=openai
LLM_MODEL=gpt-4o-mini
LLM_API_KEY=sk-your-api-key-here`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Anthropic (Claude)</h3>
    <CodeBlock>{`LLM_PROVIDER=anthropic
LLM_MODEL=claude-3-5-haiku-20241022
LLM_API_KEY=sk-ant-your-key`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Google Gemini</h3>
    <CodeBlock>{`LLM_PROVIDER=gemini
LLM_MODEL=gemini-2.5-flash
LLM_API_KEY=your-gemini-key`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Using Ollama (Local Models)</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Ollama allows you to run AI models locally on your machine:
    </p>
    <CodeBlock>{`# Install Ollama
curl https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama3.2

# Configure Mako
mako config set llm_provider ollama
mako config set llm_model llama3.2
mako config set llm_base_url http://localhost:11434

# Start Mako
mako`}</CodeBlock>

    <p className="text-muted-foreground mb-4 leading-relaxed">
      <strong className="text-foreground">Benefits of Ollama:</strong>
    </p>
    <ul className="list-disc list-inside text-muted-foreground space-y-2 mb-6 leading-relaxed">
      <li>✅ Completely free</li>
      <li>✅ Works offline</li>
      <li>✅ Privacy - data never leaves your machine</li>
      <li>✅ No API rate limits</li>
    </ul>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configuration Methods</h2>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Via CLI Commands</h3>
    <CodeBlock>{`# Set provider and credentials
mako config set llm_provider openai
mako config set llm_model gpt-4o-mini
mako config set api_key sk-your-api-key

# View current configuration
mako config list

# Check provider status
mako health`}</CodeBlock>

    <h3 className="font-mono text-lg font-semibold text-foreground mt-6 mb-3">Via Configuration File</h3>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Mako stores its configuration in <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">~/.mako/config.json</code>:
    </p>
    <CodeBlock language="json">{`{
  "llm_provider": "openai",
  "llm_model": "gpt-4o-mini",
  "api_key": "your-api-key-here",
  "llm_base_url": "https://api.openai.com/v1",
  "max_history": 1000,
  "auto_suggest": true,
  "safety_checks": true
}`}</CodeBlock>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Getting API Keys</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">OpenAI</h4>
        <a href="https://platform.openai.com/api-keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          platform.openai.com/api-keys
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Anthropic</h4>
        <a href="https://console.anthropic.com/account/keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          console.anthropic.com/account/keys
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Google Gemini</h4>
        <a href="https://ai.google.dev/" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          ai.google.dev
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">OpenRouter</h4>
        <a href="https://openrouter.ai/keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          openrouter.ai/keys
        </a>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <h4 className="font-mono text-sm font-semibold text-foreground mb-2">DeepSeek</h4>
        <a href="https://platform.deepseek.com/api_keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 text-sm">
          platform.deepseek.com/api_keys
        </a>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Configuration Options</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">max_history</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Maximum number of commands to store in history (default: 1000)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">auto_suggest</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Enable automatic command suggestions (default: true)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">safety_checks</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Enable safety validation for dangerous commands (default: true)
        </p>
      </div>
    </div>
  </motion.article>
);

const CommandsSection = ({ activeCommand }: { activeCommand: string | null }) => {
  // If no command is selected, show overview
  if (!activeCommand) {
    return (
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
          Commands
        </motion.h1>
        
        <p className="text-muted-foreground mb-6 leading-relaxed">
          Mako provides a comprehensive suite of commands to enhance your terminal experience. All commands are prefixed with <code className="font-mono text-primary bg-code px-1.5 py-0.5 rounded text-sm">mako</code>.
        </p>
        
        <motion.p 
          variants={fadeInUp}
          className="text-muted-foreground mb-6 leading-relaxed"
        >
          Select a command from the sidebar to see detailed usage, examples, and edge cases.
        </motion.p>
      </motion.article>
    );
  }

  // Render the specific command content
  if (activeCommand === "ask") {
    return (
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
              <div className="space-y-4">
                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">WARNING: Ambiguous Requests</h4>
                  <p className="text-muted-foreground text-sm leading-relaxed mb-2">
                    If your request is ambiguous, Mako will ask for clarification:
                  </p>
                  <CodeBlock>{`mako ask "delete files"
# Output: Could you be more specific? Which files should I delete?`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">ERROR: Unsupported Operations</h4>
                  <p className="text-muted-foreground text-sm leading-relaxed mb-2">
                    Some operations require manual intervention:
                  </p>
                  <CodeBlock>{`mako ask "edit file and change line 5"
# Error: Interactive editing requires manual intervention
# Suggestion: Use: nano filename.txt`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">SECURITY: Dangerous Commands</h4>
                  <p className="text-muted-foreground text-sm leading-relaxed mb-2">
                    Mako will warn about potentially destructive operations:
                  </p>
                  <CodeBlock>{`mako ask "delete everything in root"
# WARNING: This command is destructive and affects system files
# Command: rm -rf /
# Do you want to proceed? [y/N]`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">NETWORK: No Internet Connection</h4>
                  <CodeBlock>{`mako ask "find large files"
# Error: Cannot connect to AI service
# Check your internet connection or API key`}</CodeBlock>
                </div>
              </div>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Common Failure Cases</h3>
              <div className="space-y-3">
                <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
                  <p className="font-mono text-sm text-foreground mb-2">Invalid API Key</p>
                  <CodeBlock>{`Error: Authentication failed
Solution: Set your API key via: mako config set api_key your-key-here`}</CodeBlock>
                </div>

                <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
                  <p className="font-mono text-sm text-foreground mb-2">Rate Limit Exceeded</p>
                  <CodeBlock>{`Error: API rate limit exceeded
Solution: Wait a few seconds and try again`}</CodeBlock>
                </div>

                <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
                  <p className="font-mono text-sm text-foreground mb-2">Empty Prompt</p>
                  <CodeBlock>{`mako ask ""
Error: Please provide a description of what you want to do`}</CodeBlock>
                </div>
              </div>
            </div>
          </div>
        </motion.article>
      );
  }

  if (activeCommand === "history") {
    return (
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
          mako history
        </motion.h1>
        <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
          View, search, and manage your command history
        </p>
          <div className="space-y-6">
            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Usage</h3>
              <CodeBlock>{`# Show recent command history (last 50)
mako history

# Show last N commands
mako history --limit 20
mako history -n 100

# Show history with timestamps
mako history --timestamps`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Text Search</h3>
              <p className="text-muted-foreground text-sm leading-relaxed mb-3">
                Search for commands containing specific text (case-insensitive).
              </p>
              <CodeBlock>{`# Search for commands containing "docker"
mako history --search "docker"
mako history -s "docker"

# Search for multiple terms
mako history --search "git push"

# Search with wildcards
mako history --search "npm *"`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Semantic Search</h3>
              <p className="text-muted-foreground text-sm leading-relaxed mb-3">
                Find commands by meaning, not just exact text matches. Uses AI embeddings for intelligent search.
              </p>
              <CodeBlock>{`# Find database-related commands
mako history semantic "database backup"
# Finds: pg_dump, mysqldump, mongodump, etc.

# Find deployment commands
mako history semantic "deploy to production"
# Finds: git push, kubectl apply, docker push, etc.

# Find file compression commands
mako history semantic "compress files"
# Finds: tar, gzip, zip, 7z, etc.`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Filtering Options</h3>
              <CodeBlock>{`# Filter by exit code (successful commands only)
mako history --success
mako history --exit-code 0

# Filter by date range
mako history --since "2024-01-01"
mako history --since "yesterday"
mako history --since "1 week ago"

# Filter by directory
mako history --dir "/home/user/projects"

# Combine filters
mako history --success --since "1 week ago" --search "git"`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Output Formats</h3>
              <CodeBlock>{`# JSON output
mako history --format json

# CSV output
mako history --format csv

# Detailed output with metadata
mako history --verbose

# Copy command to clipboard
mako history --limit 1 --copy`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">History Management</h3>
              <CodeBlock>{`# Clear all history
mako history --clear

# Clear last N entries
mako history --clear --limit 10

# Remove specific command by ID
mako history --remove 1234

# Clear failed commands
mako history --clear --failed`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
              <div className="space-y-3">
                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Empty History</h4>
                  <CodeBlock>{`mako history
# Output: No command history found
# Run some commands to start building your history`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Corrupted Database</h4>
                  <CodeBlock>{`Error: Database error: unable to read history
Solution: Run 'mako health' to diagnose and repair`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Semantic Search Without Internet</h4>
                  <CodeBlock>{`mako history semantic "backup"
# Falls back to text search when offline
# Using text search instead (AI features unavailable)`}</CodeBlock>
                </div>
              </div>
            </div>
          </div>
        </motion.article>
      );
  }

  if (activeCommand === "stats") {
    return (
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
          mako stats
        </motion.h1>
        <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
          View detailed usage statistics and command patterns
        </p>
          <div className="space-y-6">
            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Basic Statistics</h3>
              <CodeBlock>{`mako stats

# Output:
# Mako Statistics
# ─────────────────────────────────────
# Total commands executed:    1,234
# Success rate:                94.5%
# Failed commands:             68
# Average commands per day:    42
# Most productive day:         Monday (avg 58 commands)
# First command:               2024-01-15
# Last command:                2024-02-12 (5 minutes ago)`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Detailed Breakdown</h3>
              <CodeBlock>{`# Show detailed statistics
mako stats --verbose

# Most used commands
mako stats --top-commands 10

# Command frequency by hour
mako stats --by-hour

# Command frequency by day of week
mako stats --by-day

# Command frequency by directory
mako stats --by-directory`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Time-Based Analysis</h3>
              <CodeBlock>{`# Statistics for specific period
mako stats --since "1 week ago"
mako stats --since "2024-01-01"

# Compare periods
mako stats --compare "last week" "this week"

# Monthly summary
mako stats --monthly

# Year in review
mako stats --year 2024`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Pattern Analysis</h3>
              <CodeBlock>{`# Find command patterns
mako stats --patterns

# Output:
# Common Patterns Detected:
# 1. Git workflow: commit → push (89% of git operations)
# 2. Docker workflow: build → up → logs
# 3. npm workflow: install → test → start
# 4. Most efficient time: 2pm-4pm (highest success rate)

# Productivity insights
mako stats --insights

# Command similarity analysis
mako stats --duplicates`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Export Options</h3>
              <CodeBlock>{`# Export stats as JSON
mako stats --format json --output stats.json

# Export as CSV
mako stats --format csv --output stats.csv

# Generate HTML report
mako stats --format html --output report.html`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
              <div className="space-y-3">
                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Insufficient Data</h4>
                  <CodeBlock>{`mako stats
# Output: Not enough data for meaningful statistics
# Run at least 10 commands to see statistics`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Invalid Date Range</h4>
                  <CodeBlock>{`mako stats --since "invalid date"
# Error: Invalid date format
# Use: YYYY-MM-DD or relative dates like "1 week ago"`}</CodeBlock>
                </div>
              </div>
            </div>
          </div>
        </motion.article>
      );
  }

  if (activeCommand === "export") {
    return (
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
          Export command history to various formats
        </p>
          <div className="space-y-6">
            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Supported Formats</h3>
              <CodeBlock>{`# Export to JSON (default)
mako export --output history.json
mako export -o history.json

# Export to CSV
mako export --format csv --output history.csv

# Export to plain text
mako export --format text --output history.txt

# Export to Markdown
mako export --format markdown --output history.md

# Export to SQL dump
mako export --format sql --output history.sql`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Filtering Exports</h3>
              <CodeBlock>{`# Export only successful commands
mako export --success --output success.json

# Export commands from specific date range
mako export --since "2024-01-01" --until "2024-02-01" -o jan.json

# Export commands matching pattern
mako export --search "docker" --output docker-commands.json

# Export last N commands
mako export --limit 100 --output recent.json`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Export Options</h3>
              <CodeBlock>{`# Include timestamps
mako export --timestamps -o history.json

# Include exit codes
mako export --exit-codes -o history.json

# Include working directory
mako export --include-pwd -o history.json

# Include environment variables (be careful!)
mako export --include-env --output history.json

# Compress output
mako export --compress -o history.json.gz

# Pretty print JSON
mako export --pretty -o history.json`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Backup & Restore</h3>
              <CodeBlock>{`# Create full backup
mako export --full-backup --output mako-backup-$(date +%Y%m%d).json

# Import from backup
mako import --file mako-backup-20240212.json

# Merge with existing history
mako import --merge --file history.json`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Edge Cases</h3>
              <div className="space-y-3">
                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">File Already Exists</h4>
                  <CodeBlock>{`mako export -o history.json
# Error: File already exists
# Use --force to overwrite or choose a different filename`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Permission Denied</h4>
                  <CodeBlock>{`mako export -o /root/history.json
# Error: Permission denied
# Choose a directory you have write access to`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Disk Space</h4>
                  <CodeBlock>{`# Large export fails
Error: Not enough disk space
Current history size: 2.5GB
Available space: 1.2GB
Solution: Use --compress or export specific date range`}</CodeBlock>
                </div>
              </div>
            </div>
          </div>
        </motion.article>
      );
  }

  if (activeCommand === "alias") {
    return (
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
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Creating Aliases</h3>
              <CodeBlock>{`# Simple alias
mako alias ll "ls -lah"

# Complex alias with pipes
mako alias psgrep "ps aux | grep"

# Multi-command alias
mako alias deploy "git push && ssh prod 'cd /app && ./deploy.sh'"

# Alias with parameters (use $1, $2, etc.)
mako alias findfile "find . -name '$1' -type f"

# Interactive alias with confirmation
mako alias clean "rm -rf node_modules && npm install"`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Managing Aliases</h3>
              <CodeBlock>{`# List all aliases
mako alias list

# Show specific alias
mako alias show deploy

# Edit existing alias
mako alias edit deploy

# Remove alias
mako alias remove deploy
mako alias rm deploy

# Remove all aliases
mako alias clear --confirm`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Using Aliases</h3>
              <CodeBlock>{`# After creating: mako alias ll "ls -lah"
ll  # Executes: ls -lah

# With parameters
# After creating: mako alias findfile "find . -name '$1'"
findfile "*.js"  # Executes: find . -name '*.js'

# Combine with other commands
ll | grep Documents
findfile "*.txt" | xargs cat`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Advanced Features</h3>
              <CodeBlock>{`# Add description to alias
mako alias deploy "git push && ssh prod './deploy.sh'" --description "Deploy to production"

# Make alias permanent (add to shell config)
mako alias ll "ls -lah" --permanent

# Create temporary alias (session only)
mako alias temp-cmd "echo hello" --temporary

# Export aliases
mako alias export --output aliases.json

# Import aliases
mako alias import --file aliases.json`}</CodeBlock>
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
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Alias Name Conflicts</h4>
                  <CodeBlock>{`mako alias ls "ls -lah"
# Warning: 'ls' conflicts with existing system command
# Use --force to override, or choose a different name`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Invalid Characters</h4>
                  <CodeBlock>{`mako alias my-alias! "echo hello"
# Error: Alias name contains invalid characters
# Use only letters, numbers, hyphens, and underscores`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Recursive Alias</h4>
                  <CodeBlock>{`mako alias ll "ll -a"
# Error: Recursive alias detected
# An alias cannot reference itself`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">Empty Command</h4>
                  <CodeBlock>{`mako alias test ""
# Error: Command cannot be empty`}</CodeBlock>
                </div>
              </div>
            </div>
          </div>
        </motion.article>
      );
  }

  if (activeCommand === "health") {
    return (
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
                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">ERROR: API Key Not Found</h4>
                  <CodeBlock>{`[FAIL] API Key not configured
Solution: Configure your API key
  mako config set api_key your-key-here
Or set environment variable for your provider`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">ERROR: Database Corruption</h4>
                  <CodeBlock>{`[FAIL] Database integrity check failed
Solution: Repair database with:
  mako health --repair-database
Or backup and recreate:
  cp ~/.mako/mako.db ~/.mako/mako.db.backup
  rm ~/.mako/mako.db
  mako  # Will recreate on next start`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">ERROR: Shell Integration Inactive</h4>
                  <CodeBlock>{`[FAIL] Shell integration not detected
Solution: Reinstall shell hooks:
  mako health --reinstall-shell
Then restart your shell:
  exec $SHELL`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">WARNING: Low Disk Space</h4>
                  <CodeBlock>{`[WARNING] Low disk space detected (< 1GB available)
Database size: 850MB
Recommendation: Clean up old history
  mako history --clear --before "6 months ago"`}</CodeBlock>
                </div>

                <div className="bg-code border border-border rounded-lg p-4">
                  <h4 className="font-mono text-sm font-semibold text-foreground mb-2">WARNING: Outdated Version</h4>
                  <CodeBlock>{`[INFO] Update available: v1.5.0 (current: v1.4.2)
Run: curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash`}</CodeBlock>
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
  }

  if (activeCommand === "help") {
    return (
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
          mako help
        </motion.h1>
        <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
          Display help information and command documentation
        </p>
          <div className="space-y-6">
            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">General Help</h3>
              <CodeBlock>{`# Show all available commands
mako help

# Show version information
mako --version
mako -v

# Show quick start guide
mako help quickstart`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Command-Specific Help</h3>
              <CodeBlock>{`# Get help for specific command
mako help ask
mako help history
mako help stats
mako help export
mako help alias
mako help health

# Alternative syntax
mako ask --help
mako history --help`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Examples & Tutorials</h3>
              <CodeBlock>{`# Show common examples
mako help examples

# Show interactive tutorial
mako tutorial

# Show best practices
mako help best-practices`}</CodeBlock>
            </div>

            <div>
              <h3 className="font-mono text-base font-semibold text-foreground mb-3">Documentation Links</h3>
              <CodeBlock>{`# Open online documentation
mako help --docs

# Open GitHub repository
mako help --github

# Report an issue
mako help --report-bug

# Request a feature
mako help --request-feature`}</CodeBlock>
            </div>
          </div>
        </motion.article>
      );
  }

  // Fallback
  return (
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
        Commands
      </motion.h1>
      <motion.p 
        variants={fadeInUp}
        className="text-muted-foreground mb-6 leading-relaxed"
      >
        Select a command from the sidebar to view its documentation.
      </motion.p>
    </motion.article>
  );
};

const PlatformsSection = () => (
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
      Platforms
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Mako is designed to work seamlessly across different platforms and shells.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Supported Operating Systems</h2>
    
    <div className="space-y-6">
      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-lg font-semibold text-foreground mb-3 flex items-center gap-2">
          <Monitor size={20} className="text-primary" />
          Linux
        </h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          Fully supported on all major Linux distributions.
        </p>
        <ul className="list-disc list-inside text-muted-foreground text-sm space-y-1">
          <li>Ubuntu 20.04+</li>
          <li>Debian 11+</li>
          <li>Fedora 35+</li>
          <li>Arch Linux</li>
          <li>CentOS/RHEL 8+</li>
        </ul>
      </div>

      <div className="bg-code border border-border rounded-lg p-6">
        <h3 className="font-mono text-lg font-semibold text-foreground mb-3 flex items-center gap-2">
          <Monitor size={20} className="text-primary" />
          macOS
        </h3>
        <p className="text-muted-foreground text-sm leading-relaxed mb-3">
          Native support for macOS systems.
        </p>
        <ul className="list-disc list-inside text-muted-foreground text-sm space-y-1">
          <li>macOS 11 (Big Sur) and later</li>
          <li>Intel and Apple Silicon (M1/M2/M3)</li>
          <li>Terminal.app, iTerm2, and other terminal emulators</li>
        </ul>
      </div>

      <div className="bg-code border border-border rounded-lg p-6 opacity-60">
        <h3 className="font-mono text-lg font-semibold text-foreground mb-3 flex items-center gap-2">
          <Monitor size={20} className="text-muted-foreground" />
          Windows
        </h3>
        <p className="text-muted-foreground text-sm leading-relaxed">
          Not currently supported. Consider using WSL2 (Windows Subsystem for Linux) to run Mako on Windows.
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Supported Shells</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">Bash</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Full support for Bash 4.0 and later. Most tested and recommended shell.
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">Zsh</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Full support for Zsh 5.0 and later. Compatible with Oh My Zsh and other frameworks.
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4 opacity-60">
        <code className="font-mono text-muted-foreground text-sm">Fish</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Coming soon. Fish shell support is planned for a future release.
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4 opacity-60">
        <code className="font-mono text-muted-foreground text-sm">PowerShell</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Coming soon. PowerShell support is planned for a future release.
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">System Requirements</h2>
    
    <div className="bg-code border border-border rounded-lg p-6">
      <ul className="space-y-3 text-muted-foreground text-sm">
        <li className="flex items-start gap-2">
          <span className="text-primary mt-1">•</span>
          <span><strong className="text-foreground">Architecture:</strong> x86_64 (amd64) or ARM64 (aarch64)</span>
        </li>
        <li className="flex items-start gap-2">
          <span className="text-primary mt-1">•</span>
          <span><strong className="text-foreground">RAM:</strong> Minimum 256MB free memory</span>
        </li>
        <li className="flex items-start gap-2">
          <span className="text-primary mt-1">•</span>
          <span><strong className="text-foreground">Disk:</strong> ~20MB for binaries, additional space for history database</span>
        </li>
        <li className="flex items-start gap-2">
          <span className="text-primary mt-1">•</span>
          <span><strong className="text-foreground">Network:</strong> Internet connection for AI features (API calls)</span>
        </li>
      </ul>
    </div>
  </motion.article>
);

const ReferenceSection = () => (
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
      Reference
    </motion.h1>
    
    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Architecture</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Mako uses a PTY (pseudo-terminal) architecture to intercept and process terminal I/O. Commands flow through the PTY master to your shell (bash/zsh), with a stream interceptor detecting markers for AI routing.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">File Locations</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">~/.mako/</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Main configuration directory for Mako
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">~/.mako/config.json</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Configuration file with API keys and settings
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">~/.mako/mako.db</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          SQLite database containing command history and embeddings
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">/usr/local/bin/mako</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Main Mako binary (shell orchestrator)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">/usr/local/bin/mako-menu</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Interactive menu binary for user choices
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Tech Stack</h2>
    
    <div className="overflow-x-auto">
      <table className="w-full text-sm border-collapse">
        <thead>
          <tr className="border-b border-border">
            <th className="text-left py-3 px-4 font-mono text-foreground">Component</th>
            <th className="text-left py-3 px-4 font-mono text-foreground">Technology</th>
          </tr>
        </thead>
        <tbody className="text-muted-foreground">
          <tr className="border-b border-border">
            <td className="py-3 px-4">Language</td>
            <td className="py-3 px-4 font-mono">Go 1.24+</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">PTY Handling</td>
            <td className="py-3 px-4 font-mono">creack/pty</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">Database</td>
            <td className="py-3 px-4 font-mono">SQLite with FTS5 (modernc.org/sqlite)</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">AI Providers</td>
            <td className="py-3 px-4 font-mono">OpenAI, Anthropic, Gemini, DeepSeek, OpenRouter, Ollama</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">Embeddings</td>
            <td className="py-3 px-4 font-mono">Provider-specific (Gemini, OpenAI, Ollama)</td>
          </tr>
          <tr className="border-b border-border">
            <td className="py-3 px-4">Build System</td>
            <td className="py-3 px-4 font-mono">Make</td>
          </tr>
        </tbody>
      </table>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Environment Variables</h2>
    
    <div className="space-y-4">
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">LLM_PROVIDER</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          AI provider to use (openai, anthropic, gemini, deepseek, openrouter, ollama)
        </p>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">LLM_MODEL</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Model name for your chosen provider
        </p>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">LLM_API_KEY</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Your API key for the chosen provider (not needed for Ollama)
        </p>
      </div>

      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">LLM_API_BASE</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Base URL for API calls (optional, for custom endpoints or Ollama)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">MAKO_CONFIG_DIR</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Override default config directory (default: ~/.mako)
        </p>
      </div>
      
      <div className="bg-code border border-border rounded-lg p-4">
        <code className="font-mono text-primary text-sm">MAKO_DEBUG</code>
        <p className="text-muted-foreground text-sm leading-relaxed mt-2">
          Enable debug logging (set to "1" or "true")
        </p>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Development</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      For developers interested in contributing or building from source:
    </p>
    <CodeBlock>{`# Clone repository
git clone https://github.com/fabiobrug/mako.git
cd mako/apps/cli

# Build
make build

# Run tests
make test

# Install locally
make install`}</CodeBlock>
  </motion.article>
);

const HelpSection = () => (
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
