import request from '@/utils/request'

// 获取视频列表
export function getVideoList(params) {
  return request({
    url: '/api/videos',
    method: 'get',
    params
  })
}

// 获取视频详情
export function getVideoDetail(id) {
  return request({
    url: `/api/videos/${id}`,
    method: 'get'
  })
}

// 创建视频记录
export function createVideo(data) {
  return request({
    url: '/api/videos',
    method: 'post',
    data
  })
}

// 更新视频信息
export function updateVideo(id, data) {
  return request({
    url: `/api/videos/${id}`,
    method: 'put',
    data
  })
}

// 删除视频
export function deleteVideo(id) {
  return request({
    url: `/api/videos/${id}`,
    method: 'delete'
  })
}

// 下载视频
export function downloadVideo(id) {
  return request({
    url: `/api/videos/${id}/download`,
    method: 'get',
    responseType: 'blob'
  })
}

// 开始录制
export function startRecording(data) {
  return request({
    url: '/api/videos/start-recording',
    method: 'post',
    data
  })
}

// 停止录制
export function stopRecording(workshopId) {
  return request({
    url: `/api/videos/stop-recording/${workshopId}`,
    method: 'post'
  })
}

// 获取视频预览地址
export function getVideoPreview(id) {
  return request({
    url: `/api/videos/${id}/preview`,
    method: 'get'
  })
}

// 获取车间列表
export function getWorkshopList() {
  return request({
    url: '/api/workshops',
    method: 'get'
  })
}

// 获取车间预览流
export function getWorkshopPreview(id) {
  return request({
    url: `/api/workshops/${id}/preview`,
    method: 'get'
  })
}

// 批量删除视频
export function batchDeleteVideos(ids) {
  return request({
    url: '/api/videos/batch',
    method: 'delete',
    data: { ids }
  })
}

// 获取视频统计信息
export function getVideoStats() {
  return request({
    url: '/api/videos/stats',
    method: 'get'
  })
} 