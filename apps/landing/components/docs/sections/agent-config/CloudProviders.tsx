"use client";

import { motion } from "framer-motion";
import { CodeBlock } from "../../CodeBlock";
import { fadeInUp } from "../../animations";

export const CloudProvidersConfigSection = () => (
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
      Cloud Providers
    </motion.h1>
    
    <p className="text-muted-foreground mb-6 leading-relaxed">
      Configure Mako with cloud-based AI providers for high-quality command generation. Each provider has unique strengths.
    </p>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">OpenAI Configuration</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Best for highest quality with GPT-4o models:
    </p>
    <CodeBlock>{`mako config set llm_provider openai`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Set your model:
    </p>
    <CodeBlock>{`mako config set llm_model gpt-4o-mini`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Add your API key:
    </p>
    <CodeBlock>{`mako config set api_key sk-your-api-key`}</CodeBlock>

    <div className="mt-4 border-l-2 border-primary/30 pl-4 py-2">
      <p className="text-muted-foreground text-sm leading-relaxed">
        <strong className="text-foreground">Get your API key:</strong>{" "}
        <a href="https://platform.openai.com/api-keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 underline">
          platform.openai.com/api-keys
        </a>
      </p>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Anthropic (Claude)</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Great reasoning with Claude models:
    </p>
    <CodeBlock>{`mako config set llm_provider anthropic`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Set Claude model:
    </p>
    <CodeBlock>{`mako config set llm_model claude-3-5-haiku-20241022`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Add your API key:
    </p>
    <CodeBlock>{`mako config set api_key sk-ant-your-key`}</CodeBlock>

    <div className="mt-4 border-l-2 border-primary/30 pl-4 py-2">
      <p className="text-muted-foreground text-sm leading-relaxed">
        <strong className="text-foreground">Get your API key:</strong>{" "}
        <a href="https://console.anthropic.com/account/keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 underline">
          console.anthropic.com/account/keys
        </a>
      </p>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Google Gemini</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Good balance with generous free tier (default):
    </p>
    <CodeBlock>{`mako config set llm_provider gemini`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Set Gemini model:
    </p>
    <CodeBlock>{`mako config set llm_model gemini-2.5-flash`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Add your API key:
    </p>
    <CodeBlock>{`mako config set api_key your-gemini-key`}</CodeBlock>

    <div className="mt-4 border-l-2 border-primary/30 pl-4 py-2">
      <p className="text-muted-foreground text-sm leading-relaxed">
        <strong className="text-foreground">Get your API key:</strong>{" "}
        <a href="https://ai.google.dev/" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 underline">
          ai.google.dev
        </a>
      </p>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">OpenRouter</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Access multiple models with one API key:
    </p>
    <CodeBlock>{`mako config set llm_provider openrouter`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Choose your model:
    </p>
    <CodeBlock>{`mako config set llm_model openai/gpt-4o-mini`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Add your API key:
    </p>
    <CodeBlock>{`mako config set api_key sk-or-your-key`}</CodeBlock>

    <div className="mt-4 border-l-2 border-primary/30 pl-4 py-2">
      <p className="text-muted-foreground text-sm leading-relaxed">
        <strong className="text-foreground">Get your API key:</strong>{" "}
        <a href="https://openrouter.ai/keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 underline">
          openrouter.ai/keys
        </a>
      </p>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">DeepSeek</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Cost-effective alternative:
    </p>
    <CodeBlock>{`mako config set llm_provider deepseek`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Set model:
    </p>
    <CodeBlock>{`mako config set llm_model deepseek-chat`}</CodeBlock>

    <p className="text-muted-foreground mb-4 mt-4 leading-relaxed">
      Add your API key:
    </p>
    <CodeBlock>{`mako config set api_key your-deepseek-key`}</CodeBlock>

    <div className="mt-4 border-l-2 border-primary/30 pl-4 py-2">
      <p className="text-muted-foreground text-sm leading-relaxed">
        <strong className="text-foreground">Get your API key:</strong>{" "}
        <a href="https://platform.deepseek.com/api_keys" target="_blank" rel="noopener noreferrer" className="text-primary hover:text-primary/80 underline">
          platform.deepseek.com/api_keys
        </a>
      </p>
    </div>

    <h2 className="font-mono text-2xl font-semibold text-foreground mt-8 mb-4">Environment Variables Method</h2>
    <p className="text-muted-foreground mb-4 leading-relaxed">
      Configure via .env file for any provider:
    </p>
    <CodeBlock>{`LLM_PROVIDER=openai
LLM_MODEL=gpt-4o-mini
LLM_API_KEY=sk-your-key`}</CodeBlock>
  </motion.article>
);
