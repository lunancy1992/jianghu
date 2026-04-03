#!/bin/bash
set -e

# 江湖小报部署脚本
# 使用方法: ./deploy.sh

SERVER_IP="47.95.200.127"
DOMAIN="jianghuxiaobao.cn"
SERVER_USER="root"
APP_DIR="/var/www/jianghuxiaobao"
DATA_DIR="/var/www/jianghuxiaobao/data"

echo "=========================================="
echo "江湖小报部署脚本"
echo "服务器: ${SERVER_IP}"
echo "域名: ${DOMAIN}"
echo "=========================================="
echo ""

# 检查是否已有 SSH 配置
if [ ! -f ~/.ssh/config ]; then
    echo "请确保 SSH 配置已设置好"
    exit 1
fi

echo "步骤 1: 创建服务器目录..."
ssh ${SERVER_USER}@${SERVER_IP} "mkdir -p ${APP_DIR}/dist ${DATA_DIR}"

echo "步骤 2: 上传前端静态文件..."
scp -r app/dist/build/h5/* ${SERVER_USER}@${SERVER_IP}:${APP_DIR}/dist/

echo "步骤 3: 上传后端程序..."
scp server/jianghu-server ${SERVER_USER}@${SERVER_IP}:/usr/local/bin/

echo "步骤 4: 上传配置文件..."
scp server/deploy/config.production.yaml ${SERVER_USER}@${SERVER_IP}:/root/jianghu/

echo "步骤 5: 上传 Nginx 配置..."
scp server/deploy/nginx.conf ${SERVER_USER}@${SERVER_IP}:/etc/nginx/conf.d/jianghuxiaobao.conf

echo "步骤 6: 上传 systemd 服务..."
scp server/deploy/jianghu.service ${SERVER_USER}@${SERVER_IP}:/etc/systemd/system/jianghu-server.service

echo "步骤 7: 在服务器上启动服务..."
ssh ${SERVER_USER}@${SERVER_IP} << 'EOF'
systemctl daemon-reload
systemctl enable jianghu-server
systemctl start jianghu-server
nginx -t && nginx -s reload
systemctl status jianghu-server
EOF

echo ""
echo "=========================================="
echo "部署完成!"
echo "访问地址: http://${DOMAIN}"
echo "=========================================="
