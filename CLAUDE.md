# DecisionHelper-Go

DecisionHelper 的 Go 重构版本，使用 Wails 框架实现跨平台桌面应用。目标平台：Windows 和 macOS。

原 Python 版本：https://github.com/LUNARKN1GHT/DecisionHelper

## 技术栈

- Go 1.21+
- Wails v2（Go 后端 + Web 前端）
- Vanilla JS / HTML / CSS（前端，无需框架）
- JSON 文件本地存储（格式与 Python 版本兼容）

## 项目结构

- `main.go` 入口，启动 Wails 应用
- `app.go` 后端逻辑，暴露给前端调用的方法
- `models.go` 数据结构定义
- `storage.go` 读写本地 JSON
- `frontend/` 前端界面
  - `index.html` 主页面
  - `src/main.js` 前端逻辑
  - `src/style.css` 样式

## 核心数据结构

与 Python 版本 JSON 格式保持兼容：

```json
{
  "id": "uuid",
  "title": "选 offer",
  "created_at": "2026-05-21T10:00:00",
  "options": ["公司A", "公司B"],
  "criteria": [{"id": "uuid", "name": "薪资", "weight": 3}],
  "scores": [{"option": "公司A", "criterion_id": "uuid", "value": 4}]
}
```

加权总分公式：`sum(score * weight) / sum(weight)`，结果在 1–5 区间。

## MVP 功能范围

1. 新建 / 删除决策
2. 添加 / 编辑 / 删除选项和标准（权重 1–5）
3. 矩阵评分界面（下拉框选分）
4. 结果页：加权总分和排名（手动点击查看）
5. 历史记录持久化（JSON）

## 数据存储路径

- macOS：`~/Library/Application Support/DecisionHelper/decisions.json`
- Windows：`%APPDATA%\DecisionHelper\decisions.json`

## 开发环境

- Go 1.21+（`brew install go`）
- Wails v2（`go install github.com/wailsapp/wails/v2/cmd/wails@latest`）
- Node.js（Wails 构建前端依赖）
- PATH 需包含 Go bin：`export PATH=$PATH:$(go env GOPATH)/bin`
- 开发启动：`wails dev`
- 构建：`wails build`

## Git 工作流

- 主线：`main`，保持可运行状态
- 功能分支：`feat/<name>`
- 修复分支：`fix/<name>`
- Commit message 遵循 Conventional Commits：
  - `feat: 实现矩阵评分界面`
  - `fix: 修复权重计算错误`
  - `chore: 更新依赖`
- Git 操作由 Claude 直接执行

## 当前状态

Wails 项目已初始化，尚未开始功能开发。下一步：实现 `models.go` 和 `storage.go`。
