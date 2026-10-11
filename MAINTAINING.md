# 上游同步与翻译漂移防护（维护者指南）

简体中文 · [English](MAINTAINING.en.md)

本文档是写给**维护者**的，用来说明如何在跟进上游仓库 [Autumn-27/ARTEX](https://github.com/Autumn-27/ARTEX)
变更的同时，维持本地化（韩语界面）的完整性。贡献范围·法律责任·本地化方针在
[CONTRIBUTING.md](CONTRIBUTING.md) 中说明，面向用户的说明在 [README.md](README.md) 中，
本文档只关注这些方针**具体如何落实执行**。

本地化的核心目标可以用一句话概括：**在完全保留原版判断性能的前提下，只把用户可见的
输出换成目标语言**。每次上游更新，这条边界都容易被打乱，所以用下面的流程和检查来
防止翻译漂移。

---

## 1. 本地化结构总览

本仓库是 **fork** 自上游 ARTEX，在其历史基础上叠加本地化相关的提交。上游 `main` 的
所有提交都包含在本仓库历史中，本地化提交叠加在其之上。因此，引入上游变更的工作就是
"核对与上游 `main` 的差异，分辨哪些要原样保留、哪些要翻译，然后落地"。

产出物分为三类：

- **保持原文不动的资产**（见下面第 2 节）。翻译会破坏性能或破坏与上游的对照能力。
- **写死在代码里的输出语言强制**。`agent/prompt.go` 的 `langDirective()` 会在每个
  角色 system 提示词末尾追加"面向用户的输出请用简体中文书写"的指示。
- **翻译成目标语言的用户可见字符串**。UI 文案在 `web/messages/zh.json`
  （默认·简体中文）与 `web/messages/ko.json`（可选·韩语）中维护，服务端面向用户的
  应答文案放在各个 Go 文件的具名常量里。

---

## 2. 保留原文的资产（禁止翻译）

以下资产不翻译，保持原文（中文或英文）。上游变更触及这些资产时，**原样照搬，不做
翻译**。

- **智能体内部推理提示词（大脑正文）。** `agent/promptcatalog.go` 和数据库种子
  `agent_prompts` 中的行为指令正文。必须保持原文（中文）才能延续已验证过的行为，
  翻译会让判断出现漂移。
- **同时承担"展示"与"智能体输入"双重角色的字符串。** 一部分文案（任务中止原因·
  拦截阻断消息·流量证据辅助说明等）既会显示在活动时间线上，又会作为规划器·报告器
  的输入上下文被回读，一条记录兼两种用途，因此保留原文。具体判定依据逐条记录在
  `work/DECISIONS-FOR-JIWOO.md` 的"大脑边界记录"中。
- **原始中文文档·字符串。** 文档保留在 `README.zh.md`，UI 字符串的原文保留在
  `web/messages/zh.json`，方便与上游变更对照。韩语翻译只填在 `web/messages/ko.json`。
- **命令·payload·代码·URL·标识符·日志原文。** 这些是分析所需的原始材料，不翻译。
  Go 代码注释的翻译优先级也最低，在与上游对照完成之前保留原文。

---

## 3. 上游跟踪基线

本地需要配置好如下上游远程，没有就添加：

```bash
git remote add upstream https://github.com/Autumn-27/ARTEX
git remote -v   # 确认 upstream 是否出现
```

当前本地化已经同步到的上游基线提交如下：

- **基线 = `d003372`**（上游 `main`，2026-10-03，合并 PR #189 `fix/sse-same-origin`）。

这个值的含义是"截止这个提交为止的上游变更都已经融合进本仓库"。每次引入新的上游变更，
都按第 7 节的方法更新这个基线。

---

## 4. 引入上游变更的流程

### 4.1 拉取上游并查看差异

```bash
git fetch upstream
git rev-list --count d003372..upstream/main        # 未合入的上游提交数
git log --oneline d003372..upstream/main           # 未合入提交列表
```

`git fetch` 只更新上游的远程跟踪分支，不影响工作树和 `HEAD`。如果未合入提交数为 0，
说明已经与上游同步，不需要再做什么。

### 4.2 对改动文件分类

查看未合入提交改动了哪些文件，按第 2 节分成"保留资产"和"翻译对象"。

```bash
git log --name-status --oneline d003372..upstream/main
```

分类标准如下：

- 如果改动了 `agent/promptcatalog.go`·`agent_prompts` 种子，或第 2 节里"展示兼输入"
  的字符串 → **原样合入，不翻译**。
- 如果改动了 Go 后端逻辑（`db/`·`llmrec/`·`server/` 等）→ 逻辑原样合入，但要检查是否
  出现了**新的用户应答文案**（如 `writeErr` 等），把它们翻译成具名常量。
- 如果改动了 UI（`web/src/**`）并产生了**新的界面文案** → 不要硬编码，用同一个 key
  同时加进 `web/messages/zh.json`（原文）和 `web/messages/ko.json`（翻译）。
- 如果**检测规则锚定的上游指标**（`enrich/enrich.go` 的探测 User-Agent、`selfupdate/`
  的自更新 User-Agent、`guard/guard.go` 的审计标记、`db/db.go` 的破坏性命令拒绝列表、
  `cmd/artex/main.go` 的默认监听·记录代理端口）发生变化 → 同步更新 `detections/` 下的
  Sigma·Suricata 规则、ATT&CK 层，以及 `detections/indicators/artex_indicators.csv` 里
  对应的值。这些指标不是翻译对象，而是**检测的依据**，上游改了值而规则没跟上，规则
  就会悄悄失效。5.4 节的指标一致性测试会自动抓住这种错位。

### 4.3 落地合入

按功能单元合并或挑选合入后，把 4.2 中识别出的新字符串翻译出来。合并过程中
`zh.json`·`ko.json` 的 key 容易错位，或者原文容易泄漏到用户可见的位置，所以合入后
必须跑一遍第 5 节的检查。

> **示例（截至 2026-10-05 的未合入提交）。** `git fetch upstream` 结果显示上游 `main`
> 领先到 `b55ceb1`，相对基线 `d003372` 有 2 个未合入提交（`86729b6` 模型回退授权
> token 计量功能 + 合并提交 `b55ceb1`）。这个提交改动了 `db/llm_usage.go`·
> `llmrec/llmrec.go`·`server/intercept.go`·`server/server.go` 这类 Go 逻辑，以及
> `web/src/app/(main)/system/intercept/page.tsx`·`web/src/lib/api.ts`·
> `web/src/lib/mock/handler.ts`·`web/src/lib/types.ts`。因此维护者只需原样合入 Go
> 逻辑，再把拦截配置页面新增的界面文案提取并翻译为 `zh.json`·`ko.json` 的 key 即可。
> （写这份文档时这两个提交尚未合入，所以基线仍记为 `d003372`。）

---

## 5. 翻译对称性与漂移检查

合入上游或做完翻译后，确认以下三件事。

### 5.1 zh ↔ ko 消息对称、以及用户可见处的韩文检查

`zh.json` 和 `ko.json` 的 key 必须完全一致，并且 `ko.json` 的值里不应该残留未翻译的
韩文谚文——换个角度说，真正需要盯住的回归是**用户默认看到的 `zh.json` 里不能混入
韩文谚文**。下面的脚本会打印出几个数字：

```bash
python3 - <<'PY'
import json, re
ko = json.load(open('web/messages/ko.json', encoding='utf-8'))
zh = json.load(open('web/messages/zh.json', encoding='utf-8'))
def flatten(d, p=''):
    out = {}
    if isinstance(d, dict):
        for k, v in d.items(): out.update(flatten(v, p + '/' + k))
    elif isinstance(d, list):
        for i, v in enumerate(d): out.update(flatten(v, p + '/' + str(i)))
    else: out[p] = d
    return out
fk, fz = flatten(ko), flatten(zh)
hangul = re.compile(r'[가-힣]')
print('ko leaf keys  :', len(fk))
print('zh leaf keys  :', len(fz))
print('key symdiff   :', len(set(fk) ^ set(fz)))        # 应为 0
print('zh vals w/谚文:', sum(1 for v in fz.values() if isinstance(v, str) and hangul.search(v)))  # 应为 0
PY
```

`key symdiff` 和 `zh vals w/谚文` 必须始终是 0；key 数量会随上游合入增长，但 zh 和
ko 的 key 数量应当始终相同。

### 5.2 大脑资产的原文保留确认

大脑正文要保持中文原文，所以下面这项检查里**如果汉字行数掉到 0，反而说明大脑被
误翻译污染了**。

```bash
python3 -c "import re; han=re.compile(r'[㐀-鿿豈-﫿]'); t=open('agent/promptcatalog.go', encoding='utf-8').read(); print('promptcatalog.go CJK lines =', sum(1 for l in t.splitlines() if han.search(l)))"
```

基准值（2026-10-05）：`promptcatalog.go CJK lines = 70`。如果这个数字明显下降，需要
确认大脑正文是否被误翻译。

### 5.3 构建产物里没有泄漏韩文

默认构建（不设置 `NEXT_PUBLIC_LOCALE`，即 `zh`）静态导出后，如果预渲染 HTML 里出现
韩文谚文，就是翻译泄漏（例如某处硬编码了韩文、或者忘了从 `zh.json` 切回中文值）。

```bash
cd web && npm ci && npm run build:static   # 生成 out/
python3 -I ../scripts/check-web-hangul.py  # out/ 里的用户可见 HTML 应该是 0 处韩文
```

CI 的 [`web.yml`](.github/workflows/web.yml) 已经把这一步设成合并关卡，每次 web/
变化都会自动跑一次。如果反过来要验证**可选的韩语构建**（`NEXT_PUBLIC_LOCALE=ko npm
run build:static`）没有中文泄漏，用同样的方法手动检查输出中是否还残留汉字——这个
方向目前没有自动化的 CI 关卡，需要在发布可选韩语构建前手动确认。

### 5.4 检测指标是否仍与上游源一致

`detections/` 下的规则，依据的是上游实际产出的字符串（探测 User-Agent·自更新
User-Agent·审计标记·破坏性命令拒绝列表）。上游重新同步改了这些值，翻译检查会全部
通过，但已发布的规则会悄悄停止匹配。下面的测试会双向确认每个指标在上游源和规则两边
是否仍然一致，重新同步之后要一起跑一遍。

```bash
detections/tests/indicators/run.sh   # 用 Docker 隔离运行，RESULT: PASS 即一致
```

失败时会打印出哪个指标不一致、以及不一致的方向（是上游源变了，还是规则变了），按
4.2 节最后的分类标准把规则·覆盖层调整到新值。这个测试同时也接入了仓库 CI
（[`.github/workflows/detections.yml`](.github/workflows/detections.yml)），规则树或
上述上游源文件发生变化的每次推送·PR 都会自动跑一遍，在合并关卡上抓住重新同步带来
的漂移。

新增指标时如果**锚定了新的上游源文件**（比如像给 `cmd/artex/main.go` 的端口指标那样
新增一条），必须把那个文件也加进上面工作流的 `push`·`pull_request` 的 `paths`
过滤里。漏加的话，只改了那个源文件的 PR 不会触发指标测试，漂移就会悄悄通过合并关卡。
这个同步关系本身也由指标测试自动检查（第五项检查"CI triggers this test when any
pinned source changes"）：测试读取的所有非 `detections/` 源文件，如果没有同时列在两处
`paths` 区块里，测试就会失败，这样源文件锚定和 CI 触发条件就不会在错位的情况下被合并。

### 5.5 升级检测测试工具版本时

检测测试把 `sigma-cli`·SigmaHQ 校验器插件（`pySigma-validators-sigmahq`）·Suricata
镜像锁定在固定版本上运行（各 `run.sh` 的默认值，可用环境变量覆盖）。升级这个锁定版本
可能带来的不是上游源的漂移，而是**工具本身的漂移**。SigmaHQ 校验器尤其会随版本升级
新增规范检查，`detections/tests/sigma_lint/run.sh` 可能因此冒出新的报红项。这时要把
规则改成符合新规范，或者如果某条规范不适用于某个独立规则集，就把理由写清楚后加进
[`detections/tests/sigma_lint/validators.yml`](detections/tests/sigma_lint/validators.yml)
的排除列表。后端插件改变支持范围时，`sigma_backends` 测试会给出同样的信号。

---

## 6. 用构建和测试做最终验证

合入·翻译之后，按 [CONTRIBUTING.md 的开发环境](CONTRIBUTING.md#开发环境) 流程验证
后端和前端。本地没有 Go 的话，可以用 Docker 跑同样的步骤。

```bash
docker run --rm -v "$PWD":/src -w /src \
  -v artexko-gomod:/go/pkg/mod -v artexko-gocache:/root/.cache/go-build \
  golang:1.26 sh -c 'go build ./... && go vet ./... && go test ./... -count=1'
```

翻译用户可见文案时，要同时配一个断言该文案的回归测试（`*_localized_test.go`），这样
以后上游变更再带回原文，测试能及时抓住。翻译结果的验证必须用能力足够（前沿级别）的
模型来做。廉价·小模型的输出可能又退回原文，不能只凭它自己的输出判断翻译是否生效。

---

## 7. 基线更新记录

引入上游变更并验证完毕后，把本文档第 3 节的**基线提交值更新为新的上游提交**，并把
这个改动包含在同一个提交或紧跟的提交里。这样下一个维护者能在本文档这一处就看清楚
"合入到哪里了"。

提交信息遵循 [CONTRIBUTING.md 的提交信息规则](CONTRIBUTING.md#提交信息)。例如上游
同步提交可以这样写：

```
chore(upstream): 合入上游 d003372..b55ceb1（拦截 token 计量）+ 翻译新增 UI 字符串
```

---

## 8. 检查·验证注意事项（常见陷阱）

核对上游合入·翻译·文档补充时，维护者反复踩到的两个陷阱记录在这里。两个都是
"检查方法本身用错了，把正常状态误判为出问题"的情况，为了避免不必要的回退，固定成
注意事项。

### 8.1 确认仓库 CI 状态时要显式指定仓库

本仓库是上游 ARTEX 的 fork，本地 `git remote` 里同时注册了 `origin`
（wisemonkey1990/artex）和 `upstream`（Autumn-27/ARTEX）（见第 3 节）。在这种状态下，
`gh` 命令如果不显式指定仓库，**会默认选中上游仓库**，显示的是上游那边、没有我们工作流
的运行结果。这样容易**把上游 CI 绿色误以为是我们自己的 CI 通过了**，或者把我们自己的
工作流（`ci.yml`·`detections.yml`）误判成"HTTP 404 … not found"。

所以确认 CI 时始终显式指定仓库：

```bash
gh run list -R wisemonkey1990/artex --workflow ci.yml --limit 5
gh run list -R wisemonkey1990/artex --workflow detections.yml --limit 5
```

设置一次之后，也可以让 `gh` 在省略 `-R` 时默认看我们自己的仓库。不过这是**本地 gh
配置**，不会提交进仓库，换新机器或新的 checkout 要重新设置一次。

```bash
gh repo set-default wisemonkey1990/artex
gh repo set-default --view   # 确认显示的是 wisemonkey1990/artex
```

### 8.2 文档里的外部链接要像浏览器一样用 GET 确认

防御指南（[`docs/defense-ko.md`](docs/defense-ko.md)·[`defense-en.md`](docs/defense-en.md)）
第 7 节列出了韩国官方渠道（boho.or.kr·fsec.or.kr·pipc.go.kr）的链接。确认这些链接是否
存活时，如果只用 `curl -I`（HEAD 请求）或默认 User-Agent 检测，**会把正常链接误判为
失效**。原因是这几家韩国公共·安全机构的站点出于以下三个原因会拒绝简单检测：

- **拒绝 HEAD 请求。** 例如 fsec.or.kr 对 `curl -I`（HEAD）会返回 400。
- **屏蔽 `curl` 默认 User-Agent。** fsec.or.kr 和 pipc.go.kr 对默认 UA 发出的 GET 请求
  同样返回 400（换成浏览器 UA 则返回 200）。
- **会重定向到其他地址。** pipc.go.kr 会从 `www.pipc.go.kr` 两次重定向到
  `pipc.go.kr/np/`，不跟随重定向就会错过最终状态。

所以确认链接要**用浏览器 User-Agent、发 GET 请求、并跟随重定向**：

```bash
UA='Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0 Safari/537.36'
for u in https://www.boho.or.kr https://www.fsec.or.kr https://www.pipc.go.kr; do
  curl -sS -L -A "$UA" -o /dev/null -w "$u -> %{http_code} %{url_effective}\n" "$u"
done
```

最终状态码是 200 就说明链接有效。如果状态码是 400·403，先不要急着认定链接失效，而是
怀疑**检测方法被服务器的访问策略挡住了**，逐一排除 HEAD·默认 UA·不跟随重定向这几个
因素后再检测一次。（截至 2026-10-06 确认，上面三个链接用这个方法都是 200，其中
pipc.go.kr 经过两次重定向后是 200。）

这套手动流程已经被 `scripts/check-external-links.py` 自动化了。它会从所有受跟踪的
`.md` 中收集代码围栏·行内代码之外的外部链接（排除预留·占位主机），用上面同样的方法
（浏览器 UA·GET·跟随重定向）确认状态，网络错误·5xx·429 会重试，用来区分瞬时抖动和
真正故障。结果分四类：OK（2xx·3xx）· RESTRICTED（401·403·405·429，主机存活只是检测
方式被挡）· ALLOWED（`scripts/external-links-allowlist.txt` 中记录的、我们无法修复的
上游继承失效链接）· DOWN（404·410·5xx·连接错误，很可能真的失效了）。

- 不发网络请求，只预览待检查目标：`python3 -I scripts/check-external-links.py --list`
- 发布·周期检查（出现新的失效链接就异常退出）：`python3 -I scripts/check-external-links.py --strict`

外部链接的存活状态不稳定，所以**不放进合并关卡**。而是由不挡合并的工作流
[`external-links`](.github/workflows/external-links.yml) 在每周一和手动触发时跑
`--strict`，allowlist 之外一旦出现新的 DOWN 就会变红曝出来。上游原文保留文件继承来的
失效链接（例如 `CHANGELOG.zh.md` 致谢的、已经注销的贡献者账号）我们无法修复，记录
在 allowlist 里并写明理由，从 strict 检查中排除。

---

## 9. 发布流水线

推送版本 tag（`v*`）会触发 [`.github/workflows/release.yml`](.github/workflows/release.yml)，
构建五个平台的二进制，并在满足条件时构建多架构 Docker 镜像。这个 fork 还没有打过
正式的发布 tag，这个工作流一次都没有真正运行过，所以这一节记录的是这条流水线预设了
什么前提、产出什么，以及这些前提是否与当前仓库结构相符。推送 tag 会在公开仓库上创建
GitHub Release，所以真正打 tag 要等发布决定确定之后再做。

### 9.1 如何打发布 tag

推送以 `v` 开头的 tag 就会触发工作流。

```bash
git tag v0.3.15
git push origin v0.3.15
```

### 9.2 流水线做了什么

工作流分成五个 job：

- **frontend。** 把前端静态导出一次（`web/out`），产物作为 `web-dist` artifact 上传。
  下面的 binaries job 对每个目标都会重新下载这份产物并复用。
- **binaries。** 交叉编译五个目标（linux amd64·arm64、darwin amd64·arm64、
  windows amd64），每个目标打包成 zip。linux amd64 的二进制会跑一遍
  `artex -h` 的冒烟测试，确认二进制真的能运行。
- **release。** 汇总所有 zip，生成 `SHA256SUMS` 校验文件，创建 GitHub Release 并把
  zip 和校验文件附上去。
- **docker-gate。** 检查 `DOCKERHUB_USERNAME`·`DOCKERHUB_TOKEN` 密钥是否已配置，把
  结果传给下一个 job 作为运行条件。
- **docker。** 只有上述密钥存在时才运行，拿 binaries 交叉编译好的 linux 二进制构建
  多架构镜像并推送到 Docker Hub。没有密钥就跳过这个 job，发布 CI 不会带红叉结束，
  只是只发布二进制。

### 9.3 构建前提是否与当前仓库结构相符

这条流水线建立在三个前提之上，已经在本地重现 binaries job，确认这些前提目前都与
仓库结构相符。

- **前端内嵌。** frontend job 上传的 `web-dist`（即 `web/out` 的内容），binaries job
  会下载到 `server/webui/dist`，`server/webui_embed.go` 的 `//go:embed all:webui/dist`
  把这个位置内嵌进二进制。所以 binaries job 不重新构建前端，调用
  [`build.sh`](build.sh) 时带 `ARTEX_SKIP_FRONTEND=1`。
- **二进制·打包路径。** `build.sh --target <os>/<arch>` 会生成
  `dist/artex-<os>-<arch>/artex` 二进制，并在 `dist/` 下生成 zip 包。zip 里除了二进制，
  还包含启动脚本（Linux·macOS 是 `start.sh`，Windows 是 `start.bat`）、`skills/`、
  `config.example.json`、`README.md`。
- **Docker 镜像里的二进制拷贝。** binaries job 会把 linux 二进制单独上传为
  `bin-linux-<arch>` artifact，docker job 下载后放到 `dist/<arch>/artex`。
  [`Dockerfile`](Dockerfile) 的 `COPY dist/${TARGETARCH}/artex`，在多架构构建中靠
  buildx 为每个平台填入的 `TARGETARCH` 取这个路径。[`.dockerignore`](.dockerignore)
  没有排除 `dist/`，所以二进制会被包含进构建上下文。

### 9.4 尚未决定的事项：Docker 镜像命名空间

docker job 目前把镜像名写死为上游的 `autumn27/artex`，这个 fork 要发布到哪个命名空间
是另一个待决定的事项（`work/DECISIONS-FOR-JIWOO.md` 第 8 项）。在决定之前先不配置
Docker Hub 密钥，这期间发布只会产出二进制 zip 和校验文件（docker job 会被跳过）。

### 9.5 不打 tag、在本地提前验证

如果只想确认流水线前提是否成立，又不想真的打一个公开发布 tag，可以在本地重现
binaries job。本地没有 Go 的话，可以用 Docker 跑同样的步骤。

```bash
# 1) 前端静态导出（对应 release.yml 的 frontend job）
cd web && npm ci && npm run build:static && cd ..
# 2) 放到 binaries job 下载产物的同一位置
rm -rf server/webui/dist && mkdir -p server/webui/dist && cp -a web/out/. server/webui/dist/
# 3) 用与 binaries job 相同的环境只构建一个目标
docker run --rm -v "$PWD":/app -w /app \
  -e ARTEX_SKIP_FRONTEND=1 -e ARTEX_SKIP_NPM_CI=1 \
  -e ARTEX_COMPRESS=0 -e ARTEX_PACKAGE=1 -e ARTEX_PACKAGE_DIR=dist \
  -e ARTEX_BUILD_VERSION=v0.0.0-local \
  golang:1.26 bash -c 'apt-get update && apt-get install -y zip && ./build.sh --target linux/amd64'
# 4) 确认产物：dist/artex-linux-amd64/artex · dist/*.zip · dist/SHA256SUMS
```

`dist/artex-linux-amd64/artex` 是静态链接的 ELF，带 `-h` 参数会打印用法后以退出码 0
结束，这正是 binaries job 冒烟测试确认的行为。构建产物（`dist/`·`server/webui/dist/`）
不提交到仓库（已在 `.gitignore` 中排除）。
