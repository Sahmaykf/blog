# Simple Blog

这是一个基于 Go (Gin) + GORM 后端，Vue 3 + Vite 前端的简单博客系统。
仓库地址：https://github.com/Sahmaykf/blog

## 功能概览
- 文章发布/编辑/隐藏/删除
- 评论、点赞、用户关注
- 用户头像上传
- 后台管理（文章 / 评论 / 通知）

## 技术栈
- 后端：Go, Gin, GORM, MySQL
- 前端：Vue 3, Vite, Bootstrap 5

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

注意：仓库已包含 `.gitignore`，默认忽略 `.env`、`uploads/`、`frontend/dist`、构建产物等。

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

## 自动化部署（可选）
- 可以使用 GitHub Actions 将构建产物打包并通过 SSH 或 rsync 同步到服务器，或在服务器上使用 `git pull` + 重启服务的简单部署脚本。

## 致谢
如需我为你：
- 添加 GitHub Actions 自动部署示例
- 在仓库根目录添加 README 的更多细节或截图
- 生成 `systemd`/`nginx` 配置并在仓库中提供部署脚本

请回复你想让我继续完成的项。
