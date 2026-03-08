<template>
  <div class="login-container">
    <el-card class="login-card" shadow="hover">
      <h2 class="login-title">博客系统 - 登录</h2>
      <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" class="login-btn">登录</el-button>
          <el-button @click="$router.push('/register')">注册</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store'
import { login } from '@/api/auth'

const router = useRouter()
const userStore = useUserStore()

// 表单数据
const loginForm = ref({
  username: '',
  password: ''
})

// 表单验证规则
const loginRules = ref({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
})

const loginFormRef = ref(null)

// 登录处理
const handleLogin = async () => {
  try {
    // 表单验证
    await loginFormRef.value.validate()
    // 调用登录接口
    await userStore.login(loginForm.value)
    ElMessage.success('登录成功')
    // 跳转到文章列表
    router.push('/posts')
  } catch (error) {
    console.error('登录失败：', error)
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.login-card {
  width: 400px;
  padding: 20px;
}

.login-title {
  text-align: center;
  margin-bottom: 20px;
  color: #1989fa;
}

.login-btn {
  width: 100px;
  margin-right: 10px;
}
</style>