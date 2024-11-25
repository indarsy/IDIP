import request from './request'

// 获取车间列表
// 获取车间列表
export function getWorkshopList() {
    return request({
        url: '/api/workshops',
        method: 'get'
    })
}

// 获取车间详情
export function getWorkshopDetail(id) {
    return request({
        url: `/api/workshops/${id}`,
        method: 'get'
    })
}

// 添加车间
export function addWorkshop(data) {
    return request({
        url: '/api/workshops',
        method: 'post',
        data
    })
}

// 编辑车间
export function updateWorkshop(id, data) {
    return request({
        url: `/api/workshops/${id}`,
        method: 'put',
        data
    })
}

// 删除车间
export function deleteWorkshop(id) {
    return request({
        url: `/api/workshops/${id}`,
        method: 'delete'
    })
}
