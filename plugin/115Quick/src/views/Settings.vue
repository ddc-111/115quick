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

    <!-- 重命名工具 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">重命名工具</h3>
      <el-form :model="renameForm" label-position="top">
        <el-form-item label="重命名模式">
          <el-radio-group v-model="renameForm.mode">
            <el-radio :value="0">全量重命名（递归所有文件夹）</el-radio>
            <el-radio :value="1">增量重命名（只处理指定文件夹）</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="文件夹ID（增量模式）">
          <el-input v-model="renameForm.folder_id" placeholder="留空则从根目录开始" />
        </el-form-item>
        <el-form-item label="文件名匹配正则">
          <el-input v-model="renameForm.file_pattern" placeholder="默认: (?i)[a-z]+-\d+" />
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="renameForm.delete_others">删除其他文件（只保留最大的）</el-checkbox>
        </el-form-item>
        <el-form-item>
          <el-checkbox v-model="renameForm.rename_match">将文件重命名为匹配项</el-checkbox>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="renaming"
            @click="handleRename"
          >
            执行重命名
          </el-button>
        </el-form-item>
      </el-form>
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
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useServerStore } from '@/stores/server'
import { setDownloadMode, startRename } from '@/api/server'

const serverStore = useServerStore()
const downloadMode = ref(0)
const renaming = ref(false)
const renameForm = ref({
  mode: 0,
  folder_id: '',
  file_pattern: '',
  delete_others: false,
  rename_match: false
})

async function handleModeChange(mode: number) {
  try {
    await setDownloadMode(mode)
    ElMessage.success('下载模式已更新')
  } catch (error) {
    ElMessage.error('设置失败')
  }
}

async function handleRename() {
  renaming.value = true
  try {
    await startRename(renameForm.value)
    ElMessage.success('重命名任务已启动')
  } catch (error) {
    ElMessage.error('重命名失败')
  } finally {
    renaming.value = false
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
