<template>
  <!-- 
    Ps：顶层的 NConfigProvider 和 NMessageProvider 已经移动到 Root.vue
  -->
  <div class="app-container">
    <n-space vertical :size="24">
      <!-- 标题区域 -->
      <div class="title-area">
        <n-h1 style="margin-bottom: 0;">ProxyCleaner</n-h1>
        <n-p depth="3">一个用于快速清理 Windows 系统代理残留的工具</n-p>
      </div>

      <!-- 状态与操作区域 -->
      <n-grid :x-gap="24" :y-gap="24" :cols="2">
        <n-gi>
          <n-card title="当前代理状态">
            <template #header-extra>
              <n-button @click="refreshStatus" quaternary circle size="small">
                <template #icon><n-icon :component="RefreshIcon" /></template>
              </n-button>
            </template>
            <div v-if="status.error" class="error-text">获取代理状态失败: {{ status.error }}</div>
            <n-space v-else vertical>
              <n-p>
                状态:
                <n-tag :type="status.enabled ? 'error' : 'success'" round>
                  {{ status.enabled ? '已启用' : '已禁用' }}
                </n-tag>
              </n-p>
              <n-p v-if="status.enabled">地址: {{ status.server || '未设置' }}</n-p>
            </n-space>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card title="执行操作">
            <n-space vertical style="width: 100%;">
              <n-button @click="runDisableProxyDirectly" type="primary" block>1. 直接修改注册表 (推荐)</n-button>
              <n-button @click="runDisableProxyViaPS" type="info" block>2. PowerShell修改</n-button>
            </n-space>
          </n-card>
        </n-gi>
      </n-grid>

      <!-- 日志区域 -->
      <n-card title="运行日志">
        <n-log :log="logText" :rows="10" />
      </n-card>
      
      <!-- 主题切换依然可以工作，但逻辑需要调整。为简单起见，暂时移除 -->
      <!-- 如果需要主题切换，要通过 props/emit 或状态管理库在 Root 和 App 之间通信 -->

    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue';
import { 
  useMessage, NSpace, NCard, NButton, NGrid, NGi, 
  NH1, NP, NTag, NIcon, NLog
} from 'naive-ui';
import { Refresh as RefreshIcon } from '@vicons/ionicons5';
import { GetProxyStatus, DisableProxyDirectly, DisableProxyViaPowerShell } from '../wailsjs/go/main/App';

// 现在 useMessage() 的 Provider 在父组件 Root.vue 中
const message = useMessage();

const logs = ref([]);
const status = reactive({ enabled: false, server: '', error: '' });
const logText = computed(() => logs.value.join('\n'));

const addLog = (msg) => {
  const timestamp = new Date().toLocaleTimeString();
  logs.value.push(`[${timestamp}] ${msg}`);
};

const refreshStatus = async () => {
  addLog('正在获取代理状态...');
  try {
    const result = await GetProxyStatus();
    Object.assign(status, result);
    if (result.error) {
      addLog(`获取状态失败: ${result.error}`);
      message.error(`获取状态失败: ${result.error}`);
    } else {
      addLog('代理状态已更新。');
      message.success('代理状态已刷新');
    }
  } catch (e) {
    const errorMsg = `调用后端失败: ${e}`;
    status.error = errorMsg;
    addLog(errorMsg);
    message.error(errorMsg);
  }
};

const handleOperation = async (operationFunc, startMsg) => {
  addLog(startMsg);
  const loadingMessage = message.loading('正在执行操作...', { duration: 0 });
  try {
    const result = await operationFunc();
    addLog(result);
    if (result.startsWith('失败')) {
      message.error(result, { duration: 5000 });
    } else {
      message.success(result, { duration: 5000 });
    }
  } catch (e) {
    message.error(`操作执行异常: ${e}`, { duration: 5000 });
  } finally {
    loadingMessage.destroy();
    await refreshStatus();
  }
};

const runDisableProxyDirectly = () => handleOperation(DisableProxyDirectly, '正在尝试直接修改注册表...');
const runDisableProxyViaPS = () => handleOperation(DisableProxyViaPowerShell, '正在尝试通过PowerShell修改...');

onMounted(() => {
  refreshStatus();
});
</script>

<style>
@import 'vfonts/Lato.css';
@import 'vfonts/FiraCode.css';

body {
  margin: 0;
  padding: 0;
  font-family: 'Lato', sans-serif;
}

.app-container {
  max-width: 700px;
  margin: 0 auto;
  padding: 40px 20px;
}

.title-area {
  text-align: center;
}
</style>