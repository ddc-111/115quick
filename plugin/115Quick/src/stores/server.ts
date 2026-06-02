import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useServerStore = defineStore('server', () => {
  const serverUrl = ref('')
  const isConnected = computed(() => serverUrl.value !== '')

  function setServerUrl(url: string) {
    // 移除末尾的斜杠
    serverUrl.value = url.replace(/\/+$/, '')
  }

  function getBaseUrl() {
    return `${serverUrl.value}/v1/Download`
  }

  async function testConnection(url: string): Promise<boolean> {
    try {
      const response = await axios.get(`${url}/v1/Download/api/getServerInfo`, {
        timeout: 5000
      })
      return response.status === 200
    } catch {
      return false
    }
  }

  return {
    serverUrl,
    isConnected,
    setServerUrl,
    getBaseUrl,
    testConnection
  }
}, {
  persist: true
})
