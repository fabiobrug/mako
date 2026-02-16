"use client";

import { motion } from "framer-motion";
import { fadeInUp } from "../../animations";
import { AskCommand } from "./AskCommand";
import { HistoryCommand } from "./HistoryCommand";
import { StatsCommand } from "./StatsCommand";
import { ExportCommand } from "./ExportCommand";
import { AliasCommand } from "./AliasCommand";
import { HealthCommand } from "./HealthCommand";
import { HelpCommand } from "./HelpCommand";
import { ConfigCommand } from "./ConfigCommand";
import { UpdateCommand } from "./UpdateCommand";
import { ImportCommand } from "./ImportCommand";
import { SyncCommand } from "./SyncCommand";
import { ClearCommand } from "./ClearCommand";
import { CompletionCommand } from "./CompletionCommand";
import { VersionCommand } from "./VersionCommand";

interface CommandsSectionProps {
  activeCommand: string | null;
}

export const CommandsSection = ({ activeCommand }: CommandsSectionProps) => {
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
  switch (activeCommand) {
    case "ask":
      return <AskCommand />;
    case "history":
      return <HistoryCommand />;
    case "alias":
      return <AliasCommand />;
    case "config":
      return <ConfigCommand />;
    case "update":
      return <UpdateCommand />;
    case "stats":
      return <StatsCommand />;
    case "health":
      return <HealthCommand />;
    case "export":
      return <ExportCommand />;
    case "import":
      return <ImportCommand />;
    case "sync":
      return <SyncCommand />;
    case "clear":
      return <ClearCommand />;
    case "completion":
      return <CompletionCommand />;
    case "help":
      return <HelpCommand />;
    case "version":
      return <VersionCommand />;
    default:
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
  }
};
