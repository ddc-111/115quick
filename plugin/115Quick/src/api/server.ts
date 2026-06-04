import request from './request'

export function getServerInfo() {
  return request.get('/api/getServerInfo')
}

export function addDownloadLink(downloadLink: string) {
  return request.post('/api/addDownloadLink', { downloadLink })
}

export function setDownloadMode(mode: number) {
  return request.post('/api/setDownloadMode', { mode })
}

export function getServerLogs(type: string = 'stdout', lines: number = 200) {
  return request.get('/api/getServerLogs', { params: { type, lines } })
}

export function getCloudFiles(folderId: string = '0', fileType: number = 0) {
  return request.get('/api/getCloudFiles', { params: { folderId, fileType } })
}

export function createFolder(parentId: string, folderName: string) {
  return request.post('/api/createFolder', { parentId, folderName })
}

export function deleteFile(fileIds: string) {
  return request.post('/api/deleteFile', { fileIds })
}

export function renameFile(fileId: string, newName: string) {
  return request.post('/api/renameFile', { fileId, newName })
}

export function moveFile(fileIds: string, targetDir: string) {
  return request.post('/api/moveFile', { fileIds, targetDir })
}

export function setDefaultDownloadDir(folderId: string, folderName: string) {
  return request.post('/api/setDefaultDownloadDir', { folderId, folderName })
}

export function getDefaultDownloadDir() {
  return request.get('/api/getDefaultDownloadDir')
}

export function clearTaskHistory() {
  return request.post('/api/clearTaskHistory')
}

export function clearCompletedTasks() {
  return request.post('/api/clearCompletedTasks')
}

export function startRename(params: {
  mode?: number
  folder_id?: string
  file_pattern?: string
  delete_others?: boolean
  rename_match?: boolean
}) {
  return request.post('/api/StartReName', params)
}
