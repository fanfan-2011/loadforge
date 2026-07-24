#!/usr/bin/env bash
set -euo pipefail

REPO="fanfan-2011/loadforge"
BIN_NAME="loadforge"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${CYAN}╔══════════════════════════════════════╗${NC}"
echo -e "${CYAN}║     LoadForge 一键安装脚本            ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════╝${NC}"
echo ""

# --- Detect OS & Arch ---
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo -e "${YELLOW}⚠  不支持的架构: $ARCH，尝试从源码编译...${NC}"
    BUILD_FROM_SOURCE=true
    ;;
esac

case "$OS" in
  linux) ;;
  darwin) OS="darwin" ;;
  *)
    echo -e "${YELLOW}⚠  不支持的操作系统: $OS，尝试从源码编译...${NC}"
    BUILD_FROM_SOURCE=true
    ;;
esac

# --- Check if already installed ---
if command -v "$BIN_NAME" &>/dev/null; then
  EXISTING_PATH=$(command -v "$BIN_NAME")
  echo -e "${YELLOW}⚠  已检测到 LoadForge: $EXISTING_PATH${NC}"
  echo -e "   如需覆盖安装，重新运行本脚本即可。"
  echo ""
fi

# --- Check write permission ---
if [ ! -w "$INSTALL_DIR" ]; then
  USE_SUDO=true
else
  USE_SUDO=false
fi

install_binary() {
  local src="$1"
  local dst="$2"

  if [ "$USE_SUDO" = true ]; then
    sudo mv "$src" "$dst"
    sudo chmod +x "$dst"
  else
    mv "$src" "$dst"
    chmod +x "$dst"
  fi
}

# --- Try pre-built binary first (bypass DNS pollution) ---
if [ "${BUILD_FROM_SOURCE:-false}" = false ]; then
  # Use direct IP to bypass DNS pollution (GFW)
  GITHUB_IP="20.205.243.166"
  LATEST_URL="https://github.com/$REPO/releases/latest/download/loadforge-${OS}-${ARCH}.tar.gz"
  echo -e "${CYAN}📡 正在下载 LoadForge (${OS}-${ARCH})...${NC}"

  TMP_DIR=$(mktemp -d)
  trap "rm -rf '$TMP_DIR'" EXIT

  # Try direct connection first, fallback to IP with Host header
  if curl -fsSL --connect-timeout 10 --max-time 60 "$LATEST_URL" -o "$TMP_DIR/loadforge.tar.gz" 2>/dev/null; then
    : # download succeeded
  elif curl -fsSL --connect-timeout 10 --max-time 60 --resolve "github.com:443:$GITHUB_IP" \
    "https://github.com/$REPO/releases/latest/download/loadforge-${OS}-${ARCH}.tar.gz" \
    -o "$TMP_DIR/loadforge.tar.gz" 2>/dev/null; then
    : # download succeeded via IP
  else
    echo -e "${YELLOW}⚠  预编译二进制下载失败，尝试从源码编译...${NC}"
    BUILD_FROM_SOURCE=true
  fi

  if [ "${BUILD_FROM_SOURCE:-false}" = false ]; then
    tar -xzf "$TMP_DIR/loadforge.tar.gz" -C "$TMP_DIR" 2>/dev/null || {
      mv "$TMP_DIR/loadforge.tar.gz" "$TMP_DIR/loadforge" 2>/dev/null || true
    }

    if [ -f "$TMP_DIR/loadforge" ]; then
      echo -e "${CYAN}📦 正在安装到 $INSTALL_DIR/${BIN_NAME} ...${NC}"
      install_binary "$TMP_DIR/loadforge" "$INSTALL_DIR/$BIN_NAME"
      echo -e "${GREEN}✅ LoadForge 安装成功！${NC}"
      rm -rf "$TMP_DIR"
      trap - EXIT
      echo ""
      echo -e "  运行: ${CYAN}loadforge bench -n 1000 -c 10 https://example.com${NC}"
      echo ""
      exit 0
    fi
  fi

  trap - EXIT
  rm -rf "$TMP_DIR"
fi

# --- Build from source ---
echo -e "${CYAN}🔧 正在从源码编译 LoadForge ...${NC}"

# Check prerequisites
if ! command -v go &>/dev/null; then
  echo -e "${RED}❌ 需要安装 Go 1.18+${NC}"
  echo -e "   安装方法: https://go.dev/dl/"
  exit 1
fi

GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' || echo "0.0")
if [ "$(echo "$GO_VERSION" | cut -d. -f1)" -eq 0 ] || { [ "$(echo "$GO_VERSION" | cut -d. -f1)" -eq 1 ] && [ "$(echo "$GO_VERSION" | cut -d. -f2)" -lt 18 ]; }; then
  echo -e "${RED}❌ Go 版本过低: 需要 1.18+，当前 $GO_VERSION${NC}"
  exit 1
fi

if ! command -v npm &>/dev/null && ! command -v node &>/dev/null; then
  echo -e "${YELLOW}⚠  Node.js/npm 未安装，Web 报告功能将不可用${NC}"
  echo -e "   编译仍会继续（仅 CLI 模式）..."
  SKIP_UI=true
fi

# Clone or use existing
TMP_DIR=$(mktemp -d)
trap "rm -rf '$TMP_DIR'" EXIT

echo -e "${CYAN}📥 克隆仓库...${NC}"
# Bypass DNS pollution: use IP directly
git clone --depth 1 "https://github.com/$REPO.git" "$TMP_DIR/loadforge" \
  --config http.https://github.com.proxy="" 2>/dev/null || {
  # Fallback: use IP with Host header
  GIT_TERMINAL_PROMPT=0 GIT_SSL_NO_VERIFY=0 \
  git -c http.proxy= \
    clone --depth 1 \
    "https://20.205.243.166/$REPO.git" \
    "$TMP_DIR/loadforge" 2>/dev/null || {
    echo -e "${RED}❌ 克隆仓库失败（网络不通）${NC}"
    echo -e "   请检查网络连接或手动克隆:"
    echo -e "   git clone https://github.com/$REPO.git"
    exit 1
  }
}

cd "$TMP_DIR/loadforge"

# Build frontend
if [ "${SKIP_UI:-false}" = false ]; then
  echo -e "${CYAN}🎨 构建 Web 前端...${NC}"
  cd report/ui
  npm install --silent 2>/dev/null
  npm run build 2>/dev/null
  cd ../..
fi

# Build Go binary
echo -e "${CYAN}⚙️  编译 Go 二进制...${NC}"
go build -ldflags="-s -w" -o "$BIN_NAME" . 2>&1

if [ ! -f "$BIN_NAME" ]; then
  echo -e "${RED}❌ 编译失败${NC}"
  exit 1
fi

echo -e "${CYAN}📦 正在安装到 $INSTALL_DIR/${BIN_NAME} ...${NC}"
install_binary "$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

cd /
rm -rf "$TMP_DIR"
trap - EXIT

echo -e "${GREEN}✅ LoadForge 安装成功！${NC}"
echo ""
echo -e "  运行: ${CYAN}loadforge bench -n 1000 -c 10 https://example.com${NC}"
echo ""
