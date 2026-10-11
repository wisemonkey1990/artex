# 贡献指南 (Contributing)

简体中文 · [English](CONTRIBUTING.en.md)

感谢你关注 ARTEX 中文版（`wisemonkey1990/artex`）。本文档整理了贡献前需要了解的范围·
方针·流程。提交贡献前，请务必先阅读[使用范围与法律责任](#使用范围与法律责任)和
[本地化方针](#本地化方针)。

- 报告 Bug 或提出功能建议 → 使用[issue 模板](https://github.com/wisemonkey1990/artex/issues/new/choose)。
- 发现翻译·本地化问题 → 使用"翻译·本地化问题"的 issue 模板。
- 发现安全漏洞 → **不要提交为公开 issue**，请按 [SECURITY.md](SECURITY.md) 的流程处理。
- 所有参与者都要遵守[行为准则（CODE_OF_CONDUCT.md）](CODE_OF_CONDUCT.md)。

---

## 使用范围与法律责任

ARTEX 是让 LLM 多智能体**自主**执行渗透测试的攻击性安全工具。贡献者受到与用户相同的
使用范围限制。

- 验证代码时，只能对**自己拥有或获得书面明确授权的目标**，或者**本地隔离环境**
  （例如用 Docker 启动的 OWASP Juice Shop·DVWA 这类有意留有漏洞、自己拥有的靶场）
  运行工具。
- 不接受针对授权范围之外的真实·生产·远程系统进行扫描·检测·利用的代码，也不接受
  鼓励这类用途的改动。
- 在韩国，未经授权侵入他人信息通信网络或造成故障，违反《信息通信网利用促进及信息
  保护等相关法律》，收集·泄露的个人信息还会受到《个人信息保护法》约束；具体法律
  提示详见 [README.ko.md 的"请先阅读"一节](README.ko.md#️-먼저-읽어-주세요--사용-범위와-국내법-고지)。
  使用者所在的其他司法辖区也可能有类似的网络安全·个人信息保护法规，使用前请自行
  确认当地法律。

提交的贡献代码·文档被如何使用，法律责任由实际运行它的用户自行承担。本仓库按
"现状（AS IS）"提供。

---

## 本地化方针

本仓库存在的意义，是在完全保留原版 [Autumn-27/ARTEX](https://github.com/Autumn-27/ARTEX)
**判断性能的前提下，只把用户可见的输出换成中文**。偏离这个方针的翻译贡献可能会拉低
性能，因此不予接受。

- **不要翻译智能体的内部推理提示词（行为指令正文）。** 必须保持用原文（中文）验证过
  的行为。这部分正文在 `agent/promptcatalog.go` 和数据库种子（`agent_prompts`）中。
  翻译会让智能体的判断出现漂移。
- **只强制用户可见的输出使用中文。** 检测结果（`report_finding`）、事实摘要
  （`record_fact`）、最终报告、对话回复属于这一类。这个强制是通过写死在代码里的尾部
  `agent/prompt.go` 的 `langDirective()` 实现的，会附加在每个角色 system 提示词的
  末尾。要改输出语言，改这个函数。
- **命令·payload·代码·URL·日志原文不翻译。** 这些是分析所需的原始材料，原样保留。
- **保留原始中文。** 文档保留在 `README.zh.md`，UI 字符串的原文保留在
  `web/messages/zh.json`，方便与上游（upstream）仓库的变更对照。韩语翻译则填在
  `web/messages/ko.json`。
- 翻译新的 UI 字符串时不要硬编码，要以消息文件的 key 形式添加。
- **输出语言强制是提示词引导，不是硬性上限。** `langDirective()` 只是**指示**模型用
  中文输出，并不会强制锁定。所以中文的输出完整度会随模型能力·角色·上下文而变化。
  验证本地化改动时请用能力足够（前沿级别）的模型。廉价·小模型的报告·摘要可能会退回
  原文（中文），不能只凭廉价模型的输出判断翻译是否生效。使用 OpenAI 系列模型时
  `max_tokens` 设置的陷阱，整理在 [README.ko.md 的"模型选择与输出语言"一节](README.ko.md#모델-선택과-출력-언어)。
- **跟进上游（upstream）变更的流程在维护者指南里。** 上游 ARTEX 更新时如何区分保留
  资产与翻译对象、如何检查翻译对称性与漂移，整理在 [MAINTAINING.md](MAINTAINING.md)。

---

## 开发环境

本项目由 **Go 后端**（单一二进制内嵌前端）+ **Next.js 前端**组成。

主要功能的设计意图整理在 `docs/` 的设计文档中。涉及把漏洞和流量证据关联起来的功能
（报告 Agent 的自动绑定、`report_finding` 的 `traffic_refs` 等）时，请先读
[漏洞多流量证据设计文档](docs/finding-traffic-evidence-zh.md)。韩语版在同一目录下的
`finding-traffic-evidence-ko.md`。

### 版本要求

- Go 1.26 及以上（以 `go.mod` 为准）
- Node.js 22 及以上（以发布工作流为准）
- Docker 与 Docker Compose（本地运行·验证用）

### 后端（Go）

本地装好 Go 的话，在仓库根目录执行：

```bash
go build ./...
go vet ./agent/
go test ./agent/
```

本地没有 Go 的话，可以用 Docker 做同样的验证。把模块·构建缓存放进 named volume 能让
重复运行更快。

```bash
docker run --rm -v "$PWD":/src -w /src \
  -v artexko-gomod:/go/pkg/mod -v artexko-gocache:/root/.cache/go-build \
  golang:1.26 sh -c 'go build ./... && go vet ./agent/ && go test ./agent/'
```

#### DB 集成测试（需要 postgres）

上面的 `go test ./agent/` 会**悄悄跳过需要连接 PostgreSQL 才能跑的 DB 集成测试**。
`agent`·`config`·`db`·`evidence`·`llmrec`·`server` 这六个包里有一些测试需要真实
数据库才能跑，如果没有设置环境变量 `ARTEX_PG_DSN`，配置文件里也没有 `database` 项，
这些测试就会以 `--- SKIP` 跳过，包仍然显示 `ok`。所以改了这六个包之后，不带 DSN
验证可能会**本地显示通过（ok），但 PR 的 `go-db` 检查却失败**。

要在本地跑这些测试，需要先起一个 PostgreSQL，再传入 `ARTEX_PG_DSN`。下面是起一个与
CI 一致的 `postgres:16-alpine`、放在隔离网络里跑的示例，复用上面同样的 named volume。

```bash
# 1) 起一个隔离网络和空的 postgres（与 CI 相同的镜像·账号）。
docker network create artexko-db 2>/dev/null || true
docker run -d --name artexko-pg --network artexko-db \
  -e POSTGRES_USER=artex -e POSTGRES_PASSWORD=artex -e POSTGRES_DB=artex \
  postgres:16-alpine
until docker exec artexko-pg pg_isready -U artex -d artex >/dev/null 2>&1; do sleep 1; done

# 2) 传入 DSN 跑 DB 集成包（DSN 里的 host 就是容器名）。
#    只想跑改动的包的话，把 ./agent/ 换成 config·db·evidence·llmrec·server 即可。
docker run --rm --network artexko-db -v "$PWD":/src -w /src \
  -v artexko-gomod:/go/pkg/mod -v artexko-gocache:/root/.cache/go-build \
  -e ARTEX_PG_DSN='postgres://artex:artex@artexko-pg:5432/artex?sslmode=disable' \
  golang:1.26 sh -c 'go test ./agent/ -count=1'

# 3) 清理。
docker rm -f artexko-pg && docker network rm artexko-db
```

CI 的 `go-db` 任务会把这六个包**各自用独立的 postgres 隔离开**，在合并前强制跑一遍
（`.github/workflows/ci.yml`）。如果改动了 DB 集成相关的包，建议在提 PR 前先按上面的
方法亲自验证一遍对应的包。

### 前端（web）

```bash
cd web
npm ci
npm run dev          # 开发服务器
npm run build        # 生产构建
npm run build:static # 静态导出构建（合并关卡 · 含 TypeScript 类型检查）
npm run check        # Biome lint·格式检查（仅供参考 · 历史欠账尚未作为合并关卡）
npm run check:fix    # 自动修复
```

提交前的格式化·lint 由 Biome 管理。`lint-staged` 会对已 stage 的文件自动跑一遍
`biome check --write`。

### 整套运行（Docker Compose）

```bash
cp .env.example .env     # 设置 POSTGRES_PASSWORD
docker compose up -d     # 启动 artex + postgres → http://localhost:8787
```

---

## 贡献流程

1. 先**开一个 issue。** 较大的改动建议先开 issue 对齐方向后再动手。小修改（拼写·
   链接·明显的 bug）可以直接提 PR。
2. **Fork** 仓库并创建主题分支。分支名以改动性质开头，例如 `feat/...`、`fix/...`、
   `docs/...`、`i18n/...`。
3. 编写改动，并**自行跑一遍相应范围的验证。** Go 改动要通过上面的
   `build`·`vet`·`test`。web 改动要通过 `npm run build:static`（合并关卡 · 同时做
   TypeScript 类型检查）。`npm run check`（Biome）因为还有上游带来的历史 lint
   欠账，暂时不是合并关卡，`web.yml` 里只是信息性的一步，不需要让它全部通过，只要
   **确认自己的改动没有新增错误**即可（提交时 `lint-staged` 只会对 stage 的文件自动
   跑 `biome check --write`）。如果改了文档（`.md`），用
   `python3 -I scripts/check-doc-links.py` 确认仓库内的链接·图片引用和文档锚点
   （`#标题`）没有失效。锚点是按 GitHub 同样的规则从标题生成 slug 后对照的，所以改了
   标题文字却没同步改指向它的锚点链接，会在这里被抓住（CI 的 `docs` 工作流会把同样的
   检查作为合并关卡强制执行）。这个检查也以 `docs` hook 的形式写进了仓库根目录的
   [`.pre-commit-config.yaml`](.pre-commit-config.yaml)，装好 `pre-commit install` 后
   提交时会自动跑（只用 Python 标准库，不访问网络，不需要 Docker）。
4. **开 PR。** 标题·说明按 [PR 模板](.github/PULL_REQUEST_TEMPLATE.md)填写，写清楚
   改了什么、为什么改，以及怎么验证的。改了 UI 的话附上截图。
5. 如果是用户可见的改动（功能·本地化·文档·检测规则等），在[变更记录（CHANGELOG.md）](CHANGELOG.md)
   的 `[Unreleased]` 一节加一行。内部重构或仅涉及测试的改动可以省略。

### 提交信息

遵循现有提交历史的惯例。格式是 `type(scope): 说明`，说明用中文写。

- `type`：`feat` · `fix` · `docs` · `chore` · `refactor` · `test` · `i18n` 等
- `scope`：改动涉及的区域（`agent` · `web` · `server` 等），可省略

示例：

```
feat(agent): 把用户可见输出强制为中文（langDirective）
docs: 编写中文 README，原文保留在 README.zh.md
i18n(web): 翻译仪表盘导航标签
```

---

## 贡献检测规则·检测测试

本仓库在 [`detections/`](detections/) 中同时维护用于**防御·检测**类似 ARTEX 这种
自主 AI 攻击的规则。包括可部署的 [Sigma](https://sigmahq.io) 规则
（[`detections/sigma/`](detections/sigma/)）、面向网络流量的
[Suricata](https://suricata.io) 规则（[`detections/suricata/`](detections/suricata/)）、
[MITRE ATT&CK](https://attack.mitre.org/) 覆盖层（[`detections/attack/`](detections/attack/)），
以及以可复现方式证明这些规则确实会触发的测试
（[`detections/tests/`](detections/tests/)）。提交或修改检测规则时，请遵守以下约定。
八个测试套件会机械式地强制执行其中大部分约定，只改规则而不同步更新测试·覆盖层会导致
测试失败。

- **所有指标都要有可观测的事实依据。** 规则使用的字符串·User-Agent·行为阈值，必须是
  本仓库源码中能实际确认的内容，不能凭猜测编造。请在规则里标明依据的源文件（例如
  `artex-enrich/1.0` 这个指标可以在 `enrich/enrich.go` 中确认）。指标一致性测试
  （[`detections/tests/indicators/`](detections/tests/indicators/)）会检查每个指标在
  上游源和规则两边是否仍然一致，上游重新同步改了源字符串而规则没跟上，测试就会失败。
  修改机器可读的指标列表
  （[`detections/indicators/artex_indicators.csv`](detections/indicators/artex_indicators.csv)）
  时，要同步更新承载同一批指标的 MISP 事件
  （[`detections/indicators/artex_indicators.misp.json`](detections/indicators/artex_indicators.misp.json)）。
  MISP 导出测试（[`detections/tests/misp/`](detections/tests/misp/)）会强制这两个文件
  逐行一致，并确认该事件是能被 pymisp 加载的合法 MISP 文档。
- **诚实写明局限。** Sigma 规则在 `description` 里、Suricata 规则在注释里，写明这条
  规则抓不到的情况和可能的误报。如果不是 ARTEX 专属特征而是通用的狩猎线索（例如
  破坏性命令），要明确说明，避免单次命中就武断地认定攻击者是 ARTEX。
- **通过静态校验。** Sigma 规则必须以 0 个问题通过 SigmaHQ 校验器标准
  （`sigma check --validation-config detections/tests/sigma_lint/validators.yml`）。
  默认的 `sigma check` 只跑 pySigma 核心校验器，标题写法·字段/日志源分类·引用链接这类
  SigmaHQ 规范只有用这个标准才会被检出来。有四个例外是不适用于独立规则集的 SigmaHQ
  monorepo 规范，理由写在了
  [`detections/tests/sigma_lint/validators.yml`](detections/tests/sigma_lint/validators.yml)。
  Suricata 规则必须能被 `suricata -T` 干净加载。
- **同时提交可复现的测试。** 用 [`detections/tests/`](detections/tests/) 下的测试证明
  规则确实会触发（或者结构是有效的）。测试输入每次都要确定性生成，不要把二进制文件
  提交进仓库；与引擎版本无关的属性（确实触发·没有误报）要精确断言，随版本波动的数值
  按下限断言并单独记录基准值。新增或修改 Sigma 关联规则时，后端可移植性测试
  （[`detections/tests/sigma_backends/`](detections/tests/sigma_backends/)）会确认该
  规则能在多个后端转换，请保持与
  [`detections/README.md`](detections/README.md) 里的后端支持说明一致。
- **同步更新 ATT&CK 覆盖层。** 给规则新增或修改 `attack.*` 标签时，要同步更新
  [`detections/attack/artex_navigator_layer.json`](detections/attack/artex_navigator_layer.json)
  里对应的技术·评分。一致性测试会强制规则↔覆盖层双向一致，覆盖层里没有对应规则标签，
  或规则里用了覆盖层没有的技术，都会导致失败。
- **不要写成攻击指南。** 本仓库的检测资料只保持防御·检测的定位，不接受帮助执行利用、
  或绕过检测的技术描述。

八个测试套件只需要 Docker 就能直接运行，不会把生成物提交进仓库。每个脚本只要有一个
断言失败就会以非 0 退出码结束，可以直接接入 CI 或 pre-commit hook。

```bash
detections/tests/sigma/run.sh           # Sigma：sigma check + 后端转换 + 指标保留
detections/tests/sigma_match/run.sh     # Sigma：原子规则在恶意样本触发·在正常样本静默
detections/tests/sigma_lint/run.sh      # Sigma：SigmaHQ 规范完整校验器 + 文档化标准
detections/tests/sigma_backends/run.sh  # Sigma 可移植性：关联规则能否在多个后端转换
detections/tests/suricata/run.sh        # Suricata：合成 pcap → suricata -r → 断言告警数
detections/tests/attack/run.sh          # ATT&CK：覆盖层 ↔ 规则双向一致
detections/tests/indicators/run.sh      # 指标：规则锚定的指标 ↔ 上游源双向一致
detections/tests/misp/run.sh            # MISP：指标 CSV ↔ MISP 事件同步 + pymisp 有效性
```

要一次跑完这八个，用 [`detections/tests/run-all.sh`](detections/tests/run-all.sh)。它
按与 CI 相同的顺序依次执行这八个，前面的套件失败也会继续跑完剩下的，最后输出每个套件
的 PASS/FAIL 汇总，只要有一个失败就以非 0 退出码结束。仓库根目录的
[`.pre-commit-config.yaml`](.pre-commit-config.yaml) 里有直接把这个 runner 挂成
pre-commit hook 的配置示例。用 `pip install pre-commit && pre-commit install` 安装后，
只有检测规则或其锚定的上游源发生变化的提交才会触发这个 runner（范围与 CI 一致），在
push 之前就抓住规则·测试不一致的问题。同一个配置文件里还有检查文档内部链接·图片·
锚点的 `docs` hook（也就是上面贡献流程第 3 步提到的 `check-doc-links.py`）。

这八个测试由仓库 CI（[`.github/workflows/detections.yml`](.github/workflows/detections.yml)）
在 `detections/` 下有变化的每次推送·PR 中运行。指标一致性测试在那些指标锚定的上游源
文件（`enrich/`·`selfupdate/`·`guard/`·`db/`·`cmd/artex/main.go`）变化时也会运行，
用来同时抓住上游重新同步改了 User-Agent·标记·默认端口而规则悄悄过时的情况。因此，
只改规则而没更新测试·覆盖层的改动、违反 SigmaHQ 规范的规则，或者与源码脱节的规则，
都会在合并前被 CI 标红。

规则索引、各规则的依据与局限，整理在 [`detections/README.md`](detections/README.md)；
测试的断言项和运行方法，整理在 [`detections/tests/README.md`](detections/tests/README.md)。

---

## 许可

本项目以 **GNU Affero General Public License v3.0（AGPL-3.0）** 发布。提交贡献即视为
**同意该贡献也以 AGPL-3.0 公开**。特别是如果你修改本项目后通过网络（例如作为在线
服务）提供给用户，必须向这些用户公开对应的完整源代码。完整条款见 [LICENSE](LICENSE)
文件。
