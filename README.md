# DecisionHelper

一个基于加权评分矩阵的决策辅助工具，支持 AI 辅助分析。使用 Go + Wails 构建，跨平台运行于 macOS 和 Windows。

> 原 Python 版本：[LUNARKN1GHT/DecisionHelper](https://github.com/LUNARKN1GHT/DecisionHelper)

## 功能

- **新建 / 删除决策**
- **管理选项与标准**：自由添加候选选项和评价标准，每个标准可设置权重（1–5）
- **矩阵评分**：通过表格对每个「选项 × 标准」组合打分（1–5）
- **加权结果排名**：自动计算加权总分并按名次展示
- **历史持久化**：数据存储于本地 JSON 文件，重启不丢失
- **AI 辅助**（需配置 API Key）：
  - ✦ AI 建议标准：根据决策场景自动推荐评价维度
  - ✦ AI 辅助评分：对所有组合给出建议分值和理由
  - ✦ AI 分析结果：基于评分数据生成叙述性决策建议

## 下载

前往 [Releases](https://github.com/LUNARKN1GHT/DecisionHelper-Go/releases) 页面下载对应平台的最新版本：

| 平台                                          | 文件                                              |
| --------------------------------------------- | ------------------------------------------------- |
| macOS（Universal，支持 Intel & Apple Silicon） | `DecisionHelper-Go-vX.X.X-macOS-Universal.zip`   |
| Windows x64                                   | `DecisionHelper-Go-vX.X.X-Windows-x64.exe`       |

## AI 配置

程序内置对 **DeepSeek** 及所有 OpenAI 兼容接口的支持（如 OpenAI、Ollama 等）。

1. 点击主界面右上角 **⚙** 打开设置
2. 填入 API 接入地址、API Key 和模型名称
3. 保存后即可在编辑、评分、结果页使用 AI 功能

> API Key 仅存储于本地，不会上传至任何外部服务器。

**DeepSeek 默认配置：**

```text
接入地址：https://api.deepseek.com
模型：    deepseek-chat
```

## 数据存储路径

| 平台    | 路径                                             |
| ------- | ------------------------------------------------ |
| macOS   | `~/Library/Application Support/DecisionHelper/` |
| Windows | `%APPDATA%\DecisionHelper\`                     |

## 本地开发

**环境依赖：**

- Go 1.21+
- Wails v2：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Node.js 18+

```bash
# 克隆项目
git clone https://github.com/LUNARKN1GHT/DecisionHelper-Go.git
cd DecisionHelper-Go

# 启动开发模式（热重载）
wails dev

# 构建生产包
wails build
```

## 加权得分公式

$$\text{score} = \frac{\sum(\text{分值} \times \text{权重})}{\sum\text{权重}}$$

结果区间为 1–5，分数越高排名越靠前。

## License

[MIT](LICENSE)
