# error-explain 🗣️
**Give your compiler a voice.**

error-explain is a zero-latency "sidecar" for your compiler. It sits alongside `g++`, `clang`, or `cargo` and explains build failures in plain English.

Unlike AI coding agents that try to take over your workflow, error-explain respects it. It streams your compiler's raw output instantly and only chimes in when things go wrong.

## 🚀 Features

- ⚡ **Zero Latency:** Your build runs at native speed. The AI analysis happens in a background thread and never blocks your terminal.
- 🧠 **Context-Aware:** Reads the source code around the error (the "Context Miner") so the AI sees exactly what you wrote, not just the error message.
- 🛡️ **Skeptic-First Design:** By default, you see the raw compiler error first. The AI adds a "second opinion" below it.
- 💸 **Free-Tier Friendly:** Built-in "Failover Chain" supports Groq, Gemini, Anthropic and OpenAI. If one API is rate-limited, it automatically switches to the next. Groq and Gemini currently offer free tier whereas Anthropic and OpenAI require payment.

## 🎓 Multi-Mode:

- **Direct:** "Fix the semicolon." (1 sentence)
- **Deep:** "Here is how the vtable was corrupted." (Technical deep dive)
- **Teacher:** "This is called SFINAE. Here is why it exists." (Educational)

## 📦 Installation

### From Source (Go)

```bash
go install github.com/ultrakapy/error-explain/cmd/error-explain@latest
```

### Setup API Keys

error-explain is model-agnostic. You need at least one API key. We recommend Groq for speed or Gemini for the best free tier.

```bash
# Fastest (Recommended)
export GROQ_API_KEY="gsk_..."

# Best Reasoning (Free Tier)
export GEMINI_API_KEY="AIza..."

# Fallback
export OPENAI_API_KEY="sk-..."
```

## 🛠 Usage

Simply verify your build command with `error-explain --`.

### Basic Usage:

```bash
error-explain -- g++ main.cpp
```

### Select a Voice Mode:

```bash
# Just the fix (Default)
error-explain --mode direct -- make build

# Learn the concept
error-explain --mode teacher -- cargo build

# Debug complex templates
error-explain --mode deep -- g++ -std=c++20 complex_templates.cpp
```

## ⚙️ Configuration

Error-Explain works out-of-the-box with zero configuration, using a built-in fallback chain (Groq -> Gemini -> Anthropic -> OpenAI).

However, you can fully customize the provider chain, use your own models (including local LLMs like Ollama), and manage API keys via a configuration file.

### Config File Location
The tool looks for a file named `config` (e.g., `config.yaml`, `config.json`, or `config.toml`) in your operating system's standard configuration directory:

| OS | Config Path |
| :--- | :--- |
| **Linux** | `~/.config/error-explain/config.yaml` |
| **macOS** | `~/Library/Application Support/error-explain/config.yaml` |
| **Windows** | `%AppData%\error-explain\config.yaml` |

### Configuration Structure
You can define a list of providers. The tool will try them **in order**. If one fails (rate limit, network error), it automatically moves to the next.

**Example `config.yaml`:**

```yaml
providers:
  # 1. First priority: Local Ollama (Free, Private)
  - name: "Local Llama3"
    type: "openai"
    model: "llama3"
    base_url: "http://localhost:11434/v1"
    api_key_env: "OLLAMA_API_KEY" # Optional for Ollama

  # 2. Second priority: Claude 3.5 Sonnet (Best Reasoning)
  - name: "Claude 3.5"
    type: "anthropic"
    model: "claude-3-5-sonnet-20240620"
    api_key_env: "ANTHROPIC_API_KEY"

  # 3. Fallback: Groq (Ultra Fast)
  - name: "Groq"
    type: "openai"
    model: "llama-3.3-70b-versatile"
    base_url: "https://api.groq.com/openai/v1"
    api_key_env: "GROQ_API_KEY"
```

### Supported Provider Types

| Type | Description | Required Fields |
| :--- | :--- | :--- |
| `openai` | Any OpenAI-compatible API (OpenAI, Groq, DeepSeek, Ollama, vLLM) | `model`, `api_key_env`, `base_url` |
| `anthropic` | Anthropic's Claude models | `model`, `api_key_env` |
| `gemini` | Google's Gemini models | `model`, `api_key_env` |

## 🔐 Privacy Notice

`error-explain` transmits **compiler output**, **nearby source code**, and **extra context** you provide to external LLM services to generate explanations.

- **⚠ Sensitive Data:** This transmission can unintentionally expose proprietary code or credentials.
- **🔒 Local Mode:** For private repositories, we strongly recommend using a **Local LLM** (like Ollama) via the configuration file to keep data offline.
- **🛡 Responsibility:** Ensure secrets (API keys, private keys) are redacted before running. By using hosted providers, you accept their data-handling policies. Consult your organization’s security guidelines if working with confidential code.

## 🏗 Architecture

error-explain is built in Go to be a single, lightweight binary.

- **The Transparent Pipe:** Uses `io.MultiWriter` to stream stderr to your screen and an internal memory buffer simultaneously.
- **The Context Miner:** If an error is detected, the miner parses the file:line location and reads a ±5 line window from your source code.
- **The Failover Brain:** It attempts to call the fastest AI provider first. If it hits a rate limit (429), it silently switches to the next provider in the chain.

## 🤝 Contributing

This tool is designed to be Compiler Agnostic. Currently, the "Context Miner" is optimized for GCC/Clang formats.

We welcome PRs for:

- Rust (rustc) support.
- Java (javac) support.
- Python traceback parsing.

## 📜 License

MIT License. Built for the community.
