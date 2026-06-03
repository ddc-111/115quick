const d=/magnet:\?xt=urn:[a-zA-Z0-9]+:[a-zA-Z0-9]{32,}/i,s="quick-115-push-btn",u="data-115-processed",x=`
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
`,h=`
.${s} {
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
.${s}:hover {
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  box-shadow: 0 4px 8px rgba(99, 102, 241, 0.4);
  transform: translateY(-1px);
}
.${s}:active {
  transform: translateY(0);
}
.${s}.loading {
  opacity: 0.7;
  cursor: not-allowed;
}
.${s} .icon {
  font-size: 14px;
}
`;let c=null;function y(){if(document.getElementById("quick-115-styles"))return;const e=document.createElement("style");e.id="quick-115-styles",e.textContent=x+h,document.head.appendChild(e)}function i(e,t="info",n=3e3){c||(c=document.createElement("div"),c.style.cssText="position:fixed;top:0;right:0;z-index:2147483647;pointer-events:none;",document.body.appendChild(c));const o=document.createElement("div");o.className=`quick-115-toast ${t}`,o.textContent=e,o.style.pointerEvents="auto",c.appendChild(o),setTimeout(()=>{o.style.animation="quick-115-slide-out 0.3s ease-in forwards",setTimeout(()=>o.remove(),300)},n)}function b(){return new Promise(e=>{chrome.storage.local.get(["serverUrl"],t=>{e(t.serverUrl||"")})})}async function k(e,t){const n=await b();if(!n){i("请先在115Quick插件中设置服务器地址","error");return}t.classList.add("loading"),t.textContent="推送中...";try{const o=await fetch(`${n}/v1/Download/api/addDownloadLink`,{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({downloadLink:e})}),r=await o.json();if(o.ok)i("磁力链接已推送到115下载","success"),t.textContent="已推送",t.style.background="linear-gradient(135deg, #10b981, #059669)",setTimeout(()=>{t.textContent="推送下载",t.style.background="",t.classList.remove("loading")},2e3);else throw new Error(r.message||"推送失败")}catch(o){i(`推送失败: ${o.message}`,"error"),t.textContent="重试",t.classList.remove("loading")}}function p(e){const t=document.createElement("button");return t.className=s,t.innerHTML='<span class="icon">⬇</span> 推送下载',t.title="推送到115网盘离线下载",t.addEventListener("click",n=>{n.preventDefault(),n.stopPropagation(),k(e,t)}),t}function E(e){const t=e.textContent||"",n=e.href||"";return d.test(t)||d.test(n)}function C(e){const t=e.href||"";if(d.test(t)){const r=t.match(d);return r?r[0]:""}const o=(e.textContent||"").match(d);return o?o[0]:""}function S(e){if(e.hasAttribute(u)||!E(e))return;const t=C(e);if(!t)return;e.setAttribute(u,"true");const n=e.parentElement;if(!n)return;const o=p(t);n.style&&(n.style.display=n.style.display||""),e.insertAdjacentElement("afterend",o)}function m(){const e=['a[href^="magnet:"]','a[href*="magnet:?xt="]','[onclick*="magnet:"]',"td a",".magnet","[data-magnet]","a"],t=/magnet:\?xt=urn:[a-zA-Z0-9]+:[a-zA-Z0-9]{32,}/i;document.querySelectorAll(e.join(",")).forEach(r=>{const a=r.textContent||"",l=r.href||"";(t.test(l)||t.test(a))&&S(r)});const n=document.createTreeWalker(document.body,NodeFilter.SHOW_TEXT,null),o=[];for(;n.nextNode();){const r=n.currentNode;t.test(r.textContent||"")&&o.push(r)}o.forEach(r=>{const a=r.parentElement;if(!a||a.hasAttribute(u)||a.querySelector(`.${s}`))return;const l=(r.textContent||"").match(t);if(!l)return;a.setAttribute(u,"true");const g=p(l[0]);a.appendChild(g)})}function T(){new MutationObserver(t=>{let n=!1;for(const o of t)if(o.addedNodes.length>0){n=!0;break}n&&setTimeout(m,500)}).observe(document.body,{childList:!0,subtree:!0})}function f(){y(),m(),T(),chrome.runtime.onMessage.addListener(e=>{e.type==="MAGNET_DETECTED"?i(`检测到磁力链接: ${e.magnet.substring(0,50)}...`,"info"):e.type==="PUSH_SUCCESS"?i("磁力链接已推送到115下载","success"):e.type==="PUSH_ERROR"&&i(`推送失败: ${e.error}`,"error")})}document.readyState==="loading"?document.addEventListener("DOMContentLoaded",f):f();
