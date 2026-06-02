import './assets/main.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";

const app = createApp(App)
const pinia = createPinia();





app.use(createPinia())
app.use(router)
app.use(ElementPlus)
pinia.use(piniaPluginPersistedstate);
app.use(pinia);
app.mount('#app')
