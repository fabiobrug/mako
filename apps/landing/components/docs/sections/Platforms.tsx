"use client";

import { motion } from "framer-motion";
import { CheckCircle2, XCircle, Clock } from "lucide-react";
import { fadeInUp } from "../animations";

export const PlatformsSection = () => (
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
    
    <p className="text-muted-foreground mb-10 leading-relaxed">
      Mako is designed to work seamlessly across different platforms and shells.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-12 mb-6">Operating Systems</h2>
    
    <div className="space-y-3">
      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-success">
        <div className="flex items-start gap-3">
          <CheckCircle2 size={20} className="text-success mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <h3 className="font-mono text-base font-semibold text-foreground mb-2">
              Linux
            </h3>
            <p className="text-muted-foreground text-sm leading-relaxed mb-2">
              Fully supported on all major distributions
            </p>
            <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
              <span className="px-2 py-1 bg-muted-foreground/10 rounded">Ubuntu 20.04+</span>
              <span className="px-2 py-1 bg-muted-foreground/10 rounded">Debian 11+</span>
              <span className="px-2 py-1 bg-muted-foreground/10 rounded">Fedora 35+</span>
              <span className="px-2 py-1 bg-muted-foreground/10 rounded">Arch Linux</span>
              <span className="px-2 py-1 bg-muted-foreground/10 rounded">CentOS/RHEL 8+</span>
            </div>
          </div>
        </div>
      </div>

      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-success">
        <div className="flex items-start gap-3">
          <CheckCircle2 size={20} className="text-success mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <h3 className="font-mono text-base font-semibold text-foreground mb-2">
              macOS
            </h3>
            <p className="text-muted-foreground text-sm leading-relaxed mb-2">
              Native support for Intel and Apple Silicon
            </p>
            <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
              <span className="px-2 py-1 bg-muted-foreground/10 rounded">macOS 11+</span>
              <span className="px-2 py-1 bg-muted-foreground/10 rounded">Intel & M1/M2/M3</span>
              <span className="px-2 py-1 bg-muted-foreground/10 rounded">Terminal.app</span>
              <span className="px-2 py-1 bg-muted-foreground/10 rounded">iTerm2</span>
            </div>
          </div>
        </div>
      </div>

      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-muted-foreground/30 opacity-60">
        <div className="flex items-start gap-3">
          <XCircle size={20} className="text-muted-foreground mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <h3 className="font-mono text-base font-semibold text-foreground mb-2">
              Windows
            </h3>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Not supported. Use WSL2 (Windows Subsystem for Linux) to run Mako on Windows.
            </p>
          </div>
        </div>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-12 mb-6">Shells</h2>
    
    <div className="space-y-3">
      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-success">
        <div className="flex items-start gap-3">
          <CheckCircle2 size={20} className="text-success mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-2">
              <code className="font-mono text-base text-foreground">bash</code>
              <span className="text-xs px-2 py-0.5 bg-success/20 text-success rounded">Recommended</span>
            </div>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Full support for Bash 4.0+. Most tested and stable shell integration.
            </p>
          </div>
        </div>
      </div>

      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-success">
        <div className="flex items-start gap-3">
          <CheckCircle2 size={20} className="text-success mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <code className="font-mono text-base text-foreground mb-2 block">zsh</code>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Full support for Zsh 5.0+. Compatible with Oh My Zsh and other frameworks.
            </p>
          </div>
        </div>
      </div>

      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-muted-foreground/30 opacity-60">
        <div className="flex items-start gap-3">
          <Clock size={20} className="text-muted-foreground mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-2">
              <code className="font-mono text-base text-muted-foreground">fish</code>
              <span className="text-xs px-2 py-0.5 bg-muted-foreground/20 text-muted-foreground rounded">Coming Soon</span>
            </div>
            <p className="text-muted-foreground text-sm leading-relaxed">
              Fish shell support planned for future release.
            </p>
          </div>
        </div>
      </div>

      <div className="group hover:bg-muted-foreground/5 transition-colors rounded-lg p-4 border-l-2 border-muted-foreground/30 opacity-60">
        <div className="flex items-start gap-3">
          <Clock size={20} className="text-muted-foreground mt-0.5 shrink-0" />
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-2">
              <code className="font-mono text-base text-muted-foreground">powershell</code>
              <span className="text-xs px-2 py-0.5 bg-muted-foreground/20 text-muted-foreground rounded">Coming Soon</span>
            </div>
            <p className="text-muted-foreground text-sm leading-relaxed">
              PowerShell support planned for future release.
            </p>
          </div>
        </div>
      </div>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-12 mb-6">System Requirements</h2>
    
    <div className="border-l-2 border-primary/30 pl-6 py-2">
      <dl className="space-y-4">
        <div>
          <dt className="font-mono text-sm font-semibold text-foreground mb-1">Architecture</dt>
          <dd className="text-muted-foreground text-sm">x86_64 (amd64) or ARM64 (aarch64)</dd>
        </div>
        <div>
          <dt className="font-mono text-sm font-semibold text-foreground mb-1">Memory</dt>
          <dd className="text-muted-foreground text-sm">Minimum 256MB free RAM</dd>
        </div>
        <div>
          <dt className="font-mono text-sm font-semibold text-foreground mb-1">Storage</dt>
          <dd className="text-muted-foreground text-sm">~20MB for binaries, additional space for history database</dd>
        </div>
        <div>
          <dt className="font-mono text-sm font-semibold text-foreground mb-1">Network</dt>
          <dd className="text-muted-foreground text-sm">Internet connection required for AI features (unless using Ollama locally)</dd>
        </div>
      </dl>
    </div>
  </motion.article>
);
