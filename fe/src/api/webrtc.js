import request from '@/utils/request'

// WebRTC 相关的 API
export function startWebRTC(data) {
    //console.log('Starting WebRTC with data:', data)
    return request({
        url: '/api/webrtc',
        method: 'post',
        data,
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(response => {
        console.log('WebRTC response:', response)
        // 确保响应格式正确
        if (!response.data || !response.data.success || !response.data.sdp) {
            throw new Error('Invalid response format')
        }
        //console.log("Answer SDP:", response.data.sdp);
        return response
    })
}