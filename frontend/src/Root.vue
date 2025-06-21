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
      <n-button @click="toggleLocale" size="small" class="fixed-button language-toggle-button">
        <n-icon size="20">
          <LanguageOutline />
        </n-icon>
      </n-button>
      <!-- 主题切换按钮 -->
      <n-button @click="toggleTheme" size="small" class="fixed-button theme-toggle-button">
        <n-icon size="20">
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
import { SetLocale } from '../wailsjs/go/main/App';

// 获取i18n实例
const { locale: currentLocale } = useI18n();

// 切换语言函数
const toggleLocale = async () => {
  const newLocale = currentLocale.value === 'zh' ? 'en' : 'zh';
  // 通过API设置后端语言
  const result = await SetLocale(newLocale);
  // 设置前端语言
  currentLocale.value = result;
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



