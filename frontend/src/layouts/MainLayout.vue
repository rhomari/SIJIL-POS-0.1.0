<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useQuasar } from 'quasar'
import { useI18n } from 'vue-i18n';

const { locale } = useI18n();
const $q = useQuasar();
const langs = [
  { value: 'en', label: 'English' },
  { value: 'fr', label: 'Français' },
  { value: 'ar', label: 'العربية' }
];
const currentLang = ref(locale.value);
function setLang(lang: string) {
  locale.value = lang;
  currentLang.value = lang;
  localStorage.setItem('app-lang', lang);
}
onMounted(() => {
  const savedLang = localStorage.getItem('app-lang');
  if (savedLang && langs.some(l => l.value === savedLang)) {
    locale.value = savedLang;
    currentLang.value = savedLang;
  }
});

const drawer = ref(false)
const miniState = ref(true)
const themeColors = [
  { name: 'Classic Blue', color: 'primary' },
  { name: 'Teal Green', color: 'secondary' },
  { name: 'Violet Accent', color: 'accent' },
  { name: 'Dark Mode', color: 'dark' },
  { name: 'Success Green', color: 'positive' },
  { name: 'Error Red', color: 'negative' },
  { name: 'Info Cyan', color: 'info' },
  { name: 'Warning Yellow', color: 'warning' }
]
const colorMap: Record<string, string> = {
  primary: '#1976d2',
  secondary: '#26a69a',
  accent: '#9c27b0',
  dark: '#1d1d1d',
  positive: '#21ba45',
  negative: '#c10015',
  info: '#31ccec',
  warning: '#f2c037',
}
function setThemeColor(color: string) {
  const hex = colorMap[color] || colorMap.primary;
  document.documentElement.style.setProperty('--q-primary', hex || '');
  localStorage.setItem('app-theme', color);
}

onMounted(() => {
  const savedTheme = localStorage.getItem('app-theme');
  if (savedTheme && colorMap[savedTheme]) {
    setThemeColor(savedTheme);
  }
});


onMounted(() => {
  const savedDrawer = localStorage.getItem('app-drawer');
  if (savedDrawer !== null) {
    drawer.value = savedDrawer === '1';
  } else {
    
    try {
      drawer.value = window.innerWidth >= 500;
    } catch {
      drawer.value = true; 
    }
  }
});


watch(drawer, (val) => {
  try { localStorage.setItem('app-drawer', val ? '1' : '0'); } catch (e) { void e; }
});


const isOffline = ref(!navigator.onLine);
function handleOnline() { isOffline.value = false; }
function handleOffline() { isOffline.value = true; }
onMounted(() => {
  window.addEventListener('online', handleOnline);
  window.addEventListener('offline', handleOffline);
});
onUnmounted(() => {
  window.removeEventListener('online', handleOnline);
  window.removeEventListener('offline', handleOffline);
});


const showLoginDialog = ref(false);
const loginUsername = ref('');
const loginPassword = ref('');
function doLogin() {
  
  try { localStorage.setItem('auth-token', 'demo'); } catch (e) { void e; }
  showLoginDialog.value = false;
}
onMounted(() => {
  const handler = () => { showLoginDialog.value = true; };
  window.addEventListener('app:show-login', handler as EventListener);
  onUnmounted(() => { window.removeEventListener('app:show-login', handler as EventListener); });
});

function footerLogout() {
  try { localStorage.removeItem('auth-token'); } catch (e) { void e; }
  showLoginDialog.value = true;
}

</script>
<template>
  <q-layout view="lHh lpr lFf" :class="{ 'blur-content': isOffline, 'blur-content-logout': showLoginDialog }">
    <q-header elevated>
      <q-bar>
        <q-btn dense flat icon="menu" @click="drawer = !drawer" />
        <q-icon name="laptop_chromebook" />
        <div>SIJIL-POS 0.1.0</div>
        <q-space />
        <q-btn dense flat icon="minimize" />
        <q-btn dense flat icon="crop_square" />
        <q-btn dense flat icon="close" />
      </q-bar>
    </q-header>

    <q-drawer v-model="drawer" show-if-above :mini="miniState" @mouseenter="miniState = false" @mouseleave="miniState = true" mini-to-overlay :width="300" :breakpoint="500" bordered>
      <q-list>
  <q-item clickable v-ripple :to="{ name: 'pos' }" exact>
          <q-item-section avatar>
            <q-icon name="point_of_sale" />
          </q-item-section>
    <q-item-section>{{ $t('POS') }}</q-item-section>
        </q-item>
        <q-item clickable v-ripple :to="{ name: 'products' }">
          <q-item-section avatar>
            <q-icon name="category" />
          </q-item-section>
    <q-item-section>{{ $t('ARTICLES') }}</q-item-section>
        </q-item>
        <q-item clickable v-ripple :to="{ name: 'customers' }">
          <q-item-section avatar>
            <q-icon name="people" />
          </q-item-section>
          <q-item-section>{{ $t('customers') }}</q-item-section>
        </q-item>
        <q-item clickable v-ripple>
          <q-item-section avatar>
            <q-icon name="widgets" />
          </q-item-section>
    <q-item-section>{{ $t('STOCK') }}</q-item-section>
        </q-item>
        <q-separator />
        <q-item clickable v-ripple>
          <q-item-section avatar>
            <q-icon name="settings" />
          </q-item-section>
    <q-item-section>{{ $t('PARAMS') }}</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
  <router-view />
    </q-page-container>

    <q-footer elevated class="bg-grey-2 text-black">
      <div class="row items-center justify-between q-pa-sm">
        <div class="row items-center">
          <q-btn flat :label="$q.screen.lt.md ? undefined : $t('theme')" icon="palette" class="q-mr-sm">
            <q-menu>
              <q-list style="min-width: 150px;">
                <q-item v-for="theme in themeColors" :key="theme.color" clickable @click="setThemeColor(theme.color)">
                  <q-item-section avatar>
                    <q-icon :color="theme.color === 'primary' ? '#1976d2' : theme.color" name="lens" />
                  </q-item-section>
                  <q-item-section>{{ theme.name }}</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
          <q-btn flat dense icon="language" class="q-ml-xs">
            <q-menu anchor="top right" self="bottom right">
              <q-list style="min-width: 100px;">
                <q-item v-for="lang in langs" :key="lang.value" clickable @click="setLang(lang.value)">
                  <q-item-section>
                    <span :style="{ fontWeight: currentLang === lang.value ? 'bold' : 'normal' }">{{ lang.label }}</span>
                  </q-item-section>
                  <q-item-section side v-if="currentLang === lang.value">
                    <q-icon name="check" color="primary" size="xs" />
                  </q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>
        <div class="row items-center">
          <q-btn flat :label="$q.screen.lt.md ? undefined : $t('help')" icon="help_outline" class="q-mr-sm" />
          <q-btn flat :label="$q.screen.lt.md ? undefined : $t('logout')" icon="logout" class="q-ml-xs" @click="footerLogout" />
        </div>
      </div>
    </q-footer>

    <!-- Offline persistent dialog -->
    <q-dialog v-model="isOffline" persistent class="offline-modal" transition-show="fade" transition-hide="fade">
      <q-card style="width: 420px; max-width: 90vw;" class="q-pa-md offline-card">
        <q-card-section class="row items-center q-gutter-md">
          <q-icon name="signal_wifi_off" color="negative" size="48px" />
          <div>
            <div class="text-h6">You are offline</div>
            <div class="text-body2 text-grey-7 q-mt-xs">Waiting for network connection…</div>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn color="primary" :disable="true" label="Retry" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Login dialog (shown on logout) -->
    <q-dialog v-model="showLoginDialog" persistent transition-show="scale" transition-hide="scale">
      <q-card style="width: 420px; max-width: 92vw;" class="q-pa-md">
        <q-card-section class="row items-center q-gutter-md">
          <q-icon name="lock" color="primary" size="40px" />
          <div>
            <div class="text-h6">Login</div>
            <div class="text-caption text-grey-7">Please sign in to continue</div>
          </div>
        </q-card-section>
        <q-card-section class="q-pt-none">
          <q-input v-model="loginUsername" dense outlined label="Username" autofocus @keyup.enter="doLogin" />
          <q-input v-model="loginPassword" type="password" dense outlined label="Password" class="q-mt-sm" @keyup.enter="doLogin" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Cancel" @click="showLoginDialog = false" />
          <q-btn color="primary" label="Login" @click="doLogin" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-layout>
</template>

<style scoped>

.blur-content {
  filter: blur(3px) saturate(0.9);
  transition: filter .15s ease;
  pointer-events: none; 
  user-select: none;
}

.blur-content-logout {
  filter: blur(3px) grayscale(1) saturate(0.9);
  transition: filter .15s ease;
  pointer-events: none;
  user-select: none;
}

.offline-modal :deep(.q-dialog__backdrop) {
  backdrop-filter: blur(4px);
  background: rgba(0, 0, 0, 0.18);
}
.offline-card {
  border-radius: 14px;
}
</style>