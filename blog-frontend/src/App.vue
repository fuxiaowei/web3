<template>
  <router-view></router-view> <!-- 路由出口，所有页面都在这里渲染 -->
</template>

<script setup>
// 这里可以留空，或者添加全局逻辑
</script>

<style>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
}
</style>