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

    <!-- SMB 网络存储配置 -->
    <div class="card">
      <h3 style="margin-bottom: 16px">SMB 网络存储</h3>
      <el-form :model="smbForm" label-position="top">
        <el-form-item>
          <el-switch v-model="smbForm.enabled" active-text="启用 SMB" inactive-text="禁用" />
        </el-form-item>
        
        <template v-if="smbForm.enabled">
          <el-form-item label="SMB 服务器地址">
            <el-input v-model="smbForm.host" placeholder="例如: 192.168.1.100 或 NAS.local" />
          </el-form-item>
          <el-form-item label="共享名称">
            <el-input v-model="smbForm.share" placeholder="例如: downloads" />
          </el-form-item>
          <el-form-item label="用户名">
            <el-input v-model="smbForm.username" placeholder="留空表示匿名访问" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="smbForm.password" type="password" show-password placeholder="留空表示无密码" />
          </el-form-item>
          <el-form-item label="挂载点">
            <el-input v-model="smbForm.mountPoint" placeholder="Linux/Mac: /mnt/smb, Windows: Z:" />
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" :loading="savingSMB" @click="handleSaveSMB">
            保存配置
          </el-button>
          <el-button v-if="smbForm.enabled" :loading="testingSMB" @click="handleTestSMB">
            测试连接
          </el-button>
        </el-form-item>

        <el-form-item v-if="smbStatus">
          <el-tag :type="smbStatus.isMounted ? 'success' : 'info'" size="large">
            {{ smbStatus.isMounted ? '已挂载' : smbStatus.enabled ? '已启用未挂载' : '未启用' }}
          </el-tag>
        </el-form-item>
      </el-form>
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
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useServerStore } from '@/stores/server'
import { setDownloadMode, startRename, getSMBConfig, setSMBConfig, testSMBConnection, getServerInfo } from '@/api/server'

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

const smbForm = ref({
  enabled: false,
  host: '',
  share: '',
  username: '',
  password: '',
  mountPoint: ''
})

const smbStatus = ref<{
  enabled: boolean
  isMounted: boolean
} | null>(null)

const savingSMB = ref(false)
const testingSMB = ref(false)

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
    const [smbData, serverData] = await Promise.all([
      getSMBConfig(),
      getServerInfo()
    ])

    // SMB配置
    smbForm.value = {
      enabled: smbData.enabled,
      host: smbData.host,
      share: smbData.share,
      username: smbData.username,
      password: smbData.password,
      mountPoint: smbData.mountPoint
    }
    smbStatus.value = {
      enabled: smbData.enabled,
      isMounted: smbData.isMounted
    }

    // 服务器版本信息
    if (serverData.version) {
      serverVersion.value = serverData.version
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

async function handleSaveSMB() {
  savingSMB.value = true
  try {
    await setSMBConfig(smbForm.value)
    ElMessage.success('SMB 配置已保存')
    const { data } = await getSMBConfig()
    smbStatus.value = {
      enabled: data.enabled,
      isMounted: data.isMounted
    }
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    savingSMB.value = false
  }
}

async function handleTestSMB() {
  testingSMB.value = true
  try {
    await testSMBConnection({
      host: smbForm.value.host,
      share: smbForm.value.share,
      username: smbForm.value.username,
      password: smbForm.value.password
    })
    ElMessage.success('SMB 连接测试成功')
  } catch (error: any) {
    ElMessage.error(error.message || '连接测试失败')
  } finally {
    testingSMB.value = false
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
