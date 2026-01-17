/**
 * 获取响应头中的文件名
 * @param contentDisposition 响应头内容
 * @returns 文件名
 */
export function getFilenameFromContentDisposition(contentDisposition?: string): string | undefined {
  if (!contentDisposition) return undefined

  // 兼容 RFC 5987: filename*=UTF-8''xxx
  const filenameStarMatch = contentDisposition.match(/filename\*=(?:UTF-8'')?([^;]+)/i)
  if (filenameStarMatch?.[1]) {
    const raw = filenameStarMatch[1].trim().replace(/^"|"$/g, '')
    try {
      return decodeURIComponent(raw)
    } catch {
      return raw
    }
  }

  const filenameMatch = contentDisposition.match(/filename=([^;]+)/i)
  if (filenameMatch?.[1]) {
    return filenameMatch[1].trim().replace(/^"|"$/g, '')
  }

  return undefined
}

/**
 * 下载 Blob 对象
 * @param blob Blob 对象
 * @param filename 文件名
 */
export function downloadBlob(blob: Blob, filename: string) {
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  a.remove()
  window.URL.revokeObjectURL(url)
}
