<template>
  <div class="post-detail-container">
    <el-card shadow="hover">
      <h2>{{ post.title }}</h2>
      <div class="post-meta">
        <span>作者：{{ post.user?.username || '未知' }}</span>
        <span>创建时间：{{ formatTime(post.CreatedAt) }}</span>
      </div>
      <el-divider></el-divider>
      <div class="post-content">{{ post.content }}</div>
      <el-divider></el-divider>
      <el-button @click="$router.back()">返回列表</el-button>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getPostDetail } from '@/api/post'

const route = useRoute()
const post = ref({})

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString()
}

// 获取文章详情
const fetchPostDetail = async () => {
  try {
    const res = await getPostDetail(route.params.id)
    post.value = res.data
  } catch (error) {
    console.error('获取文章详情失败：', error)
  }
}

onMounted(() => {
  fetchPostDetail()
})
</script>

<style scoped>
.post-detail-container {
  padding: 20px;
}
.post-meta {
  margin: 10px 0;
  color: #666;
  display: flex;
  gap: 20px;
}
.post-content {
  margin: 20px 0;
  line-height: 1.8;
  font-size: 16px;
}
</style>