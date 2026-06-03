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
          <el-tag :type="smbStatus.isConnected ? 'success' : 'info'" size="large">
            {{ smbStatus.isConnected ? '已连接' : smbStatus.enabled ? '已启用未连接' : '未启用' }}
          </el-tag>
        </el-form-item>
      </el-form>
    </div>

    <!-- SMB 文件浏览器 -->
    <div class="card" v-if="smbStatus?.isConnected">
      <h3 style="margin-bottom: 16px">SMB 文件浏览器</h3>
      <div class="smb-browser">
        <!-- 当前路径 -->
        <div class="smb-path">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item @click="smbNavigateTo('')">根目录</el-breadcrumb-item>
            <el-breadcrumb-item v-for="(part, index) in smbPathParts" :key="index" 
              @click="smbNavigateToPart(index)">
              {{ part }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <!-- 文件列表 -->
        <el-table :data="smbFiles" style="width: 100%" v-loading="smbLoading" max-height="400">
          <el-table-column label="文件名" min-width="300">
            <template #default="{ row }">
              <div class="file-name" @click="row.isDir ? smbNavigateTo(row.name) : null">
                <el-icon v-if="row.isDir" color="#409eff"><Folder /></el-icon>
                <el-icon v-else color="#909399"><Document /></el-icon>
                <span :class="{ 'is-dir': row.isDir }">{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="大小" width="120">
            <template #default="{ row }">
              {{ row.isDir ? '-' : formatFileSize(row.size) }}
            </template>
          </el-table-column>
          <el-table-column label="修改时间" width="180">
            <template #default="{ row }">
              {{ row.modTime }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="{ row }">
              <el-button v-if="!row.isDir" size="small" type="primary" @click="handleSMBDownload(row)">
                下载
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
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
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder, Document } from '@element-plus/icons-vue'
import { useServerStore } from '@/stores/server'
import { 
  setDownloadMode, 
  startRename, 
  getSMBConfig, 
  setSMBConfig, 
  testSMBConnection, 
  getServerInfo,
  smbBrowse,
  smbDownload
} from '@/api/server'

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
  password: ''
})

const smbStatus = ref<{
  enabled: boolean
  isConnected: boolean
} | null>(null)

const savingSMB = ref(false)
const testingSMB = ref(false)

// SMB 文件浏览器
const smbLoading = ref(false)
const smbFiles = ref<any[]>([])
const smbCurrentPath = ref('')

const smbPathParts = computed(() => {
  if (!smbCurrentPath.value || smbCurrentPath.value === '.') return []
  return smbCurrentPath.value.split('/').filter(p => p)
})

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
      getSMBConfig().catch(() => null),
      getServerInfo().catch(() => null)
    ])

    // SMB配置
    if (smbData) {
      smbForm.value = {
        enabled: (smbData as any).enabled || false,
        host: (smbData as any).host || '',
        share: (smbData as any).share || '',
        username: (smbData as any).username || '',
        password: (smbData as any).password || ''
      }
      smbStatus.value = {
        enabled: (smbData as any).enabled || false,
        isConnected: (smbData as any).isConnected || false
      }

      // 如果已连接，加载文件列表
      if ((smbData as any).isConnected) {
        loadSMBFiles()
      }
    }

    // 服务器版本信息
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

async function handleSaveSMB() {
  savingSMB.value = true
  try {
    await setSMBConfig(smbForm.value)
    ElMessage.success('SMB 配置已保存')
    const data = await getSMBConfig() as any
    smbStatus.value = {
      enabled: data.enabled,
      isConnected: data.isConnected
    }
    // 如果连接成功，加载文件列表
    if (data.isConnected) {
      loadSMBFiles()
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

async function loadSMBFiles() {
  smbLoading.value = true
  try {
    const data = await smbBrowse(smbCurrentPath.value) as any
    smbFiles.value = data.files || []
  } catch (error: any) {
    ElMessage.error(error.message || '加载SMB文件失败')
  } finally {
    smbLoading.value = false
  }
}

function smbNavigateTo(path: string) {
  smbCurrentPath.value = path
  loadSMBFiles()
}

function smbNavigateToPart(index: number) {
  const parts = smbPathParts.value.slice(0, index + 1)
  smbCurrentPath.value = parts.join('/')
  loadSMBFiles()
}

async function handleSMBDownload(file: any) {
  const remotePath = smbCurrentPath.value 
    ? `${smbCurrentPath.value}/${file.name}` 
    : file.name
  
  try {
    await ElMessageBox.confirm(`确定要下载 "${file.name}" 吗？`, '确认下载', {
      type: 'info'
    })
    
    const data = await smbDownload(remotePath) as any
    ElMessage.success(`文件已下载到: ${data.filePath}`)
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '下载失败')
    }
  }
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
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

<style lang="scss" scoped>
.smb-browser {
  .smb-path {
    margin-bottom: 16px;
    padding: 8px;
    background: #f5f7fa;
    border-radius: 4px;
  }
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;

  .is-dir {
    color: #409eff;
    font-weight: bold;
  }
}
</style>
