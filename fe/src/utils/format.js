export const formatFileSize = (size) => {
    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    let index = 0
    let fileSize = size

    while (fileSize >= 1024 && index < units.length - 1) {
        fileSize /= 1024
        index++
    }

    return `${fileSize.toFixed(2)} ${units[index]}`
}

export const formatDuration = (seconds) => {
    const h = Math.floor(seconds / 3600)
    const m = Math.floor((seconds % 3600) / 60)
    const s = Math.floor(seconds % 60)
    return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
}

export function formatDateTime(dateStr) {
    if (!dateStr) return ''
    try {
        const date = new Date(dateStr)
        if (isNaN(date.getTime())) return '' // 检查日期是否有效

        return date.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false,
            timeZone: 'Asia/Shanghai' // 明确指定时区为中国时区
        })
    } catch (error) {
        console.error('日期格式化错误:', error)
        return ''
    }
}

export function formatDate(dateStr) {
    if (!dateStr) return ''
    const date = new Date(dateStr)
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
    })
}