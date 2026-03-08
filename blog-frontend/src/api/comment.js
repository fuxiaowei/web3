import service from './index'

// 获取文章评论
export function getCommentList(postId) {
    return service({
        url: `/comments/post/${postId}`,
        method: 'get'
    })
}

// 发布评论
export function createComment(data) {
    return service({
        url: '/comments',
        method: 'post',
        data
    })
}


// 获取某篇文章的评论列表
export function getCommentsByPostId(postId) {
    return service({
        url: `/comments/post/${postId}`,
        method: 'get'
    })
}

// 发布新评论
export function addComment(data) {
    return service({
        url: '/comments',
        method: 'post',
        data
    })
}