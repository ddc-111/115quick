<template>
  <div class="dashboard page-container">
    <div class="page-header">
      <h1 class="page-title">主面板</h1>
      <el-button @click="refreshData" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <!-- 状态卡片 -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="8">
        <div class="card stat-card">
          <div class="stat-icon" style="background: rgba(103, 194, 58, 0.2)">
            <el-icon color="#67c23a"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">服务器状态</div>
            <div class="stat-value" style="color: #67c23a">已连接</div>
          </div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="card stat-card">
          <div class="stat-icon" :style="{ background: tokenBgColor }">
            <el-icon :color="tokenColor"><Key /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">Token状态</div>
            <div class="stat-value" :style="{ color: tokenColor }">{{ tokenStatus }}</div>
          </div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="card stat-card">
          <div class="stat-icon" style="background: rgba(64, 158, 255, 0.2)">
            <el-icon color="#409eff"><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">待下载任务</div>
            <div class="stat-value">{{ pendingCount }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 快速添加 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">快速添加下载链接</h3>
      <el-input
        v-model="downloadLink"
        placeholder="粘贴磁力链接或HTTP链接..."
        size="large"
        @keyup.enter="handleAddLink"
      >
        <template #append>
          <el-button type="primary" @click="handleAddLink" :loading="adding">
            添加
          </el-button>
        </template>
      </el-input>
    </div>

    <!-- 当前任务 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">当前任务 ({{ cloudTasks.length }})</h3>
      <div v-if="cloudTasks.length === 0" class="empty-state">
        <div class="empty-icon">📭</div>
        <p>暂无下载任务</p>
      </div>
      <div v-else>
        <div v-for="task in cloudTasks" :key="task.infoHash" class="task-item">
          <div class="task-header">
            <span class="task-name">{{ task.name || task.url }}</span>
            <el-tag :type="getStatusType(task.status)" size="small">
              {{ getStatusText(task.status) }}
            </el-tag>
          </div>
          <el-progress
            :percentage="task.percentDone"
            :status="task.status === 2 ? 'success' : task.status === -1 ? 'exception' : ''"
          />
          <div class="task-meta">
            <span v-if="task.size > 0">大小: {{ formatFileSize(task.size) }}</span>
            <span v-if="task.rateDownload > 0">速度: {{ formatSpeed(task.rateDownload) }}</span>
            <span v-if="task.addTime">添加: {{ task.addTime }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, CircleCheck, Key, Document } from '@element-plus/icons-vue'
import { getServerInfo, addDownloadLink } from '@/api/server'
import { getTokenStatus, getCloudTasks } from '@/api/token'
import { formatFileSize, formatSpeed, getStatusText, getStatusType } from '@/utils/format'

const loading = ref(false)
const adding = ref(false)
const downloadLink = ref('')
const tokenInfo = ref<any>({})
const cloudTasks = ref<any[]>([])
const pendingCount = ref(0)

const tokenStatus = computed(() => {
  if (!tokenInfo.value.configured) return '未配置'
  if (!tokenInfo.value.valid) return '已过期'
  return '有效'
})

const tokenColor = computed(() => {
  if (!tokenInfo.value.configured) return '#909399'
  if (!tokenInfo.value.valid) return '#f56c6c'
  return '#67c23a'
})

const tokenBgColor = computed(() => {
  if (!tokenInfo.value.configured) return 'rgba(144, 147, 153, 0.2)'
  if (!tokenInfo.value.valid) return 'rgba(245, 108, 108, 0.2)'
  return 'rgba(103, 194, 58, 0.2)'
})

async function refreshData() {
  loading.value = true
  try {
    const [serverData, tokenData, tasksData] = await Promise.all([
      getServerInfo(),
      getTokenStatus(),
      getCloudTasks()
    ])
    tokenInfo.value = tokenData
    cloudTasks.value = (tasksData as any).tasks || []
    pendingCount.value = (serverData as any).downFileInfoList?.length || 0
  } catch (error) {
    console.error('刷新数据失败:', error)
    ElMessage.error('刷新数据失败')
  } finally {
    loading.value = false
  }
}

async function handleAddLink() {
  if (!downloadLink.value.trim()) {
    ElMessage.warning('请输入下载链接')
    return
  }

  adding.value = true
  try {
    await addDownloadLink(downloadLink.value.trim())
    ElMessage.success('添加成功')
    downloadLink.value = ''
    refreshData()
  } catch (error) {
    ElMessage.error('添加失败')
  } finally {
    adding.value = false
  }
}

onMounted(() => {
  refreshData()
})
</script>

<style lang="scss" scoped>
.stats-row {
  margin-bottom: 16px;
}
</style>
