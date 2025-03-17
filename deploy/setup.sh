#!/bin/bash

# 创建服务用户
sudo useradd -r -s /bin/false dev-helper

# 创建应用目录
sudo mkdir -p /opt/dev-helper
sudo chown -R dev-helper:dev-helper /opt/dev-helper

# 复制服务文件
sudo cp dev-helper.service /etc/systemd/system/

# 重新加载systemd
sudo systemctl daemon-reload

# 启用并启动服务
sudo systemctl enable dev-helper
sudo systemctl start dev-helper

# 检查服务状态
sudo systemctl status dev-helper 