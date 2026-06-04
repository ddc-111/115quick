<template>
  <div class="settings page-container">
    <div class="page-header">
      <h1 class="page-title">设置</h1>
    </div>

    <!-- 服务器信息 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">服务器信息</h3>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="服务器地址">
          {{ serverStore.serverUrl }}
        </el-descriptions-item>
        <el-descriptions-item label="连接状态">
          <el-tag type="success">已连接</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="服务器版本">
          <span>{{ serverVersion.version || '获取中...' }}</span>
          <el-tag v-if="serverVersion.hasUpdate" type="warning" size="small" style="margin-left: 8px">
            有新版本 {{ serverVersion.latestVersion }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="serverVersion.gitCommit" label="Git Commit">
          {{ serverVersion.gitCommit }}
        </el-descriptions-item>
        <el-descriptions-item v-if="serverVersion.buildTime" label="构建时间">
          {{ serverVersion.buildTime }}
        </el-descriptions-item>
        <el-descriptions-item v-if="serverVersion.hasUpdate" label="更新链接">
          <el-link type="primary" :href="serverVersion.updateUrl" target="_blank">
            前往下载新版本
          </el-link>
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <!-- 下载模式 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">下载模式</h3>
      <el-radio-group v-model="downloadMode" @change="handleModeChange">
        <el-radio :value="0">仅下载视频文件</el-radio>
        <el-radio :value="1">下载所有文件</el-radio>
      </el-radio-group>
    </div>

    <!-- 断开连接 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">其他操作</h3>
      <el-button type="danger" @click="handleDisconnect">
        断开服务器连接
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useServerStore } from '@/stores/server'
import { setDownloadMode, getServerInfo } from '@/api/server'

const serverStore = useServerStore()
const downloadMode = ref(0)

const serverVersion = ref<{
  version: string
  gitCommit: string
  buildTime: string
  latestVersion: string
  updateUrl: string
  hasUpdate: boolean
}>({
  version: '',
  gitCommit: '',
  buildTime: '',
  latestVersion: '',
  updateUrl: '',
  hasUpdate: false
})

onMounted(async () => {
  try {
    const serverData = await getServerInfo().catch(() => null)
    if (serverData && (serverData as any).version) {
      serverVersion.value = (serverData as any).version
    }
  } catch (error) {
    console.error('Failed to load config:', error)
  }
})

async function handleModeChange(mode: number) {
  try {
    await setDownloadMode(mode)
    ElMessage.success('下载模式已更新')
  } catch (error) {
    ElMessage.error('设置失败')
  }
}

async function handleDisconnect() {
  try {
    await ElMessageBox.confirm('确定要断开服务器连接吗？', '确认', {
      type: 'warning'
    })
    serverStore.setServerUrl('')
    ElMessage.success('已断开连接')
  } catch {
    // 取消
  }
}
</script>
