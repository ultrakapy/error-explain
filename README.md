# error-explain 🗣️
**Give your compiler a voice.**

error-explain is a zero-latency "sidecar" for your compiler. It sits alongside `g++`, `clang`, or `cargo` and explains build failures in plain English.

Unlike AI coding agents that try to take over your workflow, error-explain respects it. It streams your compiler's raw output instantly and only chimes in when things go wrong.

🚀 Features

- ⚡ **Zero Latency:** Your build runs at native speed. The AI analysis happens in a background thread and never blocks your terminal.
- 🧠 **Context-Aware:** Reads the source code around the error (the "Context Miner") so the AI sees exactly what you wrote, not just the error message.
- 🛡️ **Skeptic-First Design:** By default, you see the raw compiler error first. The AI adds a "second opinion" below it.
- 💸 **Free-Tier Friendly:** Built-in "Failover Chain" supports Groq, Gemini Flash, and OpenAI. If one API is rate-limited, it automatically switches to the next.

🎓 Multi-Mode:

- **Direct:** "Fix the semicolon." (1 sentence)
- **Deep:** "Here is how the vtable was corrupted." (Technical deep dive)
- **Teacher:** "This is called SFINAE. Here is why it exists." (Educational)

📦 Installation

**From Source (Go)**

```bash
go install github.com/yourusername/error-explain/cmd/error-explain@latest
```

**Setup API Keys**

error-explain is model-agnostic. You need at least one API key. We recommend Groq for speed or Gemini for the best free tier.

```bash
# Fastest (Recommended)
export GROQ_API_KEY="gsk_..."

# Best Reasoning (Free Tier)
export GEMINI_API_KEY="AIza..."

# Fallback
export OPENAI_API_KEY="sk-..."
```

🛠 Usage

Simply verify your build command with `error-explain --`.

**Basic Usage:**

```bash
error-explain -- g++ main.cpp
```

**Select a Voice Mode:**

```bash
# Just the fix (Default)
error-explain --mode direct -- make build

# Learn the concept
error-explain --mode teacher -- cargo build

# Debug complex templates
error-explain --mode deep -- g++ -std=c++20 complex_templates.cpp
```

⚙️ Configuration

You can configure defaults in `~/.config/error-explain/config.toml` (optional):

```toml
[voice]
default_mode = "direct"  # direct, deep, teacher
timeout = "10s"

[ai]
# The tool tries these in order until one succeeds
chain = ["groq", "gemini", "openai"]
```

🏗 Architecture

error-explain is built in Go to be a single, lightweight binary.

- **The Transparent Pipe:** Uses `io.MultiWriter` to stream stderr to your screen and an internal memory buffer simultaneously.
- **The Context Miner:** If an error is detected, the miner parses the file:line location and reads a ±5 line window from your source code.
- **The Failover Brain:** It attempts to call the fastest AI provider first. If it hits a rate limit (429), it silently switches to the next provider in the chain.

🤝 Contributing

This tool is designed to be Compiler Agnostic. Currently, the "Context Miner" is optimized for GCC/Clang formats.

We welcome PRs for:

- Rust (rustc) support.
- Java (javac) support.
- Python traceback parsing.

📜 License

MIT License. Built for the community.
