# oc-go-cc

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/github/go-mod/go-version/newmancai/oc-go-cc-optim)](https://github.com/newmancai/oc-go-cc-optim)

A Go CLI proxy that lets you use your [OpenCode Go](https://opencode.ai/docs/go/) subscription with [Claude Code](https://docs.anthropic.com/en/docs/claude-code).

`oc-go-cc` sits between Claude Code and OpenCode Go, intercepting Anthropic API requests, transforming them to OpenAI format, and forwarding them to OpenCode Go's endpoint. Claude Code thinks it's talking to Anthropic — but your requests go to affordable open models instead.

[English](#english) | [中文](#chinese)

---

## English

### Why?

OpenCode Go gives you access to powerful open coding models for **$5/month** (then $10/month). This proxy makes those models work seamlessly with Claude Code's interface — no patches, no forks, just set two environment variables and go.

### Features

- **Transparent Proxy** — Claude Code sends Anthropic-format requests, proxy transforms to OpenAI format and back
- **Model Routing** — Automatically routes to different models based on context (default, thinking, long context, background)
- **Manual Model Selection** — Override routing with Claude Code's `/model` command — use any scenario key from `config.json` (e.g., `/model think`, `/model complex`)
- **Fallback Chains** — If a model fails, automatically tries the next one in your configured chain
- **Circuit Breaker** — Tracks model health and skips failing models to avoid latency spikes
- **Real-time Streaming** — Full SSE streaming with live OpenAI -> Anthropic format transformation
- **Tool Calling** — Proper Anthropic tool_use/tool_result <-> OpenAI function calling translation
- **Token Counting** — Uses tiktoken (cl100k_base) for accurate token counting and context threshold detection
- **JSON Configuration** — Flexible config file with environment variable overrides and `${VAR}` interpolation
- **Hot Reload** — Watch config file for changes and reload automatically (off by default)
- **Background Mode** — Run as daemon detached from terminal
- **Auto-start on Login** — Launch on system startup via launchd (macOS)

### Quick Start

#### 1. Install

```bash
# macOS / Linux
# Build from source (see INSTALLATION.md for Homebrew/Scoop)
git clone https://github.com/newmancai/oc-go-cc-optim.git
cd oc-go-cc-optim
make build
export PATH=$PWD/bin:$PATH
```

Or see [INSTALLATION.md](INSTALLATION.md) for more options.

#### 2. Initialize Configuration

```bash
oc-go-cc init
```

Creates a default config at `~/.config/oc-go-cc/config.json`. Edit it to add your API key, or set the environment variable:

```bash
export OC_GO_CC_API_KEY=sk-opencode-your-key-here
```

See [configs/config.oc-go-cc.json](configs/config.oc-go-cc.json) for a fully annotated example.

#### 3. Start the Proxy

```bash
oc-go-cc serve
```

#### 4. Configure Claude Code

```bash
export ANTHROPIC_BASE_URL=http://127.0.0.1:3456
export ANTHROPIC_AUTH_TOKEN=unused
```

Or add these to Claude Code's `~/.claude/settings.local.json`. See [configs/config.cc.json](configs/config.cc.json) for a complete example.

#### 5. Run Claude Code

```bash
claude
```

### CLI Commands

```
oc-go-cc serve              Start the proxy server
oc-go-cc serve -b           Start in background (detached from terminal)
oc-go-cc serve --port 8080  Start on a custom port
oc-go-cc stop               Stop the running proxy server
oc-go-cc status             Check if the proxy is running
oc-go-cc init               Create default configuration file
oc-go-cc validate           Validate configuration file
oc-go-cc models             List available OpenCode Go models
oc-go-cc autostart enable   Enable auto-start on login
oc-go-cc autostart disable  Disable auto-start on login
oc-go-cc autostart status   Check autostart status
oc-go-cc --version          Show version
```

### Manual Model Selection

In Claude Code, use `/model <scenario>` to override automatic routing and force a specific scenario from `config.json`:

```
/model default     → kimi-k2.6 (default scenario)
/model think       → glm-5 (thinking/reasoning)
/model complex     → glm-5.1 (complex architecture/tools)
/model background  → qwen3.5-plus (simple read-only ops)
/model long_context → minimax-m2.5 (1M context window)
/model fast        → qwen3.6-plus (fast streaming)
```

This works because the proxy checks if the Claude model name matches a key in `config.json`'s `models` section. You can also add custom aliases to the config file for any name you prefer.

### Documentation

| Document | Description |
| -------- | ----------- |
| [INSTALLATION.md](INSTALLATION.md) | Homebrew, Scoop, build from source, release binaries |
| [CONFIGURATION.md](CONFIGURATION.md) | Config file reference, env vars, model routing, fallback chains |
| [MODELS.md](MODELS.md) | Model capabilities, costs, and routing recommendations |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Development setup, architecture, how it works |
| [TROUBLESHOOTING.md](TROUBLESHOOTING.md) | Common issues and debug mode |

### License

[AGPL-3.0](LICENSE) — see the [LICENSE](LICENSE) file for details.

---

## Chinese

### 这是什么？

`oc-go-cc` 是一个 Go 语言编写的 CLI 代理工具，让你可以将 [OpenCode Go](https://opencode.ai/docs/go/) 订阅与 [Claude Code](https://docs.anthropic.com/en/docs/claude-code) 配合使用。

它位于 Claude Code 和 OpenCode Go 之间，拦截 Anthropic API 请求，转换为 OpenAI 格式，然后转发到 OpenCode Go。Claude Code 以为自己在和 Anthropic 通信，但实际上你的请求被路由到了更经济的开源模型。

### 为什么需要它？

OpenCode Go 以 **$5/月**（之后 $10/月）的价格提供强大的开源编程模型。这个代理让这些模型与 Claude Code 的界面无缝协作 —— 无需补丁，无需分支，只需设置两个环境变量即可。

### 快速开始

```bash
# 1. 安装（详见 INSTALLATION.md）
git clone https://github.com/newmancai/oc-go-cc-optim.git
cd oc-go-cc-optim
make build
export PATH=$PWD/bin:$PATH

# 2. 初始化配置
oc-go-cc init
export OC_GO_CC_API_KEY=sk-opencode-your-key-here

# 3. 启动代理
oc-go-cc serve

# 4. 配置 Claude Code
export ANTHROPIC_BASE_URL=http://127.0.0.1:3456
export ANTHROPIC_AUTH_TOKEN=unused

# 5. 运行 Claude Code
claude
```

### 手动选择模型

在 Claude Code 中使用 `/model` 命令覆盖自动路由，直接选择 `config.json` 中定义的场景：

```
/model default     → kimi-k2.6（默认场景）
/model think       → glm-5（思考/推理）
/model complex     → glm-5.1（复杂架构/工具调用）
/model background  → qwen3.5-plus（简单只读操作）
/model long_context → minimax-m2.5（长上下文）
/model fast        → qwen3.6-plus（快速流式响应）
```

工作原理：当 Claude 的模型名匹配到 `config.json` 中 `models` 段的某个 key 时，代理会直接使用该 key 对应的配置，跳过自动场景检测。

