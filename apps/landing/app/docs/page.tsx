"use client";

import { useState } from "react";
import { DocsLayout } from "@/components/docs/DocsLayout";
import { GetStartedSection } from "@/components/docs/sections/GetStarted";
import { InstallSection } from "@/components/docs/sections/Install";
import { AgentConfigSection } from "@/components/docs/sections/AgentConfig";
import { CommandsSection } from "@/components/docs/sections/commands";
import { PlatformsSection } from "@/components/docs/sections/Platforms";
import { ReferenceSection } from "@/components/docs/sections/Reference";
import { HelpSection } from "@/components/docs/sections/Help";

export default function Docs() {
  const [activeSection, setActiveSection] = useState("get-started");
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [commandsExpanded, setCommandsExpanded] = useState(false);
  const [activeCommand, setActiveCommand] = useState<string | null>(null);

  return (
    <DocsLayout
      activeSection={activeSection}
      setActiveSection={setActiveSection}
      sidebarOpen={sidebarOpen}
      setSidebarOpen={setSidebarOpen}
      commandsExpanded={commandsExpanded}
      setCommandsExpanded={setCommandsExpanded}
      activeCommand={activeCommand}
      setActiveCommand={setActiveCommand}
    >
      {activeSection === "get-started" && <GetStartedSection />}
      {activeSection === "install" && <InstallSection />}
      {activeSection === "agent-config" && <AgentConfigSection />}
      {activeSection === "commands" && <CommandsSection activeCommand={activeCommand} />}
      {activeSection === "platforms" && <PlatformsSection />}
      {activeSection === "reference" && <ReferenceSection />}
      {activeSection === "help" && <HelpSection />}
    </DocsLayout>
  );
}
