import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const service = axios.create({
    baseURL: '/api',
    timeout: 30000
})

// 请求拦截器
service.interceptors.request.use(
    config => {
        // 可以在这里添加 token
        // const token = localStorage.getItem('token')
        // if (token) {
        //   config.headers['Authorization'] = `Bearer ${token}`
        // }
        return config
    },
    error => {
        console.error('请求错误:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
service.interceptors.response.use(
    response => {
        const res = response.data

        // 如果是下载文件，直接返回
        if (response.config.responseType === 'blob') {
            return response.data
        }

        // 这里可以根据后端的响应结构进行调整
        if (res.code !== 0) {
            ElMessage.error(res.message || '请求失败')
            return Promise.reject(new Error(res.message || '请求失败'))
        }

        return res
    },
    error => {
        console.error('响应错误:', error)
        const message = error.response?.data?.message || error.message
        ElMessage.error(message)
        return Promise.reject(error)
    }
)

export default service
