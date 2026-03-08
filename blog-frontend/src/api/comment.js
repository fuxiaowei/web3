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