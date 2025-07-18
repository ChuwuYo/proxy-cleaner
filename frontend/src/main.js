import { createApp } from 'vue';
import Root from './Root.vue';
import naive from 'naive-ui';
import i18n from './i18n';
import './assets/styles/index.css';

const app = createApp(Root);

app.use(naive);
app.use(i18n);

app.mount('#app');