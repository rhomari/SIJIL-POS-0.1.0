<template>
  <div class="scanner-wrapper q-pa-md">
    <div class="video-container" :class="{ 'inactive': !active }">
      <video ref="videoEl" autoplay playsinline muted class="scanner-video"></video>
      <canvas ref="canvasEl" class="hidden-canvas"></canvas>
      <div v-if="!active && !error" class="overlay center">{{ t('scanning') }}...</div>
      <div v-if="error" class="overlay error">{{ error }}</div>
      <div v-if="lastCode" class="last-code badge">{{ lastCode }}</div>
    </div>
    <div class="row q-gutter-sm q-mt-sm">
      <q-btn dense color="primary" :label="t('close')" @click="$emit('close')" />
      <q-btn dense flat color="secondary" :label="t('reset')" @click="resetScan" />
      <q-space />
      <q-btn dense outline color="primary" :label="paused ? t('resume') : t('pause')" @click="togglePause" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { BrowserMultiFormatReader } from '@zxing/browser';
import { useI18n } from 'vue-i18n';

const emit = defineEmits<{ (e: 'decoded', code: string): void; (e: 'close'): void }>();

const { t } = useI18n();

const videoEl = ref<HTMLVideoElement | null>(null);
const canvasEl = ref<HTMLCanvasElement | null>(null);
const reader = new BrowserMultiFormatReader();
let stream: MediaStream | null = null;
const active = ref(false);
const paused = ref(false);
const error = ref<string | null>(null);
const lastCode = ref<string | null>(null);
let lastEmitTime = 0;
const DUP_INTERVAL = 1500; // ms

async function init() {
  try {
    error.value = null;
    const devices = await navigator.mediaDevices.enumerateDevices();
    const videoDevices = devices.filter(d => d.kind === 'videoinput');
    const preferred = videoDevices.find(d => /back|rear|environment/i.test(d.label)) || videoDevices[0];
    stream = await navigator.mediaDevices.getUserMedia({
      video: preferred ? { deviceId: { exact: preferred.deviceId } } : { facingMode: 'environment' },
      audio: false
    });
    if (videoEl.value) {
      videoEl.value.srcObject = stream;
      await videoEl.value.play();
    }
  active.value = true;
  void tick();
  } catch {
    error.value = t('cameraPermissionDenied');
  }
}

function resetScan() {
  lastCode.value = null;
  lastEmitTime = 0;
}

function togglePause() {
  paused.value = !paused.value;
  if (!paused.value) void tick();
}

async function tick() {
  if (paused.value) return;
  try {
    const result = await reader.decodeOnceFromVideoElement(videoEl.value!);
    const text = result.getText();
    const now = Date.now();
    if (text && now - lastEmitTime > DUP_INTERVAL) {
      lastCode.value = text;
      lastEmitTime = now;
      emit('decoded', text);
    }
  } catch (err) {
    void err; // ignore
  }
  if (active.value && !paused.value) {
    requestAnimationFrame(() => { void tick(); });
  }
}

onMounted(init);

onBeforeUnmount(() => {
  try {
    // Some builds expose stopContinuousDecode; guard before calling
    if (typeof (reader as unknown as { stopContinuousDecode?: () => void }).stopContinuousDecode === 'function') {
      (reader as unknown as { stopContinuousDecode: () => void }).stopContinuousDecode();
    }
  } catch {
    // ignore cleanup errors
  }
  if (stream) {
    for (const track of stream.getTracks()) track.stop();
    stream = null;
  }
});
</script>

<style scoped>
.scanner-wrapper { width: 100%; max-width: 480px; }
.video-container { position: relative; width: 100%; aspect-ratio: 3/4; background: #000; border-radius: 12px; overflow: hidden; }
.scanner-video { width: 100%; height: 100%; object-fit: cover; }
.overlay { position: absolute; inset: 0; display:flex; align-items:center; justify-content:center; color:#fff; font-weight:600; background:rgba(0,0,0,0.25); backdrop-filter: blur(2px); }
.overlay.error { background: rgba(200,0,0,0.45); }
.last-code { position:absolute; bottom:8px; left:8px; background:rgba(0,0,0,0.6); color:#fff; padding:4px 10px; border-radius:8px; font-size:12px; }
.hidden-canvas { display:none; }
.badge { font-family: monospace; letter-spacing: 1px; }
</style>
