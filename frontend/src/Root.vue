<template>
  <!-- 
      提供 Naive UI 的全局上下文。
      把 NConfigProvider 和 NMessageProvider 放在这里。
    -->
  <n-config-provider :theme="theme">
    <n-message-provider>
      <!-- 在 Provider 内部，加载 App 组件 -->
      <App />
      <!-- 语言切换按钮 -->
      <n-button @click="toggleLocale" size="small"
        style="position: fixed; top: 20px; right: 60px; z-index: 999; width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; background-color: transparent; border-radius: 12; padding: 0;">
        <svg width="20" height="20" aria-hidden="true" class="iconLanguage_nlXk" viewBox="0 0 24 24">
          <path fill="currentColor"
            d="m12.87 15.07-2.54-2.51.03-.03A17.5 17.5 0 0 0 14.07 6H17V4h-7V2H8v2H1v1.99h11.17C11.5 7.92 10.44 9.75 9 11.35 8.07 10.32 7.3 9.19 6.69 8h-2c.73 1.63 1.73 3.17 2.98 4.56l-5.09 5.02L4 19l5-5 3.11 3.11zM18.5 10h-2L12 22h2l1.12-3h4.75L21 22h2zm-2.62 7 1.62-4.33L19.12 17z" />
        </svg>
      </n-button>
      <!-- 主题切换按钮 -->
      <n-button @click="toggleTheme" size="small"
        style="position: fixed; top: 20px; right: 20px; z-index: 999; width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; background-color: transparent; border-radius: 12; padding: 0;">
        <n-icon size="20" style="display: flex; align-items: center; justify-content: center;">
          <Sunny v-if="themeName === 'light'" />
          <Moon v-else />
        </n-icon>
      </n-button>
    </n-message-provider>
    <n-global-style />
  </n-config-provider>
</template>

<script setup>
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { NConfigProvider, NMessageProvider, NGlobalStyle, lightTheme, darkTheme, NButton, NIcon } from 'naive-ui';
import { Sunny, Moon } from '@vicons/ionicons5';
import App from './App.vue';

// 获取i18n实例
const { locale: currentLocale } = useI18n();

// 切换语言函数
const toggleLocale = () => {
  currentLocale.value = currentLocale.value === 'zh' ? 'en' : 'zh';
}; // 引入主应用组件

// 定义主题名称，默认为浅色主题
const themeName = ref('light');

// 切换主题的函数
const toggleTheme = () => {
  themeName.value = themeName.value === 'light' ? 'dark' : 'light';
};

// 根据主题名称计算当前主题
const theme = computed(() => {
  return themeName.value === 'light' ? lightTheme : darkTheme;
});
</script>

<style scoped>
.theme-toggle-button {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 999;
  cursor: pointer;
  transition: background-color 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.theme-toggle-button:hover {
  background-color: rgba(255, 255, 255, 0.2);
}
</style>

<style>
@font-face {
  font-family: 'OPlusSans30';
  src: url('./assets/fonts/OPlusSans3.woff2') format('woff2');
  font-weight: normal;
  font-style: normal;
}

body,
n-button,
n-card,
n-log,
n-text,
n-p,
n-h1 {
  font-family: 'OPlusSans30', sans-serif !important;
}
</style>