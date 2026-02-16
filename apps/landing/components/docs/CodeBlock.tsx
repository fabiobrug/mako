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
