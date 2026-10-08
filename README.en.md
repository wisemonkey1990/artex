<div align="center">

# ARTEX Chinese Edition

**An LLM-powered autonomous penetration testing system** (Go backend + Next.js frontend)

[简体中文](README.md) · [한국어](README.ko.md) · English

[![license: AGPL-3.0](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)](LICENSE)

</div>

---

> **Authorized security testing only.** ARTEX can autonomously perform reconnaissance, invoke tools, and validate security findings. Use it only in isolated environments that you own or are explicitly authorized to test in writing. Unauthorized scanning, access, or exploitation may be illegal and cause harm. Read the safety scope below and [LICENSE](LICENSE) before use.

ARTEX coordinates multiple LLM agents to break down goals, run tools, and record results. It presents findings through an asset graph and task activity view. The project includes a Go service, a Next.js web interface, and PostgreSQL storage. This repository defaults to Simplified Chinese and retains Korean as an optional interface language.

## Quick start

Requires Docker and Docker Compose. The first launch starts ARTEX and PostgreSQL. Visit `http://localhost:8787` and follow the page to set the administrator password.

```bash
git clone https://github.com/wisemonkey1990/artex.git
cd artex
cp .env.example .env
docker compose up -d
```

Set the database password in `.env` and configure an LLM API key if needed. To build from source, build the static frontend and embed it in the Go service:

```bash
cd web
npm ci
npm run build:static
cd ..
rm -rf server/webui/dist
mkdir -p server/webui/dist
cp -a web/out/. server/webui/dist/
CGO_ENABLED=0 go build -tags embedui -o artex ./cmd/artex
```

## Localization

- Simplified Chinese is the default interface language. Set `NEXT_PUBLIC_LOCALE=ko` at build time to use Korean.
- Dates and times use the `zh-CN` locale.
- Agent-authored natural language shown to users is requested in Simplified Chinese. Commands, code, URLs, request/response bodies, and raw evidence remain unchanged.
- `README.zh.md` preserves the upstream Chinese documentation; this page describes the localized edition in this repository.

## Safety scope

Run ARTEX only against targets covered by explicit authorization and agreed scope and time windows. Prefer isolated, intentionally vulnerable labs; do not use real user data in tests. Handle any data and credentials generated during testing according to the agreed rules and remove them promptly. See [SECURITY.md](SECURITY.md) for reporting security issues.

## Documentation and license

- [Simplified Chinese project overview](README.md)
- [Upstream Chinese documentation](README.zh.md)
- [Korean project overview](README.ko.md)
- [Contribution guide](CONTRIBUTING.md)
- [AGPL-3.0 license](LICENSE)
