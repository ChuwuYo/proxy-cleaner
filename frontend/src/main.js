import { createApp } from 'vue';
// 入口点为 Root
import Root from './Root.vue'; 
import naive from 'naive-ui';

// 创建 Vue 应用实例，以 Root 为根
const app = createApp(Root);

// naive 的 .use() 仍然需要，它注册了所有组件的指令等
app.use(naive);

// 挂载应用
app.mount('#app');