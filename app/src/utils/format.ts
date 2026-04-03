/**
 * Format a timestamp into relative time string.
 * e.g., 刚刚, 3分钟前, 2小时前, 昨天 14:30, 2025-03-20
 */
export function formatTime(time: string | number | Date): string {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (seconds < 60) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`

  const yesterday = new Date(now)
  yesterday.setDate(yesterday.getDate() - 1)
  if (
    date.getDate() === yesterday.getDate() &&
    date.getMonth() === yesterday.getMonth() &&
    date.getFullYear() === yesterday.getFullYear()
  ) {
    const h = String(date.getHours()).padStart(2, '0')
    const m = String(date.getMinutes()).padStart(2, '0')
    return `昨天 ${h}:${m}`
  }

  if (days < 7) return `${days}天前`

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  if (year === now.getFullYear()) {
    return `${month}-${day}`
  }
  return `${year}-${month}-${day}`
}

/**
 * Format coin amount with +/- prefix.
 */
export function formatCoin(amount: number): string {
  if (amount > 0) return `+${amount}`
  return String(amount)
}

/**
 * Format a date for headlines display.
 * e.g., 3月25日 星期二
 */
export function formatHeadlineDate(time?: string | number | Date): string {
  const date = time ? new Date(time) : new Date()
  const weekdays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  const month = date.getMonth() + 1
  const day = date.getDate()
  const weekday = weekdays[date.getDay()]
  return `${month}月${day}日 ${weekday}`
}
