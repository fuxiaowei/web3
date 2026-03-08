import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path' // 导入Node.js内置的path模块

// https://vite.dev/config/
export default defineConfig({
    plugins: [vue()],
    // 配置路径别名
    resolve: {
        alias: {
            // 把 @ 指向 src 目录（关键！）
            '@': path.resolve(__dirname, './src')
        }
    },
    // 配置端口、跨域代理等
    server: {
        port: 3000, // 永久指定端口为3000
        host: '0.0.0.0', // 允许局域网访问（可选）
        // 跨域代理（对接后端8080端口）
        proxy: {
            '/api': {
                target: 'http://localhost:8080', // 后端地址
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/api/, '/api')
            }
        }
    }
})
