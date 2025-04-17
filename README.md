# ✅ Minimal Open Source AI Chat (LibreAIChat Community)

LibreAIChat Community is a self-hostable version of LibreAIChat -- an AI chat UI that runs open-source models via Ollama. It's designed to be simple, fast, and privacy-respecting -- no billing, no tracking, and no complex setup.

## 🚀 Features

- ✅ Clean Go Fiber + HTMX frontend
- ✅ GitHub OAuth login (optional)
- ✅ Docker Compose setup with:
  - PostgreSQL
  - Ollama
  - Caddy (reverse proxy)
- ✅ Automatically pulls models on boot
- ✅ No Stripe, no billing logic
- ✅ Fully local and lightweight

## 📦 Requirements

- Docker + Docker Compose
- GitHub OAuth app (optional)

## ⚙️ Quickstart

Clone the repo:

```bash
git clone https://github.com/LibreAIChat/community.git
cd community

Copy the example env:

cp .env.example .env

Edit .env with your GitHub OAuth keys or leave them blank to skip OAuth.
🐳 Start with Docker Compose

docker compose up --build

This will:

    Start PostgreSQL at localhost:15432

    Start Ollama and preload models via pull-models.sh

    Run the LibreAIChat app at http://localhost:3000

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

The pull-models.sh script will auto-pull your preferred models like phi, gemma, or mistral.

## 💾 Adding Models

To add a new model to your database, use the following SQL:

```sql
INSERT INTO models (name, identifier, description, category, required_tier, is_active, created_at, updated_at) 
VALUES ('Qwen', 'qwen2.5:0.5b', 'Quick & efficient responses', 'small', 'free', true, NOW(), NOW());


You can edit it to pull more models.
🌐 Caddy (Reverse Proxy)

Caddy is used for easy HTTPS and reverse proxying.

Edit Caddyfile to customize ports or enable TLS.
🤝 Contributing

PRs welcome! Want to add:

    New models?

    Local chat history?

    Usage stats?

Fork and go wild 🛠


📘 License

MIT. Use it freely for commercial or personal projects.