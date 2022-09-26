import {createApp} from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import RootProps from "./props";

const app = createApp(App, RootProps);

app.use(ElementPlus)
createApp(App).mount('#app')
