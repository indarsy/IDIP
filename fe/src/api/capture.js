import request from '@/utils/request'

export function createCapture(data) {
    return request({
        url: '/api/captures',
        method: 'post',
        data
    })
}

export function getCaptureList(workshopId) {
    return request({
        url: '/api/captures',
        method: 'get',
        params: { workshopId }
    })
}

// 获取单个采集任务详情
export function getCaptureDetail(id) {
    return request({
        url: `/api/captures/${id}`,
        method: 'get'
    })
}

// 取消采集任务
export function cancelCapture(id) {
    return request({
        url: `/api/captures/${id}/cancel`,
        method: 'post'
    })
}