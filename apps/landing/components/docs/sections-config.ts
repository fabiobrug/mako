import { Rocket, Terminal, Settings, Wrench, Monitor, BookOpen, HelpCircle } from "lucide-react";

export const sections = [
  { id: "get-started", label: "Get Started", icon: Rocket },
  { id: "install", label: "Install", icon: Terminal },
  { id: "agent-config", label: "Agent Config", icon: Settings },
  { id: "commands", label: "Commands", icon: Wrench },
  { id: "platforms", label: "Platforms", icon: Monitor },
  { id: "reference", label: "Reference", icon: BookOpen },
  { id: "help", label: "Help", icon: HelpCircle },
];

export const commandsList = [
  { id: "ask", label: "mako ask" },
  { id: "history", label: "mako history" },
  { id: "alias", label: "mako alias" },
  { id: "config", label: "mako config" },
  { id: "update", label: "mako update" },
  { id: "stats", label: "mako stats" },
  { id: "health", label: "mako health" },
  { id: "export", label: "mako export" },
  { id: "import", label: "mako import" },
  { id: "sync", label: "mako sync" },
  { id: "clear", label: "mako clear" },
  { id: "completion", label: "mako completion" },
  { id: "help", label: "mako help" },
  { id: "version", label: "mako version" },
];
