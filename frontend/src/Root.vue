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
        <n-icon size="20" style="display: flex; align-items: center; justify-content: center;">
          <LanguageOutline />
        </n-icon>
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
import { Sunny, Moon, LanguageOutline } from '@vicons/ionicons5';
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