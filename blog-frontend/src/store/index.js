import { defineStore } from 'pinia'
import { setToken, getToken, removeToken } from '@/utils/auth'
import { login as apiLogin, getUserInfo } from '@/api/auth'

// 定义用户Store
export const useUserStore = defineStore('user', {
    state: () => ({
        token: getToken() || '', // 从本地存储获取Token
        username: '',
        userId: 0
    }),
    actions: {
        // 登录
        async login(userInfo) {
            const res = await apiLogin(userInfo)
            this.token = res.data.token
            setToken(res.data.token) // 保存到本地存储
            // 可选：获取用户信息
            // const userRes = await getUserInfo()
            // this.username = userRes.data.username
            // this.userId = userRes.data.id
        },
        // 退出登录
        logout() {
            this.token = ''
            this.username = ''
            this.userId = 0
            removeToken() // 清除本地存储
        },
        // 初始化用户信息（页面刷新后）
        initUserInfo() {
            const token = getToken()
            if (token) {
                this.token = token
                // 可选：重新获取用户信息
                // getUserInfo().then(res => {
                //   this.username = res.data.username
                //   this.userId = res.data.id
                // })
            }
        }
    }
})