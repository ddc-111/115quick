// Content script: 检测页面磁力链并添加悬浮推送按钮

const MAGNET_REGEX = /magnet:\?xt=urn:[a-zA-Z0-9]+:[a-zA-Z0-9]{32,}/i;
const BUTTON_CLASS = 'quick-115-push-btn';
const PROCESSED_ATTR = 'data-115-processed';

// Toast 提示样式
const TOAST_STYLES = `
.quick-115-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  padding: 12px 24px;
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  z-index: 2147483647;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  animation: quick-115-slide-in 0.3s ease-out;
  max-width: 400px;
  word-break: break-all;
}
.quick-115-toast.success {
  background: linear-gradient(135deg, #67C23A, #4CAF50);
}
.quick-115-toast.error {
  background: linear-gradient(135deg, #F56C6C, #E74C3C);
}
.quick-115-toast.info {
  background: linear-gradient(135deg, #409EFF, #2196F3);
}
@keyframes quick-115-slide-in {
  from { transform: translateX(100%); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}
@keyframes quick-115-slide-out {
  from { transform: translateX(0); opacity: 1; }
  to { transform: translateX(100%); opacity: 0; }
}
`;

// 悬浮按钮样式
const BUTTON_STYLES = `
.${BUTTON_CLASS} {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-left: 8px;
  padding: 4px 10px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  cursor: pointer;
  transition: all 0.2s ease;
  vertical-align: middle;
  white-space: nowrap;
  box-shadow: 0 2px 4px rgba(99, 102, 241, 0.3);
}
.${BUTTON_CLASS}:hover {
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  box-shadow: 0 4px 8px rgba(99, 102, 241, 0.4);
  transform: translateY(-1px);
}
.${BUTTON_CLASS}:active {
  transform: translateY(0);
}
.${BUTTON_CLASS}.loading {
  opacity: 0.7;
  cursor: not-allowed;
}
.${BUTTON_CLASS} .icon {
  font-size: 14px;
}
`;

let toastContainer: HTMLDivElement | null = null;

function injectStyles() {
  if (document.getElementById('quick-115-styles')) return;
  
  const style = document.createElement('style');
  style.id = 'quick-115-styles';
  style.textContent = TOAST_STYLES + BUTTON_STYLES;
  document.head.appendChild(style);
}

function showToast(message: string, type: 'success' | 'error' | 'info' = 'info', duration = 3000) {
  if (!toastContainer) {
    toastContainer = document.createElement('div');
    toastContainer.style.cssText = 'position:fixed;top:0;right:0;z-index:2147483647;pointer-events:none;';
    document.body.appendChild(toastContainer);
  }

  const toast = document.createElement('div');
  toast.className = `quick-115-toast ${type}`;
  toast.textContent = message;
  toast.style.pointerEvents = 'auto';
  toastContainer.appendChild(toast);

  setTimeout(() => {
    toast.style.animation = 'quick-115-slide-out 0.3s ease-in forwards';
    setTimeout(() => toast.remove(), 300);
  }, duration);
}

function getServerUrl(): Promise<string> {
  return new Promise((resolve) => {
    chrome.storage.local.get(['serverUrl'], (result) => {
      resolve(result.serverUrl || '');
    });
  });
}

async function pushMagnet(magnetUrl: string, button: HTMLButtonElement) {
  const serverUrl = await getServerUrl();
  if (!serverUrl) {
    showToast('请先在115Quick插件中设置服务器地址', 'error');
    return;
  }

  button.classList.add('loading');
  button.textContent = '推送中...';

  try {
    const response = await fetch(`${serverUrl}/v1/Download/api/addDownloadLink`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ downloadLink: magnetUrl })
    });

    const data = await response.json();
    
    if (response.ok) {
      showToast('磁力链接已推送到115下载', 'success');
      button.textContent = '已推送';
      button.style.background = 'linear-gradient(135deg, #10b981, #059669)';
      setTimeout(() => {
        button.textContent = '推送下载';
        button.style.background = '';
        button.classList.remove('loading');
      }, 2000);
    } else {
      throw new Error(data.message || '推送失败');
    }
  } catch (error: any) {
    showToast(`推送失败: ${error.message}`, 'error');
    button.textContent = '重试';
    button.classList.remove('loading');
  }
}

function createPushButton(magnetUrl: string): HTMLButtonElement {
  const btn = document.createElement('button');
  btn.className = BUTTON_CLASS;
  btn.innerHTML = '<span class="icon">⬇</span> 推送下载';
  btn.title = '推送到115网盘离线下载';
  btn.addEventListener('click', (e) => {
    e.preventDefault();
    e.stopPropagation();
    pushMagnet(magnetUrl, btn);
  });
  return btn;
}

function isMagnetLink(element: Element): boolean {
  const text = element.textContent || '';
  const href = (element as HTMLAnchorElement).href || '';
  return MAGNET_REGEX.test(text) || MAGNET_REGEX.test(href);
}

function getMagnetUrl(element: Element): string {
  const href = (element as HTMLAnchorElement).href || '';
  if (MAGNET_REGEX.test(href)) {
    const match = href.match(MAGNET_REGEX);
    return match ? match[0] : '';
  }
  
  const text = element.textContent || '';
  const match = text.match(MAGNET_REGEX);
  return match ? match[0] : '';
}

function processElement(element: Element) {
  if (element.hasAttribute(PROCESSED_ATTR)) return;
  if (!isMagnetLink(element)) return;

  const magnetUrl = getMagnetUrl(element);
  if (!magnetUrl) return;

  element.setAttribute(PROCESSED_ATTR, 'true');
  
  const parent = element.parentElement;
  if (!parent) return;

  const btn = createPushButton(magnetUrl);
  
  if (parent.style) {
    parent.style.display = parent.style.display || '';
  }
  
  element.insertAdjacentElement('afterend', btn);
}

function scanPage() {
  const selectors = [
    'a[href^="magnet:"]',
    'a[href*="magnet:?xt="]',
    '[onclick*="magnet:"]',
    'td a',
    '.magnet',
    '[data-magnet]',
    'a'
  ];

  const magnetPattern = /magnet:\?xt=urn:[a-zA-Z0-9]+:[a-zA-Z0-9]{32,}/i;
  
  document.querySelectorAll(selectors.join(',')).forEach(el => {
    const text = el.textContent || '';
    const href = (el as HTMLAnchorElement).href || '';
    
    if (magnetPattern.test(href) || magnetPattern.test(text)) {
      processElement(el);
    }
  });

  const walker = document.createTreeWalker(
    document.body,
    NodeFilter.SHOW_TEXT,
    null
  );

  const textNodes: Text[] = [];
  while (walker.nextNode()) {
    const node = walker.currentNode as Text;
    if (magnetPattern.test(node.textContent || '')) {
      textNodes.push(node);
    }
  }

  textNodes.forEach(node => {
    const parent = node.parentElement;
    if (!parent || parent.hasAttribute(PROCESSED_ATTR) || parent.querySelector(`.${BUTTON_CLASS}`)) {
      return;
    }

    const match = (node.textContent || '').match(magnetPattern);
    if (!match) return;

    parent.setAttribute(PROCESSED_ATTR, 'true');
    const btn = createPushButton(match[0]);
    parent.appendChild(btn);
  });
}

function initObserver() {
  const observer = new MutationObserver((mutations) => {
    let shouldScan = false;
    for (const mutation of mutations) {
      if (mutation.addedNodes.length > 0) {
        shouldScan = true;
        break;
      }
    }
    if (shouldScan) {
      setTimeout(scanPage, 500);
    }
  });

  observer.observe(document.body, {
    childList: true,
    subtree: true
  });
}

function init() {
  injectStyles();
  scanPage();
  initObserver();

  chrome.runtime.onMessage.addListener((message) => {
    if (message.type === 'MAGNET_DETECTED') {
      showToast(`检测到磁力链接: ${message.magnet.substring(0, 50)}...`, 'info');
    } else if (message.type === 'PUSH_SUCCESS') {
      showToast('磁力链接已推送到115下载', 'success');
    } else if (message.type === 'PUSH_ERROR') {
      showToast(`推送失败: ${message.error}`, 'error');
    }
  });
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
