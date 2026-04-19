# GPT-Load Mini

基于 [tbphp/gpt-load](https://github.com/tbphp/gpt-load) 开发的轻量级 AI API 透明代理，采用 MIT License。

## 功能特性

- **API Key 管理** - 添加、编辑、删除、导入导出多个 OpenAI API Key
- **Key 轮换** - 基于 Redis 原子操作实现分组自动轮换
- **自动故障转移** - Key 连续失败 5 次后自动禁用
- **透明代理** - 转发请求到上游 API，支持重试逻辑
- **管理后台** - Vue 3 Web UI，用于管理 Key 和分组

## API 使用说明

### URL

```
http://localhost:8080/proxy/<分组名称>
```


### 请求格式

```bash
# Chat Completions
curl -X POST http://localhost:8080/proxy/mygroup/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MY_API_KEY" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "Hello"}]
  }'

# Embeddings
curl -X POST http://localhost:8080/proxy/mygroup/v1/embeddings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MY_API_KEY" \
  -d '{
    "model": "text-embedding-3-small",
    "input": "Hello world"
  }'
```

### 路径映射

路径转换逻辑：
- 移除 UpstreamURL 末尾的 `/v1` 后缀和所有尾部斜杠
- 移除请求路径开头的斜杠
- 拼接：`UpstreamURL(去/v1) + "/" + 请求路径`

示例：假设分组 UpstreamURL 为 `https://api.openai.com/v1`

| 请求路径 | 上游路径 |
|---------|---------|
| `/proxy/mygroup/v1/chat/completions` | `https://api.openai.com/v1/chat/completions` |
| `/proxy/mygroup/v1/embeddings` | `https://api.openai.com/v1/embeddings` |

### 认证说明

- 如果分组未设置 ProxyAPIKey：无需认证
- 如果分组设置了 ProxyAPIKey：需要 header `Authorization: Bearer <ProxyAPIKey>`

## 项目结构

```
gpt-load-mini/
├── main.go                     # Go 后端入口
├── internal/
│   ├── api/
│   │   ├── handler/           # HTTP 处理器
│   │   └── router/            # 路由和中间件
│   ├── core/
│   │   ├── keypool/           # Key 轮换
│   │   └── proxy/              # 透明代理
│   ├── data/
│   │   ├── model/              # GORM 模型
│   │   ├── db/                 # MySQL 连接
│   │   └── store/              # 存储接口 + Redis
│   └── pkg/                    # 工具函数
└── frontend/                   # Vue 3 管理后台
```

## 环境要求

- Go 1.24+
- MySQL 8.0+
- Redis 7+

## 部署

### 手动部署

```bash
# 后端
go run main.go

# 前端
cd frontend
npm install
npm run dev      # 开发模式
npm run build    # 生产构建
```

## 配置说明

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `SERVER_PORT` | 后端服务端口 | `8080` |
| `AUTH_KEY` | 管理后台认证密钥 | `change-me-secret` |
| `MYSQL_HOST` | MySQL 主机 | `localhost` |
| `MYSQL_PORT` | MySQL 端口 | `3306` |
| `MYSQL_USER` | MySQL 用户名 | `root` |
| `MYSQL_PASSWORD` | MySQL 密码 | `123456` |
| `MYSQL_DATABASE` | MySQL 数据库名 | `gptload` |
| `REDIS_HOST` | Redis 主机 | `localhost` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `ENCRYPTION_KEY` | 32 字符 AES-256-CFB 加密密钥 | - |

## 管理接口

所有接口需要 header `X-Auth-Key: <AUTH_KEY>`

### 分组管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/groups` | 创建分组 |
| GET | `/api/groups` | 获取分组列表 |
| GET | `/api/groups/:id` | 获取分组详情 |
| PUT | `/api/groups/:id` | 更新分组 |
| DELETE | `/api/groups/:id` | 删除分组 |

### Key 管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/keys` | 添加 Key |
| GET | `/api/keys` | 获取 Key 列表 |
| PUT | `/api/keys/:id` | 更新 Key |
| DELETE | `/api/keys/:id` | 删除 Key |
| POST | `/api/keys/:id/restore` | 恢复已禁用 Key |
| GET | `/api/keys/export` | 导出 Keys |
| POST | `/api/keys/import` | 导入 Keys |

### 系统

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/health` | 健康检查（无需认证） |
| GET | `/api/stats` | 获取统计信息 |
| GET | `/api/logs` | 获取请求日志 |
| POST | `/api/admin/reload-config` | 重载配置 |

## 数据模型

### Group (分组)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 唯一标识名 |
| display_name | string | 显示名称 |
| channel_type | string | 渠道类型 (openai/azure) |
| upstream_url | string | 上游 API 地址 |
| test_model | string | 测试用模型名 |
| sort | int | 排序权重 |

### APIKey

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| group_id | uint | 所属分组 |
| key_value | string | 加密后的 Key |
| key_hash | string | Key 哈希（用于校验） |
| status | string | active/invalid |
| failure_count | int | 连续失败次数 |
| last_used_at | datetime | 最后使用时间 |

### RequestLog

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 主键 |
| timestamp | datetime | 请求时间 |
| group_id | int | 分组 ID |
| group_name | string | 分组名称 |
| key_id | int | 使用的 Key ID |
| model | string | 请求模型 |
| is_success | bool | 是否成功 |
| source_ip | string | 请求来源 IP |
| status_code | int | 响应状态码 |
| request_path | string | 请求路径 |
| duration_ms | int | 耗时（毫秒） |
| error_message | string | 错误信息 |

## 核心特性

- **Key 轮换**: 基于 Redis LMOVE 的原子轮换，保证高并发下的线程安全
- **失败追踪**: 连续失败 5 次自动禁用 Key，防止无效请求浪费资源
- **重试逻辑**: 失败后自动重试最多 3 次，每次选择不同的 Key
- **加密存储**: AES-256-CFB 加密 API Key，确保安全性
- **管理后台**: 支持深色/浅色主题切换，提供直观的统计面板

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/your-repo/gpt-load-mini.git
cd gpt-load-mini
```

### 2. 配置环境

复制 `.env.example` 为 `.env`，修改必要的配置：

```bash
cp .env.example .env
```

关键配置项：

| 变量 | 说明 | 必填 |
|------|------|------|
| `AUTH_KEY` | 管理后台认证密钥 | 是 |
| `ENCRYPTION_KEY` | 32 字符加密密钥 | 是 |
| `MYSQL_*` | MySQL 连接信息 | 是 |
| `REDIS_*` | Redis 连接信息 | 是 |

### 3. 启动服务

```bash
go run main.go
```

访问 `http://localhost:8080` 打开管理后台。

## License

本项目基于 [MIT License](https://opensource.org/licenses/MIT) 开源。

```
MIT License

Copyright (c) 2024 GPT-Load Mini Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
