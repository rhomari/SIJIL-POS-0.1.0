<template>
  <q-page class="q-pa-md">
  <div class="row q-col-gutter-none">
      <!-- Articles Section -->
  <div class="col-12 col-md-7 col-xl-8 products-pane">
        <div class="row items-center no-wrap q-gutter-sm">
          <div class="col">
            <q-tabs v-model="selectedCategory" class="text-primary" :class="{ 'is-dragging': isDragging }" align="left" dense shrink>
              <q-tab :name="MOST_USED" :label="$t('pos.mostUsed')" />
              <q-tab 
                v-for="cat in visibleCategories" 
                :key="cat" 
                :name="cat" 
                :label="cat" 
                :class="{
                  'drag-over': dragOverCategory === cat,
                  'dragging': draggedCategory === cat,
                  'drop-cursor-left': dropCursorPosition?.category === cat && dropCursorPosition?.side === 'left',
                  'drop-cursor-right': dropCursorPosition?.category === cat && dropCursorPosition?.side === 'right'
                }"
                draggable="true"
                @dragstart="onCategoryDragStart(cat, $event)"
                @dragend="onCategoryDragEnd"
                @dragover="onCategoryDragOver(cat, $event)"
                @dragleave="onCategoryDragLeave"
                @drop="onCategoryDrop(cat, $event)"
              />
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
                  debounce="300"
                  :placeholder="$t('pos.searchCategoriesOrProducts')"
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
          <q-btn-dropdown v-if="overflowCategories.length" dense flat icon="more_horiz" :label="$t('pos.more')" class="text-primary" content-class="text-primary">
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
          :label="$t('receipt') + (totalItemsCount ? ' (' + totalItemsCount + ')' : '')"
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
              <q-btn dense flat icon="person" size="sm" :label="selectedCustomer ? selectedCustomer.name : $t('customer')" @click="openCustomerDialog" :color="selectedCustomer ? 'primary' : 'grey-7'" class="q-mr-sm" />
              <q-btn dense flat icon="playlist_add_check" size="sm" :label="$t('holdList') + ' (' + holds.length + ')'" @click="openHoldList" />
            </div>
          </q-card-section>
          <q-separator />
          <!-- Mobile Delete Hint -->
          <div v-if="isTouchDevice && orderedLines.length > 0" class="mobile-hint q-pa-sm bg-blue-1 text-blue-8 text-center text-caption">
            <q-icon name="swipe" size="xs" class="q-mr-xs" />
            {{ t('pos.swipeLeftToDelete') }}
          </div>
          <div class="relative-position receipt-body">
            <q-list ref="receiptListRef" class="receipt-list">
              <q-item v-for="item in orderedLines" :key="item.id" :data-line-id="item.id"
                draggable="true"
                @dragstart="onDragStart(item.id)"
                @dragend="onDragEnd"
                @dragenter="onReorderDragEnter(item.id)"
                @dragover.prevent
                @drop.prevent="onItemDrop(item.id, $event)"
                @touchstart="onTouchStart(item.id, $event)"
                @touchmove="onTouchMove($event)"
                @touchend="onTouchEnd"
                :class="{ 
                  'dragging': draggingId === item.id || (touchDragState.isDragging && touchDragState.itemId === item.id), 
                  'receipt-added': lastAddedId === item.id, 
                  'reorder-over': reorderOverId === item.id, 
                  'reorder-dragging': reorderDragId === item.id,
                  'touch-dragging': touchDragState.isDragging && touchDragState.itemId === item.id
                }"
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
                <q-item-section side class="q-pl-xs">
                  <q-icon name="drag_indicator" class="reorder-handle" draggable="true"
                    @dragstart.stop="onReorderHandleStart(item.id, $event)" />
                </q-item-section>
              </q-item>
            </q-list>
          </div>
          
          <!-- Total Reduction Dialog -->
          <q-dialog v-model="totalReductionDialog">
            <q-card style="width:420px; max-width:95vw;" class="q-pa-sm">
              <q-card-section class="row items-center q-pb-none">
                <div class="text-h6">{{ $t('pos.totalDiscount') }}</div>
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
                  :placeholder="totalReductionKind === 'percent' ? $t('pos.enterDiscountPercent') : $t('pos.enterDiscountAmount')"
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
                <div class="text-h6">{{ $t('pos.itemDiscount') }}</div>
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
                  :placeholder="reductionKind === 'percent' ? $t('pos.enterDiscountPercent') : $t('pos.enterDiscountAmount')"
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
                <div class="text-h6">{{ $t('editQuantity') }}</div>
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
                <div class="text-h6">{{ $t('pos.itemNote') }}</div>
                <q-space />
                <q-btn icon="close" flat round dense v-close-popup />
              </q-card-section>
              <q-card-section class="q-pt-sm">
                <div class="text-subtitle2 text-primary">{{ noteLineName }}</div>
                <q-input v-model="noteText" type="textarea" autogrow outlined dense :placeholder="$t('pos.addNotePlaceholder')" />
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
              <q-chip v-if="selectedCustomer" dense color="primary" text-color="white" class="q-ml-sm" removable @remove="clearSelectedCustomer">
                <q-icon name="person" size="16px" class="q-mr-xs" /> {{ selectedCustomer.name }}
              </q-chip>
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
  <!-- Mobile Receipt Dialog (bottom sheet) - FIRST -->
  <q-dialog v-model="receiptDialog" position="bottom" transition-show="slide-up" transition-hide="slide-down">
    <q-card style="width: 100vw; max-width: 100vw; max-height: 85vh;" class="q-pa-none">
        <q-card-section class="bg-grey-2">
          <div class="row items-center no-wrap">
            <div class="text-h6">{{ $t('receipt') }}</div>
            <q-space />
            <q-btn dense flat icon="person" size="sm" :label="selectedCustomer ? selectedCustomer.name : $t('customer')" @click="openCustomerDialog" :color="selectedCustomer ? 'primary' : 'grey-7'" class="q-mr-sm" />
            <q-btn dense flat icon="playlist_add_check" size="sm" :label="$t('holdList') + ' (' + holds.length + ')'" @click="openHoldFromMobile" />
            <q-btn flat round dense icon="close" v-close-popup class="q-ml-sm" />
          </div>
        </q-card-section>
        <q-separator />
        <div class="relative-position" style="overflow-y: auto; max-height: calc(85vh - 160px);">
          <!-- Mobile Delete Hint -->
          <div v-if="orderedLines.length > 0" class="mobile-hint q-pa-sm bg-blue-1 text-blue-8 text-center text-caption">
            <q-icon name="swipe" size="xs" class="q-mr-xs" />
            {{ t('pos.swipeLeftToDelete') }}
          </div>
          <q-list class="receipt-list q-px-sm q-pt-sm">
            <q-item v-for="item in orderedLines" :key="item.id" :data-line-id="item.id"
              @dragenter="onReorderDragEnter(item.id)" @dragover.prevent @drop.stop.prevent="onReorderDrop(item.id)"
              @touchstart="onTouchStart(item.id, $event)"
              @touchmove="onTouchMove($event)"
              @touchend="onTouchEnd"
              :class="{ 
                'reorder-over': reorderOverId === item.id, 
                'reorder-dragging': reorderDragId === item.id,
                'touch-dragging': touchDragState.isDragging && touchDragState.itemId === item.id
              }"
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
              <q-item-section side class="q-pl-xs">
                <q-icon name="drag_indicator" class="reorder-handle" draggable="true"
                  @dragstart.stop="onReorderHandleStart(item.id, $event)" />
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
            <q-chip v-if="selectedCustomer" dense color="primary" text-color="white" class="q-ml-sm" removable @remove="clearSelectedCustomer">
              <q-icon name="person" size="16px" class="q-mr-xs" /> {{ selectedCustomer.name }}
            </q-chip>
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
  <!-- Hold List Dialog -->
  <q-dialog v-model="holdListDialog">
    <q-card style="width:520px; max-width:95vw;" class="q-pa-sm">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">{{ $t('holdList') }}</div>
        <q-space />
        <q-btn flat round dense icon="close" v-close-popup />
      </q-card-section>
      <q-card-section class="q-pt-sm">
        <div v-if="holds.length === 0" class="text-grey-7">{{ $t('noHolds') }}</div>
        <q-list v-else bordered separator style="max-height:300px; overflow-y:auto;">
          <q-item v-for="h in holds" :key="h.id" clickable @click="resumeHold(h.id)">
            <q-item-section>
              <div class="row items-center justify-between">
                <span class="text-subtitle2">#{{ h.id }}</span>
                <span class="text-caption">{{ h.total.toFixed(2) }} MAD</span>
              </div>
              <div class="row items-center justify-between text-caption q-mt-xs">
                <span>{{ h.lines.length }} {{ $t('items') }}</span>
                <span>{{ formatDate(h.createdAt) }}</span>
              </div>
            </q-item-section>
            <q-item-section side>
              <q-btn dense flat color="negative" icon="delete" @click.stop="removeHold(h.id)" />
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
      <q-separator inset />
      <q-card-actions align="right">
        <q-btn flat :label="$t('close')" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>
  <!-- Mobile Receipt Dialog (bottom sheet) - SECOND -->
  <q-dialog v-model="receiptDialog" position="bottom" transition-show="slide-up" transition-hide="slide-down">
    <q-card style="width: 100vw; max-width: 100vw; max-height: 85vh;" class="q-pa-none">
        <q-card-section class="bg-grey-2">
          <div class="row items-center no-wrap">
            <div class="text-h6">{{ $t('receipt') }}</div>
            <q-space />
            <q-btn dense flat icon="person" size="sm" :label="selectedCustomer ? selectedCustomer.name : $t('customer')" @click="openCustomerDialog" :color="selectedCustomer ? 'primary' : 'grey-7'" class="q-mr-sm" />
            <q-btn dense flat icon="playlist_add_check" size="sm" :label="$t('holdList') + ' (' + holds.length + ')'" @click="openHoldFromMobile" />
            <q-btn flat round dense icon="close" v-close-popup class="q-ml-sm" />
          </div>
        </q-card-section>
        <q-separator />
        <div class="relative-position" style="overflow-y: auto; max-height: calc(85vh - 160px);">
          <!-- Mobile Delete Hint -->
          <div v-if="orderedLines.length > 0" class="mobile-hint q-pa-sm bg-blue-1 text-blue-8 text-center text-caption">
            <q-icon name="swipe" size="xs" class="q-mr-xs" />
            {{ t('pos.swipeLeftToDelete') }}
          </div>
          <q-list class="receipt-list q-px-sm q-pt-sm">
            <q-item v-for="item in orderedLines" :key="item.id" :data-line-id="item.id"
              @touchstart="onTouchStart(item.id, $event)"
              @touchmove="onTouchMove($event)"
              @touchend="onTouchEnd"
              @dragenter="onReorderDragEnter(item.id)" @dragover.prevent @drop.stop.prevent="onReorderDrop(item.id)"
              :class="{ 
                'reorder-over': reorderOverId === item.id, 
                'reorder-dragging': reorderDragId === item.id,
                'touch-dragging': touchDragState.isDragging && touchDragState.itemId === item.id
              }"
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
              <q-item-section side class="q-pl-xs">
                <q-icon name="drag_indicator" class="reorder-handle" draggable="true"
                  @dragstart.stop="onReorderHandleStart(item.id, $event)" />
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
            <q-chip v-if="selectedCustomer" dense color="primary" text-color="white" class="q-ml-sm" removable @remove="clearSelectedCustomer">
              <q-icon name="person" size="16px" class="q-mr-xs" /> {{ selectedCustomer.name }}
            </q-chip>
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
  <!-- Customer Selection Dialog -->
  <q-dialog v-model="customerDialog">
      <q-card style="width:520px; max-width:95vw;" class="q-pa-sm">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{ $t('customer') }}</div>
          <q-space />
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>
        <q-card-section class="q-pt-sm">
          <q-input v-model="customerSearch" dense outlined clearable debounce="150" :placeholder="$t('search') + '...'" />
          <div class="row q-mt-sm q-col-gutter-sm">
            <div class="col">
              <q-input v-model="newCustomerName" dense outlined :placeholder="$t('name')" />
            </div>
            <div class="col">
              <q-input v-model="newCustomerPhone" dense outlined :placeholder="$t('phone')" />
            </div>
            <div class="col-auto flex items-center">
              <q-btn dense color="primary" icon="person_add" :disable="!canAddCustomer" @click="addCustomer" />
            </div>
          </div>
          <q-list bordered separator class="q-mt-md" style="max-height:260px; overflow-y:auto;">
            <q-item v-for="c in filteredCustomers" :key="c.id" clickable @click="selectCustomer(c.id)" :active="c.id === selectedCustomerId" active-class="bg-primary text-white">
              <q-item-section>
                <div class="row items-center justify-between">
                  <span class="text-subtitle2">{{ c.name }}</span>
                  <span class="text-caption">{{ c.phone }}</span>
                </div>
              </q-item-section>
              <q-item-section side v-if="c.id === selectedCustomerId">
                <q-icon name="check" />
              </q-item-section>
            </q-item>
            <div v-if="filteredCustomers.length === 0" class="text-caption text-grey-7 q-pa-sm">{{ $t('noResults') }}</div>
          </q-list>
        </q-card-section>
        <q-separator inset />
        <q-card-actions align="right">
          <q-btn flat :label="$t('clear')" @click="clearSelectedCustomer" :disable="!selectedCustomer" />
          <q-btn flat :label="$t('close')" v-close-popup />
        </q-card-actions>
      </q-card>
  </q-dialog>
  <!-- Barcode / QR Scanner Dialog -->
  <q-dialog v-model="scannerDialog">
    <q-card style="width:640px; max-width:95vw;" class="q-pa-sm">
      <q-card-section class="row items-center q-pb-none">
  <div class="text-h6">{{ $t('pos.scanBarcode') }}</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>
      <q-separator />
      <q-card-section class="q-pt-sm">
        <div class="scanner-frame">
          <CameraScanner @decoded="onScannerDecoded" />
        </div>
      </q-card-section>
      <q-separator inset />
      <q-card-actions align="right">
        <q-btn flat :label="$t('close')" v-close-popup />
      </q-card-actions>
    </q-card>
  </q-dialog>

  <!-- Desktop Drag-to-Delete Trash Overlay -->
  <Teleport to="body" v-if="showTrashOverlay">
    <div class="trash-overlay">
      <div 
        class="trash-drop-zone"
        :class="{ 'trash-active': draggingId !== null, 'trash-over': overTrash }"
        @dragenter.prevent="onTrashDragEnter"
        @dragover.prevent
        @dragleave="onTrashDragLeave"
        @drop.prevent="onDropTrash"
      >
        <q-icon name="delete" size="32px" color="negative" />
        <div class="trash-instruction">
          <div class="text-negative text-weight-bold">{{ $t('pos.dragHereToDelete') }}</div>
        </div>
      </div>
    </div>
  </Teleport>

  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch, type ComponentPublicInstance } from 'vue';
import { useQuasar } from 'quasar';
import type { QList } from 'quasar';
import { useI18n } from 'vue-i18n';
import CameraScanner from '../components/CameraScanner.vue';

// Quasar & i18n
const $q = useQuasar();
const { t } = useI18n();

// Reactive viewport width for physical monitor size-aware computations
const viewportWidth = ref(typeof window !== 'undefined' ? window.innerWidth : 1200);

// Debounced resize handler for better INP performance
let resizeTimeout: number | null = null;
function handleViewportResize() { 
  if (resizeTimeout) clearTimeout(resizeTimeout);
  resizeTimeout = window.setTimeout(() => {
    viewportWidth.value = window.innerWidth;
    resizeTimeout = null;
  }, 150); // Debounce resize events
}

if (typeof window !== 'undefined') {
  window.addEventListener('resize', handleViewportResize, { passive: true });
  onBeforeUnmount(() => {
    window.removeEventListener('resize', handleViewportResize);
    if (resizeTimeout) clearTimeout(resizeTimeout);
  });
}

// Drag delete state
const draggingId = ref<number|null>(null);
function onDragStart(id: number) { draggingId.value = id; }
function onDragEnd() { 
  draggingId.value = null; 
  overTrash.value = false; 
}

// Desktop drag-to-delete trash overlay state
const overTrash = ref(false);
const showTrashOverlay = computed(() => draggingId.value !== null);

// Desktop trash zone drag handlers
function onTrashDragEnter(event: DragEvent) {
  event.preventDefault();
  overTrash.value = true;
}

function onTrashDragLeave() {
  overTrash.value = false;
}

function onDropTrash(event: DragEvent) {
  event.preventDefault();
  const draggedItemId = draggingId.value;
  
  // Reset drag state
  overTrash.value = false;
  draggingId.value = null;
  
  // Delete the item if there was one being dragged
  if (draggedItemId !== null) {
    deleteLine(draggedItemId);
  }
}

// Touch-based drag state for mobile devices
const touchDragState = ref<{
  itemId: number | null;
  startX: number;
  startY: number;
  currentX: number;
  currentY: number;
  isDragging: boolean;
  swipeProgress: number;
}>({
  itemId: null,
  startX: 0,
  startY: 0,
  currentX: 0,
  currentY: 0,
  isDragging: false,
  swipeProgress: 0
});

// Touch event handlers for mobile drag-to-delete
function onTouchStart(id: number, event: TouchEvent) {
  
  if (event.touches.length === 0) {
    return;
  }
  
  const firstTouch = event.touches.item(0);
  if (!firstTouch) {
    return;
  }
  
  // Prepare item for smooth color transitions
  const itemElement = document.querySelector<HTMLElement>(`.q-item[data-line-id="${id}"]`);
  
  if (itemElement) {
    // Add a simple blue border to show touch started
    itemElement.style.border = '2px solid blue';
  } else {
    // Element not found - could try alternative selectors if needed
  }
  
  touchDragState.value = {
    itemId: id,
    startX: firstTouch.clientX,
    startY: firstTouch.clientY,
    currentX: firstTouch.clientX,
    currentY: firstTouch.clientY,
    isDragging: false,
    swipeProgress: 0
  };
  
  // Remove long-press timer - no longer needed
  // longPressTimer = window.setTimeout(() => {
  //   console.log('⏰ Long press timer triggered for item', id);
  //   if (touchDragState.value.itemId === id && !touchDragState.value.isDragging) {
  //     // Long press detected - show delete confirmation
  //     if (confirm('Delete this item from receipt?')) {
  //       console.log('🗑️ Deleting item via long press:', id);
  //       deleteLine(id);
  //     }
  //   }
  // }, 800); // 800ms long press
}

function onTouchMove(event: TouchEvent) {
  
  if (!touchDragState.value.itemId || event.touches.length === 0) {
    return;
  }
  
  const firstTouch = event.touches.item(0);
  if (!firstTouch) {
    return;
  }
  
  // Long press timer handling removed
  
  const deltaX = firstTouch.clientX - touchDragState.value.startX;
  const deltaY = Math.abs(firstTouch.clientY - touchDragState.value.startY);
  const absDeltaX = Math.abs(deltaX);
  
  // Check for horizontal swipe (left to delete)
  if (absDeltaX > 50 && deltaY < 30) {
    event.preventDefault(); // Prevent scrolling
    
    if (!touchDragState.value.isDragging) {
      touchDragState.value.isDragging = true;
    }
    
    // Calculate swipe progress and apply color phases
    const swipeProgress = Math.min(Math.abs(deltaX) / 120, 1);
    touchDragState.value.swipeProgress = swipeProgress;
    
    // Apply color classes based on swipe distance for orange-to-red progression
    const itemElement = document.querySelector<HTMLElement>(`.q-item[data-line-id="${touchDragState.value.itemId}"]`);
    
    if (itemElement && deltaX < 0) { // Only for left swipes
      
      // Reset any previous styles first
      itemElement.style.backgroundColor = '';
      itemElement.style.border = '';
      itemElement.style.boxShadow = '';
      itemElement.style.transform = '';
      
      if (Math.abs(deltaX) >= 40 && Math.abs(deltaX) < 60) {
        // Orange phase: 40-60px - VERY STRONG ORANGE
        itemElement.style.backgroundColor = 'orange !important';
        itemElement.style.border = '5px solid #ff6600';
        itemElement.style.boxShadow = '0 0 20px orange';
      } else if (Math.abs(deltaX) >= 60 && Math.abs(deltaX) < 80) {
        // Red phase: 60-80px - VERY STRONG RED
        itemElement.style.backgroundColor = 'red !important';
        itemElement.style.border = '5px solid #cc0000';
        itemElement.style.boxShadow = '0 0 20px red';
      } else if (Math.abs(deltaX) >= 80) {
        // Delete ready: 80px+ - VERY STRONG DARK RED
        itemElement.style.backgroundColor = 'darkred !important';
        itemElement.style.border = '5px solid #990000';
        itemElement.style.boxShadow = '0 0 20px darkred';
        itemElement.style.transform = 'translateX(-10px)';
      }
    }
    
    touchDragState.value.currentX = firstTouch.clientX;
    touchDragState.value.currentY = firstTouch.clientY;
  }
}

function onTouchEnd() {
  
  if (!touchDragState.value.itemId) {
    return;
  }
  
  // Long press timer handling removed
  
  const itemId = touchDragState.value.itemId;
  const wasDragging = touchDragState.value.isDragging;
  const swipeDistance = touchDragState.value.currentX - touchDragState.value.startX;
  
  // Clean up color classes from the item
  const itemElement = document.querySelector<HTMLElement>(`.q-item[data-line-id="${itemId}"]`);
  if (itemElement) {
    // Only clean up inline styles - no CSS classes needed
    itemElement.style.backgroundColor = '';
    itemElement.style.border = '';
    itemElement.style.boxShadow = '';
    itemElement.style.transform = '';
  }
  
  // Reset touch state
  touchDragState.value = {
    itemId: null,
    startX: 0,
    startY: 0,
    currentX: 0,
    currentY: 0,
    isDragging: false,
    swipeProgress: 0
  };
  draggingId.value = null;
  
  // Delete item if swiped left far enough (more than 80px)
  if (wasDragging && swipeDistance < -80) {
    
    // Add final deletion animation before removing
    if (itemElement) {
      itemElement.style.transition = 'transform 0.3s ease, opacity 0.3s ease';
      itemElement.style.transform = 'translateX(-100%)';
      itemElement.style.opacity = '0';
      
      // Delete after animation
      setTimeout(() => {
        deleteLine(itemId);
      }, 300);
    } else {
      deleteLine(itemId);
    }
  }
}

// import { useRouter } from 'vue-router';

// Special pseudo-category for Most Used products (local usage stats for now; can be replaced with server stats later)
const MOST_USED = '__MOST_USED__' as const;

// Default categories order
const DEFAULT_CATEGORIES = [
  'Stationery', 'Books', 'Office', 'Gadgets', 'Supplies',
  'Arts', 'Crafts', 'Paper', 'Ink', 'Hardware',
  'Software', 'Snacks', 'Beverages', 'Cleaning', 'Storage'
] as const;

// Load categories order from localStorage or use default
const loadCategoriesOrder = (): readonly string[] => {
  try {
    const stored = localStorage.getItem('sijil-pos-categories-order');
    if (stored) {
      const parsed = JSON.parse(stored);
      if (Array.isArray(parsed) && parsed.length === DEFAULT_CATEGORIES.length) {
        // Validate that all default categories are present
        const hasAllCategories = DEFAULT_CATEGORIES.every(cat => parsed.includes(cat));
        if (hasAllCategories) return parsed;
      }
    }
  } catch (e) {
    console.warn('Failed to load categories order:', e);
  }
  return DEFAULT_CATEGORIES;
};

// Save categories order to localStorage
const saveCategoriesOrder = (order: readonly string[]) => {
  try {
    localStorage.setItem('sijil-pos-categories-order', JSON.stringify(order));
  } catch (e) {
    console.warn('Failed to save categories order:', e);
  }
};

// Reactive categories that can be reordered
const categories = ref<readonly string[]>(loadCategoriesOrder());
type Category = string; // Any category name including MOST_USED
const selectedCategory = ref<Category>(categories.value[0] ?? 'Stationery');
const categorySearch = ref<string>('');
const searchInputRef = ref<ComponentPublicInstance | null>(null);
// No explicit state needed for popup proxy

// Drag and drop state for category reordering
const draggedCategory = ref<string | null>(null);
const dragOverCategory = ref<string | null>(null);
const isDragging = ref(false);
const dropCursorPosition = ref<{ category: string; side: 'left' | 'right' } | null>(null);

// Drag and drop handlers
function onCategoryDragStart(category: string, event: DragEvent) {
  if (!event.dataTransfer) return;
  draggedCategory.value = category;
  isDragging.value = true;
  event.dataTransfer.effectAllowed = 'move';
  event.dataTransfer.setData('text/plain', category);
  
  // Add drag styling
  if (event.target instanceof HTMLElement) {
    event.target.classList.add('dragging');
  }
}

function onCategoryDragEnd(event: DragEvent) {
  draggedCategory.value = null;
  dragOverCategory.value = null;
  isDragging.value = false;
  dropCursorPosition.value = null;
  
  // Remove drag styling
  if (event.target instanceof HTMLElement) {
    event.target.classList.remove('dragging');
  }
}

function onCategoryDragOver(category: string, event: DragEvent) {
  if (draggedCategory.value && draggedCategory.value !== category) {
    event.preventDefault();
    dragOverCategory.value = category;
    
    // Calculate drop cursor position based on mouse position relative to tab
    const tabElement = event.currentTarget as HTMLElement;
    const rect = tabElement.getBoundingClientRect();
    const mouseX = event.clientX;
    const tabCenterX = rect.left + rect.width / 2;
    
    // Show cursor on left or right side based on mouse position
    const side = mouseX < tabCenterX ? 'left' : 'right';
    dropCursorPosition.value = { category, side };
    
    if (event.dataTransfer) {
      event.dataTransfer.dropEffect = 'move';
    }
  }
}

function onCategoryDragLeave() {
  dragOverCategory.value = null;
  dropCursorPosition.value = null;
}

function onCategoryDrop(targetCategory: string, event: DragEvent) {
  event.preventDefault();
  
  const sourceCategory = draggedCategory.value;
  if (!sourceCategory || sourceCategory === targetCategory) return;
  
  const currentOrder = [...categories.value];
  const sourceIndex = currentOrder.indexOf(sourceCategory);
  const targetIndex = currentOrder.indexOf(targetCategory);
  
  if (sourceIndex === -1 || targetIndex === -1) return;
  
  // Use drop cursor position for more precise insertion
  const insertPosition = dropCursorPosition.value;
  let newTargetIndex = targetIndex;
  
  if (insertPosition?.side === 'right') {
    // Insert after the target
    newTargetIndex = sourceIndex < targetIndex ? targetIndex : targetIndex + 1;
  } else {
    // Insert before the target  
    newTargetIndex = sourceIndex < targetIndex ? targetIndex - 1 : targetIndex;
  }
  
  // Remove source and insert at calculated position
  currentOrder.splice(sourceIndex, 1);
  currentOrder.splice(newTargetIndex, 0, sourceCategory);
  
  // Update categories and save to localStorage
  categories.value = currentOrder;
  saveCategoriesOrder(currentOrder);
  
  // Clear drag state
  dragOverCategory.value = null;
  dropCursorPosition.value = null;
}
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
  const cat = categories.value[id % categories.value.length] ?? 'Stationery';
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

// --- Usage tracking (local) ---
const usageCounts = ref<Record<number, number>>({});
function recordUsage(article: Article) {
  usageCounts.value = {
    ...usageCounts.value,
    [article.id]: (usageCounts.value[article.id] || 0) + 1
  };
}

const mostUsedArticles = computed<Article[]>(() => {
  const entries = Object.entries(usageCounts.value)
    .filter(([, count]) => count > 0)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 60); // cap to 60 tiles
  if (!entries.length) return [];
  const byId = new Map(articles.value.map(a => [a.id, a] as const));
  return entries
    .map(([id]) => byId.get(Number(id)))
    .filter((a): a is Article => Boolean(a));
});

const filteredArticles = computed(() => {
  const q = categorySearch.value.trim().toLowerCase();
  if (q) {
    return articles.value.filter((a: Article) =>
      a.name.toLowerCase().includes(q) || a.category.toLowerCase().includes(q)
    );
  }
  if (selectedCategory.value === MOST_USED) {
    return mostUsedArticles.value;
  }
  return articles.value.filter((a: Article) => a.category === selectedCategory.value);
});

// Category filtering for tabs/dropdown (define before watchers using it)
const filteredCategories = computed<Category[]>(() => {
  const base = Array.from(categories.value);
  const q = categorySearch.value.trim().toLowerCase();
  let list: string[];
  if (!q) list = base; else {
    const nameMatches = new Set(base.filter(c => c.toLowerCase().includes(q)));
    const productMatches = new Set(
      articles.value
        .filter((a: Article) => a.name.toLowerCase().includes(q) || a.category.toLowerCase().includes(q))
        .map(a => a.category)
    );
    const merged = new Set([...Array.from(nameMatches), ...Array.from(productMatches)]);
    list = base.filter(c => merged.has(c));
  }
  // Always put MOST_USED at the start
  return [MOST_USED, ...list];
});

// Show a limited number of category tabs based on viewport width (optimized for large monitors)
// Note: These numbers are for regular categories only, "Most Used" tab is always shown additionally
const maxTabs = computed(() => {
  const width = viewportWidth.value;
  
  // Use raw viewport width - works better for large desktop monitors
  if (width < 640) return 1;       // Mobile/small tablets: 1 + Most Used = 2 total tabs
  if (width < 768) return 2;       // Large phones/small tablets: 2 + Most Used = 3 total tabs  
  if (width < 1024) return 3;      // Tablets/small laptops: 3 + Most Used = 4 total tabs
  if (width < 1280) return 4;      // Standard laptops: 4 + Most Used = 5 total tabs
  if (width < 1440) return 5;      // Large laptops/small desktops: 5 + Most Used = 6 total tabs
  if (width < 1920) return 6;      // Full HD desktops: 6 + Most Used = 7 total tabs
  if (width < 2560) return 7;      // 24" monitors/QHD: 7 + Most Used = 8 total tabs
  return 8;                        // 27"+ monitors/4K: 8 + Most Used = 9 total tabs
});
// filteredCategories already includes MOST_USED at index 0, but we don't want it duplicated in overflow logic.
const visibleCategories = computed(() => filteredCategories.value.filter(c => c !== MOST_USED).slice(0, maxTabs.value));
const overflowCategories = computed(() => filteredCategories.value.filter(c => c !== MOST_USED).slice(maxTabs.value));
watch([filteredCategories, maxTabs], () => {
  const list = filteredCategories.value;
  if (selectedCategory.value === MOST_USED) return; // always valid
  // Exclude MOST_USED from fallback logic ordering
  const withoutMost = list.filter(c => c !== MOST_USED);
  if (!withoutMost.includes(selectedCategory.value)) {
    selectedCategory.value = withoutMost[0] ?? MOST_USED;
  }
});

const receipt = ref<Article[]>([]);
function addToReceipt(article: Article) { receipt.value.push(article); recordUsage(article); }
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
  $q.notify({ type: 'warning', message: t('pos.barcodeInvalidChecksum', { code: primary }) });
      return false;
    }
  }
  $q.notify({ type: 'warning', message: t('pos.barcodeNotFound', { code: primary }) });
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
    const hasCoarsePointer = matchMedia('(pointer: coarse)').matches;
    const hasTouchStart = 'ontouchstart' in window;
    const hasMaxTouchPoints = navigator.maxTouchPoints > 0;
    const isMobileUserAgent = /Android|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    
    isTouchDevice.value = hasCoarsePointer || hasTouchStart || hasMaxTouchPoints || isMobileUserAgent;
  } catch (e) { 
    void e; // ignore
    isTouchDevice.value = false; 
  }
  try {
    prefersReducedMotion.value = matchMedia('(prefers-reduced-motion: reduce)').matches;
  } catch { prefersReducedMotion.value = false; }
});
function handleAdd(article: Article, ev: MouseEvent) {
  // Instant path for touch devices or when reduced motion is preferred
  if (isTouchDevice.value || prefersReducedMotion.value) {
    addToReceipt(article);
    lastAddedId.value = article.id;
    
    // Defer non-critical operations to improve INP
    if ('requestIdleCallback' in window) {
      requestIdleCallback(() => playAddBeep(), { timeout: 50 });
    } else {
      setTimeout(() => playAddBeep(), 0);
    }
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
      transition: 'transform .65s cubic-bezier(.55,.06,.27,.99), opacity .65s ease',
      willChange: 'transform, opacity',
      backfaceVisibility: 'hidden',
      contain: 'layout style paint'
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
    const cleanup = () => {
      fly.style.willChange = 'auto'; // Clean up will-change
      fly.remove();
    };
    fly.addEventListener('transitionend', cleanup, { once: true });
    setTimeout(cleanup, 900);
    setTimeout(() => { if (lastAddedId.value === article.id) lastAddedId.value = null; }, 1000);
  });
}
// const router = useRouter();
// Mobile receipt dialog state
const receiptDialog = ref(false);
// Optional per-line quantity overrides (allow decimals)
const qtyOverrides = ref<Record<number, number>>({});
// Optional per-line notes
const notesById = ref<Record<number, string>>({});
// Optional per-line reductions: percent or fixed amount (union for backward compatibility)
// Quantity / line helpers
function setLineQuantity(id: number, newQty: number) {
  if (newQty <= 0) {
    receipt.value = receipt.value.filter(a => a.id !== id);
    const rest = { ...qtyOverrides.value }; delete rest[id]; qtyOverrides.value = rest; return;
  }
  qtyOverrides.value = { ...qtyOverrides.value, [id]: newQty };
  if (!receipt.value.some(a => a.id === id)) {
    const article = articles.value.find(a => a.id === id);
    if (article) receipt.value.push(article);
  }
}
const openFabId = ref<number|null>(null);
function toggleFab(id: number, open: boolean) { openFabId.value = open ? id : (openFabId.value === id ? null : openFabId.value); }

// Debounced beep to prevent audio overlap during rapid clicks
let beepTimeout: number | null = null;
function playDebouncedBeep() {
  if (beepTimeout) clearTimeout(beepTimeout);
  beepTimeout = window.setTimeout(() => {
    void playAddBeep();
    beepTimeout = null;
  }, 50);
}

function increment(id: number) {
  if (qtyOverrides.value[id] !== undefined) {
    qtyOverrides.value = { ...qtyOverrides.value, [id]: (qtyOverrides.value[id] || 0) + 1 };
  } else {
    const article = articles.value.find(a => a.id === id); if (article) receipt.value.push(article);
  }
  playDebouncedBeep();
  void nextTick(() => { openFabId.value = id; });
}
function decrement(id: number) {
  if (qtyOverrides.value[id] !== undefined) {
    const next = Math.max(0, (qtyOverrides.value[id] || 0) - 1);
    if (next <= 0) { deleteLine(id); } else { qtyOverrides.value = { ...qtyOverrides.value, [id]: next }; playDebouncedBeep(); void nextTick(() => { openFabId.value = id; }); }
  } else {
    const index = receipt.value.findIndex(a => a.id === id); if (index !== -1) receipt.value.splice(index, 1);
    const stillExists = receipt.value.some(a => a.id === id); if (stillExists) playDebouncedBeep(); void nextTick(() => { openFabId.value = stillExists ? id : null; });
  }
}
function deleteLine(id: number) {
  receipt.value = receipt.value.filter(a => a.id !== id);
  if (qtyOverrides.value[id] !== undefined) { const rest = { ...qtyOverrides.value }; delete rest[id]; qtyOverrides.value = rest; }
  if (reductionsById.value[id] !== undefined) { const restR = { ...reductionsById.value }; delete restR[id]; reductionsById.value = restR; }
  if (notesById.value[id] !== undefined) { const restN = { ...notesById.value }; delete restN[id]; notesById.value = restN; }
  if (openFabId.value === id) openFabId.value = null;
}

type Reduction = { kind: 'percent' | 'amount'; value: number };
const reductionsById = ref<Record<number, number | Reduction>>({});

// Note / reduction / edit dialogs state
const editDialog = ref(false);
const editLineId = ref<number|null>(null);
const tempQtyStr = ref('');
const editLineName = computed(() => {
  if (editLineId.value == null) return '';
  const line = groupedReceipt.value.find(l => l.id === editLineId.value);
  return line ? line.name : '';
});
function openEdit(line: GroupedLine) { editLineId.value = line.id; tempQtyStr.value = String(line.qty ?? 0); editDialog.value = true; }
function applyEdit() {
  if (editLineId.value == null) return; const n = parseFloat(tempQtyStr.value.replace(',', '.'));
  const newQty = Number.isFinite(n) ? Math.max(0, n) : 0; setLineQuantity(editLineId.value, newQty); editDialog.value = false;
}
function tap(ch: string) { if (ch === '.') { if (tempQtyStr.value.includes('.')) return; tempQtyStr.value = tempQtyStr.value === '' ? '0.' : tempQtyStr.value + '.'; return; } tempQtyStr.value += ch; }
function backspace() { tempQtyStr.value = tempQtyStr.value.slice(0,-1); }
function clearNumpad() { tempQtyStr.value = ''; }

const noteDialog = ref(false);
const noteLineId = ref<number|null>(null);
const noteText = ref('');
const noteLineName = computed(() => { if (noteLineId.value == null) return ''; const line = groupedReceipt.value.find(l => l.id === noteLineId.value); return line ? line.name : ''; });
const hasSelectedNote = computed(() => noteLineId.value != null && Boolean(notesById.value[noteLineId.value]));
function openNote(line: GroupedLine) { noteLineId.value = line.id; noteText.value = notesById.value[line.id] ?? ''; noteDialog.value = true; }
function applyNote() { if (noteLineId.value == null) return; const id = noteLineId.value; const text = noteText.value.trim(); if (text) { notesById.value = { ...notesById.value, [id]: text }; } else { const rest = { ...notesById.value }; delete rest[id]; notesById.value = rest; } noteDialog.value = false; }
function deleteNote() { if (noteLineId.value == null) return; const id = noteLineId.value; if (notesById.value[id] !== undefined) { const rest = { ...notesById.value }; delete rest[id]; notesById.value = rest; } noteDialog.value = false; }

const reductionDialog = ref(false);
const reductionLineId = ref<number|null>(null);
const reductionKind = ref<'percent' | 'amount'>('percent');
const reductionValueStr = ref('');
const reductionLineName = computed(() => { if (reductionLineId.value == null) return ''; const line = groupedReceipt.value.find(l => l.id === reductionLineId.value); return line ? line.name : ''; });
const hasSelectedReduction = computed(() => reductionLineId.value != null && hasReduction(reductionLineId.value));
function openReduction(line: GroupedLine) {
  reductionLineId.value = line.id; const current = reductionsById.value[line.id];
  if (current == null) { reductionKind.value = 'percent'; reductionValueStr.value=''; }
  else if (typeof current === 'number') { reductionKind.value='percent'; reductionValueStr.value=String(current); }
  else { reductionKind.value=current.kind; reductionValueStr.value=String(current.value ?? ''); }
  reductionDialog.value = true;
}
function applyReduction() {
  if (reductionLineId.value == null) return; const id = reductionLineId.value; const n = parseFloat(String(reductionValueStr.value).replace(',', '.'));
  if (reductionKind.value === 'percent') { const pct = Number.isFinite(n)? Math.max(0, Math.min(100,n)):0; if (pct>0){ reductionsById.value={...reductionsById.value,[id]:pct}; } else { const rest={...reductionsById.value}; delete rest[id]; reductionsById.value=rest; } }
  else { const amt = Number.isFinite(n)? Math.max(0,n):0; if (amt>0){ reductionsById.value={...reductionsById.value,[id]:{kind:'amount', value:amt}}; } else { const rest={...reductionsById.value}; delete rest[id]; reductionsById.value=rest; } }
  reductionDialog.value=false;
}
function deleteReduction() { if (reductionLineId.value == null) return; const id = reductionLineId.value; if (reductionsById.value[id] !== undefined) { const rest={...reductionsById.value}; delete rest[id]; reductionsById.value=rest; } reductionDialog.value=false; }

// Total reduction dialog
const totalReductionDialog = ref(false);
const totalReductionKind = ref<'percent' | 'amount'>('percent');
const totalReductionValueStr = ref('');
function openTotalReduction() {
  const r = totalReduction.value; if (r == null) { totalReductionKind.value='percent'; totalReductionValueStr.value=''; }
  else if (typeof r === 'number') { totalReductionKind.value='percent'; totalReductionValueStr.value=String(r); }
  else { totalReductionKind.value=r.kind; totalReductionValueStr.value=String(r.value ?? ''); }
  totalReductionDialog.value = true;
}
function applyTotalReduction() {
  const n = parseFloat(String(totalReductionValueStr.value).replace(',', '.'));
  if (totalReductionKind.value === 'percent') { const pct = Number.isFinite(n)? Math.max(0, Math.min(100,n)):0; totalReduction.value = pct>0 ? pct : null; }
  else { const amt = Number.isFinite(n)? Math.max(0,n):0; totalReduction.value = amt>0 ? { kind:'amount', value:amt } : null; }
  totalReductionDialog.value=false;
}
function deleteTotalReduction() { totalReduction.value = null; totalReductionDialog.value=false; }

// Scanner dialog
const scannerDialog = ref(false);
function openScanner() { scannerDialog.value = true; }
function onScannerDecoded(code: string) { void handleDecodedBarcode(code); }

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

// --- Reordering state (depends on groupedReceipt) ---
const lineOrder = ref<number[]>([]);
const reorderDragId = ref<number|null>(null);
const reorderOverId = ref<number|null>(null);
const orderedLines = computed(() => {
  const orderIndex = new Map(lineOrder.value.map((id, idx) => [id, idx] as const));
  return groupedReceipt.value.slice().sort((a,b) => (orderIndex.get(a.id) ?? 0) - (orderIndex.get(b.id) ?? 0));
});
function syncLineOrder() {
  const ids = groupedReceipt.value.map(l => l.id);
  const existing = new Set(ids);
  let changed = false;
  const filtered = lineOrder.value.filter(id => existing.has(id));
  if (filtered.length !== lineOrder.value.length) changed = true;
  for (const id of ids) if (!filtered.includes(id)) { filtered.push(id); changed = true; }
  if (changed) lineOrder.value = [...filtered];
}
watch(groupedReceipt, syncLineOrder);
onMounted(syncLineOrder);
function onReorderHandleStart(id: number, ev: DragEvent) {
  reorderDragId.value = id;
  reorderOverId.value = null;
  try { ev.dataTransfer?.setData('text/plain', String(id)); } catch { /* ignore */ }
}
function onReorderDragEnter(id: number) {
  if (reorderDragId.value != null && reorderDragId.value !== id) reorderOverId.value = id;
}
function onReorderDrop(targetId: number) {
  // Only handle reorder drops, not whole-item drags (which are for deletion)
  if (reorderDragId.value != null && reorderDragId.value !== targetId) {
    const from = lineOrder.value.indexOf(reorderDragId.value);
    const to = lineOrder.value.indexOf(targetId);
    if (from > -1 && to > -1) {
      const updated = [...lineOrder.value];
      const removed = updated.splice(from, 1);
      if (removed.length) {
        updated.splice(to, 0, removed[0]!);
        lineOrder.value = updated;
      }
    }
  }
  reorderDragId.value = null;
  reorderOverId.value = null;
}

// Handle drops on receipt items - distinguish between reordering vs deletion
function onItemDrop(targetId: number, event: DragEvent) {
  // If this is a reorder drag (from the handle), handle it and stop propagation
  if (reorderDragId.value !== null) {
    event.stopPropagation();
    onReorderDrop(targetId);
  }
  // If this is a whole-item drag (draggingId set), let it bubble to trash zone
  // Don't stop propagation in this case
}

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
const totalItemsCount = computed(() => groupedReceipt.value.reduce((sum, item) => sum + item.qty, 0));
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
// --- Customers ---
interface Customer { id: number; name: string; phone?: string }
const customers = ref<Customer[]>([
  { id: 1, name: 'Walk-in', phone: '' },
  { id: 2, name: 'Alice Martin', phone: '0612-345-678' },
  { id: 3, name: 'Bob Smith', phone: '0611-111-111' }
]);
const selectedCustomerId = ref<number | null>(1);
const customerDialog = ref(false);
const customerSearch = ref('');
const newCustomerName = ref('');
const newCustomerPhone = ref('');
const selectedCustomer = computed(() => customers.value.find(c => c.id === selectedCustomerId.value) || null);
const filteredCustomers = computed(() => {
  const q = customerSearch.value.trim().toLowerCase();
  if (!q) return customers.value;
  return customers.value.filter(c => c.name.toLowerCase().includes(q) || (c.phone ?? '').toLowerCase().includes(q));
});
const canAddCustomer = computed(() => newCustomerName.value.trim().length >= 2);
function openCustomerDialog() { customerDialog.value = true; }
function selectCustomer(id: number) { selectedCustomerId.value = id; customerDialog.value = false; }
function clearSelectedCustomer() { selectedCustomerId.value = 1; }
function addCustomer() {
  if (!canAddCustomer.value) return;
  const id = Date.now();
  customers.value.push({ id, name: newCustomerName.value.trim(), phone: newCustomerPhone.value.trim() });
  newCustomerName.value = '';
  newCustomerPhone.value = '';
  selectedCustomerId.value = id;
}
// ...existing code...
function resetReceipt() {
  receipt.value = [];
  qtyOverrides.value = {};
  reductionsById.value = {};
  notesById.value = {};
  totalReduction.value = null;
}
// ...existing code...
const holds = ref<Array<{ id: number; lines: GroupedLine[]; total: number; createdAt: string; notes?: Record<number,string>; reductions?: Record<number, number|Reduction>; totalReduction?: number|Reduction|null; customerId?: number | undefined }>>([]);
const holdListDialog = ref(false);
function openHoldList() { draggingId.value = null; openFabId.value = null; searchOpen.value = false; holdListDialog.value = true; }
function openReceiptDialog() { draggingId.value = null; openFabId.value = null; searchOpen.value = false; receiptDialog.value = true; }
function putOnHold() {
  if (groupedReceipt.value.length === 0) return;
  const snapshot: GroupedLine[] = groupedReceipt.value.map(l => ({ ...l }));
  holds.value.push({ id: Date.now(), lines: snapshot, total: total.value, createdAt: new Date().toISOString(), notes: { ...notesById.value }, reductions: { ...reductionsById.value }, totalReduction: totalReduction.value ?? null, customerId: selectedCustomerId.value || undefined });
  resetReceipt();
  $q.notify({ type: 'warning', message: t('holdPlaced') });
}
function formatDate(iso: string) {
  try {
    const d = new Date(iso);
    return new Intl.DateTimeFormat(undefined, { year: '2-digit', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).format(d);
  } catch { return iso; }
}
function resumeHold(id: number) {
  const h = holds.value.find(x => x.id === id);
  if (!h) return;
  if (groupedReceipt.value.length > 0) {
    $q.dialog({ title: t('holdList'), message: t('replaceCurrentOrderConfirm'), cancel: true, persistent: true }).onOk(() => loadHold(h));
  } else {
    loadHold(h);
  }
}
function loadHold(h: { id: number; lines: GroupedLine[]; total: number; createdAt: string; notes?: Record<number,string>; reductions?: Record<number, number|Reduction>; totalReduction?: number|Reduction|null; customerId?: number | undefined }) {
  resetReceipt();
  for (const line of h.lines) {
    const base = articles.value.find(a => a.id === line.id) || { id: line.id, name: line.name, price: line.price, image: '/images/1.png', category: '' } as Article;
    receipt.value.push(base);
    qtyOverrides.value[line.id] = line.qty;
  }
  if (h.notes) notesById.value = { ...h.notes };
  if (h.reductions) reductionsById.value = { ...h.reductions };
  if (h.totalReduction !== undefined) {
    // h.totalReduction may be number | Reduction | null | undefined
    const tr = h.totalReduction;
    if (tr === null || typeof tr === 'number') {
      totalReduction.value = tr;
    } else if (typeof tr === 'object' && (tr.kind === 'percent' || tr.kind === 'amount')) {
      totalReduction.value = { kind: tr.kind, value: tr.value };
    }
  }
  if (h.customerId) selectedCustomerId.value = h.customerId;
  removeHold(h.id);
  holdListDialog.value = false;
  $q.notify({ type: 'info', message: t('resumedFromHold') });
}
function removeHold(id: number) {
  holds.value = holds.value.filter(h => h.id !== id);
}
// ...existing code...
function checkout() {
  if (groupedReceipt.value.length === 0) return;
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
  const sc = selectedCustomer.value;
  const receiptJson = {
    createdAt: new Date().toISOString(),
    currency: 'MAD',
    customer: sc ? { id: sc.id, name: sc.name, phone: sc.phone } : null,
    subtotal: subtotal.value,
    totalReduction: normalizeReduction(totalReduction.value),
    total: total.value,
    lines
  };
  console.log('Final receipt:', JSON.stringify(receiptJson, null, 2));
  $q.notify({ type: 'positive', message: t('checkoutComplete') });
  resetReceipt();
}
// ...existing code...
// (removed unused removeHold helper to satisfy ESLint)

// New order simply clears receipt and leaves selected customer (design choice)
function launchNewOrder() { resetReceipt(); }
function openHoldFromMobile() { openHoldList(); }
// ...existing code...
</script>

<style scoped>
/* Disable text selection for better touch/POS experience */
* {
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}
/* Allow text selection in input fields and editable areas */
:deep(input), :deep(textarea), :deep([contenteditable="true"]) {
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
}

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
  position: relative; /* inside flex-centered overlay */
  transform: scale(.85);
  width: 140px;
  height: 140px;
  border-radius: 50%;
  background: #fff;
  border: 4px solid var(--q-grey-4);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 24px rgba(0,0,0,0.18);
  transition: border-color .25s, box-shadow .25s, transform .25s, background .25s, opacity .25s;
  opacity: 0.8; /* Slightly more visible when first appears */
  overflow: hidden; /* Prevent content from going outside circle */
  padding: 8px; /* Small padding to keep content within bounds */
  pointer-events: all; /* Allow the drop zone to receive drag events */
  animation: trashAppear 0.3s ease-out; /* Smooth appearance animation */
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
.trash-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10050;
  pointer-events: none;
}

.trash-drop-zone.trash-over {
  background-color: rgba(244, 67, 54, 0.15);
  transform: scale(1.15);
  border-color: var(--q-negative);
  box-shadow: 0 8px 32px rgba(0,0,0,0.25), 0 0 0 6px rgba(244,67,54,0.3), 0 0 50px rgba(244,67,54,0.6);
  animation: trashGlow 0.8s ease-in-out infinite alternate;
}

.trash-instruction {
  user-select: none;
  pointer-events: none;
  text-align: center;
  font-size: 10px; /* Smaller text to fit in circle */
  line-height: 1.1;
  max-width: 100px; /* Constrain width to prevent overflow */
  margin-top: 4px; /* Small spacing from icon */
}

/* Mobile touch dragging styles */
.touch-dragging {
  opacity: 0.7;
  transform: scale(0.95);
  z-index: 1000;
}

/* Improve touch targets for receipt items on mobile */
@media (max-width: 599px) {
  .receipt-list .q-item {
    min-height: 60px;
    touch-action: pan-y; /* Allow vertical scrolling but enable custom touch handling */
  }
  
  .trash-overlay {
    background: rgba(0, 0, 0, 0.3); /* Darker overlay on mobile for better visibility */
  }
}
.trash-drop-zone.trash-active { 
  border-color: var(--q-negative); 
  opacity: 1; 
  transform: scale(1); 
  animation: trashScale 0.35s ease-out;
  box-shadow: 0 8px 24px rgba(0,0,0,0.18), 0 0 0 4px rgba(244,67,54,0.2), 0 0 30px rgba(244,67,54,0.4);
}
.trash-drop-zone.trash-active::before {
  animation: trashPulseRing 1.4s ease-out infinite;
}
.trash-drop-zone.trash-active::after {
  animation: trashPulseRing 1.4s ease-out .7s infinite;
}
  @keyframes trashScale {
    0% { transform: scale(.6); }
    100% { transform: scale(1); }
  }
@keyframes trashAppear {
  0% { 
    opacity: 0; 
    transform: scale(.5); 
  }
  100% { 
    opacity: 0.8; 
    transform: scale(.85); 
  }
}
@keyframes trashPulseRing {
  0% { box-shadow: 0 0 0 0 rgba(244,67,54,0.55); opacity: .9; }
  70% { box-shadow: 0 0 0 28px rgba(244,67,54,0); opacity: 0; }
  100% { box-shadow: 0 0 0 0 rgba(244,67,54,0); opacity: 0; }
}
@keyframes trashGlow {
  0% { 
    box-shadow: 0 8px 32px rgba(0,0,0,0.25), 0 0 0 6px rgba(244,67,54,0.3), 0 0 50px rgba(244,67,54,0.6);
  }
  100% { 
    box-shadow: 0 8px 32px rgba(0,0,0,0.25), 0 0 0 8px rgba(244,67,54,0.5), 0 0 60px rgba(244,67,54,0.8);
  }
}
  /* Reorder styles */
  .reorder-handle { cursor: grab; opacity: .6; transition: opacity .2s; }
  .reorder-handle:hover { opacity: 1; }
  .reorder-dragging { opacity: .45; }
  .reorder-over { outline: 2px dashed var(--q-primary); outline-offset: -2px; }
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

/* Improve touch responsiveness and INP performance */
:deep(.q-btn), :deep(.q-tab), :deep(.q-fab), :deep(.q-item), .product-card {
  touch-action: manipulation;
  contain: layout style;
  will-change: transform;
}

/* Performance optimizations for heavy interaction elements */
.products-scroll {
  contain: layout style paint;
  will-change: scroll-position;
  transform: translateZ(0); /* Create composite layer for smooth scrolling */
}

.receipt-list {
  contain: layout style;
}

.product-card {
  backface-visibility: hidden;
  transform: translateZ(0);
}

.product-card:active, :deep(.q-btn):active, :deep(.q-fab):active {
  will-change: transform, background-color;
}

/* Reduce paint complexity */
.product-img {
  contain: layout style paint;
  backface-visibility: hidden;
  transform: translateZ(0);
}

/* Optimize category tabs for performance */
:deep(.q-tabs) {
  contain: layout style;
}

:deep(.q-tab) {
  backface-visibility: hidden;
  transition: background-color 0.2s ease, transform 0.2s ease;
}

/* Drag and drop styles for category tabs */
:deep(.q-tabs.is-dragging) {
  user-select: none;
}

:deep(.q-tab.dragging) {
  opacity: 0.6;
  transform: scale(0.95);
  cursor: grabbing;
}

:deep(.q-tab.drag-over) {
  background-color: rgba(25, 118, 210, 0.1);
  transform: scale(1.02);
}

:deep(.q-tab[draggable="true"]:not(.dragging)) {
  cursor: grab;
}

:deep(.q-tab[draggable="true"]:hover:not(.dragging)) {
  background-color: rgba(25, 118, 210, 0.05);
}

/* Drop cursor indicators */
:deep(.q-tab.drop-cursor-left::before) {
  content: '';
  position: absolute;
  left: -2px;
  top: 10%;
  bottom: 10%;
  width: 3px;
  background: linear-gradient(to bottom, transparent, #1976d2, transparent);
  border-radius: 2px;
  animation: dropCursorPulse 1s ease-in-out infinite alternate;
  z-index: 10;
}

:deep(.q-tab.drop-cursor-right::after) {
  content: '';
  position: absolute;
  right: -2px;
  top: 10%;
  bottom: 10%;
  width: 3px;
  background: linear-gradient(to bottom, transparent, #1976d2, transparent);
  border-radius: 2px;
  animation: dropCursorPulse 1s ease-in-out infinite alternate;
  z-index: 10;
}

@keyframes dropCursorPulse {
  0% { opacity: 0.6; transform: scaleY(0.8); }
  100% { opacity: 1; transform: scaleY(1); }
}

/* Optimize receipt animations */
.receipt-added {
  animation: receiptPulse 0.6s ease-out;
  will-change: transform, background-color;
}

@keyframes receiptPulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.05); background-color: rgba(25, 118, 210, 0.1); }
  100% { transform: scale(1); }
}

/* Fix mobile FAB positioning to prevent overflow */
.receipt-fab {
  margin-left: -4px; /* reduce left gap */
}
@media (max-width: 599px) {
  .receipt-list {
    padding-left: 1px !important; /* reduce list left padding on mobile */
    padding-right: 8px !important;
  }
  .receipt-fab {
    margin-left: -10px; /* further reduce left gap on mobile */
  }
}

</style>
