"use client";

import { useState } from "react";
import { DocsLayout } from "@/components/docs/DocsLayout";
import { GetStartedSection } from "@/components/docs/sections/GetStarted";
import { OneLineInstallSection, FromSourceSection, PostInstallSection } from "@/components/docs/sections/install";
import { AgentConfigOverviewSection, OllamaConfigSection, CloudProvidersConfigSection, AdvancedConfigSection } from "@/components/docs/sections/agent-config";
import { CommandsSection } from "@/components/docs/sections/commands";
import { PlatformsSection } from "@/components/docs/sections/Platforms";
import { ReferenceSection } from "@/components/docs/sections/Reference";
import { HelpSection } from "@/components/docs/sections/Help";

export default function Docs() {
  const [activeSection, setActiveSection] = useState("get-started");
  const [sidebarOpen, setSidebarOpen] = useState(false);
  
  // Install section state
  const [installExpanded, setInstallExpanded] = useState(false);
  const [activeInstall, setActiveInstall] = useState<string | null>(null);
  
  // Agent Config section state
  const [agentConfigExpanded, setAgentConfigExpanded] = useState(false);
  const [activeAgentConfig, setActiveAgentConfig] = useState<string | null>(null);
  
  // Commands section state
  const [commandsExpanded, setCommandsExpanded] = useState(false);
  const [activeCommand, setActiveCommand] = useState<string | null>(null);

  return (
    <DocsLayout
      activeSection={activeSection}
      setActiveSection={setActiveSection}
      sidebarOpen={sidebarOpen}
      setSidebarOpen={setSidebarOpen}
      installExpanded={installExpanded}
      setInstallExpanded={setInstallExpanded}
      activeInstall={activeInstall}
      setActiveInstall={setActiveInstall}
      agentConfigExpanded={agentConfigExpanded}
      setAgentConfigExpanded={setAgentConfigExpanded}
      activeAgentConfig={activeAgentConfig}
      setActiveAgentConfig={setActiveAgentConfig}
      commandsExpanded={commandsExpanded}
      setCommandsExpanded={setCommandsExpanded}
      activeCommand={activeCommand}
      setActiveCommand={setActiveCommand}
    >
      {activeSection === "get-started" && <GetStartedSection />}
      {activeSection === "install" && activeInstall === "quick-install" && <OneLineInstallSection />}
      {activeSection === "install" && activeInstall === "from-source" && <FromSourceSection />}
      {activeSection === "install" && activeInstall === "post-install" && <PostInstallSection />}
      {activeSection === "agent-config" && activeAgentConfig === "overview" && <AgentConfigOverviewSection />}
      {activeSection === "agent-config" && activeAgentConfig === "ollama" && <OllamaConfigSection />}
      {activeSection === "agent-config" && activeAgentConfig === "cloud-providers" && <CloudProvidersConfigSection />}
      {activeSection === "agent-config" && activeAgentConfig === "advanced" && <AdvancedConfigSection />}
      {activeSection === "commands" && <CommandsSection activeCommand={activeCommand} />}
      {activeSection === "platforms" && <PlatformsSection />}
      {activeSection === "reference" && <ReferenceSection />}
      {activeSection === "help" && <HelpSection />}
    </DocsLayout>
  );
}
