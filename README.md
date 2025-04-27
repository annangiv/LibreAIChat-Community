# ✅ Minimal, Privacy-First AI Chat UI (LibreAIChat Community)

**LibreAIChat Community** is a **self-hosted, open-source AI chat UI** for running **open models** like **Mistral**, **Phi**, and **LLaMA 3** via **Ollama**.  

Built with **Go Fiber** and **HTMX**, designed to be:
- **Simple** (easy Docker setup)
- **Fast** (lightweight backend)
- **Privacy-respecting** (no tracking, no billing)

🚫 No cloud lock-in. 🚫 No complex setup. 🚫 No nonsense.

---

## 🚀 Features

- ✅ **Clean UI** with **Go Fiber** + **HTMX** + **TailwindCSS**
- ✅ **Real-time streaming responses**
- ✅ **GitHub OAuth login** (optional)
- ✅ **Docker Compose** setup with:
  - PostgreSQL
  - Ollama
  - Caddy (reverse proxy with HTTPS)
- ✅ **Auto-pulls models** on boot (Mistral, Phi, etc.)
- ✅ **Fully local** – No Stripe, no billing, **your data stays yours**
- ✅ **Lightweight & fast** – No GPUs required (CPU-compatible models available)

---

## 📦 Requirements

- Docker + Docker Compose
- (Optional) GitHub OAuth app for login

---

## ⚙️ Quickstart

1. Clone the repo:

   ```bash
   git clone https://github.com/LibreAIChat/community.git
   cd community
