import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { config } from 'md-editor-v3'
import hljs from 'highlight.js/lib/common'
import App from './App.vue'
import router from './router'
import './style.css'

// Use the bundled highlight.js in the editor preview instead of loading it from a CDN
config({ editorExtensions: { highlight: { instance: hljs } } })

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
