<template>
  <div class="cloud-files page-container">
    <div class="page-header">
      <h1 class="page-title">115云文件管理</h1>
      <div>
        <el-button @click="refreshFiles" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button type="primary" @click="showCreateFolder">
          新建文件夹
        </el-button>
        <el-button v-if="currentFolder !== '0'" @click="goBack">
          返回上级
        </el-button>
      </div>
    </div>

    <!-- 当前路径 -->
    <div class="card">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item @click="navigateTo('0')">根目录</el-breadcrumb-item>
        <el-breadcrumb-item v-for="path in pathHistory" :key="path.id" @click="navigateTo(path.id)">
          {{ path.name }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 设置为默认下载目录 -->
    <div class="card" v-if="currentFolder !== '0'">
      <el-button type="success" @click="setAsDefaultDir">
        设为默认下载目录
      </el-button>
      <span v-if="defaultDir" style="margin-left: 16px; color: #67c23a;">
        当前默认: {{ defaultDir.folderName }}
      </span>
    </div>

    <!-- 文件列表 -->
    <div class="card">
      <el-table :data="files" style="width: 100%" v-loading="loading">
        <el-table-column label="文件名" min-width="300">
          <template #default="{ row }">
            <div class="file-name" @click="row.isDir ? navigateTo(row.fileId) : null">
              <el-icon v-if="row.isDir" color="#409eff"><Folder /></el-icon>
              <el-icon v-else color="#909399"><Document /></el-icon>
              <span :class="{ 'is-dir': row.isDir }">{{ row.fileName }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="showRename(row)">重命名</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="files.length === 0 && !loading" description="空文件夹" />
    </div>

    <!-- 创建文件夹对话框 -->
    <el-dialog v-model="createFolderVisible" title="新建文件夹" width="400">
      <el-input v-model="newFolderName" placeholder="请输入文件夹名称" />
      <template #footer>
        <el-button @click="createFolderVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateFolder">确定</el-button>
      </template>
    </el-dialog>

    <!-- 重命名对话框 -->
    <el-dialog v-model="renameVisible" title="重命名" width="400">
      <el-input v-model="newFileName" placeholder="请输入新名称" />
      <template #footer>
        <el-button @click="renameVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRename">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Folder, Document } from '@element-plus/icons-vue'
import { getCloudFiles, createFolder, deleteFile, renameFile, setDefaultDownloadDir, getDefaultDownloadDir } from '@/api/server'

interface CloudFile {
  fileId: string
  fileName: string
  fileSize: number
  isDir: boolean
  parentId: string
  updateTime: string
  pickCode: string
}

const loading = ref(false)
const files = ref<CloudFile[]>([])
const currentFolder = ref('0')
const pathHistory = ref<{ id: string; name: string }[]>([])

const createFolderVisible = ref(false)
const newFolderName = ref('')

const renameVisible = ref(false)
const newFileName = ref('')
const currentFile = ref<CloudFile | null>(null)

const defaultDir = ref<{ folderId: string; folderName: string } | null>(null)

async function refreshFiles() {
  loading.value = true
  try {
    const data = await getCloudFiles(currentFolder.value) as any
    files.value = data.files || []
  } catch (error) {
    console.error('获取文件列表失败:', error)
    ElMessage.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

async function loadDefaultDir() {
  try {
    const data = await getDefaultDownloadDir() as any
    if (data.folderId) {
      defaultDir.value = data
    }
  } catch (error) {
    console.error('获取默认目录失败:', error)
  }
}

function navigateTo(folderId: string) {
  if (folderId === '0') {
    currentFolder.value = '0'
    pathHistory.value = []
  } else {
    const file = files.value.find(f => f.fileId === folderId)
    if (file) {
      const existingIndex = pathHistory.value.findIndex(p => p.id === folderId)
      if (existingIndex >= 0) {
        pathHistory.value = pathHistory.value.slice(0, existingIndex + 1)
      } else {
        pathHistory.value.push({ id: folderId, name: file.fileName })
      }
      currentFolder.value = folderId
    }
  }
  refreshFiles()
}

function goBack() {
  if (pathHistory.value.length > 0) {
    pathHistory.value.pop()
    currentFolder.value = pathHistory.value.length > 0 ? pathHistory.value[pathHistory.value.length - 1].id : '0'
    refreshFiles()
  }
}

function showCreateFolder() {
  newFolderName.value = ''
  createFolderVisible.value = true
}

async function handleCreateFolder() {
  if (!newFolderName.value.trim()) {
    ElMessage.warning('请输入文件夹名称')
    return
  }
  try {
    await createFolder(currentFolder.value, newFolderName.value.trim())
    ElMessage.success('创建成功')
    createFolderVisible.value = false
    refreshFiles()
  } catch (error) {
    ElMessage.error('创建失败')
  }
}

function showRename(file: CloudFile) {
  currentFile.value = file
  newFileName.value = file.fileName
  renameVisible.value = true
}

async function handleRename() {
  if (!currentFile.value || !newFileName.value.trim()) {
    ElMessage.warning('请输入新名称')
    return
  }
  try {
    await renameFile(currentFile.value.fileId, newFileName.value.trim())
    ElMessage.success('重命名成功')
    renameVisible.value = false
    refreshFiles()
  } catch (error) {
    ElMessage.error('重命名失败')
  }
}

async function handleDelete(file: CloudFile) {
  try {
    await ElMessageBox.confirm(`确定要删除 "${file.fileName}" 吗？`, '确认删除', {
      type: 'warning'
    })
    await deleteFile(file.fileId)
    ElMessage.success('删除成功')
    refreshFiles()
  } catch {
    // 取消
  }
}

async function setAsDefaultDir() {
  const currentPathName = pathHistory.value.length > 0 ? pathHistory.value[pathHistory.value.length - 1].name : '根目录'
  try {
    await setDefaultDownloadDir(currentFolder.value, currentPathName)
    ElMessage.success('已设为默认下载目录')
    loadDefaultDir()
  } catch (error) {
    ElMessage.error('设置失败')
  }
}

onMounted(() => {
  refreshFiles()
  loadDefaultDir()
})
</script>

<style lang="scss" scoped>
.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;

  .is-dir {
    color: #409eff;
  }
}
</style>
