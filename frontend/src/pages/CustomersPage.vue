<template>
  <q-page class="q-pa-sm q-pt-xs customers-page">
    <div class="column full-height customers-page-inner">
      <div class="row items-center justify-between q-mb-xs header-bar">
        <div class="row items-center q-gutter-x-sm q-ml-xs title-block">
          <div class="text-subtitle1">{{ $t('customers') }}</div>
          <q-badge v-if="filteredRows.length" color="primary" align="middle" class="q-ml-md">{{ filteredRows.length }}</q-badge>
        </div>
        <div class="row items-center no-wrap q-gutter-x-xs actions-block">
          <q-btn dense unelevated color="primary" icon="person_add" :label="$t('add')" @click="openAdd" class="add-btn" />
        </div>
      </div>
      <div class="row q-col-gutter-xs q-mb-xs filters-row">
        <div class="col-12 col-sm-6 col-md-4">
          <q-input v-model="search" dense outlined clearable debounce="200" :placeholder="$t('search')">
            <template #prepend><q-icon name="search" /></template>
          </q-input>
        </div>
        <div class="col-12 col-sm-6 col-md-8 flex items-center justify-end gap-xs meta-slot">
          <!-- future quick filters / export buttons -->
        </div>
      </div>
      <div class="table-responsive-wrapper">
        <q-table
          class="col customers-table"
          flat
          :rows="filteredRows"
          :columns="columnsDisplay"
          row-key="id"
          :pagination="pagination"
          @update:pagination="val => Object.assign(pagination, val)"
          :rows-per-page-options="[10,20,50]"
        >
          <template #body-cell-actions="props">
            <q-td :props="props">
              <q-btn dense flat size="sm" icon="edit" color="primary" @click="openEdit(props.row)" />
              <q-btn dense flat size="sm" icon="delete" color="negative" @click="confirmDelete(props.row)" />
            </q-td>
          </template>
          <template #no-data>
            <div class="text-grey-7 q-pa-md flex flex-center column">
              <q-icon name="groups" size="42px" class="q-mb-sm" />
              <div>{{ $t('noResults') }}</div>
            </div>
          </template>
        </q-table>
      </div>
    </div>
    <!-- Add/Edit Dialog -->
    <q-dialog v-model="dialog">
      <q-card style="width:420px; max-width:95vw;" class="q-pa-sm">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{ isEdit ? $t('edit') : $t('add') }} {{ $t('customers') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-pt-sm">
          <q-input v-model="form.name" :label="$t('name')" dense outlined autofocus :error="nameError" :error-message="$t('required')" />
          <q-input v-model="form.phone" :label="$t('phone')" dense outlined class="q-mt-sm" />
          <q-input v-model="form.email" :label="$t('email')" dense outlined class="q-mt-sm" />
          <q-input v-model="form.notes" :label="$t('pos.itemNote')" type="textarea" autogrow dense outlined class="q-mt-sm" />
          <div class="row q-col-gutter-sm q-mt-sm">
            <div class="col-6">
              <q-select v-model="form.customerType" :options="customerTypeOptions" option-value="value" option-label="label" dense outlined emit-value map-options :label="$t('customerType')" />
            </div>
            <div class="col-6">
              <q-select v-model="form.priceCategory" :options="priceCategoryOptions" option-value="value" option-label="label" dense outlined emit-value map-options :label="$t('priceCategory')" />
            </div>
          </div>
        </q-card-section>
        <q-separator inset />
        <q-card-actions align="right">
          <q-btn flat :label="$t('cancel')" v-close-popup />
          <q-btn color="primary" :disable="!canSave" :label="$t('save')" @click="save" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Delete Confirm -->
    <q-dialog v-model="deleteDialog">
      <q-card style="width:400px; max-width:90vw;" class="q-pa-sm">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{ $t('delete') }} {{ $t('customers') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-pt-sm">
          <div>{{ $t('confirmDeleteProduct', { name: deleteTarget?.name || '' }) }}</div>
        </q-card-section>
        <q-separator />
        <q-card-actions align="right">
          <q-btn flat :label="$t('cancel')" v-close-popup />
          <q-btn color="negative" :label="$t('delete')" @click="doDelete" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue';
import { useQuasar } from 'quasar';
import type { QTableColumn } from 'quasar';
import { useI18n } from 'vue-i18n';

const $q = useQuasar();
const { t } = useI18n();

type CustomerType = 'individual' | 'company';
type PriceCategory = 'standard' | 'preferential' | 'wholesale';
interface Customer { id: number; name: string; phone?: string; email?: string; notes?: string; createdAt: string; customerType: CustomerType; priceCategory: PriceCategory }

const customers = ref<Customer[]>([]);
const search = ref('');
const customerTypeOptions = [
  { value: 'individual', label: t('individual') },
  { value: 'company', label: t('company') }
];
const priceCategoryOptions = [
  { value: 'standard', label: t('standard') },
  { value: 'preferential', label: t('preferential') },
  { value: 'wholesale', label: t('wholesale') }
];

const columns = [
  { name: 'name', field: 'name', label: t('name'), align: 'left', sortable: true },
  { name: 'customerType', field: 'customerType', label: t('customerType') || 'Type', align: 'left', sortable: true, format: (val: CustomerType) => t(val) },
  { name: 'priceCategory', field: 'priceCategory', label: t('priceCategory') || 'Price Cat', align: 'left', sortable: true, format: (val: PriceCategory) => t(val) },
  { name: 'phone', field: 'phone', label: t('phone'), align: 'left', sortable: true },
  { name: 'email', field: 'email', label: t('email'), align: 'left', sortable: true },
  { name: 'createdAt', field: 'createdAt', label: t('date') || 'Date', align: 'left', sortable: true, format: (val: string) => new Date(val).toLocaleDateString() },
  { name: 'actions', field: 'actions', label: '', align: 'right' }
] as QTableColumn[];

const filteredRows = computed(() => {
  const q = search.value.trim().toLowerCase();
  if (!q) return customers.value;
  return customers.value.filter(c =>
    c.name.toLowerCase().includes(q) ||
    (c.phone || '').toLowerCase().includes(q) ||
    (c.email || '').toLowerCase().includes(q) ||
    c.customerType.toLowerCase().includes(q) ||
    c.priceCategory.toLowerCase().includes(q)
  );
});

// Responsive columns: hide email & createdAt and maybe phone on very narrow screens
const viewportWidth = ref(typeof window !== 'undefined' ? window.innerWidth : 1200);
const columnsDisplay = computed(() => {
  const w = viewportWidth.value;
  if (w >= 900) return columns; // full
  if (w >= 600) return columns.filter(c => c.name !== 'createdAt');
  // small: keep name, priceCategory, actions
  return columns.filter(c => ['name','priceCategory','actions'].includes(c.name));
});
function handleResize() { viewportWidth.value = window.innerWidth; }
window.addEventListener('resize', handleResize, { passive: true });
onBeforeUnmount(() => window.removeEventListener('resize', handleResize));

// Quasar QTable accepts a plain reactive object; using ref wrapper caused template .value confusion, use reactive-like via ref unwrap in script
const pagination = { page: 1, rowsPerPage: 20, sortBy: 'name', descending: false };

// Dialog state
const dialog = ref(false);
const isEdit = ref(false);
const form = ref<{ id?: number; name: string; phone?: string; email?: string; notes?: string; customerType: CustomerType; priceCategory: PriceCategory }>({ name: '', customerType: 'individual', priceCategory: 'standard' });
const nameError = computed(() => form.value.name.trim().length === 0 && dialog.value);
const canSave = computed(() => !nameError.value);
function openAdd() {
  isEdit.value = false;
  form.value = { name: '', phone: '', email: '', notes: '', customerType: 'individual', priceCategory: 'standard' };
  dialog.value = true;
}
function openEdit(row: Customer) {
  isEdit.value = true;
  form.value = { id: row.id, name: row.name, phone: row.phone || '', email: row.email || '', notes: row.notes || '', customerType: row.customerType, priceCategory: row.priceCategory };
  dialog.value = true;
}
function save() {
  if (!canSave.value) return;
  if (isEdit.value && form.value.id != null) {
    const idx = customers.value.findIndex(c => c.id === form.value.id);
    if (idx > -1) {
      customers.value[idx] = { ...customers.value[idx], ...form.value } as Customer;
    }
    $q.notify({ type: 'positive', message: t('save') });
  } else {
  customers.value.push({ id: Date.now(), name: form.value.name.trim(), phone: form.value.phone?.trim() || '', email: form.value.email?.trim() || '', notes: form.value.notes?.trim() || '', createdAt: new Date().toISOString(), customerType: form.value.customerType, priceCategory: form.value.priceCategory });
    $q.notify({ type: 'positive', message: t('save') });
  }
  persist();
  dialog.value = false;
}

// Delete
const deleteDialog = ref(false);
const deleteTarget = ref<Customer | null>(null);
function confirmDelete(row: Customer) {
  deleteTarget.value = row;
  deleteDialog.value = true;
}
function doDelete() {
  if (deleteTarget.value) {
    customers.value = customers.value.filter(c => c.id !== deleteTarget.value!.id);
    persist();
    $q.notify({ type: 'warning', message: t('delete') });
  }
  deleteDialog.value = false;
}

function persist() {
  try { localStorage.setItem('customers', JSON.stringify(customers.value)); } catch { /* ignore */ }
}
function load() {
  try {
    const raw = localStorage.getItem('customers');
    if (raw) {
      const arr = JSON.parse(raw);
      if (Array.isArray(arr)) {
        customers.value = arr.map(r => ({
          ...r,
          createdAt: r.createdAt || new Date().toISOString(),
          customerType: (r.customerType === 'company' ? 'company' : 'individual') as CustomerType,
          priceCategory: (['standard','preferential','wholesale'].includes(r.priceCategory) ? r.priceCategory : 'standard') as PriceCategory
        }));
      }
    }
  } catch { /* ignore */ }
  if (customers.value.length === 0) {
    // seed demo data
    customers.value = [
      { id: 1, name: 'Walk-in', phone: '', email: '', notes: '', createdAt: new Date().toISOString(), customerType: 'individual', priceCategory: 'standard' },
      { id: 2, name: 'Alice Martin', phone: '0612-345-678', email: 'alice@example.com', notes: '', createdAt: new Date().toISOString(), customerType: 'individual', priceCategory: 'preferential' },
      { id: 3, name: 'Bob Smith', phone: '0611-111-111', email: 'bob@example.com', notes: '', createdAt: new Date().toISOString(), customerType: 'company', priceCategory: 'wholesale' }
    ];
  }
}

onMounted(load);
watch(customers, persist, { deep: true });
</script>

<style scoped>
.full-height { height: 100%; }
.customers-page { margin: 0 auto; }
.customers-table :deep(thead th) { font-weight: 600; background: var(--q-grey-2); }
.customers-table :deep(tbody tr:hover) { background: rgba(25,118,210,0.06); transition: background .15s ease; }
.ellipsis-2-lines { display: -webkit-box; line-clamp: 2; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
/* tighter spacing overrides */
.q-page.customers-page :deep(.q-table) { font-size: 0.85rem; }
@media (max-width: 599px) { .q-page.customers-page { padding-left: 4px; padding-right: 4px; } }
@media (max-width: 599px) { .customers-page { padding-bottom: 80px; } }
/* responsive enhancements */
.table-responsive-wrapper { width: 100%; overflow-x: auto; overscroll-behavior: contain; }
.table-responsive-wrapper :deep(table) { min-width: 640px; }
.header-bar { flex-wrap: wrap; }
.header-bar .actions-block { gap: 4px; }
@media (max-width: 820px) {
  .header-bar { row-gap: 4px; }
  .header-bar .title-block { width: 100%; }
  .header-bar .actions-block { width: 100%; justify-content: flex-start; }
  .add-btn { width: 100%; }
}
@media (max-width: 500px) {
  .table-responsive-wrapper :deep(table) { min-width: 480px; }
}
</style>
