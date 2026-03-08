import service from './index'

// 获取文章列表
export function getPostList() {
    return service({
        url: '/posts',
        method: 'get'
    })
}

// 获取单篇文章
export function getPostDetail(id) {
    return service({
        url: `/posts/${id}`,
        method: 'get'
    })
}

// 创建文章
export function createPost(data) {
    return service({
        url: '/posts',
        method: 'post',
        data
    })
}

// 更新文章
export function updatePost(id, data) {
    return service({
        url: `/posts/${id}`,
        method: 'put',
        data
    })
}

// 删除文章
export function deletePost(id) {
    return service({
        url: `/posts/${id}`,
        method: 'delete'
    })
}