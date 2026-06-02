<template>
  <div class="task-history page-container">
    <div class="page-header">
      <h1 class="page-title">任务历史</h1>
    </div>

    <div class="card">
      <el-table
        :data="historyItems"
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
import { ref, onMounted } from 'vue'
import { getTaskHistory } from '@/api/token'
import { formatFileSize, getStatusText, getStatusType } from '@/utils/format'

const historyItems = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

async function loadHistory() {
  try {
    const data = await getTaskHistory(currentPage.value, pageSize.value) as any
    historyItems.value = data.items || []
    total.value = data.total || 0
  } catch (error) {
    console.error('加载历史记录失败:', error)
  }
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

onMounted(() => {
  loadHistory()
})
</script>

<style lang="scss" scoped>
.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
