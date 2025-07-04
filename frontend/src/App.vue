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
            <div v-if="status.error" class="error-text">
              {{ $t('status.error', { msg: status.error }) }}
              <n-tag v-if="status.isUnknown" type="warning" size="small" style="margin-left: 8px;">
                {{ $t('status.unknown') }}
              </n-tag>
            </div>
            <n-space v-else vertical>
              <n-p>
                {{ $t('status.title') }}:
                <n-tag :type="status.isUnknown ? 'warning' : (status.enabled ? 'error' : 'success')" round>
                  {{ status.isUnknown ? $t('status.unknown') : 
                     (status.enabled ? $t('status.enabled') : $t('status.disabled')) }}
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
              <n-button type="warning" block @click="runReleaseRenewIP">{{ $t('basicRepair.releaseRenewIP') }}</n-button>
            </n-space>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card :title="$t('advancedReset.title')">
            <n-space vertical style="width: 100%;">
              <n-button type="warning" block @click="runResetTCPIP">{{ $t('advancedReset.resetTCPIP') }}</n-button>
              <n-button type="warning" block @click="runResetWinsock">{{ $t('advancedReset.resetWinsock') }}</n-button>
              <n-button type="warning" block @click="runRestartDNSService">{{ $t('advancedReset.restartDNS') }}</n-button>
              <n-button type="warning" block @click="runResetFirewall">{{ $t('advancedReset.resetFirewall') }}</n-button>
            </n-space>
          </n-card>
        </n-gi>
      </n-grid>

      <!-- 连通性测试 -->
      <n-card :title="$t('connectivity.title')">
        <template #header-extra>
          <n-button @click="getCurrentIP" quaternary circle size="small">
            <template #icon><n-icon :component="RefreshIcon" /></template>
          </n-button>
        </template>
        <n-space vertical style="width: 100%;">
          <n-p v-if="currentIP">
            {{ $t('connectivity.currentIP') }}: <n-text type="info">{{ currentIP }}</n-text>
          </n-p>
          <n-input v-model:value="pingHost" :placeholder="$t('connectivity.placeholder')" />
          <n-button type="primary" block @click="runPingTest" :loading="pingTesting">
            {{ $t('connectivity.test') }}
          </n-button>
          <n-text v-if="pingResult" :type="pingSuccess ? 'success' : 'error'">
            {{ pingResult }}
          </n-text>
        </n-space>
      </n-card>

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
  NH1, NP, NTag, NIcon, NLog, NText, NInput
} from 'naive-ui';
import { Refresh as RefreshIcon } from '@vicons/ionicons5';
import { GetProxyStatus, DisableProxyDirectly, DisableProxyViaPowerShell, ResetSystemProxy,
  FlushDNSCache, ResetTCPIP, ResetWinsock, RestartDNSService, GetCurrentLocale, PingTest, GetCurrentIP, ReleaseRenewIP, ResetFirewall } from '../wailsjs/go/main/App';

// useMessage() 的 Provider 在父组件 Root.vue 中
const message = useMessage();
const { t, locale } = useI18n();

const logs = ref([]);
const status = reactive({ enabled: false, server: '', error: '', isUnknown: false });
const logText = computed(() => logs.value.join('\n'));
const currentLocale = computed(() => locale.value);

// 连通性测试相关状态
const pingHost = ref('bing.com');
const pingResult = ref('');
const pingSuccess = ref(false);
const pingTesting = ref(false);
const currentIP = ref('');

// 初始化时获取后端语言设置
onMounted(async () => {
  try {
    const backendLocale = await GetCurrentLocale();
    if (backendLocale && (backendLocale === 'zh' || backendLocale === 'en')) {
      locale.value = backendLocale;
    } else {
      locale.value = 'zh'; // 默认中文
    }
  } catch (e) {
    console.warn('获取后端语言设置失败，使用默认中文:', e);
    locale.value = 'zh';
  }
  refreshStatus();
  getCurrentIP();
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
      addLog(t('status.error', { msg: result.error }));
      message.error(t('status.error', { msg: result.error }));
    } else {
      status.error = '';
      status.isUnknown = false;
      status.enabled = result.enabled;
      status.server = result.server;
      addLog(t('logs.updateSuccess'));
      message.success(t('logs.statusRefreshed'));
    }
  } catch (e) {
    const errorMsg = t('logs.backendError', { msg: e });
    status.error = errorMsg;
    status.isUnknown = true;
    addLog(errorMsg);
    message.error(errorMsg);
  }
};

const handleOperation = async (operationFunc, startMsg) => {
  addLog(startMsg);
  const loadingMessage = message.loading(t('logs.executing'), { duration: 0 });
  try {
    const result = await operationFunc();
    addLog(result.message);
    if (result.success) {
      message.success(result.message, { duration: 5000 });
    } else {
      message.error(result.message, { duration: 5000 });
    }
  } catch (e) {
    message.error(t('logs.backendError', { msg: e }), { duration: 5000 });
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
const runResetFirewall = () => handleOperation(ResetFirewall, t('logs.resettingFirewall'));
const runReleaseRenewIP = () => handleOperation(ReleaseRenewIP, t('logs.releaseRenewingIP'));

// 执行ping连通性测试
const runPingTest = async () => {
  if (!pingHost.value.trim()) {
    message.warning(t('connectivity.emptyHost'));
    return;
  }
  
  pingTesting.value = true;
  pingResult.value = '';
  addLog(t('logs.pingTesting', { host: pingHost.value }));
  
  try {
    const result = await PingTest(pingHost.value.trim());
    pingResult.value = result.message;
    pingSuccess.value = result.success;
    addLog(result.message);
    
    if (pingSuccess.value) {
      message.success(result.message, { duration: 3000 });
    } else {
      message.error(result.message, { duration: 5000 });
    }
  } catch (e) {
    const errorMsg = t('logs.backendError', { msg: e });
    pingResult.value = errorMsg;
    pingSuccess.value = false;
    addLog(errorMsg);
    message.error(errorMsg, { duration: 5000 });
  } finally {
    pingTesting.value = false;
  }
};

// 获取当前IP地址
const getCurrentIP = async () => {
  addLog(t('logs.gettingIP'));
  try {
    const result = await GetCurrentIP();
    addLog(result.message);
    if (result.success) {
      // 从消息中提取IP地址
      const ip = result.message.split(': ')[1];
      currentIP.value = ip;
      message.success(t('logs.ipUpdated'));
    } else {
      message.error(result.message);
    }
  } catch (e) {
    const errorMsg = t('logs.backendError', { msg: e });
    addLog(errorMsg);
    message.error(errorMsg);
  }
};

</script>



