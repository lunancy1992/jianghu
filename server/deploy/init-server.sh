#!/bin/bash
# 服务器初始化脚本
# 在阿里云 ECS 上运行此脚本

set -e

echo "=== 江湖小报服务器初始化 ==="

# 安装必要软件
echo "安装必要软件..."
yum install -y nginx git

# 创建目录
echo "创建目录..."
mkdir -p /var/www/jianghuxiaobao/dist
mkdir -p /var/www/jianghuxiaobao/data
mkdir -p /root/jianghu
mkdir -p /var/log/jianghu

# 设置 Nginx 开机自启
echo "配置 Nginx..."
systemctl enable nginx
systemctl start nginx

# 创建默认配置
echo "创建默认配置..."
cat > /root/jianghu/config.yaml << 'EOF'
server:
  port: 8080
  data_dir: /var/www/jianghuxiaobao/data

database:
  path: /var/www/jianghuxiaobao/data/jianghu.db

ai:
  deepseek_api_key: ""
  deepseek_base_url: "https://api.deepseek.com"
  model: "deepseek-chat"

auth:
  jwt_secret: "$(openssl rand -hex 32)"
  jwt_expire_hours: 168
  sms_provider: "stub"

crawl:
  interval_minutes: 30
  feeds:
    - name: "36kr"
      url: "https://36kr.com/feed"
      enabled: true

cache:
  max_cost: 268435456
  num_counters: 10000000
EOF

# 配置防火墙
echo "配置防火墙..."
firewall-cmd --permanent --add-port=80/tcp 2>/dev/null || true
firewall-cmd --permanent --add-port=443/tcp 2>/dev/null || true
firewall-cmd --reload 2>/dev/null || true

echo ""
echo "=== 初始化完成 ==="
echo "接下来请在 GitHub 仓库设置 Secrets:"
echo "  SERVER_HOST: 47.95.200.127"
echo "  SERVER_USER: root"
echo "  SSH_PRIVATE_KEY: (你的 SSH 私钥)"
echo ""
echo "然后推送代码到 main 分支即可自动部署"
