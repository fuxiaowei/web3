import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/utils/auth'

// 导入页面组件
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import PostList from '@/views/PostList.vue'
import PostDetail from '@/views/PostDetail.vue'
import PostEdit from '@/views/PostEdit.vue'

// 路由规则
const routes = [
    {
        path: '/',
        redirect: '/posts' // 默认跳转到文章列表
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: { requiresAuth: false }
    },
    {
        path: '/register',
        name: 'Register',
        component: Register,
        meta: { requiresAuth: false }
    },
    {
        path: '/posts',
        name: 'PostList',
        component: PostList,
        meta: { requiresAuth: false }
    },
    {
        path: '/posts/:id',
        name: 'PostDetail',
        component: PostDetail,
        meta: { requiresAuth: false }
    },
    {
        path: '/post/edit/:id?',
        name: 'PostEdit',
        component: PostEdit,
        meta: { requiresAuth: true }
    }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
    const requiresAuth = to.meta.requiresAuth
    const token = getToken()

    if (requiresAuth && !token) {
        next('/login')
    } else {
        next()
    }
})

export default router