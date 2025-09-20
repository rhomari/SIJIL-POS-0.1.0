<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';

// Disable context menu globally
function disableContextMenu(event: Event) {
  event.preventDefault();
  return false;
}

// Disable F12, Ctrl+Shift+I, Ctrl+U, and other developer shortcuts
function disableDevTools(event: KeyboardEvent) {
  // F12 or Ctrl+Shift+I or Ctrl+Shift+J or Ctrl+U
  if (
    event.key === 'F12' ||
    (event.ctrlKey && event.shiftKey && (event.key === 'I' || event.key === 'J')) ||
    (event.ctrlKey && event.key === 'U')
  ) {
    event.preventDefault();
    return false;
  }
}

onMounted(() => {
  // Disable right-click context menu
  document.addEventListener('contextmenu', disableContextMenu);
  
  // Disable developer shortcuts
  document.addEventListener('keydown', disableDevTools);
  
  // Disable drag and drop of external files
  document.addEventListener('dragover', (e) => e.preventDefault());
  document.addEventListener('drop', (e) => e.preventDefault());
});

onUnmounted(() => {
  document.removeEventListener('contextmenu', disableContextMenu);
  document.removeEventListener('keydown', disableDevTools);
});
</script>
