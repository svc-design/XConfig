#!/bin/bash

# 获取参数
DOMAIN="$1"
VALID_DAYS="$2"
OUTPUT_DIR="$3"

# 确保参数不为空
if [[ -z "$DOMAIN" || -z "$VALID_DAYS" || -z "$OUTPUT_DIR" ]]; then
  echo "Usage: $0 <domain_name> <valid_days> <output_dir>"
  exit 1
fi

# 确保输出目录存在
mkdir -p "$OUTPUT_DIR"

CERT_FILE="$DOMAIN.cert"
KEY_FILE="$DOMAIN.key"

echo "Generating certificate for domain: $DOMAIN with validity: $VALID_DAYS days"

# 生成 CA 私钥
openssl genrsa -out "$OUTPUT_DIR/ca.key" 2048

# 生成 CA 证书
openssl req -x509 -new -nodes -key "$OUTPUT_DIR/ca.key" -sha256 -days "$VALID_DAYS" -out "$OUTPUT_DIR/ca.cert" -subj "/C=CN/ST=State/L=City/O=Company/OU=Org/CN=Custom-CA"

# 生成服务器私钥
openssl genrsa -out "$OUTPUT_DIR/$KEY_FILE" 2048

# 生成 CSR（证书签名请求）
openssl req -new -key "$OUTPUT_DIR/$KEY_FILE" -out "$OUTPUT_DIR/$DOMAIN.csr" -subj "/C=CN/ST=State/L=City/O=Company/OU=Org/CN=$DOMAIN"

# 生成服务器证书
openssl x509 -req -in "$OUTPUT_DIR/$DOMAIN.csr" -CA "$OUTPUT_DIR/ca.cert" -CAkey "$OUTPUT_DIR/ca.key" -CAcreateserial -out "$OUTPUT_DIR/$CERT_FILE" -days "$VALID_DAYS" -sha256

# 清理 CSR 文件
rm -f "$OUTPUT_DIR/$DOMAIN.csr"

echo "SSL Certificates for $DOMAIN generated successfully in $OUTPUT_DIR!"

