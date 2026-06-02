<template>
  <div class="app-layout">
    <div class="logo">
      <h1>115Quick</h1>
    </div>
    <el-menu
      :default-active="activeMenu"
      router
      class="sidebar-menu"
    >
      <el-menu-item index="/dashboard">
        <el-icon><Monitor /></el-icon>
        <span>主面板</span>
      </el-menu-item>
      <el-menu-item index="/token">
        <el-icon><Key /></el-icon>
        <span>Token配置</span>
      </el-menu-item>
      <el-menu-item index="/downloads">
        <el-icon><Download /></el-icon>
        <span>下载列表</span>
      </el-menu-item>
      <el-menu-item index="/history">
        <el-icon><Clock /></el-icon>
        <span>任务历史</span>
      </el-menu-item>
      <el-menu-item index="/settings">
        <el-icon><Setting /></el-icon>
        <span>设置</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Monitor, Key, Download, Clock, Setting } from '@element-plus/icons-vue'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style lang="scss" scoped>
.app-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.logo {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);

  h1 {
    font-size: 24px;
    font-weight: 700;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
}

.sidebar-menu {
  flex: 1;
  border-right: none;
  background-color: transparent;

  .el-menu-item {
    color: var(--el-text-color-regular);

    &:hover {
      background-color: rgba(255, 255, 255, 0.05);
    }

    &.is-active {
      color: var(--el-color-primary);
      background-color: rgba(15, 52, 96, 0.3);
    }
  }
}
</style>
