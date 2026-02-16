"use client";

import { motion } from "framer-motion";
import { Monitor } from "lucide-react";
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
