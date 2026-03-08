# Vue 3 + Vite

This template should help get you started developing with Vue 3 in Vite. The template uses Vue 3 `<script setup>` SFCs, check out the [script setup docs](https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup) to learn more.

Learn more about IDE Support for Vue in the [Vue Docs Scaling up Guide](https://vuejs.org/guide/scaling-up/tooling.html#ide-support).

一、项目设计思路
技术栈：Vue3 (Composition API) + Vite + Element Plus (UI 组件) + Axios (请求) + Pinia (状态管理，存储 Token)；
核心功能：
用户认证：注册 / 登录，Token 持久化；
文章管理：列表展示、详情、新增 / 编辑 / 删除；
评论功能：发布评论、查看文章评论；
路由守卫：未登录用户无法访问私有页面；
项目结构：遵循 Vue 官方最佳实践，分模块管理页面、组件、请求、状态。

二、完整前端项目实现
1. 项目初始化（先执行命令）
   bash
   运行
# 创建Vue3项目（命名为blog-frontend）
npm create vite@latest blog-frontend -- --template vue
cd blog-frontend

# 安装核心依赖
npm install element-plus axios pinia vue-router@4
npm install @element-plus/icons-vue # 图标依赖

# 启动项目（测试）
npm run dev
2. 项目最终结构
   plaintext
   blog-frontend/
   ├── package.json
   ├── vite.config.js       # Vite配置
   ├── index.html
   ├── src/
   │   ├── main.js          # 入口文件（注册插件、路由、Pinia）
   │   ├── router/index.js  # 路由配置
   │   ├── store/index.js   # Pinia状态管理（Token）
   │   ├── api/             # 请求封装
   │   │   ├── index.js     # Axios实例配置
   │   │   ├── auth.js      # 用户认证接口
   │   │   ├── post.js      # 文章接口
   │   │   └── comment.js   # 评论接口
   │   ├── views/           # 页面组件
   │   │   ├── Login.vue    # 登录页
   │   │   ├── Register.vue # 注册页
   │   │   ├── PostList.vue # 文章列表页
   │   │   ├── PostDetail.vue # 文章详情页
   │   │   └── PostEdit.vue # 新增/编辑文章页
   │   ├── components/      # 公共组件
   │   │   ├── NavBar.vue   # 顶部导航栏
   │   │   └── CommentList.vue # 评论列表组件
   │   └── utils/           # 工具函数
   │       └── auth.js      # Token存储/获取工具
   └── public/
   └── favicon.ico