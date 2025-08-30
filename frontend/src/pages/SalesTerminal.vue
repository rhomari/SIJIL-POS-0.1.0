<template>
  <q-page class="q-pa-md">
    <div class="row q-col-gutter-xl">
      <!-- Articles Section -->
      <div class="col-8">
        <q-tabs v-model="selectedCategory" class="text-primary" align="left" dense>
          <q-tab v-for="cat in categories" :key="cat" :name="cat" :label="cat" />
        </q-tabs>
        <div class="row q-gutter-x-xs q-mt-sm">
          <div
            v-for="article in filteredArticles"
            :key="article.id"
            class="col-6 col-sm-4 col-md-3 col-lg-2 q-mb-xs"
          >
            <q-card
              bordered
              class="shadow-4 cursor-pointer"
              style="width:100%; height: 120px;"
              @click="addToReceipt(article)"
            >
              <!--LOCK START-->
              <q-img
                :src="article.image"
                :alt="article.name"
                style="width: 100%; height:119px;"
                :img-style="{ objectFit: 'contain', width: '100%', height: '100%' }"
                fit="contain"
                img-class="q-pa-xs bg-white"
              >
                <div
                  class="absolute-bottom text-subtitle1 text-center"
                  style="height: 30%; word-wrap: break-word; padding: 0px; font-size: 0.8rem;"
                >
                  {{ article.name }} - {{ article.price.toFixed(2) }} MAD
                </div>
              </q-img>
              <!--LOCK END-->
            </q-card>
          </div>
        </div>
      </div>
      <!-- Receipt Section -->
      <div class="col-4">
        <q-card>
          <q-card-section class="bg-grey-2">
              <div class="text-h6">{{ $t('receipt') }}</div>
          </q-card-section>
          <q-separator />
          <div class="relative-position">
            <q-list>
              <q-item v-for="item in groupedReceipt" :key="item.id"
                draggable="true"
                @dragstart="onDragStart(item.id)"
                @dragend="onDragEnd"
                :class="{ 'dragging': draggingId === item.id }"
              >
                <q-item-section avatar>
                  <q-fab
                    color="primary"
                    text-color="white"
                    :label="item.qty.toString()"
                    padding="xs"
                    dense
                    direction="left"
                    icon="more_vert"
                    class="receipt-fab"
                    :model-value="openFabId === item.id"
                    @update:model-value="val => toggleFab(item.id, val)"
                  >
                    <q-fab-action color="primary" icon="add" @click.stop="increment(item.id)" />
                    <q-fab-action color="primary" icon="remove" :disable="item.qty <= 1" @click.stop="decrement(item.id)" />
                    <q-fab-action color="secondary" icon="edit" @click.stop="openEdit(item)" />
                    <q-fab-action color="negative" icon="delete" @click.stop="deleteLine(item.id)" />
                  </q-fab>
                </q-item-section>
                <q-item-section>
                  <div class="row items-center justify-between w-100">
                    <span>{{ item.name }}</span>
                    <span class="text-weight-bold">{{ (item.price * item.qty).toFixed(2) }} MAD</span>
                  </div>
                </q-item-section>
              </q-item>
            </q-list>
            <!-- Floating trash can (old style: red glow while dragging) -->
            <div
              v-show="draggingId !== null"
              class="trash-drop-zone fixed-center"
              :class="{ 'trash-active': draggingId !== null }"
              @dragover.prevent
              @drop="onDropTrash"
              :style="draggingId !== null ? 'z-index:9999;pointer-events:auto;' : 'pointer-events:none;'"
            >
              <div
                class="bg-white q-pa-lg shadow-10 rounded-borders flex flex-center transition-all"
                :style="draggingId !== null ? 'width:140px;height:140px;box-shadow:0 0 24px 8px #f44336;' : 'width:100px;height:100px;'"
              >
                <q-icon
                  name="delete"
                  :color="draggingId !== null ? 'negative' : 'grey-6'"
                  size="64px"
                  class="transition-all"
                  :class="draggingId !== null ? 'q-animate-bounce' : ''"
                  :style="draggingId !== null ? 'filter: drop-shadow(0 0 8px #f44336);' : ''"
                />
              </div>
            </div>
          </div>
          <q-dialog v-model="editDialog">
            <q-card style="width:320px; max-width:90vw;" class="q-pa-sm">
              <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Edit Quantity</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
              </q-card-section>
              <q-card-section class="q-pt-sm">
                <div class="text-subtitle2 text-primary">{{ editLineName }}</div>
                <div class="q-mt-sm row items-center no-wrap">
                  <q-btn dense round icon="remove" color="primary" :disable="tempQty <= 0" @click="tempQty = +(Math.max(0, tempQty - qtyStep).toFixed(2))" />
                  <q-input v-model.number="tempQty" type="number" dense outlined class="q-mx-sm col" :step="qtyStep" min="0" />
                  <q-btn dense round icon="add" color="primary" @click="tempQty = +(tempQty + qtyStep).toFixed(2)" />
                </div>
                  <div class="text-caption text-grey-7 q-mt-xs">{{ $t('editHint') }}</div>
              </q-card-section>
              <q-separator inset />
              <q-card-actions align="right">
                  <q-btn flat :label="$t('cancel')" v-close-popup />
                  <q-btn color="primary" :label="$t('apply')" @click="applyEdit" />
              </q-card-actions>
            </q-card>
          </q-dialog>
          <q-separator />
          <q-card-section class="text-right">
              <div class="text-subtitle1">{{ $t('total') }}: {{ total.toFixed(2) }} MAD</div>
            <q-btn color="positive" :label="$t('checkout')" class="q-mt-sm" @click="cycleLang" />
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
const draggingId = ref<number|null>(null);

function onDragStart(id: number) {
  draggingId.value = id;
  // Optionally set drag image
  // event.dataTransfer?.setDragImage(event.target as HTMLElement, 0, 0);
}
function onDragEnd() {
  draggingId.value = null;
}
function onDropTrash() {
  if (draggingId.value !== null) {
    deleteLine(draggingId.value);
    draggingId.value = null;
  }
}
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue';

const categories = ['Stationery', 'Books', 'Office'];
const selectedCategory = ref(categories[0]);
type Article = {
  id: number;
  name: string;
  price: number;
  image: string;
  category: string;
};

const articles = ref<Article[]>([
  { id: 1, name: 'Blue Ballpoint Pen', price: 1.2, image: '/images/1.png', category: 'Stationery' },
  { id: 2, name: 'A5 Spiral Notebook', price: 3.5, image: '/images/2.png', category: 'Stationery' },
  { id: 3, name: 'Highlighter Set', price: 4.0, image: '/images/3.png', category: 'Stationery' },
  { id: 4, name: 'Stapler', price: 5.5, image: '/images/4.png', category: 'Office' },
  { id: 5, name: 'Desk Organizer', price: 9.9, image: '/images/5.png', category: 'Office' },
  { id: 6, name: 'Mystery Novel', price: 7.99, image: '/images/6.png', category: 'Books' },
  { id: 7, name: 'Science Textbook', price: 24.5, image: '/images/7.png', category: 'Books' },
]);

const filteredArticles = computed(() =>
  articles.value.filter((a: Article) => a.category === selectedCategory.value)
);

const receipt = ref<Article[]>([]);
function addToReceipt(article: Article) {
  receipt.value.push(article);
}
import { useI18n } from 'vue-i18n';
const { locale } = useI18n();
const langOrder = ['en', 'fr', 'ar'];
function cycleLang() {
  const idx = langOrder.indexOf(locale.value);
  const nextLang = langOrder[(idx + 1) % langOrder.length] ?? 'en';
  locale.value = nextLang;
}
type GroupedLine = Article & { qty: number };
const groupedReceipt = computed<GroupedLine[]>(() => {
  const map = new Map<number, GroupedLine>();
  for (const item of receipt.value) {
    const existing = map.get(item.id);
    if (existing) {
      existing.qty += 1;
    } else {
      map.set(item.id, { ...item, qty: 1 });
    }
  }
  return Array.from(map.values());
});

const total = computed(() => groupedReceipt.value.reduce((sum, item) => sum + item.price * item.qty, 0));

// --- FAB open/close with outside click ---
const openFabId = ref<number|null>(null);
function toggleFab(id: number, open: boolean) {
  openFabId.value = open ? id : (openFabId.value === id ? null : openFabId.value);
}
function handleDocumentClick(ev: MouseEvent) {
  const target = ev.target as HTMLElement;
  if (!target.closest('.receipt-fab')) {
    openFabId.value = null;
  }
}
onMounted(() => document.addEventListener('click', handleDocumentClick));
onBeforeUnmount(() => document.removeEventListener('click', handleDocumentClick));

// --- Edit quantity dialog state & logic ---
const editDialog = ref(false);
const editLineId = ref<number|null>(null);
const tempQty = ref<number>(0);
const qtyStep = 1; // adjust to 0.25 or 0.1 for finer control
const editLineName = computed(() => {
  if (editLineId.value == null) return '';
  const line = groupedReceipt.value.find(l => l.id === editLineId.value);
  return line ? line.name : '';
});

function openEdit(line: GroupedLine) {
  editLineId.value = line.id;
  tempQty.value = line.qty;
  editDialog.value = true;
}

function applyEdit() {
  if (editLineId.value == null) return;
  setLineQuantity(editLineId.value, tempQty.value);
  editDialog.value = false;
}

function setLineQuantity(id: number, newQty: number) {
  const currentQty = receipt.value.filter(a => a.id === id).length;
  if (newQty <= 0) {
    // remove completely
    receipt.value = receipt.value.filter(a => a.id !== id);
    return;
  }
  if (newQty === currentQty) return; // nothing to do
  if (newQty > currentQty) {
    const article = articles.value.find(a => a.id === id);
    if (!article) return;
    const toAdd = newQty - currentQty;
    for (let i = 0; i < toAdd; i++) receipt.value.push(article);
  } else {
    let toRemove = currentQty - newQty;
    for (let i = receipt.value.length - 1; i >= 0 && toRemove > 0; i--) {
      const rec = receipt.value[i];
      if (rec && rec.id === id) {
        receipt.value.splice(i, 1);
        toRemove--;
      }
    }
  }
}

// QFab quantity controls
function increment(id: number) {
  const article = articles.value.find(a => a.id === id);
  if (article) receipt.value.push(article);
  // keep FAB open
  void nextTick(() => { openFabId.value = id; });
}
function decrement(id: number) {
  const index = receipt.value.findIndex(a => a.id === id);
  if (index !== -1) receipt.value.splice(index, 1);
  // if still has items for this id, keep FAB open; else close
  const stillExists = receipt.value.some(a => a.id === id);
  void nextTick(() => { openFabId.value = stillExists ? id : null; });
}
function deleteLine(id: number) {
  receipt.value = receipt.value.filter(a => a.id !== id);
  if (openFabId.value === id) openFabId.value = null;
}

</script>


<style scoped>
.trash-drop-zone {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  pointer-events: auto;
  transition: background 0.2s;
}
</style>
