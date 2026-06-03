// Background script: 处理消息传递、点击打开面板、右键菜单

interface PushMessage {
  type: 'PUSH_MAGNET';
  magnet: string;
  tabId?: number;
}

interface StorageData {
  serverUrl?: string;
}

const PANEL_URL = 'index.html';
const PANEL_TAB_KEY = '115quick_panel_tab_id';

// 点击插件图标打开面板标签页
chrome.action.onClicked.addListener(async () => {
  await openPanel();
});

// 打开或聚焦面板标签页
async function openPanel() {
  const data = await chrome.storage.local.get([PANEL_TAB_KEY]);
  const panelTabId = data[PANEL_TAB_KEY];

  if (panelTabId) {
    try {
      const tab = await chrome.tabs.get(panelTabId);
      if (tab) {
        await chrome.tabs.update(panelTabId, { active: true });
        await chrome.windows.update(tab.windowId!, { focused: true });
        return;
      }
    } catch {
      // 标签页已关闭
    }
  }

  const tab = await chrome.tabs.create({
    url: chrome.runtime.getURL(PANEL_URL)
  });
  await chrome.storage.local.set({ [PANEL_TAB_KEY]: tab.id });
}

// 监听标签页关闭，清理存储
chrome.tabs.onRemoved.addListener(async (tabId) => {
  const data = await chrome.storage.local.get([PANEL_TAB_KEY]);
  if (data[PANEL_TAB_KEY] === tabId) {
    await chrome.storage.local.remove(PANEL_TAB_KEY);
  }
});

// 创建右键菜单
chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.create({
    id: 'open-115quick-panel',
    title: '打开115Quick管理面板',
    contexts: ['action']
  });

  chrome.contextMenus.create({
    id: 'open-115quick-settings',
    title: '打开115Quick设置',
    contexts: ['action']
  });

  console.log('115Quick extension installed');
});

// 右键菜单点击处理
chrome.contextMenus.onClicked.addListener(async (info) => {
  if (info.menuItemId === 'open-115quick-panel') {
    await openPanel();
  } else if (info.menuItemId === 'open-115quick-settings') {
    await openPanel();
    // 延迟发送消息跳转到设置页
    setTimeout(async () => {
      const data = await chrome.storage.local.get([PANEL_TAB_KEY]);
      if (data[PANEL_TAB_KEY]) {
        chrome.tabs.sendMessage(data[PANEL_TAB_KEY], { type: 'NAVIGATE_TO', path: '/settings' });
      }
    }, 500);
  }
});

// 消息处理
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
