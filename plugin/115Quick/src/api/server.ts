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

export function startRename(params: {
  mode?: number
  folder_id?: string
  file_pattern?: string
  delete_others?: boolean
  rename_match?: boolean
}) {
  return request.post('/api/StartReName', params)
}
