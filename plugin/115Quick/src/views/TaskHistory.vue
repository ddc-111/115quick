<template>
  <div class="task-history page-container">
    <div class="page-header">
      <h1 class="page-title">任务历史</h1>
      <div>
        <el-button type="danger" @click="handleClearAll">
          清空历史
        </el-button>
      </div>
    </div>

    <div class="card">
      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-select v-model="statusFilter" placeholder="筛选状态" clearable style="width: 120px" @change="handleFilterChange">
          <el-option label="全部" value="" />
          <el-option label="已完成" value="2" />
          <el-option label="下载中" value="1" />
          <el-option label="等待中" value="0" />
          <el-option label="失败" value="3" />
        </el-select>
      </div>

      <el-table
        :data="filteredHistoryItems"
        style="width: 100%"
        :header-cell-style="{ background: 'var(--el-bg-color-overlay)', color: 'var(--el-text-color-primary)' }"
        :cell-style="{ background: 'var(--el-bg-color-overlay)', color: 'var(--el-text-color-primary)' }"
      >
        <el-table-column prop="name" label="任务名称" min-width="200">
          <template #default="{ row }">
            <span>{{ row.name || row.url }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="120">
          <template #default="{ row }">
            {{ row.size > 0 ? formatFileSize(row.size) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="progress" label="进度" width="120">
          <template #default="{ row }">
            <el-progress
              v-if="row.status === 1"
              :percentage="Math.round(row.progress)"
              :stroke-width="8"
            />
            <span v-else-if="row.status === 2">100%</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="addTime" label="添加时间" width="180" />
        <el-table-column prop="completeTime" label="完成时间" width="180">
          <template #default="{ row }">
            {{ row.completeTime || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="errorMsg" label="错误信息" min-width="150">
          <template #default="{ row }">
            <span v-if="row.errorMsg" style="color: var(--el-color-danger)">{{ row.errorMsg }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTaskHistory } from '@/api/token'
import { clearTaskHistory } from '@/api/server'
import { formatFileSize, getStatusText, getStatusType } from '@/utils/format'

const historyItems = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const statusFilter = ref('')

// 筛选后的历史记录
const filteredHistoryItems = computed(() => {
  if (!statusFilter.value) return historyItems.value
  const filterStatus = parseInt(statusFilter.value)
  return historyItems.value.filter(item => item.status === filterStatus)
})

async function loadHistory() {
  try {
    const data = await getTaskHistory(currentPage.value, pageSize.value) as any
    historyItems.value = data.items || []
    total.value = data.total || 0
  } catch (error) {
    console.error('加载历史记录失败:', error)
  }
}

function handleFilterChange() {
  // 筛选变化时不需要重新加载，使用 computed 自动过滤
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  loadHistory()
}

function handleCurrentChange(page: number) {
  currentPage.value = page
  loadHistory()
}

async function handleClearAll() {
  try {
    await ElMessageBox.confirm('确定要清空所有任务历史吗？此操作不可恢复！', '确认', {
      type: 'warning'
    })
    await clearTaskHistory()
    ElMessage.success('任务历史已清空')
    loadHistory()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清空失败')
    }
  }
}

onMounted(() => {
  loadHistory()
})
</script>

<style lang="scss" scoped>
.filter-bar {
  margin-bottom: 16px;
  display: flex;
  gap: 8px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
