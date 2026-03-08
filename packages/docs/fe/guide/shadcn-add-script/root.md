# Shadcn添加依赖脚本

因为 Shadcn 添加依赖太不稳定了，编写一个shell脚本，重复尝试add。

下载完之后的会将记录写入到 `shadcn-installed-urls.log` 文件中。

也可以配置`.npmrc`，让命令走代理。

```txt
proxy=http://127.0.0.1:29758
https-proxy=http://127.0.0.1:29758
```

---

`retry-shadcn.sh`

给脚本权限

```bash
chmod +x ./retry-shadcn.sh
```

---

```bash
#!/usr/bin/env bash
set -uo pipefail

SLEEP_SECONDS="${SLEEP_SECONDS:-5}"
SUCCESS_FILE="${SUCCESS_FILE:-.shadcn-installed-urls.log}"

# 1) 不传参数时，走这里的默认列表（每一项都应是完整命令）
# chmod +x retry-shadcn.sh
COMMANDS=(
  "pnpx shadcn@latest add @plate/editor-ai"
  "pnpx shadcn@latest add https://platejs.org/r/ai-kit"
  "pnpx shadcn@latest add https://platejs.org/r/copilot-kit"
  "pnpx shadcn@latest add https://platejs.org/r/comment-kit"
  "pnpx shadcn@latest add https://platejs.org/r/suggestion-kit"
  "pnpx shadcn@latest add https://platejs.org/r/basic-blocks-kit"
  "pnpx shadcn@latest add https://platejs.org/r/callout-kit"
  "pnpx shadcn@latest add https://platejs.org/r/code-block-kit"
  "pnpx shadcn@latest add https://platejs.org/r/column-kit"
  "pnpx shadcn@latest add https://platejs.org/r/date-kit"
  "pnpx shadcn@latest add https://platejs.org/r/math-kit"
  "pnpx shadcn@latest add https://platejs.org/r/link-kit"
  "pnpx shadcn@latest add https://platejs.org/r/list-classic-kit"
  "pnpx shadcn@latest add https://platejs.org/r/media-kit"
  "pnpx shadcn@latest add https://platejs.org/r/mention-kit"
  "pnpx shadcn@latest add https://platejs.org/r/table-kit"
  "pnpx shadcn@latest add https://platejs.org/r/toc-kit"
  "pnpx shadcn@latest add https://platejs.org/r/toggle-kit"
  "pnpx shadcn@latest add https://platejs.org/r/basic-marks-kit"
  "pnpx shadcn@latest add https://platejs.org/r/basic-marks-kit"
  "pnpx shadcn@latest add https://platejs.org/r/font-kit"
  "pnpx shadcn@latest add https://platejs.org/r/line-height-kit"
  "pnpx shadcn@latest add https://platejs.org/r/align-kit"
  "pnpx shadcn@latest add https://platejs.org/r/indent-kit"
  "pnpx shadcn@latest add https://platejs.org/r/list-kit"
  "pnpx shadcn@latest add https://platejs.org/r/exit-break-kit"
  "pnpx shadcn@latest add https://platejs.org/r/autoformat-kit"
  "pnpx shadcn@latest add https://platejs.org/r/emoji-kit"
  "pnpx shadcn@latest add https://platejs.org/r/dnd-kit"
  "pnpx shadcn@latest add https://platejs.org/r/ai-menu"
  "pnpx shadcn@latest add https://platejs.org/r/ai-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/align-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/block-context-menu"
  "pnpx shadcn@latest add https://platejs.org/r/block-selection"
  "pnpx shadcn@latest add https://platejs.org/r/import-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/export-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/caption"
  "pnpx shadcn@latest add https://platejs.org/r/font-color-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/comment-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/block-discussion"
  "pnpx shadcn@latest add https://platejs.org/r/cursor-overlay"
  "pnpx shadcn@latest add https://platejs.org/r/block-draggable"
  "pnpx shadcn@latest add https://platejs.org/r/editor"
  "pnpx shadcn@latest add https://platejs.org/r/select-editor"
  "pnpx shadcn@latest add https://platejs.org/r/emoji-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/fixed-toolbar-buttons"
  "pnpx shadcn@latest add https://platejs.org/r/fixed-toolbar-classic-buttons"
  "pnpx shadcn@latest add https://platejs.org/r/fixed-toolbar"
  "pnpx shadcn@latest add https://platejs.org/r/floating-toolbar-buttons"
  "pnpx shadcn@latest add https://platejs.org/r/floating-toolbar-classic-buttons"
  "pnpx shadcn@latest add https://platejs.org/r/floating-toolbar"
  "pnpx shadcn@latest add https://platejs.org/r/ghost-text"
  "pnpx shadcn@latest add https://platejs.org/r/history-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/list-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/indent-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/inline-combobox"
  "pnpx shadcn@latest add https://platejs.org/r/insert-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/insert-toolbar-classic-button"
  "pnpx shadcn@latest add https://platejs.org/r/line-height-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/link-toolbar"
  "pnpx shadcn@latest add https://platejs.org/r/link-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/list-classic-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/mark-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/media-toolbar"
  "pnpx shadcn@latest add https://platejs.org/r/media-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/media-upload-toast"
  "pnpx shadcn@latest add https://platejs.org/r/mode-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/more-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/resize-handle"
  "pnpx shadcn@latest add https://platejs.org/r/table-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/toggle-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/turn-into-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/turn-into-toolbar-classic-button"
  "pnpx shadcn@latest add https://platejs.org/r/remote-cursor-overlay"
  "pnpx shadcn@latest add https://platejs.org/r/toolbar"
  "pnpx shadcn@latest add https://platejs.org/r/suggestion-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/ai-node"
  "pnpx shadcn@latest add https://platejs.org/r/block-list"
  "pnpx shadcn@latest add https://platejs.org/r/blockquote-node"
  "pnpx shadcn@latest add https://platejs.org/r/code-block-node"
  "pnpx shadcn@latest add https://platejs.org/r/code-drawing-node"
  "pnpx shadcn@latest add https://platejs.org/r/code-node"
  "pnpx shadcn@latest add https://platejs.org/r/column-node"
  "pnpx shadcn@latest add https://platejs.org/r/comment-node"
  "pnpx shadcn@latest add https://platejs.org/r/suggestion-node"
  "pnpx shadcn@latest add https://platejs.org/r/date-node"
  "pnpx shadcn@latest add https://platejs.org/r/equation-node"
  "pnpx shadcn@latest add https://platejs.org/r/equation-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/emoji-node"
  "pnpx shadcn@latest add https://platejs.org/r/excalidraw-node"
  "pnpx shadcn@latest add https://platejs.org/r/font-size-toolbar-button"
  "pnpx shadcn@latest add https://platejs.org/r/heading-node"
  "pnpx shadcn@latest add https://platejs.org/r/highlight-node"
  "pnpx shadcn@latest add https://platejs.org/r/hr-node"
  "pnpx shadcn@latest add https://platejs.org/r/media-image-node"
  "pnpx shadcn@latest add https://platejs.org/r/media-preview-dialog"
  "pnpx shadcn@latest add https://platejs.org/r/kbd-node"
  "pnpx shadcn@latest add https://platejs.org/r/link-node"
  "pnpx shadcn@latest add https://platejs.org/r/list-classic-node"
  "pnpx shadcn@latest add https://platejs.org/r/media-audio-node"
  "pnpx shadcn@latest add https://platejs.org/r/media-embed-node"
  "pnpx shadcn@latest add https://platejs.org/r/media-file-node"
  "pnpx shadcn@latest add https://platejs.org/r/media-placeholder-node"
  "pnpx shadcn@latest add https://platejs.org/r/media-video-node"
  "pnpx shadcn@latest add https://platejs.org/r/mention-node"
  "pnpx shadcn@latest add https://platejs.org/r/paragraph-node"
  "pnpx shadcn@latest add https://platejs.org/r/search-highlight-node"
  "pnpx shadcn@latest add https://platejs.org/r/slash-node"
  "pnpx shadcn@latest add https://platejs.org/r/table-node"
  "pnpx shadcn@latest add https://platejs.org/r/tag-node"
  "pnpx shadcn@latest add https://platejs.org/r/toc-node"
  "pnpx shadcn@latest add https://platejs.org/r/toggle-node"
)

# 2) 如果传了参数，就用参数覆盖默认列表
# 用法示例:
# ./retry-shadcn.sh "pnpx shadcn@latest add @plate/editor-ai" "pnpm add @plate/editor-ai"
if [ "$#" -gt 0 ]; then
  COMMANDS=("$@")
fi

touch "$SUCCESS_FILE"

is_recorded_success() {
  local command="$1"
  grep -Fqx "${command}" "${SUCCESS_FILE}"
}

record_success() {
  local command="$1"
  if ! is_recorded_success "${command}"; then
    echo "${command}" >> "${SUCCESS_FILE}"
  fi
}

retry_until_success() {
  local command="$1"
  local attempt=1

  if is_recorded_success "${command}"; then
    echo "⏭️ 已跳过（历史成功）: ${command}"
    return 0
  fi

  echo "开始: ${command}"
  while true; do
    echo "[$(date '+%F %T')] 第 ${attempt} 次尝试..."
    bash -lc "${command}"
    local exit_code=$?

    if [ "${exit_code}" -eq 0 ]; then
      record_success "${command}"
      echo "✅ 成功（第 ${attempt} 次）: ${command}"
      return 0
    fi

    echo "❌ 失败（exit=${exit_code}），${SLEEP_SECONDS}s 后重试: ${command}"
    attempt=$((attempt + 1))
    sleep "${SLEEP_SECONDS}"
  done
}

SEEN_FILE="$(mktemp -t shadcn-seen.XXXXXX)"
trap 'rm -f "${SEEN_FILE}"' EXIT

is_seen_in_current_run() {
  local command="$1"
  grep -Fqx "${command}" "${SEEN_FILE}"
}

mark_seen_in_current_run() {
  local command="$1"
  if ! is_seen_in_current_run "${command}"; then
    echo "${command}" >> "${SEEN_FILE}"
  fi
}

for command in "${COMMANDS[@]}"; do
  if is_seen_in_current_run "${command}"; then
    echo "⏭️ 已跳过（本次重复）: ${command}"
    continue
  fi
  mark_seen_in_current_run "${command}"
  retry_until_success "${command}"
done

echo "🎉 全部命令已成功执行。成功记录文件: ${SUCCESS_FILE}"
```
