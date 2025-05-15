#!/bin/bash
set -e

REMOTE_HOST="root@10.1.3.179"

if [ -z "$1" ]; then
  echo "❌ 用法: $0 <image>"
  echo "示例: $0 dfcloud-image-registry-vpc.cn-beijing.cr.aliyuncs.com/dev/df-web-ai:v6.6.18839"
  exit 1
fi

IMAGE="$1"

# 提取镜像名和版本号
NAME_TAG="${IMAGE##*/}"             # df-web-ai:v6.6.18839
NAME="${NAME_TAG%%:*}"              # df-web-ai
TAG="${NAME_TAG##*:}"               # v6.6.18839
FILE_NAME="${NAME}-${TAG//v/}.tar"  # df-web-ai-6.6.18839.tar

echo "📦 镜像: $IMAGE"
echo "📁 导出文件名: $FILE_NAME"

echo "🚀 在远程拉取镜像..."
ssh $REMOTE_HOST docker pull "$IMAGE"

echo "💾 在远程保存镜像为 /tmp/$FILE_NAME..."
ssh $REMOTE_HOST "docker save $IMAGE > /tmp/$FILE_NAME"

echo "📥 拷贝镜像回本地 ~/Desktop..."
scp $REMOTE_HOST:/tmp/$FILE_NAME ~/Desktop

echo "✅ 完成！镜像保存于：~/Desktop/$FILE_NAME"
