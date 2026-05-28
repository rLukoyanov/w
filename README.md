<div align="center">

```
██╗    ██╗
██║    ██║
██║ █╗ ██║
██║███╗██║
╚███╔███╔╝
 ╚══╝╚══╝
```

# W

**Self-hosted team communication. One binary. Zero compromise.**

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)
[![Single Binary](https://img.shields.io/badge/Deploy-Single%20Binary-blueviolet?style=flat-square)]()
[![Status](https://img.shields.io/badge/Status-WIP-orange?style=flat-square)]()

_Inspired by [PocketBase](https://pocketbase.io) — the idea that powerful software doesn't need to be complex to deploy._

</div>

---

## What is W?

**W** is an open-source, self-hosted real-time communication platform — packaged as a **single executable binary** written in Go.

No Docker. No Kubernetes. No ops team. No cloud vendor lock-in.  
Download → run → done. Your data stays yours.

W brings servers, channels, roles, voice, and threads — and packages it all into one portable file you can run on a $5 VPS, a Raspberry Pi, or your laptop.

---

## Что такое W?

**W** — это open-source платформа для общения в реальном времени с возможностью самостоятельного хостинга, упакованная в **один исполняемый файл** на Go.

Никакого Docker. Никакого Kubernetes. Никакой облачной зависимости.  
Скачай → запусти → готово. Твои данные остаются у тебя.

W предоставляет серверы, каналы, роли, голосовые комнаты и треды — и всё это в одном портативном файле, который можно запустить на VPS за $5, Raspberry Pi или ноутбуке.

```bash
# That's it. Seriously.
./w serve
```

---

## Why W?

|                          | Discord | Self-hosted alternatives | **W** |
| ------------------------ | ------- | ------------------------ | ----- |
| Data ownership           | ❌      | ✅                       | ✅    |
| Single binary deploy     | ❌      | ❌                       | ✅    |
| No external dependencies | ❌      | ❌                       | ✅    |
| Scales to small teams    | ✅      | ⚠️                       | ✅    |
| Free & open source       | ❌      | ✅                       | ✅    |
| Built-in admin UI        | ❌      | ⚠️                       | ✅    |

---

## Key Features

### 🏗️ Single Binary Architecture

Everything — web server, WebSocket hub, database (SQLite), file storage, admin panel, API — is compiled into one file. No runtime dependencies. Works offline. Backup = copy one file.

### 💬 Servers & Channels

Organize your community into servers (guilds), with text channels, announcement channels, and threads. Full Markdown support with syntax highlighting.

### 🔴 Real-Time Everything

WebSocket-powered live updates. Messages, reactions, typing indicators, online status — all instant. Works behind NAT with built-in connection recovery.

### 🎙️ Voice & Video Rooms

WebRTC-based voice and video channels built in. No external media server needed for small teams (embedded SFU for larger deployments).

### 👥 Roles & Permissions

Granular permission system: per-server roles, per-channel overrides, admin groups. Works like Discord's model but you control the rules.

### 🤖 Bot API

Full REST + WebSocket API compatible with a simple bot framework. Write bots in any language. Event-driven webhooks as an alternative.

### 🔍 Full-Text Search

Message search across all channels powered by SQLite FTS5. No Elasticsearch. No extra infra.

### 📁 File Sharing

Drag-and-drop file uploads with image previews, audio players, video embeds. Files stored locally or on S3-compatible storage.

### 🎨 Themes & Customization

Server-level branding: custom icons, banners, emoji packs. Users can apply personal themes.

### 🔐 Auth That Works

Local accounts, SSO via OAuth2 (GitHub, Google, custom IdP), magic links, and 2FA — all built in.

### 📊 Built-in Admin Dashboard

Web-based admin panel at `/admin`. Manage users, servers, bans, storage, and monitor real-time stats — no CLI needed.

### 📦 Zero-Config Deployment

Sensible defaults out of the box. Drop a config file next to the binary to override anything. Automatic HTTPS via Let's Encrypt.

---

## Quickstart

```bash
# Download
curl -L https://github.com/your-org/w/releases/latest/download/w_linux_amd64 -o w
chmod +x w

# Run
./w serve

# Open browser
open http://localhost:8090
```

First run creates a fresh SQLite database and opens the setup wizard.

### Configuration (optional)

```yaml
# w.yaml
host: "0.0.0.0"
port: 8090
data_dir: "./w_data"

tls:
  auto: true # Let's Encrypt
  domain: "chat.example.com"

storage:
  type: local # or "s3"
  path: "./uploads"

limits:
  max_upload_mb: 50
  max_servers_per_user: 10
```

---

## Architecture

```
w (single binary)
├── HTTP server         — chi router, serves web app + API
├── WebSocket hub       — gorilla/websocket, fan-out message broker
├── SQLite (CGo-free)   — modernc.org/sqlite, WAL mode, FTS5
├── WebRTC SFU          — pion/webrtc for voice/video
├── File storage        — local FS or S3 via minio-go
├── Auth engine         — JWT + refresh tokens, OAuth2 client
├── Admin UI            — embedded SvelteKit build
└── Bot runtime         — WebSocket gateway + REST API
```

W is designed to be **extended as a Go library**, not just used as a binary — similar to PocketBase:

```go
package main

import (
    "github.com/your-org/w/core"
    "github.com/your-org/w/apis"
)

func main() {
    app := core.New()

    // Hook into message events
    app.OnMessageCreate().Add(func(e *core.MessageEvent) error {
        if e.Message.Content == "!ping" {
            e.Channel.Send("pong!")
        }
        return nil
    })

    apis.Serve(app)
}
```

---

## Roadmap

### Phase 1 — Foundation `v0.1` _(current focus)_

- [ ] Project scaffolding, CI/CD, release pipeline
- [ ] Core data models: users, servers, channels, messages
- [ ] HTTP API (REST) with OpenAPI spec
- [ ] WebSocket hub + event system
- [ ] Basic web client (read/send messages)
- [ ] SQLite storage with migrations
- [ ] JWT auth + local accounts
- [ ] File uploads (local storage)

### Phase 2 — Community Features `v0.2`

- [ ] Roles & permissions system
- [ ] Reactions, threads, pinned messages
- [ ] Message editing & deletion with history
- [ ] Invite links & server discovery
- [ ] Direct messages (DMs)
- [ ] Typing indicators, read receipts
- [ ] Admin dashboard v1

### Phase 3 — Rich Media `v0.3`

- [ ] Voice channels (WebRTC, p2p for 2 users)
- [ ] Embedded SFU for group voice (pion)
- [ ] Video rooms
- [ ] Screen sharing
- [ ] Link previews & rich embeds
- [ ] Image galleries, video players

### Phase 4 — Power Features `v0.4`

- [ ] Bot API & gateway (Discord-compatible events)
- [ ] Webhook support
- [ ] OAuth2 SSO (GitHub, Google, custom)
- [ ] Full-text search (SQLite FTS5)
- [ ] S3-compatible external storage
- [ ] Notification system (browser push, email digest)

### Phase 5 — Scale & Polish `v1.0`

- [ ] Horizontal scaling mode (multiple W nodes)
- [ ] Federation between W instances (ActivityPub?)
- [ ] Mobile apps (React Native or Go + Wails)
- [ ] Desktop app (Wails)
- [ ] Plugin system
- [ ] Audit logging for enterprises

---

## Contributing

W is at an early stage and **actively welcoming contributors**.

```bash
git clone https://github.com/your-org/w
cd w
go mod download
go run . serve --dev
```

The codebase follows standard Go project layout. Run tests with `go test ./...`.

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines, and [good-first-issues](https://github.com/your-org/w/labels/good%20first%20issue) for where to start.

---

## Comparison with Alternatives

| Project        | Language  | Deploy        | Voice       | Bot API | License |
| -------------- | --------- | ------------- | ----------- | ------- | ------- |
| **W**          | Go        | Single binary | ✅ Built-in | ✅      | MIT     |
| Mattermost     | Go        | Docker/K8s    | ❌          | ✅      | MIT/EE  |
| Rocket.Chat    | Node.js   | Docker/K8s    | ✅          | ✅      | MIT/EE  |
| Revolt         | Rust/Node | Docker        | ✅          | ✅      | AGPL    |
| Zulip          | Python    | Docker        | ❌          | ✅      | Apache  |
| Matrix/Element | Python    | Complex       | ✅          | ✅      | Apache  |

W's unique position: **Discord UX + PocketBase simplicity**.

---

## License

MIT — use it, fork it, build on it. Commercially too.

---

<div align="center">

Made with ☕ and Go · Star ⭐ if you find it useful

</div>
