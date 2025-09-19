<template>
  <q-page class="q-pa-sm q-pt-xs customers-page">
    <div class="column full-height customers-page-inner">
      <div class="row items-center justify-between q-mb-xs header-bar">
        <div class="row items-center q-gutter-x-sm q-ml-xs title-block">
          <div class="text-subtitle1">{{ $t('customers') }}</div>
          <q-badge v-if="filteredRows.length" color="primary" align="middle" class="q-ml-md">{{ filteredRows.length }}</q-badge>
        </div>
        <div class="row items-center no-wrap q-gutter-x-xs actions-block">
          <q-btn dense flat color="primary" icon="person_add" :label="$t('add')" @click="openAdd" class="add-btn" />
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
      <div class="table-container">
        <q-table
          class="customers-table"
          flat
          :rows="filteredRows"
          :columns="columnsDisplay"
          row-key="id"
          :pagination="pagination"
          @update:pagination="val => Object.assign(pagination, val)"
          :rows-per-page-options="[10,20,50]"
        >
          <!-- Customer Type with colored chips -->
          <template #body-cell-customerType="props">
            <q-td :props="props">
              <q-chip 
                :color="props.value === 'individual' ? 'blue-grey-4' : 'indigo-5'"
                :text-color="props.value === 'individual' ? 'white' : 'white'"
                :icon="props.value === 'individual' ? 'person' : 'business'"
                dense
                size="sm"
              >
                {{ t(props.value) }}
              </q-chip>
            </q-td>
          </template>

          <!-- Price Category with colored badges -->
          <template #body-cell-priceCategory="props">
            <q-td :props="props">
              <q-badge 
                :color="props.value === 'wholesale' ? 'purple' : props.value === 'preferential' ? 'orange' : 'grey-6'"
                :label="t(props.value)"
                align="middle"
              />
            </q-td>
          </template>

          <!-- Status with colored chips -->
          <template #body-cell-isActive="props">
            <q-td :props="props">
              <q-chip 
                :color="props.value ? 'positive' : 'grey-5'"
                :text-color="props.value ? 'white' : 'white'"
                :icon="props.value ? 'check_circle' : 'cancel'"
                dense
                size="sm"
              >
                {{ props.value ? t('active') : t('inactive') }}
              </q-chip>
            </q-td>
          </template>

          <!-- Credit Limit with color coding -->
          <template #body-cell-creditLimit="props">
            <q-td :props="props">
              <span 
                v-if="props.value != null && Number(props.value) > 0"
                :class="[
                  'text-weight-medium',
                  Number(props.value) >= 5000 ? 'text-green-7' : 
                  Number(props.value) >= 1000 ? 'text-orange-7' : 
                  'text-blue-7'
                ]"
              >
                {{ Number(props.value).toFixed(2) }} MAD
              </span>
              <span v-else class="text-grey-5">-</span>
            </q-td>
          </template>

          <!-- Discount Percentage with colored badges -->
          <template #body-cell-discountPercentage="props">
            <q-td :props="props">
              <q-badge 
                v-if="props.value != null && Number(props.value) > 0"
                :color="Number(props.value) >= 15 ? 'red-5' : Number(props.value) >= 10 ? 'orange-5' : 'blue-5'"
                :label="`${Number(props.value)}%`"
                align="middle"
              />
              <span v-else class="text-grey-5">-</span>
            </q-td>
          </template>

          <!-- Enhanced name cell with avatar and financial info on mobile -->
          <template #body-cell-name="props">
            <q-td :props="props" class="name-cell">
              <div class="row items-center no-wrap">
                <q-avatar 
                  size="28px" 
                  :color="props.row.customerType === 'company' ? 'indigo-5' : 'blue-grey-4'"
                  text-color="white"
                  :icon="props.row.customerType === 'company' ? 'business' : 'person'"
                  class="q-mr-sm avatar-centered"
                />
                <div class="column justify-center">
                  <span class="text-weight-medium">{{ props.value }}</span>
                  <!-- Show financial info on very small screens -->
                  <div v-if="viewportWidth < 500" class="row items-center q-gutter-xs q-mt-xs">
                    <q-badge 
                      v-if="props.row.creditLimit && props.row.creditLimit > 0"
                      :color="props.row.creditLimit >= 5000 ? 'green-5' : props.row.creditLimit >= 1000 ? 'orange-5' : 'blue-5'"
                      :label="`${props.row.creditLimit} MAD`"
                      size="xs"
                    />
                    <q-badge 
                      v-if="props.row.discountPercentage && props.row.discountPercentage > 0"
                      :color="props.row.discountPercentage >= 15 ? 'red-5' : props.row.discountPercentage >= 10 ? 'orange-5' : 'blue-5'"
                      :label="`${props.row.discountPercentage}%`"
                      size="xs"
                    />
                  </div>
                </div>
              </div>
            </q-td>
          </template>

          <!-- Actions with enhanced styling -->
          <template #body-cell-actions="props">
            <q-td :props="props" class="actions-cell">
              <div class="row items-center no-wrap q-gutter-xs">
                <q-btn 
                  dense 
                  round 
                  flat 
                  size="sm" 
                  icon="edit" 
                  color="primary" 
                  @click="openEdit(props.row)"
                  class="hover-btn action-btn"
                >
                  <q-tooltip>{{ t('edit') }}</q-tooltip>
                </q-btn>
                <q-btn 
                  dense 
                  round 
                  flat 
                  size="sm" 
                  icon="delete" 
                  color="negative" 
                  @click="confirmDelete(props.row)"
                  class="hover-btn action-btn"
                >
                  <q-tooltip>{{ t('delete') }}</q-tooltip>
                </q-btn>
              </div>
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
      <q-card style="width:620px; max-width:95vw;" class="q-pa-sm">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">{{ isEdit ? $t('edit') : $t('add') }} {{ $t('customers') }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-pt-sm">
          <!-- Basic Information -->
          <div class="text-subtitle2 text-primary q-mb-sm">{{ $t('basicInformation') }}</div>
          <q-input v-model="form.name" :label="$t('name')" dense outlined autofocus :error="nameError" :error-message="$t('required')" />
          
          <div class="q-mt-sm">
            <q-toggle v-model="form.isActive" :label="$t('active')" dense />
          </div>
          
          <div class="row q-col-gutter-sm q-mt-sm">
            <div class="col-6">
              <q-select v-model="form.customerType" :options="customerTypeOptions" option-value="value" option-label="label" dense outlined emit-value map-options :label="$t('customerType')" />
            </div>
            <div class="col-6">
              <q-select v-model="form.priceCategory" :options="priceCategoryOptions" option-value="value" option-label="label" dense outlined emit-value map-options :label="$t('priceCategory')" />
            </div>
          </div>

          <!-- Contact Information -->
          <div class="text-subtitle2 text-primary q-mb-sm q-mt-md">{{ $t('contactInformation') }}</div>
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-md-6">
              <q-input v-model="form.phone" :label="$t('phone')" dense outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.email" :label="$t('email')" type="email" dense outlined />
            </div>
          </div>
          
          <div class="row q-col-gutter-sm q-mt-sm" v-if="form.customerType === 'company'">
            <div class="col-12 col-md-6">
              <q-input v-model="form.contactPerson" :label="$t('contactPerson')" dense outlined />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.website" :label="$t('website')" dense outlined />
            </div>
          </div>

          <!-- Address Information -->
          <div class="text-subtitle2 text-primary q-mb-sm q-mt-md">{{ $t('addressInformation') }}</div>
          <q-input v-model="form.address" :label="$t('address')" dense outlined />
          <div class="row q-col-gutter-sm q-mt-sm">
            <div class="col-8">
              <q-input v-model="form.city" :label="$t('city')" dense outlined />
            </div>
            <div class="col-4">
              <q-input v-model="form.postalCode" :label="$t('postalCode')" dense outlined />
            </div>
          </div>

          <!-- Financial Information -->
          <div class="text-subtitle2 text-primary q-mb-sm q-mt-md">{{ $t('financialInformation') }}</div>
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-md-4">
              <q-input v-model="form.taxId" :label="$t('taxId')" dense outlined />
            </div>
            <div class="col-12 col-md-4">
              <q-input v-model="form.creditLimit" :label="$t('creditLimit')" type="number" suffix="MAD" dense outlined />
            </div>
            <div class="col-12 col-md-4">
              <q-input v-model="form.discountPercentage" :label="$t('discountPercentage')" type="number" suffix="%" min="0" max="100" dense outlined />
            </div>
          </div>

          <!-- Notes -->
          <div class="text-subtitle2 text-primary q-mb-sm q-mt-md">{{ $t('additionalNotes') }}</div>
          <q-input v-model="form.notes" :label="$t('notes')" type="textarea" autogrow rows="2" dense outlined />
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
          <div>{{ $t('confirmDeleteCustomer', { name: deleteTarget?.name || '' }) }}</div>
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
interface Customer { 
  id: number; 
  name: string; 
  phone?: string; 
  email?: string; 
  notes?: string; 
  createdAt: string; 
  customerType: CustomerType; 
  priceCategory: PriceCategory;
  // Additional enriched fields
  address?: string;
  city?: string;
  postalCode?: string;
  taxId?: string;
  contactPerson?: string;
  website?: string;
  creditLimit?: number;
  discountPercentage?: number;
  isActive?: boolean;
  lastOrderDate?: string;
  totalOrders?: number;
  totalSpent?: number;
}

const customers = ref<Customer[]>([]);
const search = ref('');
const customerTypeOptions = computed(() => [
  { value: 'individual', label: t('individual') },
  { value: 'company', label: t('company') }
]);
const priceCategoryOptions = computed(() => [
  { value: 'standard', label: t('standard') },
  { value: 'preferential', label: t('preferential') },
  { value: 'wholesale', label: t('wholesale') }
]);

const columns = computed(() => [
  { name: 'name', field: 'name', label: t('name'), align: 'left', sortable: true },
  { name: 'customerType', field: 'customerType', label: t('customerType'), align: 'left', sortable: true, format: (val: CustomerType) => t(val) },
  { name: 'phone', field: 'phone', label: t('phone'), align: 'left', sortable: true },
  { name: 'email', field: 'email', label: t('email'), align: 'left', sortable: true },
  { name: 'city', field: 'city', label: t('city'), align: 'left', sortable: true },
  { name: 'priceCategory', field: 'priceCategory', label: t('priceCategory'), align: 'left', sortable: true, format: (val: PriceCategory) => t(val) },
  { name: 'creditLimit', field: 'creditLimit', label: t('creditLimit'), align: 'right', sortable: true },
  { name: 'discountPercentage', field: 'discountPercentage', label: t('discount'), align: 'right', sortable: true },
  { name: 'isActive', field: 'isActive', label: t('status'), align: 'center', sortable: true, format: (val: boolean) => val ? t('active') : t('inactive') },
  { name: 'createdAt', field: 'createdAt', label: t('date'), align: 'left', sortable: true, format: (val: string) => new Date(val).toLocaleDateString() },
  { name: 'actions', field: 'actions', label: '', align: 'right', style: 'width: 100px', headerStyle: 'width: 100px' }
] as QTableColumn[]);

const filteredRows = computed(() => {
  const q = search.value.trim().toLowerCase();
  if (!q) return customers.value;
  return customers.value.filter(c =>
    c.name.toLowerCase().includes(q) ||
    (c.phone || '').toLowerCase().includes(q) ||
    (c.email || '').toLowerCase().includes(q) ||
    (c.city || '').toLowerCase().includes(q) ||
    (c.address || '').toLowerCase().includes(q) ||
    (c.taxId || '').toLowerCase().includes(q) ||
    c.customerType.toLowerCase().includes(q) ||
    c.priceCategory.toLowerCase().includes(q)
  );
});

// Responsive columns: Credit Limit and Discount are ALWAYS visible
const viewportWidth = ref(typeof window !== 'undefined' ? window.innerWidth : 1200);
const columnsDisplay = computed(() => {
  const w = viewportWidth.value;
  const cols = columns.value; // Get the computed columns array
  let filteredColumns;
  
  // Always include creditLimit and discountPercentage - they're essential financial info
  const essentialColumns = ['name', 'customerType', 'creditLimit', 'discountPercentage', 'actions'];
  
  if (w >= 1200) {
    filteredColumns = cols; // full - show all columns
  } else if (w >= 1000) {
    filteredColumns = cols.filter(c => !['createdAt'].includes(c.name)); // hide only date
  } else if (w >= 800) {
    filteredColumns = cols.filter(c => !['createdAt', 'city'].includes(c.name)); // hide date and city
  } else if (w >= 650) {
    filteredColumns = cols.filter(c => !['createdAt', 'city', 'email'].includes(c.name)); // hide date, city, and email
  } else if (w >= 500) {
    filteredColumns = cols.filter(c => !['createdAt', 'city', 'email', 'isActive'].includes(c.name)); // hide date, city, email, and status
  } else {
    // small: keep essential columns including financial info
    filteredColumns = cols.filter(c => ['name','customerType','phone','creditLimit','discountPercentage','actions'].includes(c.name));
  }
  
  // Force include essential financial columns if somehow they got filtered out
  essentialColumns.forEach(colName => {
    const hasColumn = filteredColumns.some(c => c.name === colName);
    if (!hasColumn) {
      const column = cols.find(c => c.name === colName);
      if (column) {
        filteredColumns.push(column);
      }
    }
  });
  
  return filteredColumns;
});

function handleResize() { 
  viewportWidth.value = window.innerWidth; 
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('resize', handleResize, { passive: true });
  }
});
onBeforeUnmount(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', handleResize);
  }
});

// Quasar QTable accepts a plain reactive object; using ref wrapper caused template .value confusion, use reactive-like via ref unwrap in script
const pagination = { page: 1, rowsPerPage: 20, sortBy: 'name', descending: false };

// Dialog state
const dialog = ref(false);
const isEdit = ref(false);
const form = ref<{ 
  id?: number; 
  name: string; 
  phone?: string; 
  email?: string; 
  notes?: string; 
  customerType: CustomerType; 
  priceCategory: PriceCategory;
  address?: string;
  city?: string;
  postalCode?: string;
  taxId?: string;
  contactPerson?: string;
  website?: string;
  creditLimit?: number;
  discountPercentage?: number;
  isActive?: boolean;
}>({ name: '', customerType: 'individual', priceCategory: 'standard', isActive: true, creditLimit: 0, discountPercentage: 0 });
const nameError = computed(() => form.value.name.trim().length === 0 && dialog.value);
const canSave = computed(() => !nameError.value);
function openAdd() {
  isEdit.value = false;
  form.value = { 
    name: '', 
    customerType: 'individual', 
    priceCategory: 'standard',
    isActive: true
  };
  dialog.value = true;
}
function openEdit(row: Customer) {
  isEdit.value = true;
  Object.assign(form.value, {
    id: row.id, 
    name: row.name, 
    phone: row.phone || undefined, 
    email: row.email || undefined, 
    notes: row.notes || undefined, 
    customerType: row.customerType, 
    priceCategory: row.priceCategory,
    address: row.address || undefined,
    city: row.city || undefined,
    postalCode: row.postalCode || undefined,
    taxId: row.taxId || undefined,
    contactPerson: row.contactPerson || undefined,
    website: row.website || undefined,
    creditLimit: row.creditLimit || undefined,
    discountPercentage: row.discountPercentage || undefined,
    isActive: row.isActive ?? true
  });
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
    const newCustomer: Customer = {
      id: Date.now(),
      name: form.value.name.trim(),
      createdAt: new Date().toISOString(),
      customerType: form.value.customerType,
      priceCategory: form.value.priceCategory,
      isActive: form.value.isActive ?? true,
      totalOrders: 0,
      totalSpent: 0,
      ...(form.value.phone?.trim() && { phone: form.value.phone.trim() }),
      ...(form.value.email?.trim() && { email: form.value.email.trim() }),
      ...(form.value.notes?.trim() && { notes: form.value.notes.trim() }),
      ...(form.value.address?.trim() && { address: form.value.address.trim() }),
      ...(form.value.city?.trim() && { city: form.value.city.trim() }),
      ...(form.value.postalCode?.trim() && { postalCode: form.value.postalCode.trim() }),
      ...(form.value.taxId?.trim() && { taxId: form.value.taxId.trim() }),
      ...(form.value.contactPerson?.trim() && { contactPerson: form.value.contactPerson.trim() }),
      ...(form.value.website?.trim() && { website: form.value.website.trim() }),
      ...((form.value.creditLimit && form.value.creditLimit > 0) && { creditLimit: form.value.creditLimit }),
      ...((form.value.discountPercentage && form.value.discountPercentage > 0) && { discountPercentage: form.value.discountPercentage })
    };
    customers.value.push(newCustomer);
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
      { 
        id: 1, 
        name: 'Walk-in Customer', 
        createdAt: new Date().toISOString(), 
        customerType: 'individual', 
        priceCategory: 'standard',
        isActive: true,
        totalOrders: 0,
        totalSpent: 0
      },
      { 
        id: 2, 
        name: 'Alice Martin', 
        phone: '0612-345-678', 
        email: 'alice@example.com', 
        address: '123 Main Street',
        city: 'Casablanca',
        postalCode: '20000',
        notes: 'Regular customer, prefers morning visits', 
        createdAt: new Date().toISOString(), 
        customerType: 'individual', 
        priceCategory: 'preferential',
        discountPercentage: 5,
        creditLimit: 1000,
        isActive: true,
        totalOrders: 25,
        totalSpent: 2500.50,
        lastOrderDate: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString()
      },
      { 
        id: 3, 
        name: 'TechnoSoft Solutions', 
        phone: '0522-999-888', 
        email: 'contact@technosoft.ma', 
        contactPerson: 'Mohammed Bennani',
        website: 'www.technosoft.ma',
        address: '456 Business Avenue',
        city: 'Rabat',
        postalCode: '10000',
        taxId: 'IF12345678',
        notes: 'IT company, bulk orders for office supplies', 
        createdAt: new Date().toISOString(), 
        customerType: 'company', 
        priceCategory: 'wholesale',
        discountPercentage: 15,
        creditLimit: 5000,
        isActive: true,
        totalOrders: 8,
        totalSpent: 12400.75,
        lastOrderDate: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString()
      },
      { 
        id: 4, 
        name: 'Sara El Amrani', 
        phone: '0661-123-456', 
        email: 'sara.elamrani@gmail.com', 
        address: '789 Student Street',
        city: 'Fès',
        postalCode: '30000',
        notes: 'Art student, frequent stationery purchases', 
        createdAt: new Date().toISOString(), 
        customerType: 'individual', 
        priceCategory: 'standard',
        isActive: true,
        totalOrders: 15,
        totalSpent: 850.25,
        lastOrderDate: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
      }
    ];
  }
}

onMounted(load);
watch(customers, persist, { deep: true });
</script>

<style scoped>
.full-height { height: 100%; }
.customers-page { margin: 0 auto; }
.customers-table :deep(thead th) { 
  font-weight: 600; 
  background: var(--q-grey-2); 
  border-bottom: 2px solid var(--q-primary);
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}
.customers-table :deep(tbody tr:hover) { 
  background: rgba(25,118,210,0.08); 
  transition: background .2s ease; 
}
.customers-table :deep(tbody tr:nth-child(even)) {
  background-color: rgba(0,0,0,0.02);
}
.customers-table :deep(tbody tr:nth-child(even):hover) {
  background: rgba(25,118,210,0.08);
}
.customers-table :deep(tbody td) {
  padding: 12px 8px;
  vertical-align: middle !important;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  max-width: 150px;
}
.customers-table :deep(thead th) {
  padding: 16px 8px;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* Button hover effects */
.hover-btn {
  transition: all 0.2s ease;
}
.hover-btn:hover {
  transform: scale(1.1);
}

/* Action buttons - ensure they're always visible */
.actions-cell {
  min-width: 80px !important;
  width: 80px !important;
  max-width: 80px !important;
  padding: 8px 4px !important;
}
.action-btn {
  opacity: 1 !important;
  visibility: visible !important;
  display: inline-flex !important;
  min-width: 28px;
  min-height: 28px;
  margin: 0 1px;
}
.action-btn .q-icon {
  font-size: 14px;
}

/* Force actions column to always be visible */
.customers-table :deep(th:last-child),
.customers-table :deep(td:last-child) {
  min-width: 80px !important;
  width: 80px !important;
  max-width: 80px !important;
}

/* Make sure row actions don't get cut off */
.customers-table :deep(tbody tr) {
  min-height: 48px;
}

/* Specific column width controls */
.customers-table :deep(td:nth-child(1)) { /* Name column */
  max-width: 180px;
  min-width: 120px;
}
.customers-table :deep(td:nth-child(2)) { /* Customer Type */
  max-width: 120px;
  min-width: 100px;
}
.customers-table :deep(td:nth-child(3)) { /* Phone */
  max-width: 130px;
  min-width: 110px;
}

/* Financial columns - ensure they're visible and properly sized */
.customers-table :deep([data-label*="Credit"]) { /* Credit Limit column */
  max-width: 120px;
  min-width: 90px;
  text-align: right;
}
.customers-table :deep([data-label*="Discount"]) { /* Discount column */
  max-width: 100px;
  min-width: 80px;
  text-align: right;
}

/* Make financial data more prominent */
.customers-table :deep(.text-green-7),
.customers-table :deep(.text-orange-7),
.customers-table :deep(.text-blue-7) {
  font-weight: 600;
}

/* Ensure actions column is sticky on mobile */
@media (max-width: 599px) {
  .table-container {
    margin: 0 -4px; /* Extend to edges on mobile */
  }
  
  .actions-cell {
    position: sticky;
    right: 0;
    background: inherit;
    z-index: 1;
    box-shadow: -2px 0 4px rgba(0,0,0,0.1);
    min-width: 70px !important;
    width: 70px !important;
    max-width: 70px !important;
  }
  .customers-table :deep(thead th:last-child) {
    position: sticky;
    right: 0;
    background: var(--q-grey-2);
    z-index: 2;
    box-shadow: -2px 0 4px rgba(0,0,0,0.1);
    min-width: 70px !important;
    width: 70px !important;
    max-width: 70px !important;
  }
  .action-btn {
    min-width: 26px;
    min-height: 26px;
    margin: 0 1px;
  }
  
  /* Tighter mobile columns */
  .customers-table :deep(td),
  .customers-table :deep(th) {
    padding: 8px 6px;
    max-width: 120px;
  }
  
  .customers-table :deep(td:nth-child(1)) { /* Name column */
    max-width: 160px;
    min-width: 120px;
  }
  
  /* Allow name column to expand when showing financial info */
  .customers-table :deep(td:nth-child(1)) {
    white-space: normal;
    overflow: visible;
  }
}

@media (max-width: 400px) {
  .customers-table :deep(td),
  .customers-table :deep(th) {
    padding: 6px 4px;
    font-size: 0.8rem;
    max-width: 100px;
  }
  
  .actions-cell {
    min-width: 60px !important;
    width: 60px !important;
    max-width: 60px !important;
  }
  
  .customers-table :deep(thead th:last-child) {
    min-width: 60px !important;
    width: 60px !important;
    max-width: 60px !important;
  }
}

/* Chip animations */
.customers-table :deep(.q-chip) {
  transition: all 0.2s ease;
}
.customers-table :deep(.q-chip:hover) {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0,0,0,0.15);
}

/* Badge styling */
.customers-table :deep(.q-badge) {
  font-weight: 500;
  font-size: 0.75rem;
  padding: 4px 8px;
  border-radius: 12px;
}

/* Avatar styling - improved centering */
.customers-table :deep(.q-avatar) {
  border: 2px solid rgba(255,255,255,0.8);
  box-shadow: 0 1px 3px rgba(0,0,0,0.12);
  flex-shrink: 0;
}

.customers-table :deep(.name-cell) {
  vertical-align: middle !important;
}

.customers-table :deep(.avatar-centered) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.customers-table :deep(.name-cell .row) {
  align-items: center;
  min-height: 40px; /* Ensure minimum height for proper centering */
}

/* Color utilities */
.text-weight-medium { font-weight: 500; }

/* tighter spacing overrides */
.q-page.customers-page :deep(.q-table) { font-size: 0.85rem; }
@media (max-width: 599px) { 
  .q-page.customers-page { 
    padding-left: 4px; 
    padding-right: 4px; 
  }
  
  /* Mobile avatar styles */
  .customers-table :deep(.q-avatar) {
    width: 20px !important;
    height: 20px !important;
    font-size: 10px;
    flex-shrink: 0;
  }
  
  .customers-table :deep(.name-cell .row) {
    min-height: 32px; /* Smaller height for mobile */
  }
}
@media (max-width: 599px) { 
  .customers-page { 
    padding-bottom: 80px; 
  } 
}

/* Table container - prevent overflow */
.table-container { 
  width: 100%;
  max-width: 100%;
  overflow-x: auto;
  overflow-y: hidden;
  overscroll-behavior: contain;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.customers-table {
  width: 100%;
  max-width: 100%;
}

.customers-table :deep(.q-table__container) {
  width: 100%;
  max-width: 100%;
}

.customers-table :deep(.q-table__middle) {
  width: 100%;
  max-width: 100%;
  overflow-x: auto;
}

.customers-table :deep(table) {
  width: 100%;
  table-layout: auto;
  border-collapse: collapse;
}
.header-bar { 
  flex-wrap: wrap; 
}
.header-bar .actions-block { 
  gap: 4px; 
}
@media (max-width: 820px) {
  .header-bar { 
    row-gap: 4px; 
  }
  .header-bar .title-block { 
    width: 100%; 
  }
  .header-bar .actions-block { 
    width: 100%; 
    justify-content: flex-start; 
  }
  .add-btn { 
    width: 100%; 
  }
}
</style>
