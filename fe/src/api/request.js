import axios from 'axios'
import { ElMessage } from 'element-plus'

const service = axios.create({
    baseURL: '/api',
    timeout: 30000
})

service.interceptors.request.use(
    config => {
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

service.interceptors.response.use(
    response => {
        const res = response.data
        if (res.code !== 0) {
            ElMessage.error(res.message || '请求失败')
            return Promise.reject(new Error(res.message || '请求失败'))
        }
        return res
    },
    error => {
        ElMessage.error(error.message || '请求失败')
        return Promise.reject(error)
    }
)

export default service
