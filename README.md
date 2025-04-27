✅ Minimal, Privacy-First AI Chat UI (LibreAIChat Community)

LibreAIChat Community is a self-hosted, open-source AI chat UI for running open models like Mistral, Phi, and LLaMA 3 via Ollama.

Built with Go Fiber and HTMX, designed to be:

Simple (easy Docker setup)
Fast (lightweight backend)
Privacy-respecting (no tracking, no billing)
🚫 No cloud lock-in. 🚫 No complex setup. 🚫 No nonsense.

🚀 Features

✅ Clean UI with Go Fiber + HTMX + TailwindCSS
✅ Real-time streaming responses
✅ GitHub OAuth login (optional)
✅ Docker Compose setup with:
PostgreSQL
Ollama
Caddy (reverse proxy with HTTPS)
✅ Auto-pulls models on boot (Mistral, Phi, etc.)
✅ Fully local – No Stripe, no billing, your data stays yours
✅ Lightweight & fast – No GPUs required (CPU-compatible models available)
📦 Requirements

Docker + Docker Compose
(Optional) GitHub OAuth app for login
⚙️ Quickstart

Clone the repo:
git clone https://github.com/LibreAIChat/community.git
cd community
Copy the example env:
cp .env.example .env
Edit .env with your GitHub OAuth keys (or leave them blank to skip OAuth).
🐳 Start with Docker Compose:
docker compose up --build
This will:

Start PostgreSQL at localhost:15432
Start Ollama and preload models via pull-models.sh
Run LibreAIChat at http://localhost:3000
Serve through Caddy at http://localhost
🛠 Example .env

MAX_PROCS=4
SESSION_SECRET=your_random_secret
DATABASE_URL=postgres://postgres:password@postgres:5432/LibreAIChat_prod?sslmode=disable
BASE_URL=http://localhost:3000

# Optional GitHub OAuth
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
🔄 Auto-Pull Models

The pull-models.sh script auto-pulls your preferred models (like phi, mistral, gemma).

Add or edit models directly in the script to customize what loads.

💾 Adding Models to the Database

Insert new models into PostgreSQL:

INSERT INTO models (name, identifier, description, category, required_tier, is_active, created_at, updated_at) 
VALUES ('Qwen', 'qwen2.5:0.5b', 'Quick & efficient responses', 'small', 'free', true, NOW(), NOW());
Customize identifiers and descriptions as needed!

🌐 Caddy (Reverse Proxy & HTTPS)

Caddy handles reverse proxy and HTTPS (if needed).
Edit the Caddyfile to:

Change ports
Enable TLS
Customize domains
🤝 Contributing

PRs welcome! Want to:

Add new models?
Build local chat history?
Add usage stats?
Fork it and have fun 🛠.

📘 License

MIT License – Free to use, modify, and deploy for commercial or personal projects.