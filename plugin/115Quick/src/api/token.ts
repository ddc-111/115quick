import request from './request'

export function setToken(accessToken: string, refreshToken: string) {
  return request.post('/api/setToken', { accessToken, refreshToken })
}

export function getTokenStatus() {
  return request.get('/api/getTokenStatus')
}

export function getDownloadProgress() {
  return request.get('/api/getDownloadProgress')
}

export function getCloudTasks() {
  return request.get('/api/getCloudTasks')
}

export function getTaskHistory(page: number = 1, pageSize: number = 20) {
  return request.get('/api/getTaskHistory', { params: { page, pageSize } })
}

export function removeDownloadTask(url: string) {
  return request.post('/api/removeDownloadTask', { url })
}

export function refreshTasks() {
  return request.post('/api/refreshTasks')
}
