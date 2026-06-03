<template>
  <div class="log-viewer page-container">
    <div class="page-header">
      <h1 class="page-title">服务器日志</h1>
      <div>
        <el-button @click="refreshLogs" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button @click="clearLogs" type="danger">
          清空
        </el-button>
      </div>
    </div>

    <div class="card">
      <div class="log-controls">
        <el-radio-group v-model="logType" @change="refreshLogs">
          <el-radio-button value="stdout">标准输出</el-radio-button>
          <el-radio-button value="stderr">错误日志</el-radio-button>
        </el-radio-group>
        <el-input-number v-model="lines" :min="50" :max="1000" :step="50" @change="refreshLogs" />
        <el-tag>行数</el-tag>
      </div>

      <div class="log-content" ref="logContainer">
        <pre v-if="logs.length > 0">{{ logs.join('\n') }}</pre>
        <el-empty v-else description="暂无日志" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { getServerLogs } from '@/api/server'

const loading = ref(false)
const logType = ref('stdout')
const lines = ref(200)
const logs = ref<string[]>([])
const logContainer = ref<HTMLElement>()

async function refreshLogs() {
  loading.value = true
  try {
    const data = await getServerLogs(logType.value, lines.value) as any
    logs.value = data.logs || []
    await nextTick()
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  } catch (error) {
    console.error('获取日志失败:', error)
  } finally {
    loading.value = false
  }
}

function clearLogs() {
  logs.value = []
}

onMounted(() => {
  refreshLogs()
})
</script>

<style lang="scss" scoped>
.log-controls {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.log-content {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 16px;
  border-radius: 8px;
  max-height: 600px;
  overflow-y: auto;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.6;

  pre {
    margin: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
  }
}
</style>
