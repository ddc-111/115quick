export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

export function formatSpeed(bytesPerSecond: number): string {
  if (bytesPerSecond === 0) return '0 B/s'
  return formatFileSize(bytesPerSecond) + '/s'
}

export function formatTime(seconds: number): string {
  if (seconds < 60) return `${Math.floor(seconds)}秒`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}分${Math.floor(seconds % 60)}秒`
  return `${Math.floor(seconds / 3600)}时${Math.floor((seconds % 3600) / 60)}分`
}

export function getStatusText(status: number): string {
  switch (status) {
    case -1: return '失败'
    case 0: return '分配中'
    case 1: return '下载中'
    case 2: return '已完成'
    default: return '未知'
  }
}

export function getStatusType(status: number): '' | 'success' | 'warning' | 'danger' | 'info' {
  switch (status) {
    case -1: return 'danger'
    case 0: return 'info'
    case 1: return 'warning'
    case 2: return 'success'
    default: return ''
  }
}
