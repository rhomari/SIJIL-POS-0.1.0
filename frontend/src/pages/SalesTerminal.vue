<template>
  <q-page class="q-pa-md">
  <div class="row q-col-gutter-none">
      <!-- Articles Section -->
  <div class="col-12 col-md-7 col-xl-8 products-pane">
        <div class="row items-center no-wrap q-gutter-sm">
          <div class="col">
            <q-tabs v-model="selectedCategory" class="text-primary" align="left" dense shrink>
              <q-tab v-for="cat in visibleCategories" :key="cat" :name="cat" :label="cat" />
            </q-tabs>
            <!-- Search popup anchored to tabs container -->
            <q-menu
              v-model="searchOpen"
              target
              no-parent-event
              anchor="bottom left"
              self="top left"
              transition-show="jump-down"
              transition-hide="jump-up"
              content-class="search-menu-content"
            >
              <div class="q-pa-sm" style="min-width: 280px;">
                <q-input
                  ref="searchInputRef"
                  v-model="categorySearch"
                  outlined
                  clearable
                  debounce="200"
                  placeholder="Search categories or products"
                  class="search-glow"
                  autofocus
                >
                  <template #prepend>
                    <q-icon name="search" />
                  </template>
                </q-input>
              </div>
            </q-menu>
          </div>
          <q-btn-dropdown v-if="overflowCategories.length" dense flat icon="more_horiz" label="More" class="text-primary" content-class="text-primary">
            <q-list style="min-width: 160px;">
              <q-item v-for="cat in overflowCategories" :key="cat" clickable v-close-popup @click="selectedCategory = cat">
                <q-item-section>{{ cat }}</q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>
          
          <q-btn
            dense
            flat
            round
            icon="search"
            class="q-ml-xs"
            @click="() => { searchOpen = true; void nextTick(() => { const el = searchInputRef?.$el as HTMLElement | undefined; const native = el?.querySelector('input') as HTMLInputElement | null; native?.focus(); }); }"
          />
          <q-btn
            dense
            flat
            round
            icon="qr_code_scanner"
            class="q-ml-xs"
            @click="openScanner"
          />
        </div>
  <div class="row q-gutter-x-none q-gutter-y-sm q-mt-xs products-scroll" style="overflow-y: auto; max-height: 82vh; -webkit-overflow-scrolling: touch; overscroll-behavior: contain;">
          <div
            v-for="article in filteredArticles"
            :key="article.id"
            class="col-6 col-sm-4 five-col-md six-col-xl q-mb-sm product-tile"
          >
            <q-card
              bordered
              class="shadow-3 cursor-pointer product-card"
              style="height: 120px;"
              v-ripple="false"
              @click="handleAdd(article, $event)"
            >
              <!--LOCK START-->
              <q-img
                :src="article.image"
                :alt="article.name"
                style="width: 100%; height:119px;"
                :img-style="{ objectFit: 'contain', width: '100%', height: '100%' }"
                fit="contain"
                img-class="q-pa-xs bg-white product-img"
              >
                <div class="absolute-bottom product-legend text-center">
                  <div class="legend-name">{{ article.name }}</div>
                  <div class="legend-price">{{ article.price.toFixed(2) }} MAD</div>
                </div>
              </q-img>
              <!--LOCK END-->
            </q-card>
          </div>
        </div>
        <!-- Mobile receipt open button (full-width bar above footer) -->
        <q-btn
          v-if="$q.screen.lt.md && !receiptDialog"
          color="primary"
          unelevated
          icon="receipt_long"
          :label="$t('receipt') + (groupedReceipt.length ? ' (' + groupedReceipt.length + ')' : '')"
          class="mobile-receipt-bar"
          @click="openReceiptDialog"
        />
      </div>
      <!-- Receipt Section -->
  <div class="col-12 col-md-5 col-xl-4 receipt-pane q-pl-sm" v-if="$q.screen.gt.sm">
  <q-card class="receipt-card">
          <q-card-section class="bg-grey-2">
            <div class="row items-center no-wrap">
              <div class="text-h6">{{ $t('receipt') }}</div>
              <q-space />
              <q-btn dense flat icon="playlist_add_check" size="sm" :label="$t('holdList') + ' (' + holds.length + ')'" @click="openHoldList" />
            </div>
          </q-card-section>
          <q-separator />
          <div class="relative-position receipt-body">
            <q-list ref="receiptListRef" class="receipt-list">
              <q-item v-for="item in groupedReceipt" :key="item.id" :data-line-id="item.id"
                draggable="true"
                @dragstart="onDragStart(item.id)"
                @dragend="onDragEnd"
                :class="{ 'dragging': draggingId === item.id, 'receipt-added': lastAddedId === item.id }"
              >
                <q-item-section avatar>
                  <q-fab
                    color="primary"
                    text-color="white"
                    :label="item.qty.toString()"
                    padding="xs"
                    dense
                    direction="right"
                    icon="more_vert"
                    class="receipt-fab"
                    :model-value="openFabId === item.id"
                    @update:model-value="val => toggleFab(item.id, val)"
                  >
                    <q-fab-action color="primary" icon="add" @click.stop="increment(item.id)" />
                    <q-fab-action color="primary" icon="remove" :disable="item.qty <= 0" @click.stop="decrement(item.id)" />
                    <q-fab-action color="secondary" icon="edit" @click.stop="openEdit(item)" />
                    <q-fab-action color="info" icon="local_offer" @click.stop="openReduction(item)" />
                    <q-fab-action color="accent" icon="edit_note" @click.stop="openNote(item)" />
                    <q-fab-action color="negative" icon="delete" @click.stop="deleteLine(item.id)" />
                  </q-fab>
                </q-item-section>
                <q-item-section>
                  <div class="row items-center justify-between w-100">
                    <span>{{ item.name }}</span>
                    <span class="text-weight-bold">{{ lineTotal(item).toFixed(2) }} MAD</span>
                  </div>
                  <div v-if="hasReduction(item.id)" class="text-caption text-grey-7 q-mt-xs">
                    <q-icon name="local_offer" size="16px" class="q-mr-xs" /> {{ reductionLabel(item.id) }}
                  </div>
                  <div v-if="notesById[item.id]" class="text-caption text-grey-7 q-mt-xs ellipsis">
                    <q-icon name="notes" size="16px" class="q-mr-xs" /> {{ notesById[item.id] }}
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
              :style="draggingId !== null ? 'z-index:1100;pointer-events:auto;' : 'pointer-events:none;'"
            >
              <q-icon
                name="delete"
                :color="draggingId !== null ? 'negative' : 'grey-6'"
                size="64px"
                :class="draggingId !== null ? 'q-animate-bounce' : ''"
              />
            </div>
          </div>
          
          <!-- Total Reduction Dialog -->
          <q-dialog v-model="totalReductionDialog">
            <q-card style="width:420px; max-width:95vw;" class="q-pa-sm">
              <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Total Discount</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
              </q-card-section>
              <q-card-section class="q-pt-sm">
                <div class="row items-center q-gutter-sm q-mb-sm">
                  <q-btn-toggle
                    v-model="totalReductionKind"
                    :options="[
                      { label: '%', value: 'percent' },
                      { label: 'MAD', value: 'amount' }
                    ]"
                    unelevated
                    spread
                    size="sm"
                  />
                </div>
                <q-input
                  v-model="totalReductionValueStr"
                  type="number"
                  :min="0"
                  :max="totalReductionKind === 'percent' ? 100 : undefined"
                  :step="totalReductionKind === 'percent' ? 0.1 : 0.01"
                  outlined
                  dense
                  :suffix="totalReductionKind === 'percent' ? '%' : 'MAD'"
                  :placeholder="totalReductionKind === 'percent' ? 'Enter discount %' : 'Enter discount amount'"
                />
              </q-card-section>
              <q-separator inset />
              <q-card-actions align="right">
                <q-btn v-if="hasTotalReduction" flat color="negative" :label="$t('remove')" @click="deleteTotalReduction" />
                <q-btn flat :label="$t('cancel')" v-close-popup />
                <q-btn color="primary" :label="$t('apply')" @click="applyTotalReduction" />
              </q-card-actions>
            </q-card>
          </q-dialog>
          <!-- Reduction Dialog -->
          <q-dialog v-model="reductionDialog">
            <q-card style="width:420px; max-width:95vw;" class="q-pa-sm">
              <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Item Discount</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
              </q-card-section>
              <q-card-section class="q-pt-sm">
                <div class="text-subtitle2 text-primary">{{ reductionLineName }}</div>
                <div class="row items-center q-gutter-sm q-mt-sm">
                  <q-btn-toggle
                    v-model="reductionKind"
                    :options="[
                      { label: '%', value: 'percent' },
                      { label: 'MAD', value: 'amount' }
                    ]"
                    unelevated
                    spread
                    size="sm"
                  />
                </div>
                <q-input
                  v-model="reductionValueStr"
                  type="number"
                  :min="0"
                  :max="reductionKind === 'percent' ? 100 : undefined"
                  :step="reductionKind === 'percent' ? 0.1 : 0.01"
                  outlined
                  dense
                  :suffix="reductionKind === 'percent' ? '%' : 'MAD'"
                  :placeholder="reductionKind === 'percent' ? 'Enter discount %' : 'Enter discount amount'"
                  class="q-mt-sm"
                />
              </q-card-section>
              <q-separator inset />
              <q-card-actions align="right">
                <q-btn v-if="hasSelectedReduction" flat color="negative" :label="$t('remove')" @click="deleteReduction" />
                <q-btn flat :label="$t('cancel')" v-close-popup />
                <q-btn color="primary" :label="$t('apply')" @click="applyReduction" />
              </q-card-actions>
            </q-card>
          </q-dialog>
          <q-dialog v-model="editDialog">
            <q-card style="width:320px; max-width:90vw;" class="q-pa-sm">
              <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Edit Quantity</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
              </q-card-section>
              <q-card-section class="q-pt-sm">
                <div class="text-subtitle2 text-primary">{{ editLineName }}</div>
                <div class="q-mt-sm">
                  <div class="row items-center">
                    <q-input v-model="tempQtyStr" dense outlined class="col" inputmode="decimal" pattern="[0-9]*" :placeholder="$t('editQuantity')" />
                  </div>
                  <div class="numpad grid q-mt-sm">
                    <q-btn dense class="nkey" @click="tap('7')">7</q-btn>
                    <q-btn dense class="nkey" @click="tap('8')">8</q-btn>
                    <q-btn dense class="nkey" @click="tap('9')">9</q-btn>

                    <q-btn dense class="nkey" @click="tap('4')">4</q-btn>
                    <q-btn dense class="nkey" @click="tap('5')">5</q-btn>
                    <q-btn dense class="nkey" @click="tap('6')">6</q-btn>

                    <q-btn dense class="nkey" @click="tap('1')">1</q-btn>
                    <q-btn dense class="nkey" @click="tap('2')">2</q-btn>
                    <q-btn dense class="nkey" @click="tap('3')">3</q-btn>

                    <q-btn dense class="nkey" @click="tap('0')">0</q-btn>
                    <q-btn dense class="nkey" @click="tap('.')">.</q-btn>
                    <q-btn dense class="nkey" @click="backspace()"><q-icon name="backspace" /></q-btn>
                  </div>
                  <div class="row q-col-gutter-sm q-mt-sm">
                    <div class="col-12"><q-btn outline color="negative" class="full-width" @click="clearNumpad">CLR</q-btn></div>
                  </div>
                  <div class="text-caption text-grey-7 q-mt-xs">{{ $t('editHint') }}</div>
                </div>
              </q-card-section>
              <q-separator inset />
              <q-card-actions align="right">
                  <q-btn flat :label="$t('cancel')" v-close-popup />
                  <q-btn color="primary" :label="$t('apply')" @click="applyEdit" />
              </q-card-actions>
            </q-card>
          </q-dialog>
          <!-- Note Dialog -->
          <q-dialog v-model="noteDialog">
            <q-card style="width:420px; max-width:95vw;" class="q-pa-sm">
              <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">Item Note</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
              </q-card-section>
              <q-card-section class="q-pt-sm">
                <div class="text-subtitle2 text-primary">{{ noteLineName }}</div>
                <q-input v-model="noteText" type="textarea" autogrow outlined dense placeholder="Add a note..." />
              </q-card-section>
              <q-separator inset />
              <q-card-actions align="right">
                <q-btn v-if="hasSelectedNote" flat color="negative" :label="$t('remove')" @click="deleteNote" />
                <q-btn flat :label="$t('cancel')" v-close-popup />
                <q-btn color="primary" :label="$t('apply')" @click="applyNote" />
              </q-card-actions>
            </q-card>
          </q-dialog>
          <q-separator />
          <q-card-section>
            <div class="total-summary row items-center no-wrap q-mb-sm">
              <div class="total-chip row items-center no-wrap">
                <q-icon name="receipt_long" size="18px" class="q-mr-sm" />
                <span class="total-label">{{ $t('total') }}</span>
              </div>
              <q-btn dense flat round icon="local_offer" class="q-ml-sm" :color="hasTotalReduction ? 'info' : 'grey-7'" @click="openTotalReduction" />
              <q-space />
              <div class="total-amount row items-end no-wrap">
                <span class="amount">{{ total.toFixed(2) }}</span>
                <span class="currency q-ml-xs">MAD</span>
              </div>
            </div>
            <div v-if="hasTotalReduction" class="row items-center text-caption text-grey-7 q-mb-sm">
              <q-icon name="local_offer" size="16px" class="q-mr-xs" />
              <span class="q-mr-sm">{{ totalReductionLabel }}</span>
              <span>−{{ totalDiscountAmount.toFixed(2) }} MAD</span>
            </div>
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <q-btn
                  unelevated
                  color="positive"
                  icon="shopping_cart_checkout"
                  :label="$q.screen.lt.md ? undefined : $t('checkout')"
                  class="full-width"
                  :disable="groupedReceipt.length === 0"
                  @click="checkout"
                />
              </div>
              <div class="col-6">
                <q-btn
                  unelevated
                  color="negative"
                  icon="restart_alt"
                  :label="$q.screen.lt.md ? undefined : $t('reset')"
                  class="full-width"
                  :disable="groupedReceipt.length === 0"
                  @click="resetReceipt"
                />
              </div>
              <div class="col-6 q-mt-sm">
                <q-btn
                  unelevated
                  color="warning"
                  icon="pause_circle"
                  :label="$q.screen.lt.md ? undefined : $t('hold')"
                  class="full-width"
                  :disable="groupedReceipt.length === 0"
                  @click="putOnHold"
                />
              </div>
              <div class="col-6 q-mt-sm">
                <q-btn
                  unelevated
                  color="primary"
                  icon="add_shopping_cart"
                  :label="$q.screen.lt.md ? undefined : $t('newOrder')"
                  class="full-width"
                  @click="launchNewOrder"
                />
              </div>
              
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  <!-- Mobile Receipt Dialog (bottom sheet) -->
    <q-dialog v-model="receiptDialog" position="bottom" transition-show="slide-up" transition-hide="slide-down">
      <q-card style="width: 100vw; max-width: 100vw; max-height: 85vh;" class="q-pa-none">
        <q-card-section class="bg-grey-2">
          <div class="row items-center no-wrap">
            <div class="text-h6">{{ $t('receipt') }}</div>
            <q-space />
            <q-btn dense flat icon="playlist_add_check" size="sm" :label="$t('holdList') + ' (' + holds.length + ')'" @click="openHoldFromMobile" />
            <q-btn flat round dense icon="close" v-close-popup class="q-ml-sm" />
          </div>
        </q-card-section>
        <q-separator />
        <div class="relative-position" style="overflow-y: auto; max-height: calc(85vh - 160px);">
          <q-list class="receipt-list q-px-sm q-pt-sm">
            <q-item v-for="item in groupedReceipt" :key="item.id" :data-line-id="item.id">
              <q-item-section avatar>
                <q-fab
                  color="primary"
                  text-color="white"
                  :label="item.qty.toString()"
                  padding="xs"
                  dense
                  direction="right"
                  icon="more_vert"
                  class="receipt-fab"
                  :model-value="openFabId === item.id"
                  @update:model-value="val => toggleFab(item.id, val)"
                >
                  <q-fab-action color="primary" icon="add" @click.stop="increment(item.id)" />
                  <q-fab-action color="primary" icon="remove" :disable="item.qty <= 0" @click.stop="decrement(item.id)" />
                  <q-fab-action color="secondary" icon="edit" @click.stop="openEdit(item)" />
                  <q-fab-action color="info" icon="local_offer" @click.stop="openReduction(item)" />
                  <q-fab-action color="accent" icon="edit_note" @click.stop="openNote(item)" />
                  <q-fab-action color="negative" icon="delete" @click.stop="deleteLine(item.id)" />
                </q-fab>
              </q-item-section>
              <q-item-section>
                <div class="row items-center justify-between w-100">
                  <span>{{ item.name }}</span>
                  <span class="text-weight-bold">{{ lineTotal(item).toFixed(2) }} MAD</span>
                </div>
                <div v-if="hasReduction(item.id)" class="text-caption text-grey-7 q-mt-xs">
                  <q-icon name="local_offer" size="16px" class="q-mr-xs" /> {{ reductionLabel(item.id) }}
                </div>
                <div v-if="notesById[item.id]" class="text-caption text-grey-7 q-mt-xs ellipsis">
                  <q-icon name="notes" size="16px" class="q-mr-xs" /> {{ notesById[item.id] }}
                </div>
              </q-item-section>
            </q-item>
          </q-list>
        </div>
        <q-separator />
        <q-card-section>
          <div class="total-summary row items-center no-wrap q-mb-sm">
            <div class="total-chip row items-center no-wrap">
              <q-icon name="receipt_long" size="18px" class="q-mr-sm" />
              <span class="total-label">{{ $t('total') }}</span>
            </div>
            <q-btn dense flat round icon="local_offer" class="q-ml-sm" :color="hasTotalReduction ? 'info' : 'grey-7'" @click="openTotalReduction" />
            <q-space />
            <div class="total-amount row items-end no-wrap">
              <span class="amount">{{ total.toFixed(2) }}</span>
              <span class="currency q-ml-xs">MAD</span>
            </div>
          </div>
          <div v-if="hasTotalReduction" class="row items-center text-caption text-grey-7 q-mb-sm">
            <q-icon name="local_offer" size="16px" class="q-mr-xs" />
            <span class="q-mr-sm">{{ totalReductionLabel }}</span>
            <span>−{{ totalDiscountAmount.toFixed(2) }} MAD</span>
          </div>
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-btn
                unelevated
                color="positive"
                icon="shopping_cart_checkout"
                :label="$q.screen.lt.md ? undefined : $t('checkout')"
                class="full-width"
                :disable="groupedReceipt.length === 0"
                @click="checkout"
              />
            </div>
            <div class="col-6">
              <q-btn
                unelevated
                color="negative"
                icon="restart_alt"
                :label="$q.screen.lt.md ? undefined : $t('reset')"
                class="full-width"
                :disable="groupedReceipt.length === 0"
                @click="resetReceipt"
              />
            </div>
            <div class="col-6 q-mt-sm">
              <q-btn
                unelevated
                color="warning"
                icon="pause_circle"
                :label="$q.screen.lt.md ? undefined : $t('hold')"
                class="full-width"
                :disable="groupedReceipt.length === 0"
                @click="putOnHold"
              />
            </div>
            <div class="col-6 q-mt-sm">
              <q-btn
                unelevated
                color="primary"
                icon="add_shopping_cart"
                :label="$q.screen.lt.md ? undefined : $t('newOrder')"
                class="full-width"
                @click="launchNewOrder"
              />
            </div>
            
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
    <q-dialog v-model="scannerDialog" maximized>
      <q-card class="q-pa-md">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{$t('scan') }}</div>
          <q-space />
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>
        <q-card-section class="q-pt-sm flex flex-center">
          <CameraScanner @decoded="onScannerDecoded" @close="scannerDialog=false" />
        </q-card-section>
      </q-card>
    </q-dialog>
    
    <!-- Hold List Dialog (global, works on all sizes) -->
    <q-dialog v-model="holdListDialog">
      <q-card style="width: 520px; max-width: 95vw;">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{ $t('holdList') }}</div>
          <q-space />
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>
        <q-card-section class="q-pt-sm">
          <div v-if="holds.length === 0" class="text-grey-7">{{ $t('noHolds') }}</div>
          <q-list v-else separator>
            <q-item v-for="h in holds" :key="h.id">
              <q-item-section>
                <div class="row items-center justify-between">
                  <div class="text-subtitle2">#{{ h.id }}</div>
                  <div class="text-caption text-grey-7">{{ formatDate(h.createdAt) }}</div>
                </div>
                <div class="row items-center justify-between q-mt-xs">
                  <div class="text-body2">{{ h.lines.length }} {{ $t('items') }}</div>
                  <div class="text-weight-bold">{{ h.total.toFixed(2) }} MAD</div>
                </div>
              </q-item-section>
              <q-item-section side>
                <div class="column q-gutter-xs">
                  <q-btn dense color="primary" :label="$t('resume')" @click="resumeHold(h.id)" />
                  <q-btn dense color="negative" flat :label="$t('remove')" @click="removeHold(h.id)" />
                </div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat :label="$t('close')" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
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
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch, type ComponentPublicInstance } from 'vue';
// import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import type { QList } from 'quasar';

const categories = [
  'Stationery', 'Books', 'Office', 'Gadgets', 'Supplies',
  'Arts', 'Crafts', 'Paper', 'Ink', 'Hardware',
  'Software', 'Snacks', 'Beverages', 'Cleaning', 'Storage'
] as const;
type Category = typeof categories[number];
const selectedCategory = ref<Category>(categories[0]);
const categorySearch = ref<string>('');
const searchInputRef = ref<ComponentPublicInstance | null>(null);
// No explicit state needed for popup proxy
const searchOpen = ref(false);
type Article = {
  id: number;
  name: string;
  price: number;
  image: string;
  category: string;
  barcode: string;
};

// Base seed items
const baseArticles: Article[] = [
  { id: 1, name: 'Blue Ballpoint Pen', price: 1.2, image: '/images/1.png', category: 'Stationery', barcode: '2000000000001' },
  { id: 2, name: 'A5 Spiral Notebook', price: 3.5, image: '/images/2.png', category: 'Stationery', barcode: '887930951998' },
  { id: 3, name: 'Highlighter Set', price: 4.0, image: '/images/3.png', category: 'Stationery', barcode: '2000000000003' },
  { id: 4, name: 'Stapler', price: 5.5, image: '/images/4.png', category: 'Office', barcode: '2000000000004' },
  { id: 5, name: 'Desk Organizer', price: 9.9, image: '/images/5.png', category: 'Office', barcode: '8017486106678' },
  { id: 6, name: 'Mystery Novel', price: 7.99, image: '/images/6.png', category: 'Books', barcode: '2000000000006' },
  { id: 7, name: 'Science Textbook', price: 24.5, image: '/images/7.png', category: 'Books', barcode: '2000000000007' }
];

// Generate extra dummy products (120 more), reuse images 1..35 and cycle categories
const extraArticles: Article[] = Array.from({ length: 120 }, (_v, idx) => {
  const id = baseArticles.length + idx + 1; // start at 8
  const imgNum = ((id - 1) % 35) + 1; // 1..35 loop
  const cat = categories[id % categories.length] ?? 'Stationery';
  const price = +( (1 + ((id % 30) * 0.5)).toFixed(2) ); // 1.00 .. 16.00 range
  return {
    id,
    name: `Dummy Product ${id}`,
    price,
    image: `/images/${imgNum}.png`,
    category: cat,
    barcode: String(2000000000000 + id)
  } satisfies Article;
});

const articles = ref<Article[]>([...baseArticles, ...extraArticles]);

const filteredArticles = computed(() => {
  const q = categorySearch.value.trim().toLowerCase();
  if (q) {
    return articles.value.filter((a: Article) =>
      a.name.toLowerCase().includes(q) || a.category.toLowerCase().includes(q)
    );
  }
  return articles.value.filter((a: Article) => a.category === selectedCategory.value);
});
// Category filtering for tabs/dropdown
const filteredCategories = computed<Category[]>(() => {
  const q = categorySearch.value.trim().toLowerCase();
  if (!q) return Array.from(categories);
  const nameMatches = new Set(Array.from(categories).filter(c => c.toLowerCase().includes(q)));
  const productMatches = new Set(
    articles.value
      .filter((a: Article) => a.name.toLowerCase().includes(q) || a.category.toLowerCase().includes(q))
      .map(a => a.category as Category)
  );
  const merged = new Set<Category>([...Array.from(nameMatches), ...Array.from(productMatches)]);
  return Array.from(categories).filter(c => merged.has(c));
});

const receipt = ref<Article[]>([]);
function addToReceipt(article: Article) {
  receipt.value.push(article);
}
// Audio feedback (single shared AudioContext)
let audioCtx: AudioContext | null = null;
function playAddBeep() {
  try {
    if (!audioCtx) {
      // Some browsers (older Safari) expose webkitAudioContext; define a safe typed fallback.
      type AudioContextConstructor = { new(): AudioContext };
      const AC: AudioContextConstructor | undefined = (window as unknown as { AudioContext?: AudioContextConstructor }).AudioContext;
      const WAC: AudioContextConstructor | undefined = (window as unknown as { webkitAudioContext?: AudioContextConstructor }).webkitAudioContext;
      const Ctor = AC || WAC;
      if (!Ctor) return; // environment unsupported
      audioCtx = new Ctor();
    }
    if (audioCtx.state === 'suspended') {
      // Attempt resume (may still be blocked until user gesture)
      void audioCtx.resume().catch(() => {});
    }
    const duration = 0.09; // seconds
    const osc = audioCtx.createOscillator();
    const gain = audioCtx.createGain();
    osc.type = 'triangle';
    osc.frequency.value = 2093; // A6 high pitched retail-like beep
    gain.gain.setValueAtTime(0.001, audioCtx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.25, audioCtx.currentTime + 0.005);
    gain.gain.exponentialRampToValueAtTime(0.0001, audioCtx.currentTime + duration);
    osc.connect(gain).connect(audioCtx.destination);
    osc.start();
    osc.stop(audioCtx.currentTime + duration + 0.02);
  } catch {
    // ignore audio errors silently (no blocking)
  }
}
// Barcode scan handling
// --- Barcode normalization & validation (supports GTIN-8,12,13,14) ---
function computeGTINCheckDigit(body: string): number {
  // body: all digits except final check digit
  // Weighting from rightmost moving left: 3,1,3,1...
  let sum = 0;
  for (let i = body.length - 1, pos = 0; i >= 0; i--, pos++) {
    const d = body.charCodeAt(i) - 48;
    if (d < 0 || d > 9) return -1; // invalid char
    const weight = (pos % 2 === 0) ? 3 : 1; // start with 3 on rightmost
    sum += d * weight;
  }
  const mod = sum % 10;
  return (10 - mod) % 10;
}
function isValidGTIN(code: string): boolean {
  if (!/^[0-9]+$/.test(code)) return false;
  if (![8, 12, 13, 14].includes(code.length)) return false;
  const body = code.slice(0, -1);
  const check = computeGTINCheckDigit(body);
  return check === (code.charCodeAt(code.length - 1) - 48);
}
function generateBarcodeCandidates(raw: string): string[] {
  const candidates = new Set<string>();
  const digits = raw.trim();
  if (!digits) return [];
  candidates.add(digits);
  // UPC-A (12) => EAN-13 by prefixing 0
  if (/^\d{12}$/.test(digits)) {
    candidates.add('0' + digits);
  }
  // EAN-13 starting with 0 might correspond to UPC-A without leading 0
  if (/^0\d{12}$/.test(digits)) {
    candidates.add(digits.slice(1));
  }
  // GTIN-14: sometimes data set stores trimmed (drop leading) if leading is 0
  if (/^0\d{13}$/.test(digits)) {
    candidates.add(digits.slice(1)); // drop leading 0 -> EAN-13
  }
  return Array.from(candidates);
}
function handleDecodedBarcode(code: string): boolean {
  const candidates = generateBarcodeCandidates(code);
  let matched: Article | undefined;
  for (const c of candidates) {
    matched = articles.value.find(a => a.barcode === c);
    if (matched) {
        addToReceipt(matched);
        playAddBeep();
      $q.notify({ type: 'positive', message: `${matched.name} +1` });
      return true;
    }
  }
  // If no match and code appears to be GTIN but invalid checksum -> notify invalid
  const primary = candidates[0] ?? code;
  if (/^\d{8}$/.test(primary) || /^\d{12}$/.test(primary) || /^\d{13}$/.test(primary) || /^\d{14}$/.test(primary)) {
    if (!isValidGTIN(primary)) {
      $q.notify({ type: 'warning', message: `Invalid barcode checksum (${primary})` });
      return false;
    }
  }
  $q.notify({ type: 'warning', message: `Code ${primary} not found` });
  return false;
}
// animated add
// Use explicit QList instance type (component public instance) instead of any for ESLint compliance
const receiptListRef = ref<InstanceType<typeof QList> | null>(null); // q-list component
const lastAddedId = ref<number | null>(null);
// Input performance helpers
const isTouchDevice = ref<boolean>(false);
const prefersReducedMotion = ref<boolean>(false);
onMounted(() => {
  try {
    isTouchDevice.value = matchMedia('(pointer: coarse)').matches || 'ontouchstart' in window || navigator.maxTouchPoints > 0;
  } catch { isTouchDevice.value = false; }
  try {
    prefersReducedMotion.value = matchMedia('(prefers-reduced-motion: reduce)').matches;
  } catch { prefersReducedMotion.value = false; }
});
function handleAdd(article: Article, ev: MouseEvent) {
  // Instant path for touch devices or when reduced motion is preferred
  if (isTouchDevice.value || prefersReducedMotion.value) {
    addToReceipt(article);
    lastAddedId.value = article.id;
    playAddBeep();
    return;
  }
  const card = (ev.currentTarget as HTMLElement) || null;
  const listComp = receiptListRef.value; // QList instance or null
  const listEl: HTMLElement | null = listComp ? (listComp.$el as HTMLElement) : null;
  if (!card || !listEl) {
    addToReceipt(article);
    lastAddedId.value = article.id;
    playAddBeep();
    return;
  }
  // Capture non-null list element for use inside async callbacks without re-narrowing
  const stableListEl: HTMLElement = listEl;
  const start = card.getBoundingClientRect();
  addToReceipt(article);
  lastAddedId.value = article.id;
  playAddBeep();
  // Wait for receipt DOM update so the target line exists / updated
  void nextTick(() => {
  const targetLine = document.querySelector<HTMLElement>(`.receipt-list [data-line-id="${article.id}"]`);
  // listEl is guaranteed (early return above) but keep safe fallback
  const baseEl = targetLine ?? stableListEl;
  const targetRect = baseEl.getBoundingClientRect();
    const fly = card.cloneNode(true) as HTMLElement;
    fly.classList.add('fly-clone');
    Object.assign(fly.style, {
      position: 'fixed',
      margin: '0',
      top: start.top + 'px',
      left: start.left + 'px',
      width: start.width + 'px',
      height: start.height + 'px',
      zIndex: '99999',
      pointerEvents: 'none',
      transform: 'translate(0,0) scale(1)',
      transition: 'transform .65s cubic-bezier(.55,.06,.27,.99), opacity .65s ease'
    });
    document.body.appendChild(fly);
    // Force reflow
    void fly.offsetWidth;
    // Compute destination (center of target line)
    const endX = targetRect.left + targetRect.width * 0.1; // left padding inside receipt
    const endY = targetRect.top + targetRect.height * 0.5;
    const translateX = endX - start.left;
    const translateY = endY - start.top;
    requestAnimationFrame(() => {
      fly.style.transform = `translate(${translateX}px, ${translateY}px) scale(.28)`;
  // Keep some opacity so the flying object remains visible
  fly.style.opacity = '0.45';
    });
    const cleanup = () => fly.remove();
    fly.addEventListener('transitionend', cleanup, { once: true });
    setTimeout(cleanup, 900);
    setTimeout(() => { if (lastAddedId.value === article.id) lastAddedId.value = null; }, 1000);
  });
}
import { useI18n } from 'vue-i18n';
import CameraScanner from '../components/CameraScanner.vue';
const { t } = useI18n();
// const router = useRouter();
// Mobile receipt dialog state
const receiptDialog = ref(false);
// Optional per-line quantity overrides (allow decimals)
const qtyOverrides = ref<Record<number, number>>({});
// Optional per-line notes
const notesById = ref<Record<number, string>>({});
// Optional per-line reductions: percent or fixed amount (union for backward compatibility)
type Reduction = { kind: 'percent' | 'amount'; value: number };
const reductionsById = ref<Record<number, number | Reduction>>({});

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
  const arr = Array.from(map.values());
  // apply overrides if present
  return arr.map(line => {
    const ov = qtyOverrides.value[line.id];
    return (ov !== undefined) ? { ...line, qty: ov } : line;
  });
});

// Reset filters on Escape key
function onGlobalKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    if (categorySearch.value) {
      categorySearch.value = '';
    }
    if (searchOpen.value) {
      searchOpen.value = false;
    }
    return;
  }
  if (e.key === 'F3') {
    e.preventDefault();
    searchOpen.value = true;
    void nextTick(() => {
      const el = searchInputRef.value?.$el as HTMLElement | undefined;
      const native = el?.querySelector('input') as HTMLInputElement | null;
      if (native) native.focus();
    });
  }
}
onMounted(() => window.addEventListener('keydown', onGlobalKeydown));
onBeforeUnmount(() => window.removeEventListener('keydown', onGlobalKeydown));

// Helper: compute line total with reduction
function lineTotal(line: GroupedLine) {
  const base = line.price * line.qty;
  const r = reductionsById.value[line.id];
  if (r == null) return base;
  if (typeof r === 'number') {
    const pct = Math.max(0, Math.min(100, r));
    return base * (1 - pct / 100);
  }
  if (r.kind === 'percent') {
    const pct = Math.max(0, Math.min(100, r.value || 0));
    return base * (1 - pct / 100);
  }
  const amt = Math.max(0, r.value || 0);
  return Math.max(0, base - amt);
}
function hasReduction(id: number) {
  const r = reductionsById.value[id];
  return r != null && (typeof r === 'number' ? r > 0 : r.value > 0);
}
function reductionLabel(id: number) {
  const r = reductionsById.value[id];
  if (r == null) return '';
  if (typeof r === 'number') return `-${Math.round(r)}%`;
  return r.kind === 'percent' ? `-${Math.round(r.value)}%` : `-${(r.value || 0).toFixed(2)} MAD`;
}

const subtotal = computed(() => groupedReceipt.value.reduce((sum, item) => sum + lineTotal(item), 0));
// Total-level reduction (percent or amount)
const totalReduction = ref<number | Reduction | null>(null);
const hasTotalReduction = computed(() => {
  const r = totalReduction.value;
  if (r == null) return false;
  return typeof r === 'number' ? r > 0 : (r.value || 0) > 0;
});
const totalReductionLabel = computed(() => {
  const r = totalReduction.value;
  if (r == null) return '';
  if (typeof r === 'number') return `-${Math.round(r)}%`;
  return r.kind === 'percent' ? `-${Math.round(r.value)}%` : `-${(r.value || 0).toFixed(2)} MAD`;
});
const total = computed(() => {
  const base = subtotal.value;
  const r = totalReduction.value;
  if (r == null) return base;
  if (typeof r === 'number') {
    const pct = Math.max(0, Math.min(100, r));
    return base * (1 - pct / 100);
  }
  if (r.kind === 'percent') {
    const pct = Math.max(0, Math.min(100, r.value || 0));
    return base * (1 - pct / 100);
  }
  const amt = Math.max(0, r.value || 0);
  return Math.max(0, base - amt);
});
const totalDiscountAmount = computed(() => Math.max(0, subtotal.value - total.value));

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
const passiveClickOptions: AddEventListenerOptions = { passive: true };
onMounted(() => document.addEventListener('click', handleDocumentClick, passiveClickOptions));
onBeforeUnmount(() => document.removeEventListener('click', handleDocumentClick, passiveClickOptions));


const editDialog = ref(false);
const editLineId = ref<number|null>(null);
// string buffer for numpad input
const tempQtyStr = ref<string>('');
const editLineName = computed(() => {
  if (editLineId.value == null) return '';
  const line = groupedReceipt.value.find(l => l.id === editLineId.value);
  return line ? line.name : '';
});

function openEdit(line: GroupedLine) {
  editLineId.value = line.id;
  tempQtyStr.value = String(line.qty ?? 0);
  editDialog.value = true;
}

function applyEdit() {
  if (editLineId.value == null) return;
  const n = parseFloat(tempQtyStr.value.replace(',', '.'));
  const newQty = Number.isFinite(n) ? Math.max(0, n) : 0;
  setLineQuantity(editLineId.value, newQty);
  editDialog.value = false;
}

// Numpad helpers
function tap(ch: string) {
  if (ch === '.') {
    if (tempQtyStr.value.includes('.')) return;
    tempQtyStr.value = tempQtyStr.value === '' ? '0.' : tempQtyStr.value + '.';
    return;
  }
  // digits
  tempQtyStr.value += ch;
}
function backspace() {
  tempQtyStr.value = tempQtyStr.value.slice(0, -1);
}
function clearNumpad() {
  tempQtyStr.value = '';
}

// Note editing
const noteDialog = ref(false);
const noteLineId = ref<number|null>(null);
const noteText = ref<string>('');
const noteLineName = computed(() => {
  if (noteLineId.value == null) return '';
  const line = groupedReceipt.value.find(l => l.id === noteLineId.value);
  return line ? line.name : '';
});
const hasSelectedNote = computed(() => {
  return noteLineId.value != null && Boolean(notesById.value[noteLineId.value]);
});
function openNote(line: GroupedLine) {
  noteLineId.value = line.id;
  noteText.value = notesById.value[line.id] ?? '';
  noteDialog.value = true;
}
function applyNote() {
  if (noteLineId.value == null) return;
  const id = noteLineId.value;
  const text = noteText.value.trim();
  if (text) {
    notesById.value = { ...notesById.value, [id]: text };
  } else {
    const rest = { ...notesById.value };
    delete rest[id];
    notesById.value = rest;
  }
  noteDialog.value = false;
}
function deleteNote() {
  if (noteLineId.value == null) return;
  const id = noteLineId.value;
  if (notesById.value[id] !== undefined) {
    const rest = { ...notesById.value };
    delete rest[id];
    notesById.value = rest;
  }
  noteDialog.value = false;
}

// Reduction editing
const reductionDialog = ref(false);
const reductionLineId = ref<number|null>(null);
const reductionKind = ref<'percent' | 'amount'>('percent');
const reductionValueStr = ref<string>('');
const reductionLineName = computed(() => {
  if (reductionLineId.value == null) return '';
  const line = groupedReceipt.value.find(l => l.id === reductionLineId.value);
  return line ? line.name : '';
});
const hasSelectedReduction = computed(() => {
  return reductionLineId.value != null && hasReduction(reductionLineId.value);
});
function openReduction(line: GroupedLine) {
  reductionLineId.value = line.id;
  const current = reductionsById.value[line.id];
  if (current == null) {
    reductionKind.value = 'percent';
    reductionValueStr.value = '';
  } else if (typeof current === 'number') {
    reductionKind.value = 'percent';
    reductionValueStr.value = String(current);
  } else {
    reductionKind.value = current.kind;
    reductionValueStr.value = String(current.value ?? '');
  }
  reductionDialog.value = true;
}
function applyReduction() {
  if (reductionLineId.value == null) return;
  const id = reductionLineId.value;
  const n = parseFloat(String(reductionValueStr.value).replace(',', '.'));
  if (reductionKind.value === 'percent') {
    const pct = Number.isFinite(n) ? Math.max(0, Math.min(100, n)) : 0;
    if (pct > 0) {
      // store percent as a number for backward compatibility
      reductionsById.value = { ...reductionsById.value, [id]: pct };
    } else {
      const rest = { ...reductionsById.value } as Record<number, number | Reduction>;
      delete rest[id];
      reductionsById.value = rest;
    }
  } else {
    const amt = Number.isFinite(n) ? Math.max(0, n) : 0;
    if (amt > 0) {
      reductionsById.value = { ...reductionsById.value, [id]: { kind: 'amount', value: amt } };
    } else {
      const rest = { ...reductionsById.value } as Record<number, number | Reduction>;
      delete rest[id];
      reductionsById.value = rest;
    }
  }
  reductionDialog.value = false;
}
function deleteReduction() {
  if (reductionLineId.value == null) return;
  const id = reductionLineId.value;
  if (reductionsById.value[id] !== undefined) {
    const rest = { ...reductionsById.value };
    delete rest[id];
    reductionsById.value = rest;
  }
  reductionDialog.value = false;
}

function setLineQuantity(id: number, newQty: number) {
  if (newQty <= 0) {
    // remove completely (and clear override)
    receipt.value = receipt.value.filter(a => a.id !== id);
  const rest = { ...qtyOverrides.value };
  delete rest[id];
  qtyOverrides.value = rest;
    return;
  }
  // set override to allow decimals
  qtyOverrides.value = { ...qtyOverrides.value, [id]: newQty };
  // ensure at least one instance exists in receipt to render the line
  const exists = receipt.value.some(a => a.id === id);
  if (!exists) {
    const article = articles.value.find(a => a.id === id);
    if (article) receipt.value.push(article);
  }
}

// QFab quantity controls
function increment(id: number) {
  if (qtyOverrides.value[id] !== undefined) {
    qtyOverrides.value = { ...qtyOverrides.value, [id]: (qtyOverrides.value[id] || 0) + 1 };
  } else {
    const article = articles.value.find(a => a.id === id);
    if (article) receipt.value.push(article);
  }
  // keep FAB open
  void nextTick(() => { openFabId.value = id; });
}
function decrement(id: number) {
  if (qtyOverrides.value[id] !== undefined) {
    const next = Math.max(0, (qtyOverrides.value[id] || 0) - 1);
    if (next <= 0) {
      deleteLine(id);
    } else {
      qtyOverrides.value = { ...qtyOverrides.value, [id]: next };
      void nextTick(() => { openFabId.value = id; });
    }
  } else {
    const index = receipt.value.findIndex(a => a.id === id);
    if (index !== -1) receipt.value.splice(index, 1);
    // if still has items for this id, keep FAB open; else close
    const stillExists = receipt.value.some(a => a.id === id);
    void nextTick(() => { openFabId.value = stillExists ? id : null; });
  }
}
function deleteLine(id: number) {
  receipt.value = receipt.value.filter(a => a.id !== id);
  if (qtyOverrides.value[id] !== undefined) {
  const rest = { ...qtyOverrides.value };
  delete rest[id];
  qtyOverrides.value = rest;
  }
  if (reductionsById.value[id] !== undefined) {
    const restR = { ...reductionsById.value };
    delete restR[id];
    reductionsById.value = restR;
  }
  if (notesById.value[id] !== undefined) {
    const restN = { ...notesById.value };
    delete restN[id];
    notesById.value = restN;
  }
  if (openFabId.value === id) openFabId.value = null;
}

// --- Receipt actions (checkout, reset, hold, new order) ---
const $q = useQuasar();
// Show a limited number of category tabs based on breakpoint
const maxTabs = computed(() => {
  if ($q.screen.lt.sm) return 3; // xs
  if ($q.screen.lt.md) return 4; // sm
  if ($q.screen.lt.lg) return 6; // md
  return 8; // lg+
});
const visibleCategories = computed(() => filteredCategories.value.slice(0, maxTabs.value));
const overflowCategories = computed(() => filteredCategories.value.slice(maxTabs.value));
// Ensure selected category is part of filtered list; if not, switch to first visible
watch([filteredCategories, maxTabs], () => {
  const list = filteredCategories.value;
  const current = selectedCategory.value;
  if (!list.includes(current)) {
    const next = (list[0] ?? categories[0]);
    selectedCategory.value = next;
  }
});
// Simple ISO datetime formatter for display
function formatDate(iso: string) {
  try {
    const d = new Date(iso);
    return new Intl.DateTimeFormat(undefined, {
      year: '2-digit', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit'
    }).format(d);
  } catch {
    return iso;
  }
}
const holds = ref<Array<{ id: number; lines: GroupedLine[]; total: number; createdAt: string }>>([]);
const holdListDialog = ref(false);
const orderCounter = ref(1);

function resetReceipt() {
  receipt.value = [];
  qtyOverrides.value = {};
  reductionsById.value = {};
  notesById.value = {};
  totalReduction.value = null;
}

function putOnHold() {
  if (groupedReceipt.value.length === 0) return;
  const snapshot: GroupedLine[] = groupedReceipt.value.map(l => ({ ...l }));
  holds.value.push({
    id: Date.now(),
    lines: snapshot,
    total: total.value,
  createdAt: new Date().toISOString(),
  // @ts-expect-error: extend object at runtime to include notes (backward compatible)
  notes: { ...notesById.value },
  reductions: { ...reductionsById.value },
  totalReduction: totalReduction.value ?? null
  });
  resetReceipt();
  $q.notify({ type: 'warning', message: t('holdPlaced') });
}

function launchNewOrder() {
  orderCounter.value += 1;
  resetReceipt();
  $q.notify({ type: 'info', message: t('newOrderStarted') });
}

function checkout() {
  if (groupedReceipt.value.length === 0) return;
  // Build final receipt JSON and log it
  const normalizeReduction = (r: number | Reduction | null | undefined): Reduction | null => {
    if (r == null) return null;
    return typeof r === 'number' ? { kind: 'percent', value: r } : { kind: r.kind, value: r.value };
  };
  const lines = groupedReceipt.value.map(l => {
    const r = reductionsById.value[l.id];
    return {
      id: l.id,
      name: l.name,
      unitPrice: l.price,
      qty: l.qty,
      lineSubtotal: l.price * l.qty,
      reduction: normalizeReduction(r),
      lineTotal: lineTotal(l),
      note: notesById.value[l.id] ?? undefined
    };
  });
  const receiptJson = {
    createdAt: new Date().toISOString(),
    currency: 'MAD',
    subtotal: subtotal.value,
    totalReduction: normalizeReduction(totalReduction.value),
    total: total.value,
    lines
  };
  console.log('Final receipt:', JSON.stringify(receiptJson, null, 2));
  $q.notify({ type: 'positive', message: t('checkoutComplete') });
  resetReceipt();
}

// Holds: open list, resume, remove, persist
function openHoldList() {
  // prevent any overlay (drag zone, FAB, search) from sitting above dialog
  draggingId.value = null;
  openFabId.value = null;
  searchOpen.value = false;
  holdListDialog.value = true;
}

function openReceiptDialog() {
  // Close overlays before opening the mobile receipt
  draggingId.value = null;
  openFabId.value = null;
  searchOpen.value = false;
  receiptDialog.value = true;
}

function openHoldFromMobile() {
  receiptDialog.value = false;
  // Small delay to allow dialog to close transition before opening next
  setTimeout(() => {
    openHoldList();
  }, 200);
}

// Logout is handled globally in MainLayout via the footer button & login dialog

function removeHold(id: number) {
  holds.value = holds.value.filter(h => h.id !== id);
}

function resumeHold(id: number) {
  const h = holds.value.find(x => x.id === id);
  if (!h) return;
  const doLoad = () => {
    // reconstruct receipt from grouped lines
    receipt.value = [];
    qtyOverrides.value = {};
    notesById.value = {};
    totalReduction.value = null;
    for (const line of h.lines) {
      const base = articles.value.find(a => a.id === line.id) || {
        id: line.id,
        name: line.name,
        price: line.price,
        image: '/images/1.png',
        category: ''
      } as Article;
      // ensure a visible line; use override for exact quantity
      receipt.value.push(base);
      qtyOverrides.value[line.id] = line.qty;
    }
    // restore notes if present
    // @ts-expect-error: historical holds may not have notes
    if (h.notes && typeof h.notes === 'object') {
      // @ts-expect-error - runtime assign
      notesById.value = { ...h.notes };
    }
    // restore reductions if present
    // @ts-expect-error: historical holds may not have reductions
    if (h.reductions && typeof h.reductions === 'object') {
      // @ts-expect-error - runtime assign
      reductionsById.value = { ...h.reductions };
    }
    // restore total reduction if present
    // @ts-expect-error: historical holds may not have totalReduction
    if (h.totalReduction !== undefined) {
      // @ts-expect-error - runtime assign
      totalReduction.value = h.totalReduction;
    }
    // remove the hold after loading
    removeHold(id);
    holdListDialog.value = false;
    $q.notify({ type: 'info', message: t('resumedFromHold') });
  };
  if (groupedReceipt.value.length > 0) {
    $q.dialog({
      title: t('holdList'),
      message: t('replaceCurrentOrderConfirm'),
      cancel: true,
      persistent: true
    }).onOk(doLoad);
  } else {
    doLoad();
  }
}

// persist holds in localStorage (simple persistence)
const HOLDS_KEY = 'pos_holds_v1';
onMounted(() => {
  try {
    const raw = localStorage.getItem(HOLDS_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      if (Array.isArray(parsed)) holds.value = parsed;
    }
  } catch { /* ignore */ }
});
watch(holds, (val) => {
  try { localStorage.setItem(HOLDS_KEY, JSON.stringify(val)); } catch { /* ignore */ }
}, { deep: true });

// Total reduction dialog state and actions
const totalReductionDialog = ref(false);
const totalReductionKind = ref<'percent' | 'amount'>('percent');
const totalReductionValueStr = ref<string>('');
// Scanner dialog state
const scannerDialog = ref(false);
function openScanner() { scannerDialog.value = true; }
function onScannerDecoded(code: string) {
  void handleDecodedBarcode(code);
}

// Classic hardware (keyboard wedge) barcode listener
const wedgeBuffer = ref('');
let wedgeLast = 0;
const WEDGE_TIMEOUT_MS = 90;
function isLikelyBarcodeSequence(delta: number) { return delta < WEDGE_TIMEOUT_MS; }
function flushWedge() {
  const raw = wedgeBuffer.value.trim();
  wedgeBuffer.value = '';
  if (!raw) return;
  const candidates = generateBarcodeCandidates(raw);
  for (const c of candidates) {
    if (handleDecodedBarcode(c)) return; // stop on first success
  }
  // If none succeeded, handleDecodedBarcode already produced a notification for first candidate
}
function onGlobalWedgeKey(e: KeyboardEvent) {
  const target = e.target as HTMLElement | null;
  if (target && target.closest('input,textarea,[contenteditable=true],.q-field')) return;
  const now = performance.now();
  const delta = now - wedgeLast;
  wedgeLast = now;
  if (e.key === 'Enter') {
    flushWedge();
    return;
  }
  if (e.key.length === 1 && /[0-9A-Za-z]/.test(e.key)) {
    if (!isLikelyBarcodeSequence(delta)) {
      // New sequence
      wedgeBuffer.value = '';
    }
    wedgeBuffer.value += e.key;
    // Heuristic: if length reaches 13 and next likely terminator imminent, we still wait for Enter
    if (wedgeBuffer.value.length >= 16) {
      // Fallback flush for long raw codes without Enter (some scanners can be configured w/out terminator)
      flushWedge();
    }
  }
}
onMounted(() => window.addEventListener('keydown', onGlobalWedgeKey));
onBeforeUnmount(() => window.removeEventListener('keydown', onGlobalWedgeKey));
function openTotalReduction() {
  const r = totalReduction.value;
  if (r == null) {
    totalReductionKind.value = 'percent';
    totalReductionValueStr.value = '';
  } else if (typeof r === 'number') {
    totalReductionKind.value = 'percent';
    totalReductionValueStr.value = String(r);
  } else {
    totalReductionKind.value = r.kind;
    totalReductionValueStr.value = String(r.value ?? '');
  }
  totalReductionDialog.value = true;
}
function applyTotalReduction() {
  const n = parseFloat(String(totalReductionValueStr.value).replace(',', '.'));
  if (totalReductionKind.value === 'percent') {
    const pct = Number.isFinite(n) ? Math.max(0, Math.min(100, n)) : 0;
    totalReduction.value = pct > 0 ? pct : null;
  } else {
    const amt = Number.isFinite(n) ? Math.max(0, n) : 0;
    totalReduction.value = amt > 0 ? { kind: 'amount', value: amt } : null;
  }
  totalReductionDialog.value = false;
}
function deleteTotalReduction() {
  totalReduction.value = null;
  totalReductionDialog.value = false;
}

</script>


<style scoped>
/* Full-width mobile receipt bar (fixed above footer) */
.mobile-receipt-bar {
  position: fixed;
  left: 12px;
  right: 12px;
  bottom: calc(64px + env(safe-area-inset-bottom, 0px)); /* footer (~56px) + gap */
  z-index: 1200; /* above page content, below dialogs */
  border-radius: 999px;
  box-shadow: 0 6px 16px rgba(0,0,0,0.18);
}
/* Slightly larger touch target on very small screens */
@media (max-width: 599px) {
  .mobile-receipt-bar :deep(.q-btn__content) { padding: 0 8px; }
}
.receipt-card {
  display: flex;
  flex-direction: column;
  
  /* keep the card from exceeding the products section height */
  max-height: 82vh;
  /* allow FAB actions to extend outside card */
  overflow: visible;
}
.receipt-pane { position: relative; z-index: 100; }
.products-pane { position: relative; z-index: 0; }
.receipt-list :deep(.q-item) { overflow: visible; }
.receipt-list :deep(.q-list) { overflow: visible; }
.receipt-list :deep(.q-item__section--side),
.receipt-list :deep(.q-item__section--avatar),
.receipt-list :deep(.q-item__section--top),
.receipt-list :deep(.q-item__section--bottom) {
  overflow: visible;
}
.receipt-list :deep(.q-fab) { z-index: 300; }
.receipt-list :deep(.q-fab__actions) { z-index: 310; }
.receipt-body {
  /* let the middle list area grow and scroll */
  overflow-y: auto;
  max-height: 100%;
  position: relative;
  z-index: 200;
}
.products-scroll { position: relative; z-index: 0; }
/* Reserve space for vertical scrollbar so 6th column doesn't get squeezed */
.products-scroll { scrollbar-gutter: stable both-edges; }
.products-scroll, .receipt-body { contain: content; }
.search-menu-content {
  z-index: 250; /* above list, below dialogs */
}
.trash-drop-zone {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%) scale(.85);
  width: 140px;
  height: 140px;
  border-radius: 50%;
  background: #fff;
  border: 4px solid var(--q-grey-4);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 24px rgba(0,0,0,0.18);
  transition: border-color .25s, box-shadow .25s, transform .25s, background .25s, opacity .25s;
  opacity: 0.2;
  position: fixed;
  overflow: visible;
}
.trash-drop-zone::before,
.trash-drop-zone::after {
  content: "";
  position: absolute;
  top: 0; left: 0;
  width: 100%; height: 100%;
  border-radius: 50%;
  box-shadow: 0 0 0 0 rgba(244,67,54,0.4);
  opacity: 0;
}
.trash-drop-zone.trash-active {
  border-color: var(--q-negative);
  opacity: 1;
  transform: translate(-50%, -50%) scale(1);
  animation: trashScale 0.35s ease-out;
}
.trash-drop-zone.trash-active::before {
  animation: trashPulseRing 1.4s ease-out infinite;
}
.trash-drop-zone.trash-active::after {
  animation: trashPulseRing 1.4s ease-out .7s infinite;
}
@keyframes trashScale {
  0% { transform: translate(-50%, -50%) scale(.6); }
  100% { transform: translate(-50%, -50%) scale(1); }
}
@keyframes trashPulseRing {
  0% { box-shadow: 0 0 0 0 rgba(244,67,54,0.55); opacity: .9; }
  70% { box-shadow: 0 0 0 28px rgba(244,67,54,0); opacity: 0; }
  100% { box-shadow: 0 0 0 0 rgba(244,67,54,0); opacity: 0; }
}
.fly-clone {
  box-shadow: 0 4px 12px rgba(0,0,0,.25);
  transition: transform .55s cubic-bezier(.55,.06,.27,.99), opacity .55s ease;
  will-change: transform, opacity;
  pointer-events: none;
  border-radius: 8px;
  overflow: hidden;
}
.product-card { width: 96%; margin: 0 auto; }
.receipt-added { animation: receiptAddedFlash .8s ease-out; }
@keyframes receiptAddedFlash {
  0% { background: rgba(33,150,243,0.45); }
  60% { background: rgba(33,150,243,0); }
  100% { background: transparent; }
}
.product-legend {
  height: 30%;
  padding: 0;
  font-size: 0.8rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2px;
  white-space: normal;
  word-break: break-word;
  overflow-wrap: anywhere;
  line-height: 1.15;
}

/* Floating search glow */
.search-glow :deep(.q-field__control) {
  box-shadow: 0 0 0 2px rgba(33,150,243,0.15), 0 8px 24px rgba(33,150,243,0.18);
  border-radius: 12px;
  transition: box-shadow .2s ease, transform .12s ease;
}
.search-glow :deep(.q-field__control):hover {
  box-shadow: 0 0 0 2px rgba(33,150,243,0.25), 0 10px 28px rgba(33,150,243,0.24);
}
.search-glow :deep(.q-field__native) {
  font-weight: 600;
}
.search-glow :deep(.q-field__control--focused) {
  box-shadow: 0 0 0 2px rgba(25,118,210,0.35), 0 12px 30px rgba(25,118,210,0.35);
}
.product-card { width: 96%; margin: 0 auto; }
.product-legend .legend-name {
  /* allow up to 2 lines for name */
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.product-legend .legend-price {
  font-weight: 700;
}
.product-legend {
  height: 30%;
  padding: 0;
  font-size: 0.8rem;
  white-space: normal; /* allow wrapping */
  word-break: break-word;
  overflow-wrap: anywhere;
  line-height: 1.1;
}

/* Modern total summary styles */
.total-summary {
  background: var(--q-grey-1);
  border: 1px solid var(--q-grey-4);
  border-radius: 12px;
  padding: 10px 12px;
  box-shadow: 0 6px 18px rgba(0,0,0,0.06);
}
.total-chip {
  padding: 6px 10px;
  border-radius: 999px;
  background: color-mix(in oklab, var(--q-primary) 12%, white);
  color: var(--q-primary);
  font-weight: 600;
}
.total-amount .amount {
  font-size: 1.6rem;
  font-weight: 800;
  line-height: 1;
  letter-spacing: 0.2px;
}
.total-amount .currency {
  font-size: 0.9rem;
  opacity: 0.8;
  font-weight: 600;
}

/* Dark mode tweaks */
:deep(.body--dark) .total-summary {
  background: rgba(255,255,255,0.04);
  border-color: rgba(255,255,255,0.12);
  box-shadow: 0 6px 18px rgba(0,0,0,0.3);
}
:deep(.body--dark) .total-chip {
  background: color-mix(in oklab, var(--q-primary) 24%, transparent);
  color: #eaeaea;
}

/* Numpad styles */
.numpad {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}
.nkey {
  min-height: 40px;
  font-weight: 600;
}

/* Responsive adjustments */
/* 5 columns on md+ */
@media (min-width: 1024px) {
  .five-col-md { flex: 0 0 20%; max-width: 20%; }
}
/* 6 columns on xl+ (approx >= 1920px) */
@media (min-width: 1920px) {
  .six-col-xl { flex: 0 0 16.6667%; max-width: 16.6667%; }
}
/* Slightly tighter product card for very large screens */
@media (min-width: 1920px) {
  .product-card { width: 94%; }
}
@media (max-width: 1023px) { /* below md */
  .products-scroll {
    max-height: none !important; /* let page scroll when stacked */
    row-gap: 8px !important; /* reduced vertical spacing */
  }
  /* add a bit more vertical breathing room on small screens */
  .product-tile { margin-bottom: 8px !important; }
  .receipt-card {
    max-height: none !important;
  }
  .product-card {
    height: 110px !important;
  }
  .product-img {
    height: 109px !important;
  }
  .total-summary {
    padding: 8px 10px;
  }
  .total-amount .amount {
    font-size: 1.4rem;
  }
}
@media (max-width: 599px) { /* xs */
  /* leave space for the fixed receipt bar */
  .products-scroll { padding-bottom: 120px; row-gap: 10px !important; }
  .product-tile { margin-bottom: 10px !important; }
  .product-card {
    height: 96px !important;
  }
  .product-img {
    height: 95px !important;
  }
  .product-legend { font-size: 0.72rem; }
  .total-amount .amount {
    font-size: 1.25rem;
  }
  .nkey { min-height: 36px; }
}

/* Improve touch responsiveness */
:deep(.q-btn), :deep(.q-tab), :deep(.q-fab), :deep(.q-item), .product-card {
  touch-action: manipulation;
}

</style>
