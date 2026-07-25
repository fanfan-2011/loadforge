#!/usr/bin/env bash
set -euo pipefail

REPO="fanfan-2011/loadforge"
GITEE_REPO="fan-haoran-01/loadforge"
BIN_NAME="loadforge"
INSTALL_DIR="/usr/local/bin"
VERSION="v1.0.0"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

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

# ============================================================
#  多源下载：依次尝试，一个成功即停止
# ============================================================
download_and_install() {
  local TMP_DIR
  TMP_DIR=$(mktemp -d)
  trap "rm -rf '$TMP_DIR'" EXIT

  local FILENAME="loadforge-${OS}-${ARCH}"
  local GZ_FILE="${FILENAME}.gz"

  # ---- 源列表 ----
  # 依次尝试以下 URL，成功即安装退出
  local URLS=()

  # 源1: GitHub Releases
  URLS+=("https://github.com/$REPO/releases/download/$VERSION/$GZ_FILE")

  # 源2: jsDelivr CDN (从 repo dist/ 目录分发，国内速度快)
  URLS+=("https://cdn.jsdelivr.net/gh/$REPO@$VERSION/dist/$GZ_FILE")

  # 源3: raw.githubusercontent.com (不同 CDN IP，部分网络可用)
  URLS+=("https://raw.githubusercontent.com/$REPO/$VERSION/dist/$GZ_FILE")

  # 源4: Gitee raw 文件（直接用仓库里的 dist/ 目录）
  URLS+=("https://gitee.com/$GITEE_REPO/raw/main/dist/$GZ_FILE")

  local SUCCESS=false
  for url in "${URLS[@]}"; do
    echo -e "${CYAN}📡 尝试下载: ${url}${NC}"
    if curl -fsSL --connect-timeout 10 --max-time 60 "$url" -o "$TMP_DIR/$GZ_FILE" 2>/dev/null; then
      echo -e "${GREEN}✅ 下载成功${NC}"
      # 解压（.gz 文件用 gunzip 或 gzip -d）
      if command -v gunzip &>/dev/null; then
        gunzip -c "$TMP_DIR/$GZ_FILE" > "$TMP_DIR/$FILENAME" 2>/dev/null || cp "$TMP_DIR/$GZ_FILE" "$TMP_DIR/$FILENAME"
      else
        gzip -d -c "$TMP_DIR/$GZ_FILE" > "$TMP_DIR/$FILENAME" 2>/dev/null || cp "$TMP_DIR/$GZ_FILE" "$TMP_DIR/$FILENAME"
      fi

      if [ -f "$TMP_DIR/$FILENAME" ] && [ -s "$TMP_DIR/$FILENAME" ]; then
        echo -e "${CYAN}📦 正在安装到 $INSTALL_DIR/$BIN_NAME ...${NC}"
        install_binary "$TMP_DIR/$FILENAME" "$INSTALL_DIR/$BIN_NAME"
        echo -e "${GREEN}✅ LoadForge 安装成功！${NC}"
        rm -rf "$TMP_DIR"
        trap - EXIT
        SUCCESS=true
        break
      fi
    fi
    echo -e "${YELLOW}⚠  该源不可用，尝试下一个...${NC}"
  done

  if [ "$SUCCESS" = false ]; then
    trap - EXIT
    rm -rf "$TMP_DIR"
    return 1
  fi
  return 0
}

# --- 尝试预编译二进制下载 ---
if [ "${BUILD_FROM_SOURCE:-false}" = false ]; then
  if download_and_install; then
    echo ""
    echo -e "  运行: ${CYAN}$BIN_NAME bench -n 1000 -c 10 https://example.com${NC}"
    echo ""
    exit 0
  fi
  echo -e "${YELLOW}⚠  所有下载源均不可用，尝试从源码编译...${NC}"
fi

# ============================================================
#  源码编译（兜底）
# ============================================================
echo -e "${CYAN}🔧 正在从源码编译 LoadForge ...${NC}"

# 检查 Go
if ! command -v go &>/dev/null; then
  echo -e "${RED}❌ 需要安装 Go 1.18+${NC}"
  echo -e "   安装方法: https://go.dev/dl/"
  echo ""
  echo -e "${YELLOW}💡 也可以手动下载预编译二进制:${NC}"
  echo -e "   https://github.com/$REPO/releases/tag/$VERSION"
  echo -e "   或从 Gitee 下载:"
  echo -e "   https://gitee.com/$GITEE_REPO/blob/main/dist/loadforge-linux-${OS}-${ARCH}.gz"
  exit 1
fi

GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' || echo "0.0")
if [ "$(echo "$GO_VERSION" | cut -d. -f1)" -eq 0 ] || { [ "$(echo "$GO_VERSION" | cut -d. -f1)" -eq 1 ] && [ "$(echo "$GO_VERSION" | cut -d. -f2)" -lt 18 ]; }; then
  echo -e "${RED}❌ Go 版本过低: 需要 1.18+，当前 $GO_VERSION${NC}"
  exit 1
fi

# 设置国内代理（如果可用）
export GOPROXY=https://goproxy.cn,https://goproxy.io,direct

# 检查 Node.js
if ! command -v npm &>/dev/null && ! command -v node &>/dev/null; then
  echo -e "${YELLOW}⚠  Node.js/npm 未安装，Web 报告功能将不可用${NC}"
  SKIP_UI=true
fi

# 克隆仓库
TMP_DIR=$(mktemp -d)
trap "rm -rf '$TMP_DIR'" EXIT

echo -e "${CYAN}📥 克隆仓库 (GitHub)...${NC}"
git clone --depth 1 "https://github.com/$REPO.git" "$TMP_DIR/loadforge" 2>/dev/null || {
  echo -e "${YELLOW}⚠  GitHub 克隆失败，尝试 Gitee...${NC}"
  git clone --depth 1 "https://gitee.com/$GITEE_REPO.git" "$TMP_DIR/loadforge" 2>/dev/null || {
    echo -e "${RED}❌ 仓库克隆失败${NC}"
    exit 1
  }
}

cd "$TMP_DIR/loadforge"

# 构建前端
if [ "${SKIP_UI:-false}" = false ]; then
  echo -e "${CYAN}🎨 构建 Web 前端...${NC}"
  cd report/ui
  npm config set registry https://registry.npmmirror.com 2>/dev/null || true
  npm install --silent 2>/dev/null
  npm run build 2>/dev/null
  cd ../..
fi

# 编译
echo -e "${CYAN}⚙️  编译 Go 二进制...${NC}"
go build -ldflags="-s -w" -o "$BIN_NAME" . 2>&1

if [ ! -f "$BIN_NAME" ]; then
  echo -e "${RED}❌ 编译失败${NC}"
  exit 1
fi

echo -e "${CYAN}📦 正在安装到 $INSTALL_DIR/$BIN_NAME ...${NC}"
install_binary "$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

cd /
rm -rf "$TMP_DIR"
trap - EXIT

echo -e "${GREEN}✅ LoadForge 安装成功！${NC}"
echo ""
echo -e "   运行: ${CYAN}$BIN_NAME bench -n 1000 -c 10 https://example.com${NC}"
echo ""
echo -e "   安装已完成！具体操作可以前往仓库主页查看教程："
echo -e "   GitHub: ${CYAN}https://github.com/$REPO${NC}"
echo -e "   Gitee:  ${CYAN}https://gitee.com/$GITEE_REPO${NC}"
echo ""
