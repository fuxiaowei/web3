<template>
  <div class="post-detail-container">
    <!-- 原有文章内容区域 -->
    <el-card shadow="hover" class="post-card">
      <div class="post-title">{{ post.title }}</div>
      <el-divider />
      <div class="post-meta">
        <span>作者：{{ post.user?.username || '未知' }}</span>
        <span>发布时间：{{ formatTime(post.CreatedAt) }}</span>
      </div>
      <el-divider />
      <div class="post-content">{{ post.content }}</div>
      <el-divider />
      <el-button @click="goBack">返回列表</el-button>
    </el-card>

    <!-- 新增评论模块 -->
    <el-card shadow="hover" class="comment-card" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>评论区 ({{ comments.length }})</span>
        </div>
      </template>

      <!-- 发布评论表单 -->
      <div class="comment-form">
        <el-form :model="commentForm" @submit.prevent="submitComment">
          <el-form-item>
            <el-input
                v-model="commentForm.content"
                type="textarea"
                placeholder="请输入评论内容..."
                rows="3"
            ></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitComment">发布评论</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-divider />

      <!-- 评论列表 -->
      <div class="comment-list" v-if="comments.length > 0">
        <div class="comment-item" v-for="item in comments" :key="item.ID">
          <div class="comment-author">
            <span>{{ item.user?.username || '匿名用户' }}</span>
            <span class="comment-time">{{ formatTime(item.CreatedAt) }}</span>
          </div>
          <div class="comment-content">{{ item.content }}</div>
          <el-divider size="small" />
        </div>
      </div>

      <!-- 空评论提示 -->
      <div class="empty-comment" v-else>
        <el-empty description="暂无评论，快来抢沙发吧～"></el-empty>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getPostDetail } from '@/api/post'
// 引入评论API
import {getCommentsByPostId, addComment, createComment} from '@/api/comment'

const route = useRoute()
const router = useRouter()
const postId = route.params.id

// 文章数据
const post = ref({})
// 评论相关
const comments = ref([])
const commentForm = ref({
  post_id: postId, // 关联文章ID
  content: ''      // 评论内容
})

// 时间格式化（适配你的字段名 CreatedAt）
const formatTime = (time) => {
  if (!time) return ''
  // 兼容时间戳（秒级）转毫秒级
  const timestamp = typeof time === 'number' && time.toString().length === 10 ? time * 1000 : time
  return new Date(timestamp).toLocaleString()
}

// 获取文章详情（原有逻辑）
const fetchPostDetail = async () => {
  try {
    const res = await getPostDetail(postId)
    post.value = res.data
  } catch (err) {
    ElMessage.error('获取文章失败')
    console.error(err)
  }
}

// 获取评论列表（新增）
const fetchComments = async () => {
  try {
    const res = await getCommentsByPostId(postId)
    comments.value = res.data
  } catch (err) {
    ElMessage.error('获取评论列表失败')
    console.error(err)
  }
}

// 发布评论（新增）
const submitComment = async () => {
  if (!commentForm.value.content.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  try {
    await createComment(commentForm.value)
    ElMessage.success('评论发布成功')
    // 清空输入框并重新加载评论列表
    commentForm.value.content = ''
    await fetchComments()
  } catch (err) {
    ElMessage.error('发布评论失败')
    console.error(err)
  }
}

// 返回列表（原有逻辑）
const goBack = () => {
  router.push('/posts')
}

// 组件挂载时加载数据
onMounted(() => {
  fetchPostDetail()
  fetchComments() // 新增加载评论
})
</script>

<style scoped>
.post-detail-container {
  max-width: 1200px;
  margin: 20px auto;
  padding: 0 20px;
}

.post-card {
  padding: 20px;
}

.post-title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 10px;
}

.post-meta {
  color: #999;
  margin-bottom: 20px;
  display: flex;
  gap: 20px;
}

.post-content {
  line-height: 1.8;
  font-size: 16px;
  margin-bottom: 20px;
}

.comment-card {
  padding: 20px;
}

.comment-form {
  margin-bottom: 20px;
}

.comment-item {
  padding: 10px 0;
}

.comment-author {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}

.comment-time {
  color: #999;
  font-size: 12px;
}

.comment-content {
  line-height: 4.0;
  font-size: 14px;
  color: #666;
  background-color: rgba(245, 245, 245, 0.8); /* 最后一个值是透明度（0-1） */
}

.empty-comment {
  padding: 40px 0;
  text-align: center;
}
</style>