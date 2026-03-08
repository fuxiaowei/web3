import axios from 'axios'
import {getToken} from '@/utils/auth'
import {ElMessage} from 'element-plus'
import {useUserStore} from '@/store'

// 创建Axios实例
const service = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api', // 基础地址
    timeout: 5000 // 请求超时
})

// 请求拦截器：添加Token
service.interceptors.request.use(
    (config) => {
        const token = getToken()
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        ElMessage.error('请求异常：' + error.message)
        return Promise.reject(error)
    }
)

// 响应拦截器：统一处理错误（核心修改）
service.interceptors.response.use(
    (response) => {
        const res = response.data
        // 后端返回非200码，视为错误
        if (res.code !== 200) {
            // 这里只抛错，不弹提示（交给组件处理）
            return Promise.reject(res)
        }
        return res
    },
    (error) => {
        // Token过期/未登录：仍全局处理（必须的全局逻辑）
        if (error.response?.status === 401) {
            const userStore = useUserStore()
            userStore.logout()
            ElMessage.error('登录已过期，请重新登录')
            window.location.href = '/login'
        }
        // 其他错误（如403/404/500）：只抛错，不弹包含状态码的提示
        else {
            // 移除原有弹提示的代码，仅抛错
            // ElMessage.error(error.message || '服务器错误')
        }
        return Promise.reject(error)
    }
)

export default service