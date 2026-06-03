<template>
  <div class="download-list page-container">
    <div class="page-header">
      <h1 class="page-title">下载列表</h1>
      <div>
        <el-button @click="handleRefresh" :loading="refreshing">
          <el-icon><Refresh /></el-icon>
          刷新任务
        </el-button>
        <el-button type="warning" @click="handleClearCompleted">
          清空已完成
        </el-button>
      </div>
    </div>

    <!-- 本地下载进度 -->
    <div class="card" v-if="downloadProgress.length > 0">
      <h3 style="margin-bottom: 16px">本地下载 ({{ downloadProgress.length }})</h3>
      <div v-for="item in downloadProgress" :key="item.fileName" class="task-item">
        <div class="task-header">
          <span class="task-name">{{ item.fileName }}</span>
          <span class="speed">{{ formatSpeed(item.speed * 1024 * 1024) }}</span>
        </div>
        <el-progress :percentage="Math.round(item.percent)" />
        <div class="task-meta">
          <span>已下载: {{ formatFileSize(item.downloaded) }}</span>
          <span>总大小: {{ formatFileSize(item.fileSize) }}</span>
        </div>
      </div>
    </div>

    <!-- 云下载任务 -->
    <div class="card">
      <div class="section-header">
        <h3>云下载任务 ({{ filteredCloudTasks.length }})</h3>
        <div class="filter-bar">
          <el-select v-model="cloudTaskFilter" placeholder="筛选状态" clearable style="width: 120px">
            <el-option label="全部" value="" />
            <el-option label="下载中" value="1" />
            <el-option label="等待中" value="0" />
            <el-option label="已完成" value="2" />
            <el-option label="失败" value="-1" />
          </el-select>
        </div>
      </div>
      
      <div v-if="filteredCloudTasks.length === 0" class="empty-state">
        <div class="empty-icon">☁️</div>
        <p>暂无云下载任务</p>
      </div>
      <div v-else>
        <div v-for="task in filteredCloudTasks" :key="task.infoHash" class="task-item">
          <div class="task-header">
            <span class="task-name">{{ task.name || task.url }}</span>
            <div>
              <el-tag :type="getStatusType(task.status)" size="small" style="margin-right: 8px">
                {{ getStatusText(task.status) }}
              </el-tag>
              <el-button
                v-if="task.status === 0"
                type="danger"
                size="small"
                @click="handleRemove(task.url)"
              >
                删除
              </el-button>
            </div>
          </div>
          <el-progress
            :percentage="Math.round(task.percentDone)"
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

    <!-- 待下载队列 -->
    <div class="card">
      <div class="section-header">
        <h3>待下载队列 ({{ pendingTasks.length }})</h3>
      </div>
      
      <div v-if="pendingTasks.length === 0" class="empty-state">
        <div class="empty-icon">📋</div>
        <p>暂无待下载任务</p>
      </div>
      <div v-else>
        <div v-for="task in pendingTasks" :key="task.url" class="task-item">
          <div class="task-header">
            <span class="task-name">{{ task.url }}</span>
            <el-button type="danger" size="small" @click="handleRemove(task.url)">
              删除
            </el-button>
          </div>
          <div class="task-meta">
            <span>添加时间: {{ task.addTime }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getServerInfo, clearCompletedTasks } from '@/api/server'
import { getDownloadProgress, getCloudTasks, removeDownloadTask, refreshTasks } from '@/api/token'
import { formatFileSize, formatSpeed, getStatusText, getStatusType } from '@/utils/format'

const refreshing = ref(false)
const downloadProgress = ref<any[]>([])
const cloudTasks = ref<any[]>([])
const pendingTasks = ref<any[]>([])
const cloudTaskFilter = ref('')
let timer: ReturnType<typeof setInterval> | null = null

// 筛选后的云下载任务
const filteredCloudTasks = computed(() => {
  if (!cloudTaskFilter.value) return cloudTasks.value
  const filterStatus = parseInt(cloudTaskFilter.value)
  return cloudTasks.value.filter(task => task.status === filterStatus)
})

async function loadData() {
  try {
    const [progressData, tasksData, serverData] = await Promise.all([
      getDownloadProgress(),
      getCloudTasks(),
      getServerInfo()
    ])
    downloadProgress.value = (progressData as any).downloads || []
    cloudTasks.value = (tasksData as any).tasks || []
    pendingTasks.value = (serverData as any).downFileInfoList || []
  } catch (error) {
    console.error('加载数据失败:', error)
  }
}

async function handleRefresh() {
  refreshing.value = true
  try {
    await refreshTasks()
    ElMessage.success('已触发刷新')
    setTimeout(loadData, 2000)
  } catch (error) {
    ElMessage.error('刷新失败')
  } finally {
    refreshing.value = false
  }
}

async function handleRemove(url: string) {
  try {
    await ElMessageBox.confirm('确定要删除这个任务吗？', '确认', {
      type: 'warning'
    })
    await removeDownloadTask(url)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

async function handleClearCompleted() {
  try {
    await ElMessageBox.confirm('确定要清空所有已完成的任务吗？', '确认', {
      type: 'warning'
    })
    await clearCompletedTasks()
    ElMessage.success('已完成任务已清空')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清空失败')
    }
  }
}

onMounted(() => {
  loadData()
  timer = setInterval(loadData, 10000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style lang="scss" scoped>
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  h3 {
    margin: 0;
  }
}

.filter-bar {
  display: flex;
  gap: 8px;
}

.speed {
  color: var(--el-color-primary);
  font-weight: 600;
}
</style>
