import service from './index'

// 注册
export function register(data) {
    return service({
        url: '/register',
        method: 'post',
        data
    })
}

// 登录
export function login(data) {
    return service({
        url: '/login',
        method: 'post',
        data
    })
}

// 获取用户信息（可选）
export function getUserInfo() {
    return service({
        url: '/user/info',
        method: 'get'
    })
}