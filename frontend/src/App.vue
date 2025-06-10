<template>
  <!-- 
    Ps：顶层的 NConfigProvider 和 NMessageProvider 已经移动到 Root.vue
  -->
  <div class="app-container">
    <n-space vertical :size="24">
      <!-- 标题区域 -->
      <div class="title-area">
        <n-space vertical :size="12">
          <n-space justify="center" align="center">
            <n-h1 style="margin-bottom: 0;">{{ $t('title') }}</n-h1>
            <n-space>
              <n-button quaternary circle size="small" @click="toggleLocale">
                {{ currentLocale === 'zh' ? 'EN' : '中' }}
              </n-button>
            </n-space>
          </n-space>
          <n-p depth="3">{{ $t('subtitle') }}</n-p>
        </n-space>
      </div>

      <!-- 状态与操作区域 -->
      <n-grid :x-gap="24" :y-gap="24" :cols="2">
        <n-gi>
          <n-card :title="$t('status.title')">
            <template #header-extra>
              <n-button @click="refreshStatus" quaternary circle size="small">
                <template #icon><n-icon :component="RefreshIcon" /></template>
              </n-button>
            </template>
            <div v-if="status.error" class="error-text">{{ $t('status.error', { msg: status.error }) }}</div>
            <n-space v-else vertical>
              <n-p>
                {{ $t('status.title') }}:
                <n-tag :type="status.enabled ? 'error' : 'success'" round>
                  {{ status.enabled ? $t('status.enabled') : $t('status.disabled') }}
                </n-tag>
              </n-p>
              <n-p>
                {{ $t('status.address') }}:
                <n-text :depth="status.enabled ? 1 : 3">
                  {{ status.server || $t('status.notSet') }}
                  <n-tag v-if="!status.enabled && status.server" size="small" type="default">{{ $t('status.notEffective') }}</n-tag>
                </n-text>
              </n-p>
            </n-space>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card :title="$t('operations.title')">
            <n-space vertical style="width: 100%;">
              <n-button @click="runDisableProxyDirectly" type="primary" block>{{ $t('operations.direct') }}</n-button>
              <n-button @click="runDisableProxyViaPS" type="info" block>{{ $t('operations.powershell') }}</n-button>
            </n-space>
          </n-card>
        </n-gi>
      </n-grid>

      <!-- 系统网络修复功能 -->
      <n-grid :x-gap="24" :y-gap="24" :cols="2">
        <n-gi>
          <n-card :title="$t('basicRepair.title')">
            <n-space vertical style="width: 100%;">
              <n-button type="primary" block @click="runResetSystemProxy">{{ $t('basicRepair.resetProxy') }}</n-button>
              <n-button type="info" block @click="runFlushDNSCache">{{ $t('basicRepair.flushDNS') }}</n-button>
            </n-space>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card :title="$t('advancedReset.title')">
            <n-space vertical style="width: 100%;">
              <n-button type="warning" block @click="runResetTCPIP">{{ $t('advancedReset.resetTCPIP') }}</n-button>
              <n-button type="warning" block @click="runResetWinsock">{{ $t('advancedReset.resetWinsock') }}</n-button>
              <n-button type="warning" block @click="runRestartDNSService">{{ $t('advancedReset.restartDNS') }}</n-button>
            </n-space>
          </n-card>
        </n-gi>
      </n-grid>

      <!-- 日志区域 -->
      <n-card :title="$t('logs.title')">
        <n-log :log="logText" :rows="10" />
      </n-card>
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  useMessage, NSpace, NCard, NButton, NGrid, NGi,
  NH1, NP, NTag, NIcon, NLog, NText
} from 'naive-ui';
import { Refresh as RefreshIcon } from '@vicons/ionicons5';
import { GetProxyStatus, DisableProxyDirectly, DisableProxyViaPowerShell, ResetSystemProxy,
  FlushDNSCache, ResetTCPIP, ResetWinsock, RestartDNSService, GetCurrentLocale, SetLocale } from '../wailsjs/go/main/App';

// useMessage() 的 Provider 在父组件 Root.vue 中
const message = useMessage();
const { t, locale } = useI18n();

const logs = ref([]);
const status = reactive({ enabled: false, server: '', error: '' });
const logText = computed(() => logs.value.join('\n'));
const currentLocale = computed(() => locale.value);

const toggleLocale = async () => {
  const newLocale = locale.value === 'zh' ? 'en' : 'zh';
  // 通过API设置后端语言
  const result = await SetLocale(newLocale);
  // 设置前端语言
  locale.value = result;
};

// 初始化时获取后端语言设置
onMounted(async () => {
  const backendLocale = await GetCurrentLocale();
  locale.value = backendLocale;
  refreshStatus();
});

const addLog = (msg) => {
  const timestamp = new Date().toLocaleTimeString();
  logs.value.push(`[${timestamp}] ${msg}`);
};

const refreshStatus = async () => {
  addLog(t('logs.refreshing'));
  try {
    const result = await GetProxyStatus();
    if (result.error) {
      status.error = result.error;
      addLog(`获取状态失败: ${result.error}`);
      message.error(`获取状态失败: ${result.error}`);
    } else {
      status.error = '';  // 清除错误状态
      status.enabled = result.enabled;
      status.server = result.server;
      addLog(t('logs.updateSuccess'));
      message.success(t('logs.statusRefreshed'));
    }
  } catch (e) {
    const errorMsg = t('logs.backendError', { msg: e });
    status.error = errorMsg;
    status.enabled = false;  // 错误时设置为禁用状态
    status.server = '';     // 错误时清空服务器地址
    addLog(errorMsg);
    message.error(errorMsg);
  }
};

const handleOperation = async (operationFunc, startMsg) => {
  addLog(startMsg);
  const loadingMessage = message.loading(t('logs.executing'), { duration: 0 });
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

const runDisableProxyDirectly = () => handleOperation(DisableProxyDirectly, t('logs.directModifying'));
const runDisableProxyViaPS = () => handleOperation(DisableProxyViaPowerShell, t('logs.psModifying'));

const runResetSystemProxy = () => handleOperation(ResetSystemProxy, t('logs.resettingProxy'));
const runFlushDNSCache = () => handleOperation(FlushDNSCache, t('logs.flushingDNS'));
const runResetTCPIP = () => handleOperation(ResetTCPIP, t('logs.resettingTCPIP'));
const runResetWinsock = () => handleOperation(ResetWinsock, t('logs.resettingWinsock'));
const runRestartDNSService = () => handleOperation(RestartDNSService, t('logs.restartingDNS'));

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

.header-btn {
  transition: background-color 0.3s ease;
}

.header-btn:hover {
  background-color: rgba(24, 160, 88, 0.12);
}

.header-btn:active {
  background-color: transparent;
  transition: background-color 0.1s ease;
}
</style>