"use client";

import { ReactNode } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import { Github, Menu, X, ChevronDown } from "lucide-react";
import { sections, commandsList } from "./sections-config";

interface DocsLayoutProps {
  children: ReactNode;
  activeSection: string;
  setActiveSection: (section: string) => void;
  sidebarOpen: boolean;
  setSidebarOpen: (open: boolean) => void;
  commandsExpanded: boolean;
  setCommandsExpanded: (expanded: boolean) => void;
  activeCommand: string | null;
  setActiveCommand: (command: string | null) => void;
}

export const DocsLayout = ({
  children,
  activeSection,
  setActiveSection,
  sidebarOpen,
  setSidebarOpen,
  commandsExpanded,
  setCommandsExpanded,
  activeCommand,
  setActiveCommand,
}: DocsLayoutProps) => {
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
                        {commandsList.map((cmd) => (
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
          {children}
        </motion.main>
      </div>
    </div>
  );
};
