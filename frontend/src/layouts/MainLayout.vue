<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n';

const { locale } = useI18n();
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
</script>
<template>
  <q-layout view="lHh lpr lFf">
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
          <q-btn flat :label="$t('theme')" icon="palette" class="q-mr-sm">
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
        <div>
          <q-btn flat :label="$t('help')" icon="help_outline" class="q-mr-sm" />
          <q-btn flat :label="$t('settings')" icon="settings" class="q-mr-sm" />
          <q-btn color="primary" :label="$t('logout')" icon="logout" />
        </div>
      </div>
    </q-footer>
  </q-layout>
</template>