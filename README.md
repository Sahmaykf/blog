# Simple Blog

这是一个基于 Go (Gin) + GORM 后端，Vue 3 + Vite 前端的简单博客系统。
仓库地址：https://github.com/Sahmaykf/blog

## 功能概览
- **文章管理**：发布、编辑、隐藏、删除、Markdown 支持
- **多级评论**：支持无限层级回复、@用户提及、评论删除
- **互动功能**：点赞、收藏、用户关注/粉丝系统
- **搜索系统**：全局搜索（支持文章标题、作者用户名、标签名分类检索）
- **置顶功能**：
  - **全站置顶**：管理员可设置首页全局置顶文章
  - **个人置顶**：用户可在个人主页置顶自己的文章
- **消息通知**：点赞、评论、关注实时通知及未读提醒
- **用户系统**：头像上传、个人资料修改、博客个性化命名
- **后台管理**：文章管理、评论审核、全站置顶控制

## 技术栈
- **后端**：Go 1.25, Gin, GORM (MySQL), JWT 认证
- **前端**：Vue 3, Vite, Vue Router, Axios, Bootstrap 5, Bootstrap Icons

## 先决条件
- 本地开发：Go >= 1.20（仓库中使用 go1.25 指定），Node.js + npm
- 服务器：Ubuntu（示例为 Ubuntu 22.04），MySQL，Nginx

## 环境变量 (.env)
在项目根目录放置 `.env`（不要提交到 Git）。示例：

```
DB_DSN=blog_user:your_password@tcp(127.0.0.1:3306)/simple_blog?charset=utf8mb4&parseTime=True&loc=Local
APP_PORT=8080
UPLOAD_DIR=uploads
JWT_SECRET=your_jwt_secret
VITE_API_BASE_URL=/api
```

注意：仓库已包含 `.gitignore`，默认忽略 `.env`、`uploads/`、`frontend/dist`、构建产物等。**请务必不要将包含敏感信息的 `.env` 文件提交至公开仓库。**

## 核心页面说明
- **首页 (`/`)**：展示最新文章、全站置顶文章、热门文章及标签导航。
- **搜索页 (`/search`)**：分类展示文章、用户和标签的搜索结果。
- **文章详情 (`/post/:id`)**：支持 Markdown 渲染、多级评论互动、点赞收藏。
- **个人主页 (`/user/:id`)**：展示用户文章（含个人置顶）、关注/粉丝列表、获赞统计。
- **后台管理 (`/admin`)**：集中管理个人文章、评论及系统级置顶。

## 本地开发

1. 后端（开发模式）

```bash
# 在项目根目录
go run cmd/main.go
```

2. 前端（开发模式）

```bash
cd frontend
npm install
npm run dev
```

前端默认通过 Vite 的代理或 `VITE_API_BASE_URL` 访问后端 API（开发时可使用 `http://localhost:8080`）。

## 生产构建 & 部署（简要步骤）

1. 在本地交叉编译后端（生成 Linux 可执行文件）：

```powershell
# $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o blog-api cmd/main.go
```

2. 打包前端：

```bash
cd frontend
npm run build
# 生成 dist 文件夹
```

3. 将 `blog-api`、`.env`、`frontend/dist` 上传到服务器（示例目录 `/var/www/blog`）：

```bash
scp blog-api user@your_server:/home/you/
scp .env user@your_server:/home/you/
scp -r frontend/dist user@your_server:/home/you/
```

4. 服务器上准备目录并设置权限：

```bash
sudo mkdir -p /var/www/blog/uploads
sudo mv /home/you/blog-api /var/www/blog/
sudo mv /home/you/.env /var/www/blog/
sudo mv /home/you/dist /var/www/blog/dist
sudo chown -R www-data:www-data /var/www/blog
sudo chmod -R 755 /var/www/blog
```

5. 使用 Systemd 管理后端（示例 `/etc/systemd/system/blog.service`）：

```ini
[Unit]
Description=Simple Blog Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/blog
ExecStart=/var/www/blog/blog-api
Restart=always

[Install]
WantedBy=multi-user.target
```

启动并开机自启：

```bash
sudo systemctl daemon-reload
sudo systemctl start blog
sudo systemctl enable blog
```

6. Nginx 配置（替换 `/etc/nginx/sites-available/default` 的 `server` 块）：

```nginx
server {
    listen 80;
    server_name YOUR_SERVER_IP_OR_DOMAIN;

    location / {
        root /var/www/blog/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /uploads {
        alias /var/www/blog/uploads;
        add_header Access-Control-Allow-Origin *;
    }
}
```

检查并重启 Nginx：

```bash
sudo nginx -t
sudo systemctl restart nginx
```

## 常见问题与排查
- 403 Forbidden：通常是权限问题，执行 `sudo chown -R www-data:www-data /var/www/blog` 并检查 `root` 是否正确指向 `dist`。
- SSH 登录超时/被拒绝：检查腾讯云安全组是否放行端口 22/80，重置实例密码或使用控制台 Web 登录。
- 数据库连接失败：检查 `.env` 中的 `DB_DSN`、MySQL 是否已创建数据库与用户，并 `FLUSH PRIVILEGES`。

