"use client";

import { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { 
  Github, Search, ShieldCheck, Copy, Check, 
  Zap, Database, Brain,
  Layers, BookOpen, Users
} from "lucide-react";

export default function Home() {
  const [copied, setCopied] = useState(false);
  const installCmd = "curl -sSL https://raw.githubusercontent.com/fabiobrug/mako/dev/scripts/install.sh | bash";

  const handleCopy = () => {
    navigator.clipboard.writeText(installCmd);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="min-h-screen flex flex-col">
      {/* Underwater Bubble Background - Home page only */}
      <svg className="bubble-bg" xmlns="http://www.w3.org/2000/svg">
        <defs>
          {/* Dark blue - deep ocean */}
          <radialGradient id="bubble-dark">
            <stop offset="0%" stopColor="rgba(30, 100, 150, 0.15)" />
            <stop offset="50%" stopColor="rgba(30, 100, 150, 0.08)" />
            <stop offset="100%" stopColor="transparent" />
          </radialGradient>
          
          {/* Medium blue - current primary color */}
          <radialGradient id="bubble-medium">
            <stop offset="0%" stopColor="rgba(59, 172, 214, 0.12)" />
            <stop offset="50%" stopColor="rgba(59, 172, 214, 0.06)" />
            <stop offset="100%" stopColor="transparent" />
          </radialGradient>
          
          {/* Light blue - closer to surface */}
          <radialGradient id="bubble-light">
            <stop offset="0%" stopColor="rgba(120, 200, 230, 0.10)" />
            <stop offset="50%" stopColor="rgba(120, 200, 230, 0.05)" />
            <stop offset="100%" stopColor="transparent" />
          </radialGradient>
        </defs>
        
        {/* Background rectangle */}
        <rect width="100%" height="100%" fill="hsl(220, 30%, 5%)" />
        
        {/* Large bubbles - mixed depths */}
        <circle cx="40%" cy="60%" r="90" fill="url(#bubble-medium)">
          <animate attributeName="cy" from="110%" to="-10%" dur="25s" begin="-12s" repeatCount="indefinite"/>
        </circle>
        
        <circle cx="85%" cy="30%" r="80" fill="url(#bubble-dark)">
          <animate attributeName="cy" from="110%" to="-10%" dur="22s" begin="-8s" repeatCount="indefinite"/>
        </circle>

        <circle cx="30%" cy="50%" r="75" fill="url(#bubble-light)">
          <animate attributeName="cy" from="110%" to="-10%" dur="23s" begin="-15s" repeatCount="indefinite"/>
        </circle>

        {/* Medium bubbles - mixed depths */}
        <circle cx="15%" cy="20%" r="50" fill="url(#bubble-light)">
          <animate attributeName="cy" from="110%" to="-10%" dur="20s" begin="-10s" repeatCount="indefinite"/>
        </circle>
        
        <circle cx="70%" cy="80%" r="45" fill="url(#bubble-dark)">
          <animate attributeName="cy" from="110%" to="-10%" dur="18s" begin="-5s" repeatCount="indefinite"/>
        </circle>

        <circle cx="60%" cy="90%" r="40" fill="url(#bubble-medium)">
          <animate attributeName="cy" from="110%" to="-10%" dur="19s" begin="-13s" repeatCount="indefinite"/>
        </circle>

        {/* Small bubbles - mixed depths */}
        <circle cx="20%" cy="40%" r="25" fill="url(#bubble-dark)">
          <animate attributeName="cy" from="110%" to="-10%" dur="16s" begin="-6s" repeatCount="indefinite"/>
        </circle>

        <circle cx="50%" cy="70%" r="30" fill="url(#bubble-medium)">
          <animate attributeName="cy" from="110%" to="-10%" dur="17s" begin="-9s" repeatCount="indefinite"/>
        </circle>

        <circle cx="75%" cy="50%" r="20" fill="url(#bubble-light)">
          <animate attributeName="cy" from="110%" to="-10%" dur="15s" begin="-11s" repeatCount="indefinite"/>
        </circle>

        <circle cx="90%" cy="85%" r="28" fill="url(#bubble-dark)">
          <animate attributeName="cy" from="110%" to="-10%" dur="21s" begin="-18s" repeatCount="indefinite"/>
        </circle>

        <circle cx="8%" cy="65%" r="22" fill="url(#bubble-light)">
          <animate attributeName="cy" from="110%" to="-10%" dur="14s" begin="-7s" repeatCount="indefinite"/>
        </circle>

        <circle cx="95%" cy="45%" r="18" fill="url(#bubble-medium)">
          <animate attributeName="cy" from="110%" to="-10%" dur="13s" begin="-4s" repeatCount="indefinite"/>
        </circle>
      </svg>

      {/* Header */}
      <header className="border-b border-border sticky top-0 z-50 bg-background/95 backdrop-blur">
        <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
          <span className="font-mono font-bold text-lg tracking-tight text-foreground hover:text-primary transition-colors">
            Mako
          </span>
          <div className="flex items-center gap-6">
            <Link href="/docs" className="text-primary hover:text-primary/80 transition-colors text-sm font-medium">
              Documentation
            </Link>
            <a
              href="https://github.com/fabiobrug/mako"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary hover:text-primary/80 transition-colors"
              aria-label="GitHub"
            >
              <Github size={20} />
            </a>
          </div>
        </div>
      </header>

      <main className="flex-1">
        {/* Hero Section */}
        <section className="max-w-5xl mx-auto px-6 pt-32 pb-20 text-center">
          <h1 className="font-mono font-bold text-5xl md:text-7xl text-foreground mb-8">
            <span className="text-primary">Mako</span> AI Shell Orchestrator
          </h1>

          {/* Install Command */}
          <div className="max-w-3xl mx-auto mb-16">
            <div className="rounded-lg overflow-hidden shadow-lg">
              {/* Mac terminal header */}
              <div className="bg-[#2d2d2d] px-4 py-2 flex items-center gap-2 border-b border-[#1a1a1a]">
                <div className="flex gap-1.5">
                  <div className="w-3 h-3 rounded-full bg-[#ff5f57]"></div>
                  <div className="w-3 h-3 rounded-full bg-[#febc2e]"></div>
                  <div className="w-3 h-3 rounded-full bg-[#28c840]"></div>
                </div>
              </div>
              {/* Terminal content */}
              <div className="bg-[#1e1e1e] px-4 py-4 flex items-center justify-between gap-4">
                <code className="font-mono text-sm text-foreground flex-1 text-left">
                  <span className="text-success">$ </span>{installCmd}
                </code>
                <button
                  onClick={handleCopy}
                  className="shrink-0 p-2 rounded hover:bg-[#2d2d2d] transition-colors"
                  aria-label="Copy install command"
                >
                  {copied ? (
                    <Check size={16} className="text-success" />
                  ) : (
                    <Copy size={16} className="text-muted-foreground" />
                  )}
                </button>
              </div>
            </div>
            <p className="text-sm text-muted-foreground mt-3">
              Supports Linux and macOS
            </p>
          </div>

          {/* What It Does - 3x2 Grid */}
          <div className="mb-20">
            <h2 className="font-mono font-semibold text-2xl text-foreground mb-8">What It Does</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-3xl mx-auto">
              <FeatureBox
                icon={<Brain size={20} />}
                title="AI Command Generation"
                description="Convert natural language to shell commands instantly"
              />
              <FeatureBox
                icon={<Search size={20} />}
                title="Semantic History"
                description="Search commands by meaning, not just text"
              />
              <FeatureBox
                icon={<Zap size={20} />}
                title="Fast Performance"
                description="Sub-10ms command saves with async processing"
              />
              <FeatureBox
                icon={<Layers size={20} />}
                title="PTY Integration"
                description="Seamless terminal wrapping without disruption"
              />
              <FeatureBox
                icon={<ShieldCheck size={20} />}
                title="Safety Validation"
                description="Detects dangerous commands before execution"
              />
              <FeatureBox
                icon={<Database size={20} />}
                title="Local-First"
                description="All data stored locally with SQLite + FTS5"
              />
            </div>
          </div>

          {/* Supported LLMs Section */}
          <div className="mb-20">
            <h2 className="font-mono font-semibold text-2xl text-foreground mb-8">Supported LLMs</h2>
            <div className="flex flex-wrap justify-center gap-3 max-w-3xl mx-auto">
              <LLMPill name="Ollama" icon="/ollama-icon.svg" />
              <LLMPill name="OpenAI" icon="/openai-icon.svg" />
              <LLMPill name="Anthropic" icon="/claude-ai-icon.svg" />
              <LLMPill name="Gemini" icon="/google-gemini-icon.svg" />
              <LLMPill name="OpenRouter" icon="/openrouter-icon.svg" />
              <LLMPill name="DeepSeek" icon="/deepseek-logo-icon.svg" />
            </div>
          </div>

          {/* Action Boxes */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-3xl mx-auto">
            <ActionBox
              icon={<Github size={20} />}
              title="GitHub"
              description="View source code and contribute"
              href="https://github.com/fabiobrug/mako"
              color="white"
            />
            <ActionBox
              icon={<BookOpen size={20} />}
              title="Documentation"
              description="Learn how to use Mako"
              href="/docs"
              color="orange"
            />
            <ActionBox
              icon={<Users size={20} />}
              title="Contribute"
              description="Join the community and help"
              href="https://github.com/fabiobrug/mako/blob/dev/docs/CONTRIBUTING.md"
              color="purple"
            />
          </div>
        </section>

      </main>

      {/* Footer */}
      <footer className="border-t border-border">
        <div className="max-w-5xl mx-auto px-6 py-8">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <div className="flex items-center gap-6">
              <Link href="/docs" className="text-primary hover:text-primary/80 transition-colors text-sm font-medium">
                Documentation
              </Link>
              <a 
                href="https://github.com/fabiobrug/mako" 
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:text-primary/80 transition-colors text-sm font-medium"
              >
                GitHub
              </a>
              <a 
                href="https://github.com/fabiobrug/mako/blob/dev/LICENSE" 
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:text-primary/80 transition-colors text-sm font-medium"
              >
                License
              </a>
            </div>
            <div className="text-sm text-muted-foreground">
              Built by <a href="https://github.com/fabiobrug" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80">Fabio Brug</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}

const FeatureBox = ({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) => (
  <div className="bg-code rounded-lg border border-border p-4 hover:scale-[1.02] transition-all duration-300 ease-out cursor-default hover:shadow-lg hover:shadow-primary/10">
    <div className="flex justify-center text-primary mb-3">{icon}</div>
    <h3 className="font-mono font-semibold text-foreground mb-1.5 text-sm text-left">{title}</h3>
    <p className="text-xs text-muted-foreground leading-relaxed text-left">{description}</p>
  </div>
);

const ActionBox = ({ icon, title, description, href, color }: { icon: React.ReactNode; title: string; description: string; href: string; color: 'white' | 'orange' | 'purple' }) => {
  const isExternal = href.startsWith('http');
  const Component = isExternal ? 'a' : Link;
  const extraProps = isExternal ? { target: "_blank", rel: "noopener noreferrer" } : {};

  const colorClasses = {
    white: 'hover:border-white hover:shadow-white/20',
    orange: 'hover:border-[#f59e0b] hover:shadow-[#f59e0b]/20',
    purple: 'hover:border-[#8b5cf6] hover:shadow-[#8b5cf6]/20'
  };

  const iconColorClasses = {
    white: 'group-hover:text-white',
    orange: 'group-hover:text-[#f59e0b]',
    purple: 'group-hover:text-[#8b5cf6]'
  };

  return (
    <Component
      href={href}
      {...extraProps}
      className={`bg-code rounded-lg border border-border p-4 hover:scale-[1.02] transition-all duration-300 ease-out group cursor-pointer block hover:shadow-xl ${colorClasses[color]}`}
    >
      <div className={`flex justify-center text-primary mb-3 transition-colors duration-300 ${iconColorClasses[color]}`}>{icon}</div>
      <h3 className="font-mono font-semibold text-foreground mb-1.5 text-sm text-left">{title}</h3>
      <p className="text-xs text-muted-foreground leading-relaxed text-left">{description}</p>
    </Component>
  );
};

const LLMPill = ({ name, icon }: { name: string; icon: string }) => (
  <div className="bg-code rounded-full border border-border px-5 py-2.5 hover:scale-[1.02] transition-all duration-300 ease-out cursor-default hover:shadow-lg hover:shadow-primary/10 flex items-center gap-2.5">
    <Image src={icon} alt={`${name} icon`} width={20} height={20} className="shrink-0" />
    <span className="font-mono font-semibold text-foreground text-sm">{name}</span>
  </div>
);
