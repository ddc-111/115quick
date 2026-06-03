// Background script: 处理消息传递和剪贴板监控

interface PushMessage {
  type: 'PUSH_MAGNET';
  magnet: string;
  tabId?: number;
}

interface StorageData {
  serverUrl?: string;
}

const serverUrlCache: { [tabId: number]: string } = {};

chrome.runtime.onMessage.addListener((message: PushMessage, sender, sendResponse) => {
  if (message.type === 'PUSH_MAGNET') {
    handlePushMagnet(message.magnet, sender.tab?.id)
      .then(result => sendResponse(result))
      .catch(error => sendResponse({ success: false, error: error.message }));
    return true;
  }
});

async function handlePushMagnet(magnet: string, tabId?: number): Promise<{ success: boolean; message?: string; error?: string }> {
  try {
    const data = await chrome.storage.local.get(['serverUrl']);
    const serverUrl = data.serverUrl;

    if (!serverUrl) {
      return { success: false, error: '请先在115Quick插件中设置服务器地址' };
    }

    const response = await fetch(`${serverUrl}/v1/Download/api/addDownloadLink`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ downloadLink: magnet })
    });

    const result = await response.json();

    if (response.ok) {
      if (tabId) {
        chrome.tabs.sendMessage(tabId, { type: 'PUSH_SUCCESS', magnet });
      }
      return { success: true, message: '磁力链接已推送到115下载' };
    } else {
      throw new Error(result.message || '推送失败');
    }
  } catch (error: any) {
    if (tabId) {
      chrome.tabs.sendMessage(tabId, { type: 'PUSH_ERROR', error: error.message });
    }
    return { success: false, error: error.message };
  }
}

chrome.runtime.onInstalled.addListener(() => {
  console.log('115Quick extension installed');
});
