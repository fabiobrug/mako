# Mako Landing Page

Marketing website and documentation for Mako AI Shell Orchestrator.

## Tech Stack

- **Framework**: Next.js 16 (App Router)
- **Styling**: Tailwind CSS
- **Language**: TypeScript
- **Build Tool**: Turbopack

## Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Start production server
npm start
```

Visit http://localhost:3000

## Project Structure

```
apps/landing/
├── app/              # Next.js app directory
│   ├── layout.tsx    # Root layout
│   ├── page.tsx      # Home page
│   └── ...
├── components/       # React components (add your components here)
├── public/          # Static assets
├── styles/          # Global styles
└── ...
```

## Deployment

This site can be deployed to:
- Vercel (recommended)
- Netlify
- Cloudflare Pages
- Any Node.js hosting platform

## Embedding Configuration Content

If you're creating documentation or content for the landing page about Mako's embedding features:

### What are Embeddings?

Embeddings are numerical representations of text that capture semantic meaning. They enable Mako's semantic search feature, allowing users to find commands by meaning rather than exact text.

**Key Points to Communicate:**
- **Semantic search**: Find commands by describing what you want, not just exact text matches
- **Example**: Searching "show running containers" finds `docker ps`, `docker container ls`, `kubectl get pods`, etc.
- **Automatic**: Embeddings are generated in the background as you use commands
- **Cached**: Generated once and stored in the local database

### Configuration Guide for Landing Page

When documenting embedding setup for users:

1. **Default Behavior**: Embeddings use the same provider as the LLM (no extra config needed)
2. **Supported Providers**:
   - Gemini: `gemini-embedding-001` (768-dimensional, state-of-the-art)
   - OpenAI: `text-embedding-3-small` (1536-dimensional)
   - Ollama: `nomic-embed-text` (local, free, private)

3. **Advanced: Separate Embedding Provider**:
   ```bash
   # Use local Ollama for free embeddings
   EMBEDDING_PROVIDER=ollama
   EMBEDDING_MODEL=nomic-embed-text
   EMBEDDING_API_BASE=http://localhost:11434
   ```

4. **Health Check**: Users can validate their configuration with:
   ```bash
   mako health          # Check embedding provider status
   mako config list     # View current configuration
   mako history semantic "test"  # Test semantic search
   ```

### Benefits to Highlight

- **Better search**: Find commands by meaning, not memorization
- **Privacy option**: Use local Ollama embeddings for zero cloud data
- **Cost optimization**: Use free local embeddings while keeping cloud LLM
- **Zero friction**: Works automatically with default provider

## Learn More

- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS](https://tailwindcss.com)
- [Main Mako Repository](../../README.md)
- [CLI README (embedding details)](../cli/README.md)
