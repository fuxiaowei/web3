<template>
  <el-header style="text-align: right; font-size: 12px">
    <el-menu :default-active="activeIndex" mode="horizontal" background-color="#1989fa" text-color="#fff" active-text-color="#ffd04b">
      <el-menu-item index="1" @click="$router.push('/posts')">
        <el-icon><House /></el-icon>
        <span>首页</span>
      </el-menu-item>
      <el-menu-item index="2" v-if="!isLogin" @click="$router.push('/login')">
        <el-icon><User /></el-icon>
        <span>登录</span>
      </el-menu-item>
      <el-menu-item index="3" v-if="!isLogin" @click="$router.push('/register')">
        <el-icon><UserFilled /></el-icon>
        <span>注册</span>
      </el-menu-item>
      <el-menu-item index="4" v-if="isLogin" @click="handleLogout">
        <el-icon><SwitchButton /></el-icon>
        <span>退出登录</span>
      </el-menu-item>
    </el-menu>
  </el-header>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { House, User, UserFilled, SwitchButton } from '@element-plus/icons-vue'
import { useUserStore } from '@/store'
import { getToken } from '@/utils/auth'

const router = useRouter()
const userStore = useUserStore()
const activeIndex = ref('1')
const isLogin = ref(false)

// 退出登录
const handleLogout = () => {
  userStore.logout()
  ElMessage.success('退出成功')
  router.push('/login')
}

// 初始化登录状态
onMounted(() => {
  isLogin.value = !!getToken()
  userStore.initUserInfo()
})
</script>

<style scoped>
.el-header {
  background-color: #1989fa;
  color: #333;
  line-height: 60px;
  margin-bottom: 20px;
}
</style>