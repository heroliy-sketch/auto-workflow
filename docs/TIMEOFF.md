# auto-sky 多端项目交接文档

> 更新时间：2026-09-20（Asia/Shanghai）
>
> 本文写给完全没有当前会话上下文的新会话。先读取本文，再读取目标项目的真实 Git 状态；不要仅凭“目录存在”“workflow 存在”推断已经构建、部署或验收。

## 1. 当前任务

用户围绕同一个 `auto-sky` 产品建立了一组独立仓库：

1. 一个最小 Go 服务，作为后端和部署参考。
2. 一个 Next.js 网站客户端。
3. 一个 Android 客户端。
4. 一个 iOS 客户端。
5. 一个 HarmonyOS 客户端。
6. 一个专门替代 Figma 做静态产品原型和交互演示的前端工作台。
7. 之后又复制 Go 服务初始化为第二个独立 Go 服务 `auto-sky-one`。

共同边界：

- 这些仓库都是初始化工程，不是完整产品实现。
- 不复制 `frontier-server` 的业务代码、数据库、生产配置、第三方凭据或历史业务文档。
- 客户端当前使用本地静态内容和假数据，不调用真实后端 API。
- 后端部署文件保留 ECR/EKS/Helm/GitHub Actions 的结构，但没有把 AWS 生产参数写入仓库。
- 所有个人 GitHub 推送都使用个人 SSH key；不要用公司账号的 SSH key 误推送。

## 2. 项目总览

| 项目 | 本地路径 | GitHub | 当前分支 | 当前提交 | 主要职责 |
| --- | --- | --- | --- | --- | --- |
| Go backend | `/Users/devin/go/src/auto-sky` | [auto-sky](https://github.com/heroliy-sketch/auto-sky) | `dev` | `35f96a2` | 最小 Go HTTP 服务和 ECR/EKS/Helm 部署模板 |
| Go backend copy | `/Users/devin/go/src/auto-sky-one` | [auto-sky-one](https://github.com/heroliy-sketch/auto-sky-one) | `dev` | `11c7015` | 从 `auto-sky` 派生的独立最小 Go 服务 |
| Web client | `/Users/devin/frontend/auto-sky-website` | [auto-sky-website](https://github.com/heroliy-sketch/auto-sky-website) | `dev` | `008b03c` | Next.js/React/TypeScript 初始网站客户端 |
| Android client | `/Users/devin/frontend/auto-sky-ad` | [auto-sky-ad](https://github.com/heroliy-sketch/auto-sky-ad) | `dev` | `05d175f` | Kotlin/Compose Android 初始客户端 |
| iOS client | `/Users/devin/frontend/auto-sky-ios` | [auto-sky-ios](https://github.com/heroliy-sketch/auto-sky-ios) | `dev` | `ea67525` | Swift/SwiftUI iOS 初始客户端 |
| HarmonyOS client | `/Users/devin/frontend/auto-sky-hm` | [auto-sky-hm](https://github.com/heroliy-sketch/auto-sky-hm) | `dev` | `043fb43` | ArkTS/ArkUI/Stage/Hvigor 初始客户端 |
| Prototype workbench | `/Users/devin/frontend/auto-ui` | [auto-ui](https://github.com/heroliy-sketch/auto-ui) | `main` | `a7fe2b2` | 只做静态原型页面和交互，不接真实业务 |

最后一次状态核验：除 `auto-ui` 外的六个仓库均为 `dev...origin/dev`，且 `dev` 与各自 `main` 指向同一提交；`auto-ui` 仍为 `main...origin/main`，未创建 `dev`。七个工作区均无未提交的 tracked 文件改动。

## 3. Go 后端项目

### 3.1 auto-sky

路径：`/Users/devin/go/src/auto-sky`

这个 Go 仓库曾被误移动到 `/Users/devin/frontend/auto-sky`，现已按用户要求还原到 `/Users/devin/go/src/auto-sky`。移动保留了 `.git`、远程地址、提交历史和 `main` tracking；新会话不要再移动它。相关交接文档放在 Go 后端的 `/Users/devin/go/src/auto-sky/docs`。

模块：`github.com/heroliy-sketch/auto-sky`

核心内容：

- 标准库 HTTP 服务，默认端口 `8080`
- `GET /health` 返回 `{"service":"auto-sky","status":"ok"}`
- `GET /` 返回服务状态
- 未知路径返回 404
- `SIGINT`/`SIGTERM` 优雅退出，最长等待 10 秒
- YAML 配置包位于 `config/`
- `Makefile`、Dockerfile、Helm chart、GitHub Actions 已初始化

部署结构：

- GitHub Actions 构建并推送 ECR 镜像
- EKS kubeconfig 更新
- `deploy/deploy.sh` 执行 Helm upgrade/install
- Helm 包含 ConfigMap、Deployment、Service、Ingress、HPA、ServiceAccount 和 `/health` 探针
- 环境白名单只允许 `devin`、`uat`、`prod`
- 默认 namespace 是 `auto-sky-{environment}`

2026-09-18 在项目还原到 `/Users/devin/go/src/auto-sky` 后重新独立验证：

- `go mod tidy` 没有产生文件改动
- `make test` 通过
- `make build` 通过，生成 `bin/auto-sky`
- `go vet ./...` 通过
- `go test -race ./...` 通过
- `bash -n deploy/deploy.sh` 通过
- `AUTO_SKY_ENV=billy` 被环境白名单正确拒绝
- GitHub Actions、Helm values、Chart 和本地 config YAML 解析通过
- Helm 3.17.3 lint 通过，template 成功渲染 125 行 Kubernetes 资源
- 本地启动后，`GET /health` 和 `GET /` 返回 200，未知路径返回 404
- 敏感内容扫描未发现 `frontier-server` 生产配置、数据库凭据或第三方密钥
- Docker CLI 存在，但 Docker daemon 未运行，因此本次没有重新执行容器镜像构建

未完成：

- 没有真实 AWS/ECR/EKS 发布
- 没有配置个人仓库的 AWS Secrets/Variables
- 已创建远程 `dev` 和历史 `devin`；默认开发分支是 `dev`，当前 `dev` 与 `main` 指向同一提交；尚未创建 `uat`、`prod`
- 没有真实 GitHub Actions deploy 成功记录

### 3.2 auto-sky-one

路径：`/Users/devin/go/src/auto-sky-one`

模块：`github.com/heroliy-sketch/auto-sky-one`

这是从 `auto-sky` 当前最小版本复制并重命名的独立 Go 服务，不共享 Git 历史。以下标识已全部改名：

- 二进制：`auto-sky-one`
- 配置变量：`AUTO_SKY_ONE_CONFIG`
- Docker 镜像和 ECR 默认仓库：`auto-sky-one`
- Helm chart：`deploy/charts/auto-sky-one`
- Helm release：`auto-sky-one`
- namespace：`auto-sky-one-{environment}`
- Go module：`github.com/heroliy-sketch/auto-sky-one`

它保留和 `auto-sky` 一样的 `devin`、`uat`、`prod` 部署白名单和部署骨架。它同样没有真实 AWS/EKS 发布证据。

## 4. Web 前端项目

### 4.1 auto-sky-website

路径：`/Users/devin/frontend/auto-sky-website`

仓库参考：`/Users/devin/frontend/pond-website-new`

技术栈：

- TypeScript
- React 18
- Next.js 13.5.11
- Pages Router
- CSS
- npm + `package-lock.json`

已完成：

- `src/pages` 初始化结构
- `src/pages/index.tsx` 初始页面
- `src/pages/api/health.ts` 本地健康接口
- `_app.tsx`、`_document.tsx`
- 严格 TypeScript 和 `@/*` 路径别名
- ESLint、README、Next 配置

没有复制 `pond-website-new` 的 Pond 业务页面、钱包、登录、API 客户端、生产环境文件或第三方密钥。

已验证：

- `npm run typecheck` 通过
- `npm run lint` 通过
- `npm run build` 通过
- 本地页面 HTTP 200
- `/api/health` 返回 `{"service":"auto-sky-website","status":"ok"}`

安装依赖时本机 npm 证书链曾返回 `UNABLE_TO_GET_ISSUER_CERT_LOCALLY`，使用了一次临时 `--strict-ssl=false` 和临时 cache；没有修改全局 npm 配置，也没有把该选项写进项目。安装时 npm 报告 Next 13/ESLint 依赖链存在审计告警，初始化任务没有执行 `npm audit fix --force`。

### 4.2 auto-ui：静态原型工作台

路径：`/Users/devin/frontend/auto-ui`

用途：

- 替代 Figma 做产品早期静态原型和交互演示
- 后端需求讨论阶段，先在这里确认页面结构、文案、交互状态和流程
- 不接真实后端、不保存生产数据、不放 API key
- 原型页面、假数据和本地交互都集中在这个仓库

技术栈：

- React 19.3.0
- TypeScript 7.0.2
- Vite 8.3.0
- CSS

当前工作台已经包含：

- 左侧 screen 列表：`Overview`、`Details`、`Checkout`
- 中央静态预览画布
- Continue/Back 交互
- Prototype notice 可关闭
- 右侧 Inspector
- 本地状态和 fixture metadata

Vite 配置使用 `base: './'`。完成构建后可以直接打开：

```bash
cd /Users/devin/frontend/auto-ui
npm install
npm run build
open dist/index.html
```

也可以用：

```bash
npm run dev
npm run preview
```

已验证：

- `npm run typecheck` 通过
- `npm run build` 通过
- `dist/index.html` 使用 `./assets/...` 相对路径
- Vite dev server 返回 HTTP 200
- 构建产物和 `node_modules` 已忽略，未提交

## 5. Android、iOS、HarmonyOS 客户端

### 5.1 Android：auto-sky-ad

路径：`/Users/devin/frontend/auto-sky-ad`

注意：这个仓库最初在 `/Users/devin/go/src/auto-sky-ad` 创建，后来按用户要求移动到了 `/Users/devin/frontend/auto-sky-ad`。移动只改变了本地路径，`.git`、远程仓库、提交历史和 `main` tracking 都保持不变。新会话不要再使用旧路径。

技术栈：

- Kotlin `2.4.10`
- Android Gradle Plugin `9.4.0`
- Gradle `9.6.0`
- Jetpack Compose + Material 3
- Compose BOM `2026.09.00`
- compileSdk/targetSdk `37`
- minSdk `26`
- Gradle Kotlin DSL + Version Catalog

主要内容：

- `app` Android application module
- Kotlin Compose 初始页面
- Material 3 主题
- Debug/Release build types
- Unit Test 和 Compose Instrumented Test
- Gradle wrapper、README、`.gitignore`

未完成：

- 本机没有 JDK 17、Gradle 或 Android SDK
- 没有执行 `./gradlew test`、`assembleDebug` 或真机/模拟器测试
- 没有配置签名、包发布、后端 API、登录或持久化

README 中的预期验证命令：

```bash
./gradlew test
./gradlew lint
./gradlew assembleDebug
./gradlew connectedDebugAndroidTest
```

### 5.2 iOS：auto-sky-ios

路径：`/Users/devin/frontend/auto-sky-ios`

技术栈：

- Swift
- SwiftUI
- Xcode project format
- XCTest unit test 和 UI test

主要内容：

- `auto-sky-ios.xcodeproj`
- SwiftUI `@main` App entry point
- `ContentView`
- Assets catalog
- Debug/Release 配置
- 共享 scheme `auto-sky-ios`
- Unit Tests 和 UI Tests

已验证：

- Swift 文件 `swiftc -parse` 通过
- `project.pbxproj` plist 解析通过
- Git 检查通过

未完成：

- 当前 active developer directory 只有 Command Line Tools，没有完整 Xcode
- `xcodebuild` 无法执行
- 没有 iOS Simulator 构建、运行和 UI test 证据
- 没有 Team ID、签名证书、Provisioning Profile、API、认证或持久化

不要把 iOS 的“Swift 语法解析通过”写成“iOS app 已构建通过”。

### 5.3 HarmonyOS：auto-sky-hm

路径：`/Users/devin/frontend/auto-sky-hm`

技术栈：

- ArkTS
- ArkUI
- HarmonyOS NEXT Stage model
- Hvigor
- `AppScope` + `entry` HAP module

主要内容：

- 工程级 `build-profile.json5`
- 模块级 `entry/build-profile.json5`
- `hvigorfile.ts` 和 `hvigor/hvigor-config.json5`
- `AppScope/app.json5`
- `EntryAbility.ets`
- `pages/Index.ets`
- string/color/profile/media 资源
- Unit/OHOS test 占位

已验证：

- 配置文件可被标准 JSON 解析；本次配置没有依赖 JSON5 特有语法
- 入口、页面、资源和 module 配置引用存在
- Git whitespace 检查通过

未完成：

- 本机没有 DevEco Studio、HarmonyOS SDK 或真实 Hvigor 工具链
- 没有 HAP 编译、设备/模拟器运行或签名发布验证
- `hvigor/hvigor-config.json5` 的版本应与实际 DevEco Studio/SDK 配套版本核对；DevEco 可能要求同步调整

## 6. 验证状态矩阵

| 项目 | 本地验证 | 尚未验证 |
| --- | --- | --- |
| `auto-sky` | `go mod tidy`、`make test`、`make build`、`go vet`、`go test -race`、YAML、Bash、Helm lint/template、本地 `/health` | Docker 镜像实构建、AWS/ECR/EKS、真实 Actions、UAT/PROD |
| `auto-sky-one` | `make test`、`make build`、`go vet`、Bash、Helm lint/template、非法环境拒绝 | AWS/ECR/EKS、真实 Actions、UAT/PROD |
| `auto-sky-website` | npm typecheck、lint、build、页面 HTTP 200、API health | 真实后端、部署域名、登录、业务 API |
| `auto-ui` | npm typecheck、Vite build、dev server HTTP 200、相对 dist 资源 | 真实产品页面、API、用户验收流程 |
| `auto-sky-ad` | 静态文件和 Git 状态 | JDK/SDK、Gradle build、APK、设备/模拟器 |
| `auto-sky-ios` | Swift parse、pbxproj parse、Git 状态 | 完整 Xcode、xcodebuild、Simulator、签名 |
| `auto-sky-hm` | 配置 JSON parse、文件引用、Git 状态 | DevEco、Hvigor、HAP、设备/模拟器、签名 |

任何后续交接必须继续使用这个区分：代码存在和提交成功，不等于目标平台构建或部署成功。

## 7. Git、SSH 和目录规则

### 7.1 开发分支约定

2026-09-20 已基于各自最新 `main` 创建并推送默认开发分支 `dev`：

- `/Users/devin/go/src/auto-sky`
- `/Users/devin/go/src/auto-sky-one`
- `/Users/devin/frontend/auto-sky-website`
- `/Users/devin/frontend/auto-sky-ad`
- `/Users/devin/frontend/auto-sky-ios`
- `/Users/devin/frontend/auto-sky-hm`

以上六个工作区当前均停留在 `dev`，并跟踪 `origin/dev`。创建时 `dev` 与 `main` 指向完全相同的提交，没有额外 commit。后续日常开发默认在 `dev` 进行，`main` 保留为初始化稳定基线；合并前必须重新检查分支差异和目标仓库规则。

2026-09-18 创建的 `devin` 分支仍保留在这六个仓库，但不再是默认开发分支；不要在没有用户明确要求时删除或强制改写它。

`/Users/devin/frontend/auto-ui` 是纯静态原型工作台，按用户要求排除在本次分支创建之外，仍停留在 `main`，远程没有 `dev` 或 `devin` 分支。

两个 Go 仓库的 `.github/workflows/deploy.yml` 当前只监听 `devin`、`uat`、`prod`，不监听 `dev`。因此 `dev` 是代码开发分支，不是部署环境分支；创建和推送 `dev` 不代表任何环境已经部署。真实部署仍需核对 Actions 日志、AWS/ECR/EKS 配置和 Pod 状态。

### 7.2 SSH 和提交规则

所有新仓库使用个人 GitHub SSH key，并在本地仓库写入 repo-local 配置：

```text
core.sshCommand=ssh -o IdentitiesOnly=yes -i /Users/devin/.ssh/id_ed25519_github_personal
```

私钥路径只用于本机配置，不能上传、打印或写入 README：

```text
/Users/devin/.ssh/id_ed25519_github_personal
```

公开 `.pub` 公钥不等于可以公开私钥。以后在新目录初始化仓库时，先确认：

```bash
ssh -T -o IdentitiesOnly=yes -i /Users/devin/.ssh/id_ed25519_github_personal git@github.com
git status --short --branch
git remote -v
git ls-remote --heads origin
```

目录约定：

- Go 服务：`/Users/devin/go/src/auto-sky`、`/Users/devin/go/src/auto-sky-one`
- 所有前端和客户端：`/Users/devin/frontend/`
- 不要把 Android、iOS、HarmonyOS 客户端再放回 `/Users/devin/go/src/`
- 原型工作台只放 `/Users/devin/frontend/auto-ui`

安全提交流程：

```bash
git status --short --branch
git diff --check
git add <明确的文件>
git diff --cached --check
git commit -m "<message>"
git push origin main
```

不要使用 `git reset --hard` 或 `git checkout --` 清理来源不明的改动。不要复制源项目的 `.git`、`node_modules`、`.next`、`dist`、`bin`、`.idea`、`DerivedData`、`.hvigor` 或签名文件。

## 8. 旧项目规则和规范

`frontier-server` 的实际规则文件是：[AGENTS.md](/Users/devin/go/src/frontier-server/AGENTS.md)

修改旧项目时必须遵守：

- Frontier HTTP API response 禁止中文文本。
- 用户/产品侧使用 `task/tasks`，不要把 backend `bounty/bounties` 原样带入用户文案；技术引用、SQL 表列和原始数据库数据除外。
- 邮件正文要有正确换行、粗体强调和可点击导航链接。
- 邮件签名固定为 `Best,` 换行 `Pond Team`。
- 数据库集成测试按当前项目规则使用 SQLite memory，Redis 集成使用 miniredis；专项 Agent Billing 的真实 MySQL 并发要求只适用于对应专项测试，不要泛化到所有模块。
- Go 构建必须使用 `make build`，不要把直接 `go build` 当作项目构建流程。
- `internal/dao/ent/` 冲突先改 schema，再执行 `make ent`。
- 涉及 `uat`/`prod` merge 或 rebase 前，先更新目标分支到最新远程，冲突优先保留已部署行为。
- 只做用户请求对应的手术式改动，不顺手重构或删除无关代码。

对 `auto-sky` 和 `auto-sky-one`，部署环境只能是 `devin`、`uat`、`prod`。代码开发分支 `dev` 已存在，但不是部署环境；除非用户明确要求改变发布策略，不要把 `dev`、`billy`、`jerry` 或其他分支加入部署 workflow。

新项目会话还应读取当前目录的编码规范副本：[编码规范.md](/Users/devin/go/src/auto-sky/docs/编码规范.md)。该文件是从 [frontier-server/docs/编码规范.md](/Users/devin/go/src/frontier-server/docs/编码规范.md) 复制而来，当前两份内容一致；旧仓库文件是来源文件，新项目副本用于 `auto-sky` 相关工作。若来源规范发生变化，需要重新复制或明确同步差异。

## 9. 设计、接口和记忆文档位置

### 9.1 Frontier repository 文档

长期 Agent Billing 设计/API 文档目录：[docs/ai_agent/agent](/Users/devin/go/src/frontier-server/docs/ai_agent/agent/)

关键文件：

- [2026-08-11_agent_billing_api.md](/Users/devin/go/src/frontier-server/docs/ai_agent/agent/2026-08-11_agent_billing_api.md)：持续维护的 Billing HTTP/API 文档，后续接口变更优先更新这里。
- [2026-08-03_devin_agent_marketplace_payment_billing_final_design.md](/Users/devin/go/src/frontier-server/docs/ai_agent/agent/2026-08-03_devin_agent_marketplace_payment_billing_final_design.md)：Devin Agent Marketplace 支付计费最终设计。
- [2026-08-17_agent_billing_wallet_frontend_api.md](/Users/devin/go/src/frontier-server/docs/ai_agent/agent/2026-08-17_agent_billing_wallet_frontend_api.md)：Wallet 前端联调说明。
- [2026-08-19_agent_detail_plan_frontend_api.md](/Users/devin/go/src/frontier-server/docs/ai_agent/agent/2026-08-19_agent_detail_plan_frontend_api.md)：Agent 详情页套餐选择和支付联调说明。

不要把整个 `frontier-server/docs/` 目录复制到任何 `auto-sky` 项目；先确认文档是否与当前工作直接相关。

### 9.2 Codex memory

- [MEMORY.md](/Users/devin/.codex/memories/MEMORY.md)：长期记忆索引。
- [memory_summary.md](/Users/devin/.codex/memories/memory_summary.md)：汇总入口。
- [/Users/devin/.codex/memories/rollout_summaries/](/Users/devin/.codex/memories/rollout_summaries/)：历史 rollout 证据。
- [/Users/devin/.codex/memories/extensions/ad_hoc/notes/](/Users/devin/.codex/memories/extensions/ad_hoc/notes/)：增量记录。

这些记忆不是仓库代码，也不是 GitHub 内容。除非用户明确要求更新记忆，不要编辑 `MEMORY.md`；需要保留新结论时按记忆规则新增 note，并在最终回复说明。

当前 auto-sky 文档目录包含：

- [TIMEOFF.md](/Users/devin/go/src/auto-sky/docs/TIMEOFF.md)：本次多端项目交接文档。
- [编码规范.md](/Users/devin/go/src/auto-sky/docs/编码规范.md)：从 `frontier-server` 复制的编码规范副本。

本文件最终位于 `/Users/devin/go/src/auto-sky/docs/TIMEOFF.md`。`auto-sky/.gitignore` 包含 `docs/`，所以文档位于 Go 仓库目录内但不会出现在普通 `git status` 中，也不会自动提交到 GitHub；检查时使用 `ls`、`wc` 和内容读取。

## 10. 新会话启动检查清单

先确认所有工作区和提交状态：

```bash
for project in \
  /Users/devin/go/src/auto-sky \
  /Users/devin/go/src/auto-sky-one \
  /Users/devin/frontend/auto-sky-website \
  /Users/devin/frontend/auto-sky-ad \
  /Users/devin/frontend/auto-sky-ios \
  /Users/devin/frontend/auto-sky-hm \
  /Users/devin/frontend/auto-ui; do
  echo "### $project"
  git -C "$project" status --short --branch
  git -C "$project" log -1 --oneline
  git -C "$project" remote -v
done
git -C /Users/devin/go/src/frontier-server status --short --branch
```

开始新工作前必须回答：

- 用户是在修改后端、真实客户端，还是只做 `auto-ui` 静态原型？
- 当前项目的本地路径和 remote 是否正确？
- 是否有未提交的用户改动？
- 目标平台工具链是否真实存在？
- 代码通过了哪一层验证，哪些只是静态配置？
- 是否会触碰真实 API、签名、云环境或敏感数据？
- Go 部署环境是否仍只有 `devin`、`uat`、`prod`？
- 是否需要先更新对应设计/API文档，再修改实现？

## 11. 下一步建议

1. 先在 `auto-ui` 中把产品需求拆成静态 screens 和交互状态，完成评审后再决定后端 API。
2. 后端 API 确认后，更新 `auto-sky-website`、Android、iOS 和 HarmonyOS 的接口契约，再接入真实服务。
3. 为两个 Go 服务分别配置个人 GitHub Actions 的 AWS Secrets/Variables、ECR repository、EKS cluster、namespace 和 ingress host。
4. 配置 JDK/Android SDK 后执行 Android build；安装完整 Xcode 后执行 iOS build/test；安装 DevEco Studio/HarmonyOS SDK 后执行 Hvigor/HAP build。
5. 在所有客户端完成工具链验证前，不要声称“客户端可发布”或“已通过真机验收”。

## 12. 当前结论

已完成的范围是：两个最小 Go 服务、一个 Next.js 网站初始工程、一个 Android 初始工程、一个 iOS 初始工程、一个 HarmonyOS 初始工程，以及一个可直接打开静态构建产物的 `auto-ui` 原型工作台；所有代码均已提交到各自个人 GitHub 仓库。

除 `auto-ui` 外，六个产品仓库已经从各自 `main` 创建并推送默认 `dev` 开发分支，本地工作区均切换到 `dev` 并跟踪 `origin/dev`；历史 `devin` 分支保留，`auto-ui` 继续使用 `main`。

尚未完成的范围是：真实后端业务、统一 API 合约、登录/鉴权、数据持久化、签名发布、设备验收、AWS/EKS 部署和生产/UAT 验收。新会话必须把“初始化提交完成”和“产品/平台交付完成”严格分开。
