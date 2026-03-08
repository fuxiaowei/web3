<template>
  <div class="post-list-container">
    <navbar></navbar>
    <el-card class="post-list-card" shadow="never">
      <div class="post-list-header">
        <h2>文章列表</h2>
        <el-button type="primary" @click="$router.push('/post/edit')">
          <el-icon><Plus /></el-icon> 新增文章
        </el-button>
      </div>
      <el-divider></el-divider>
      <el-table :data="postList" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="title" label="标题" min-width="300"></el-table-column>
        <el-table-column prop="user.username" label="作者" width="120"></el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="200">
          <template #default="scope">
            {{ formatTime(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewPost(scope.row.id)">查看</el-button>
            <el-button type="warning" size="small" @click="editPost(scope.row.id)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDeletePost(scope.row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus' // 替换 ElConfirm 为 ElMessageBox
import { Plus } from '@element-plus/icons-vue'
import NavBar from '@/components/NavBar.vue'
import { getPostList, deletePost } from '@/api/post'

const router = useRouter()
const postList = ref([])

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString()
}

// 获取文章列表
const fetchPostList = async () => {
  try {
    const res = await getPostList()
    postList.value = res.data
  } catch (error) {
    console.error('获取文章列表失败：', error)
  }
}

// 查看文章
const viewPost = (id) => {
  router.push(`/posts/${id}`)
}

// 编辑文章
const editPost = (id) => {
  router.push(`/post/edit/${id}`)
}

// 删除文章
const handleDeletePost = async (id) => {
  try {
    // 将 ElConfirm.confirm 替换为 ElMessageBox.confirm
    await ElMessageBox.confirm(
        '确定要删除这篇文章吗？',
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
    )
    await deletePost(id)
    ElMessage.success('删除成功')
    fetchPostList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 页面挂载时加载列表
onMounted(() => {
  fetchPostList()
})
</script>

<style scoped>
.post-list-container {
  padding: 20px;
}

.post-list-card {
  padding: 20px;
}

.post-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}
</style>