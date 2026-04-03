# 江湖小报部署指南

## 架构

```
GitHub Actions ── SSH ── 阿里云 ECS
                    │
                    ├── Nginx (80)
                    │
                    └── jianghu-server (8080)
```

## 第一步：服务器初始化

SSH 登录服务器执行：

```bash
# 安装必要软件
yum install -y nginx

# 创建目录
mkdir -p /var/www/jianghuxiaobao/dist
mkdir -p /var/www/jianghuxiaobao/data
mkdir -p /root/jianghu
mkdir -p /var/log/jianghu

# 启动 Nginx
systemctl enable nginx
systemctl start nginx

# 开放防火墙端口
firewall-cmd --permanent --add-port=80/tcp
firewall-cmd --reload
```

## 第二步：配置 GitHub Secrets

在 GitHub 仓库设置 → Secrets and variables → Actions → New repository secret

| Secret 名称 | 值 |
|-------------|-----|
| `SERVER_HOST` | `47.95.200.127` |
| `SERVER_USER` | `root` |
| `SSH_PRIVATE_KEY` | 你的 SSH 私钥内容 |

### 生成 SSH 密钥（如果没有）

```bash
# 本地生成密钥
ssh-keygen -t ed25519 -C "github-actions" -f github-deploy-key

# 将公钥添加到服务器
ssh-copy-id -i github-deploy-key.pub root@47.95.200.127

# 将私钥内容添加到 GitHub Secrets
cat github-deploy-key
```

## 第三步：配置域名解析

在域名服务商添加 A 记录：

| 类型 | 主机记录 | 记录值 |
|------|----------|--------|
| A | @ | 47.95.200.127 |

## 第四步：推送代码部署

```bash
git add .
git commit -m "Setup CI/CD deployment"
git push origin main
```

## 部署文件说明

```
server/deploy/
├── config.production.yaml  # 生产环境配置
├── nginx.conf             # Nginx 配置
├── jianghu.service       # systemd 服务配置
└── init-server.sh        # 服务器初始化脚本

.github/workflows/
└── deploy.yml            # GitHub Actions 工作流
```

## 常用命令

```bash
# 查看服务状态
ssh root@47.95.200.127 'systemctl status jianghu-server'

# 查看日志
ssh root@47.95.200.127 'journalctl -u jianghu-server -f'

# 重启服务
ssh root@47.95.200.127 'systemctl restart jianghu-server'

# 手动部署
cd /Users/bytedance/code/jianghu
./deploy.sh
```

## 注意事项

1. **JWT Secret**: 首次部署时会自动生成随机 JWT secret
2. **短信服务**: 如需真实短信，需配置阿里云短信 AccessKey
3. **数据库**: SQLite 数据库文件在 `/var/www/jianghuxiaobao/data/jianghu.db`
4. **日志**: 服务日志通过 journalctl 查看
