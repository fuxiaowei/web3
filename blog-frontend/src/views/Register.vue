<template>
  <div class="register-container">
    <el-card class="register-card" shadow="hover">
      <h2 class="register-title">博客系统 - 注册</h2>
      <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="registerForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="registerForm.email" placeholder="请输入邮箱"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleRegister" class="register-btn">注册</el-button>
          <el-button @click="$router.push('/login')">返回登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { register } from '@/api/auth'

const router = useRouter()

// 表单数据
const registerForm = ref({
  username: '',
  password: '',
  email: ''
})

// 表单验证规则
const registerRules = ref({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
})

const registerFormRef = ref(null)

// 注册处理
const handleRegister = async () => {
  try {
    // 表单验证
    await registerFormRef.value.validate()
    // 调用注册接口
    await register(registerForm.value)
    ElMessage.success('注册成功，请登录')
    // 跳转到登录页
    router.push('/login')
  } catch (error) {
    console.error('注册失败：', error)
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.register-card {
  width: 400px;
  padding: 20px;
}

.register-title {
  text-align: center;
  margin-bottom: 20px;
  color: #1989fa;
}

.register-btn {
  width: 100px;
  margin-right: 10px;
}
</style>