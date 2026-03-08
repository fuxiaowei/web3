// Token 存储的key
const TOKEN_KEY = 'blog_token'

// 保存Token到localStorage
export function setToken(token) {
    localStorage.setItem(TOKEN_KEY, token)
}

// 获取Token
export function getToken() {
    return localStorage.getItem(TOKEN_KEY)
}

// 移除Token
export function removeToken() {
    localStorage.removeItem(TOKEN_KEY)
}