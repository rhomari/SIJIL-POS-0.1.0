<script setup lang="ts">
import { ref } from 'vue'
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
  const hex = colorMap[color] || colorMap.primary
  document.documentElement.style.setProperty('--q-primary', hex || '')
}
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

    <q-drawer v-model="drawer" show-if-above :mini="miniState" @mouseenter="miniState = false" @mouseleave="miniState = true" mini-to-overlay :width="200" :breakpoint="500" bordered>
      <q-list>
        <q-item clickable v-ripple>
          <q-item-section avatar>
            <q-icon name="inbox" />
          </q-item-section>
          <q-item-section>Inbox</q-item-section>
        </q-item>
        <q-item clickable v-ripple>
          <q-item-section avatar>
            <q-icon name="star" />
          </q-item-section>
          <q-item-section>Star</q-item-section>
        </q-item>
        <q-item clickable v-ripple>
          <q-item-section avatar>
            <q-icon name="send" />
          </q-item-section>
          <q-item-section>Send</q-item-section>
        </q-item>
        <q-separator />
        <q-item clickable v-ripple>
          <q-item-section avatar>
            <q-icon name="drafts" />
          </q-item-section>
          <q-item-section>Drafts</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
  <router-view />
    </q-page-container>

    <q-footer elevated class="bg-grey-2 text-black">
      <div class="row items-center justify-between q-pa-sm">
        <div>
          <q-btn flat label="Theme" icon="palette" class="q-mr-sm">
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
        </div>
        <div>
          <q-btn flat label="Help" icon="help_outline" class="q-mr-sm" />
          <q-btn flat label="Settings" icon="settings" class="q-mr-sm" />
          <q-btn color="primary" label="Logout" icon="logout" />
        </div>
      </div>
    </q-footer>
  </q-layout>
</template>