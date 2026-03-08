<template>
  <div class="post-edit-container">
    <el-card shadow="hover">
      <h2>{{ isEdit ? '编辑文章' : '新增文章' }}</h2>
      <el-form :model="postForm" :rules="postRules" ref="postFormRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="postForm.title" placeholder="请输入文章标题"></el-input>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="postForm.content" type="textarea" :rows="10" placeholder="请输入文章内容"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSubmit">提交</el-button>
          <el-button @click="$router.push('/posts')">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import {ref, onMounted} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage} from 'element-plus'
import {createPost, getPostDetail, updatePost} from '@/api/post'

const route = useRoute()
const router = useRouter()
const isEdit = !!route.params.id // 是否为编辑模式
const postFormRef = ref(null)

// 表单数据
const postForm = ref({
  title: '',
  content: ''
})

// 表单验证规则
const postRules = ref({
  title: [{required: true, message: '请输入文章标题', trigger: 'blur'}],
  content: [{required: true, message: '请输入文章内容', trigger: 'blur'}]
})

// 获取文章详情（编辑模式）
const fetchPostDetail = async () => {
  try {
    const res = await getPostDetail(route.params.id)
    postForm.value = {
      title: res.data.title,
      content: res.data.content
    }
  } catch (error) {
    ElMessage.error('获取文章详情失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    await postFormRef.value.validate()
    if (isEdit) {
      // 编辑文章
      await updatePost(route.params.id, postForm.value)
      ElMessage.success('编辑成功')
    } else {
      // 新增文章
      await createPost(postForm.value)
      ElMessage.success('新增成功')
    }
    router.push('/posts')
  } catch (error) {
    // 优先从后端响应中提取 msg
    let errMsg = '提交失败'
    if (error.response?.data?.msg) {
      // 后端返回了业务信息，如 "无权限删除该文章"
      errMsg = error.response.data.msg
    } else if (error.message) {
      // 网络错误或其他 JS 错误
      errMsg = error.message
    }

    // 显示提示信息
    ElMessage.error(errMsg)
  }
}

// 编辑模式下加载文章详情
onMounted(() => {
  if (isEdit) {
    fetchPostDetail()
  }
})
</script>

<style scoped>
.post-edit-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}
</style>