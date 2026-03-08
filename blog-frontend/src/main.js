import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { createPinia } from 'pinia'
import router from './router' // 确保路由已导入
import App from './App.vue'

const app = createApp(App)

// 注册 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

// 安装插件
app.use(ElementPlus)
app.use(createPinia())
app.use(router) // 注册路由（关键）

app.mount('#app')