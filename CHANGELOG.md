# 变更记录

简体中文 · [English](CHANGELOG.en.md) · [中文(原版·上游)](CHANGELOG.zh.md)

本文档记录 ARTEX 中文版（本 fork）相对上游仓库新增的变更。格式参考
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/)。

上游 ARTEX 项目按版本的发布历史（0.3.x 及以下）和贡献者名单，按原始中文原样保留在
[`CHANGELOG.zh.md`](CHANGELOG.zh.md) 中，与 `README.zh.md` 采用相同的做法保留原文，
方便与上游变更对照。每项变更的具体内容和依据可以在仓库的 commit 历史中查到。

## [Unreleased] · 本 fork 变更

### 本地化调整：默认语言改回简体中文

- **把默认界面和智能体输出语言从韩语改回简体中文，韩语保留为可选语言。** 构建时设置
  `NEXT_PUBLIC_LOCALE=ko` 仍可切换到韩语界面；默认（不设置该变量）构建的是简体中文。
  README·CONTRIBUTING.md·MAINTAINING.md 等说明文档，以及脚本·CI 工作流·issue 和 PR
  模板里原本以韩语为主的提示文字，都已翻译成简体中文。
- **修复了 CI 的语言混入检测关卡。** 原来的 `scripts/check-web-cjk.py` 检查静态导出的
  默认（曾是韩语）界面里"不能混入中文"，默认语言改回中文后这个关卡会在每次构建都
  报红。已重命名为 `scripts/check-web-hangul.py` 并反转检测逻辑，改为检查默认（中文）
  构建产物里"不能混入韩文谚文"，[`web.yml`](.github/workflows/web.yml) 中的调用和说明
  已同步更新。
- **管理端口 8787 只绑定本机地址，不对外暴露。**
- **自更新发布源改为指向本 fork 自身的仓库**（`wisemonkey1990/artex`），不再从上游
  （`Autumn-27/artex`）拉取发布版本。
- **修复了多处文档·模板·工作流中指向旧仓库地址的失效链接。** issue 模板·PR 模板·
  `MAINTAINING.md`·`MAINTAINING.en.md` 中引用的 `jiwoochris/artex-ko`（本仓库早期的
  地址）已更新为当前地址 `wisemonkey1990/artex`。

### 本地化 (i18n)

- **把用户可见的输出强制为韩语。** 为保留性能，经过基准验证的智能体行为指令正文
  （大脑）保留原文不动，用写死在代码里的片段（`langDirective`）只指示面向用户的
  产出（漏洞报告、事实摘要、最终摘要、聊天回复）用韩语书写。命令·payload·代码·
  日志原文保留原样。
- **把 Web UI 迁移到韩语。** 在 Next App Router 中引入 `next-intl`，把字符串拆分到
  `web/messages/ko.json` 和 `web/messages/zh.json`。原始中文保留在 `zh.json` 中，方便
  与上游更新对照。仪表盘·漏洞·对话·通知发送·拦截·LLM 设置等界面文案都迁移成了韩语。
- **把服务端 API 面向用户的错误·应答迁移到韩语。** 回传给浏览器的 HTTP 错误·应答文案
  替换成了韩语。但是会被回读进智能体大脑输入的文案，为避免基准性能漂移保留了原文，
  判定依据记录在仓库的工作文档中。
- **整理了韩语文档。** 制作了韩语 `README.md`，同时配有英语版 `README.en.md`，原始
  中文保留为 `README.zh.md`。

### 防御·检测资料

- **新增了防御·检测指南。** 整理了自主 AI 攻击与传统扫描器的区别、防御者可观测的
  指纹（IoC）、入口点与加固、检测规则、事件响应的韩语指南
  （[`docs/defense-ko.md`](docs/defense-ko.md)），以及内容相同的英文版
  （[`docs/defense-en.md`](docs/defense-en.md)）。
- **提供可直接部署的检测规则。** 把指南里的指纹落地成可用的规则：主机·日志层是
  [Sigma](https://sigmahq.io) 原子·关联规则（[`detections/sigma/`](detections/sigma/)），
  网络层是针对 enrich 探测器和 norma SDK WebFetch 的 User-Agent 的
  [Suricata](https://suricata.io) 规则（[`detections/suricata/`](detections/suricata/)）。
- **把 ATT&CK 覆盖情况可视化。** 把规则标注的技术整理成 MITRE ATT&CK Navigator 层
  （[`detections/attack/`](detections/attack/)）。
- **以标准格式提供机器可读的入侵指标（IoC）。** 把 ARTEX 对外暴露的独有指纹汇总成一份
  CSV（[`detections/indicators/artex_indicators.csv`](detections/indicators/artex_indicators.csv)），
  以及可直接导入威胁情报平台的 MISP 事件
  （[`detections/indicators/artex_indicators.misp.json`](detections/indicators/artex_indicators.misp.json)）。
  有规则支撑的指标标记为 `to_ids`，主机取证用的端口标记为分类线索加以区分。
- **加上了可复现的检测测试。** 新增八种实际运行规则来证明其有效的测试（Sigma 结构·
  编译校验、Sigma 实时事件匹配、后端可移植性、SigmaHQ 规范 lint、Suricata 加载·触发、
  ATT&CK 层一致性、指标-源一致性、MISP 导出 ↔ CSV 同步），以及一次性全部运行的批量
  runner·pre-commit 示例，并接入 CI 合并关卡。Sigma 实时事件匹配不仅确认规则能编译，
  还会在原子·关联规则上确认规则确实在恶意样本事件上触发、在正常事件上保持静默。

### 仓库维护

- **加入了安全·滥用警告和当地法律提示。** 在 README 顶部加入了使用范围、《信息通信网
  法》·《个人信息保护法》提示、禁止滥用的警告。
- **把界面预览截图换成了韩语 UI 截图。**
- **整理了维护者手册和贡献指南。** 新增了防止上游同步·翻译漂移的手册
  （[`MAINTAINING.md`](MAINTAINING.md)）和检测规则贡献约定
  （[`CONTRIBUTING.md`](CONTRIBUTING.md)）。手册里还整理了发布流水线的构建前提，以及
  不打 tag 在本地验证这些前提的方法。
- **新增了推送·PR 合并关卡 CI。** 上游仓库只在打 tag 发布时跑 CI，这个 fork 则在每次
  推送和 PR 都跑 Go 构建·静态分析（`go vet`）·单元测试
  （[`ci.yml`](.github/workflows/ci.yml)）、韩语 UI 静态构建
  （[`web.yml`](.github/workflows/web.yml)）、文档的仓库内部链接·图片引用完整性
  （[`docs.yml`](.github/workflows/docs.yml)），在合并前抓住本地化过程中产生的回归。
  需要数据库的集成测试用每个包各自隔离的 PostgreSQL 服务一起验证。文档链接检查用不
  依赖外部网络的确定性脚本（[`scripts/check-doc-links.py`](scripts/check-doc-links.py)）
  运行，防止多语言文档之间大量互相指向的相对链接和界面预览图片带着失效状态被合并。
  文档锚点（`#标题`）链接也用与 GitHub 相同的 slug 规则与标题对照，一并抓住标题文字
  改动后悄悄失效的目录·互相引用链接。检测规则套件由上面"防御·检测资料"一节描述的
  合并关卡负责。
- **定期检查外部链接的存活状态。** 防御指南指向的事件报告渠道·标准参考等外部链接依赖
  远程服务器状态，结果不稳定，所以不放进合并关卡，而是由不挡合并的工作流
  （[`external-links`](.github/workflows/external-links.yml)）在每周一和手动触发时跑
  浏览器 User-Agent·GET·跟随重定向的检查
  （[`scripts/check-external-links.py`](scripts/check-external-links.py)）。主机存活、
  只是检测方式被挡住（防爬虫·限流）的情况，以及我们无法修复的上游继承失效链接
  （allowlist）不计入失败，只有我们文档整理的外部链接出现新的失效时才会变红。
- **搭建了贡献·治理基础设施。** 新增了 Bug·功能·翻译 issue 模板
  （[`.github/ISSUE_TEMPLATE/`](.github/ISSUE_TEMPLATE/)）、Pull Request 模板
  （[`PULL_REQUEST_TEMPLATE.md`](.github/PULL_REQUEST_TEMPLATE.md)）、安全漏洞报告策略
  （[`SECURITY.md`](SECURITY.md)）、行为准则（[`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md)），
  让外部贡献者能用统一格式提交 issue·PR·安全报告。
- **补全了面向海外贡献者的英文文档层。** 本仓库以韩语为主要语言，但为了让不懂韩语的
  贡献者·安全研究者·防御者也能获取同样的信息，核心文档都配有英文版。除了英文
  `README.en.md`·防御指南（[`docs/defense-en.md`](docs/defense-en.md)），还有变更记录
  （[`CHANGELOG.en.md`](CHANGELOG.en.md)）、安全报告策略
  （[`SECURITY.en.md`](SECURITY.en.md)）、行为准则
  （[`CODE_OF_CONDUCT.en.md`](CODE_OF_CONDUCT.en.md)）、贡献指南
  （[`CONTRIBUTING.en.md`](CONTRIBUTING.en.md)）、维护者手册
  （[`MAINTAINING.en.md`](MAINTAINING.en.md)）、流量证据设计文档
  （[`docs/finding-traffic-evidence-en.md`](docs/finding-traffic-evidence-en.md)），以及
  Bug·功能·翻译 issue 模板的英文版。韩语版和英语版在页眉互相指向对方，不管从哪种语言
  进入都能切换到另一种。（Pull Request 模板目前只有中文版。）

---

上游 ARTEX 项目按版本的发布历史和贡献者名单，可以在 [`CHANGELOG.zh.md`](CHANGELOG.zh.md)
中看到原文。
