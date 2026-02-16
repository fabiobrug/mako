"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { Copy, Check } from "lucide-react";

interface CodeBlockProps {
  children: string;
  language?: string;
}

export const CodeBlock = ({ children, language = "bash" }: CodeBlockProps) => {
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
      
      // Environment variable assignments (LLM_PROVIDER=, etc.)
      if (line.match(/^[A-Z_]+=/)) {
        const [varName, ...valueParts] = line.split('=');
        const value = valueParts.join('=');
        return (
          <div key={index}>
            <span className="text-[#8b5cf6]">{varName}</span>
            <span className="text-muted-foreground">=</span>
            <span className="text-[#f59e0b]">{value}</span>
          </div>
        );
      }
      
      // Lines with commands (curl, git, make, cd, cp, nano, ollama)
      const commandRegex = /^(curl|git|make|cd|cp|nano|ollama|sudo|brew|apt|export|source|\$)\s/;
      if (commandRegex.test(line.trim()) || line.startsWith('$') || line.startsWith('mako')) {
        return (
          <div key={index}>
            {highlightCommandLine(line)}
          </div>
        );
      }
      
      // Output lines
      return (
        <div key={index} className="text-muted-foreground/80">
          {line || '\u00A0'}
        </div>
      );
    });
  };

  const highlightCommandLine = (line: string) => {
    const tokens: JSX.Element[] = [];
    let currentIndex = 0;
    
    // Commands to highlight
    const commands = ['curl', 'git', 'make', 'cd', 'cp', 'nano', 'ollama', 'mako', 'sudo', 'brew', 'apt', 'export', 'source', 'bash', 'sh'];
    
    // Split by spaces but preserve quotes
    const parts = line.match(/(?:[^\s"']+|"[^"]*"|'[^']*')+/g) || [line];
    
    parts.forEach((part, i) => {
      const key = `${currentIndex}-${i}`;
      
      // Dollar sign
      if (part === '$') {
        tokens.push(<span key={key} className="text-success mr-2">{part}</span>);
      }
      // Commands
      else if (commands.includes(part) || part.startsWith('mako')) {
        tokens.push(<span key={key} className="text-[#3bace6] font-semibold mr-1">{part}</span>);
      }
      // Flags (- or --)
      else if (part.startsWith('--') || (part.startsWith('-') && part.length > 1 && part[1] !== ' ')) {
        tokens.push(<span key={key} className="text-[#8b5cf6] mr-1">{part}</span>);
      }
      // Strings (quoted)
      else if ((part.startsWith('"') && part.endsWith('"')) || (part.startsWith("'") && part.endsWith("'"))) {
        tokens.push(<span key={key} className="text-[#f59e0b] mr-1">{part}</span>);
      }
      // URLs
      else if (part.includes('http://') || part.includes('https://')) {
        tokens.push(<span key={key} className="text-[#3bace6] underline mr-1">{part}</span>);
      }
      // Pipes and redirects
      else if (part === '|' || part === '>' || part === '>>' || part === '<' || part === '\\') {
        tokens.push(<span key={key} className="text-success mx-1">{part}</span>);
      }
      // File paths (containing / or ~)
      else if (part.includes('/') || part.startsWith('~') || part.includes('.')) {
        tokens.push(<span key={key} className="text-[#78c8e6] mr-1">{part}</span>);
      }
      // Default
      else {
        tokens.push(<span key={key} className="text-foreground mr-1">{part}</span>);
      }
      
      currentIndex++;
    });
    
    return tokens;
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
