<template>
	<q-page class="q-pa-md">
		<div class="column full-height">
			<div class="row items-center justify-between q-mb-sm">
				<div class="row items-center q-gutter-x-sm" style="margin-left:10px">
					<div class="text-subtitle1">
						{{ $t('pm.products') }}
					</div>
					<q-badge v-if="displayedProducts.length" color="primary" align="middle" class="q-ml-md" >{{ displayedProducts.length }}</q-badge>
				</div>
				<div class="row items-center no-wrap q-gutter-x-xs">
					<q-toggle v-model="onlyActive" :label="$t('pm.activeOnly')" dense />
					<q-btn dense flat round icon="refresh" @click="refreshAll" />
					<q-btn dense flat color="primary" icon="add" :label="$t('pm.addCategory')" @click="openCategoryDialog()" />
					<q-btn dense flat color="primary" icon="add" :label="$t('pm.addProduct')" @click="openProductDialog()" />
				</div>
			</div>

			<div class="row q-col-gutter-sm q-mb-sm">
				<div class="col-12 col-md-6">
					<q-input v-model="productSearch" dense outlined clearable :placeholder="$t('pm.searchProducts')">
						<template #prepend><q-icon name="search" /></template>
					</q-input>
				</div>
				<div class="col-12 col-md-6">
					<q-select
						v-model="categoryFilterId"
						:options="categorySelectOptions"
						dense outlined clearable
						map-options emit-value
						:label="$t('pm.filterByCategory')"
					>
						<template #option="scope">
							<q-item v-bind="scope.itemProps">
								<q-item-section>
									<q-item-label>{{ scope.opt.label }}</q-item-label>
								</q-item-section>
								<q-item-section side>
									<div class="row no-wrap items-center q-gutter-x-xs">
										<q-btn dense round flat size="sm" icon="edit" @click.stop.prevent="onEditCategoryOption(scope.opt.value)">
											<q-tooltip>{{ $t('pm.edit') }}</q-tooltip>
										</q-btn>
										<q-btn dense round flat size="sm" color="negative" icon="delete" @click.stop.prevent="onDeleteCategoryOption(scope.opt.value)">
											<q-tooltip>{{ $t('pm.delete') }}</q-tooltip>
										</q-btn>
									</div>
								</q-item-section>
							</q-item>
						</template>
					</q-select>
				</div>
			</div>

			<q-table
				class="col"
				flat
				:rows="displayedProducts"
				:columns="productColumns"
				row-key="product_id"
				:pagination="pagination"
				@update:pagination="val => pagination = val"
			>
				<template #body-cell-name="props">
					<q-td :props="props">
						<div class="column">
							<div class="text-weight-medium">{{ props.row.name }}</div>
							<div class="text-grey-7 text-caption ellipsis-2-lines">{{ props.row.description }}</div>
						</div>
					</q-td>
				</template>

				<template #body-cell-category="props">
					<q-td :props="props">
						{{ categoryLabelById(props.row.category_id) || '-' }}
					</q-td>
				</template>

				<template #body-cell-base_price="props">
					<q-td :props="props">
						{{ formatCurrency(props.row.base_price) }}
					</q-td>
				</template>

				<template #body-cell-cost_price="props">
					<q-td :props="props">
						{{ formatCurrency(props.row.cost_price) }}
					</q-td>
				</template>

				<template #body-cell-stock_quantity="props">
					<q-td :props="props">
						<q-badge :color="props.row.stock_quantity > props.row.min_stock ? 'positive' : 'warning'">
							{{ props.row.stock_quantity }}
						</q-badge>
					</q-td>
				</template>

				<template #body-cell-is_active="props">
					<q-td :props="props">
						<q-toggle v-model="props.row.is_active" dense @update:model-value="onToggleActive(props.row)" />
					</q-td>
				</template>

				<template #body-cell-actions="props">
					<q-td :props="props" class="q-gutter-x-xs">
						<q-btn dense round flat size="sm" icon="edit" @click="openProductDialog(props.row)">
							<q-tooltip>{{ $t('pm.edit') }}</q-tooltip>
						</q-btn>
						<q-btn dense round flat size="sm" icon="content_copy" @click="duplicateProduct(props.row)">
							<q-tooltip>{{ $t('pm.duplicate') }}</q-tooltip>
						</q-btn>
						<q-btn dense round flat size="sm" color="negative" icon="delete" @click="confirmDeleteProduct(props.row)">
							<q-tooltip>{{ $t('pm.delete') }}</q-tooltip>
						</q-btn>
					</q-td>
				</template>
			</q-table>
		</div>

		<!-- Category dialog -->
		<q-dialog v-model="categoryDialog.open" persistent>
			<q-card style="min-width: 420px; max-width: 92vw;">
				<q-card-section class="row items-center justify-between">
					<div class="text-h6">{{ categoryDialog.model?.category_id ? $t('pm.category.editTitle') : $t('pm.category.addTitle') }}</div>
					<q-btn flat round dense icon="close" v-close-popup />
				</q-card-section>

				<q-separator />

				<q-card-section>
					<q-form @submit.prevent="saveCategory">
						<div class="row q-col-gutter-md">
							<div class="col-12">
								<q-input v-model="categoryDialog.model!.name" :label="$t('pm.category.name')" dense outlined :rules="[v => !!v || $t('pm.required')]" />
							</div>
						</div>
						<div class="row justify-end q-gutter-sm q-mt-md">
							<q-btn flat :label="$t('cancel')" v-close-popup />
							<q-btn color="primary" :label="$t('save')" type="submit" />
						</div>
					</q-form>
				</q-card-section>
			</q-card>
		</q-dialog>

		<!-- Product dialog -->
		<q-dialog v-model="productDialog.open" persistent>
			<q-card style="min-width: 720px; max-width: 96vw;">
				<q-card-section class="row items-center justify-between">
					<div class="text-h6">{{ productDialog.model?.product_id ? $t('pm.product.editTitle') : $t('pm.product.addTitle') }}</div>
					<q-btn flat round dense icon="close" v-close-popup />
				</q-card-section>

				<q-separator />

				<q-card-section>
					<q-form @submit.prevent="saveProduct">
						<div class="row q-col-gutter-md">
							<div class="col-12 col-md-4">
								<q-input v-model="productDialog.model!.sku" :label="$t('pm.product.sku')" dense outlined :rules="[v => !!v || $t('pm.required')]" />
							</div>
							<div class="col-12 col-md-8">
								<q-input v-model="productDialog.model!.name" :label="$t('pm.product.name')" dense outlined :rules="[v => !!v || $t('pm.required')]" />
							</div>

							<div class="col-12">
								<q-input v-model="productDialog.model!.description" type="textarea" autogrow :label="$t('pm.product.description')" dense outlined />
							</div>

							<div class="col-12 col-md-6">
								<q-select
									v-model="productDialog.model!.category_id"
									:options="categorySelectOptions"
									:label="$t('pm.product.category')"
									dense outlined clearable
									map-options emit-value
								/>
							</div>

							<div class="col-6 col-md-3">
								<q-input v-model.number="productDialog.model!.base_price" type="number" step="0.01" :label="$t('pm.product.basePrice')" dense outlined />
							</div>
							<div class="col-6 col-md-3">
								<q-input v-model.number="productDialog.model!.cost_price" type="number" step="0.01" :label="$t('pm.product.costPrice')" dense outlined />
							</div>

							<div class="col-6 col-md-3">
								<q-input v-model.number="productDialog.model!.stock_quantity" type="number" step="1" :label="$t('pm.product.stock')" dense outlined />
							</div>
							<div class="col-6 col-md-3">
								<q-input v-model.number="productDialog.model!.min_stock" type="number" step="1" :label="$t('pm.product.minStock')" dense outlined />
							</div>

							<div class="col-6 col-md-3">
								<q-input v-model="productDialog.model!.unit" :label="$t('pm.product.unit')" dense outlined />
							</div>
							<div class="col-6 col-md-3 flex items-center">
								<q-toggle v-model="productDialog.model!.is_active" :label="$t('pm.product.active')" dense />
							</div>

							<div class="col-12 col-md-6">
								<q-select
									v-model="barcodesEdit"
									use-input use-chips multiple new-value-mode="add-unique"
									input-debounce="0"
									:label="$t('pm.product.barcodes')"
									dense outlined
									:options="[]"
								/>
								<div class="text-caption text-grey-7 q-mt-xs">{{ $t('pm.product.barcodesHint') }}</div>
							</div>

							<div class="col-12">
								<q-expansion-item dense expand-separator icon="tune" :label="$t('pm.advancedJson')">
									<div class="row q-col-gutter-md q-mt-sm">
										<div class="col-12 col-md-4">
											<q-input v-model="pricesJson" type="textarea" autogrow :label="$t('pm.pricesJson')" dense outlined />
										</div>
										<div class="col-12 col-md-4">
											<q-input v-model="discountPolicyJson" type="textarea" autogrow :label="$t('pm.discountPolicyJson')" dense outlined />
										</div>
										<div class="col-12 col-md-4">
											<q-input v-model="merchJson" type="textarea" autogrow :label="$t('pm.merchJson')" dense outlined />
										</div>
									</div>
								</q-expansion-item>
							</div>
						</div>

						<div class="row justify-end q-gutter-sm q-mt-md">
							<q-btn flat :label="$t('cancel')" v-close-popup />
							<q-btn color="primary" :label="$t('save')" type="submit" />
						</div>
					</q-form>
				</q-card-section>
			</q-card>
		</q-dialog>

		<!-- Confirm dialogs -->
		<q-dialog v-model="confirm.open">
			<q-card>
				<q-card-section class="row items-center q-gutter-sm">
					<q-icon name="warning" color="warning" />
					<div class="text-body1">{{ confirm.message }}</div>
				</q-card-section>
				<q-card-actions align="right">
					<q-btn flat :label="$t('cancel')" v-close-popup />
					<q-btn color="negative" :label="$t('pm.delete')" @click="confirm.onConfirm?.()" v-close-popup />
				</q-card-actions>
			</q-card>
		</q-dialog>
	</q-page>
	</template>

	<script lang="ts">
		import { defineComponent, ref, computed, onMounted } from 'vue';
		import { useI18n } from 'vue-i18n';
		import type { QTableProps } from 'quasar';
  
	type Nullable<T> = T | null;
  
	interface Category {
		category_id: number;
		name: string;
		created_at?: string;
	}
  
	interface Product {
		product_id: number;
		sku: string;
		name: string;
		description: string | null;
		category_id: number | null;
		brand_id: number | null;
		base_price: number;
		cost_price: number;
		stock_quantity: number;
		min_stock: number;
		unit: string;
		prices: unknown;
		price_template_id: number | null;
		discount_policy: unknown;
		merch: unknown;
		barcodes: string[];
		is_active: boolean;
		created_at?: string;
		updated_at?: string;
		sku_generated?: string | null;
	}
  
	export default defineComponent({
		name: 'ProductsManagement',
		setup() {
            const { t } = useI18n();
			const categories = ref<Category[]>([]);
  			const products = ref<Product[]>([
				{ product_id: 1, sku: 'PEN-BLUE-001', name: 'Blue Ballpoint Pen', description: 'Smooth writing pen', category_id: 2, brand_id: null, base_price: 1.5, cost_price: 0.6, stock_quantity: 120, min_stock: 20, unit: 'pcs', prices: {}, price_template_id: null, discount_policy: {}, merch: {}, barcodes: ['0123456789012'], is_active: true },
				{ product_id: 2, sku: 'NOTE-A5-001', name: 'A5 Notebook', description: 'Ruled, 100 pages', category_id: 3, brand_id: null, base_price: 3.99, cost_price: 1.5, stock_quantity: 80, min_stock: 15, unit: 'pcs', prices: {}, price_template_id: null, discount_policy: {}, merch: {}, barcodes: ['0123456789013'], is_active: true },
				{ product_id: 3, sku: 'BOOK-NOVEL-001', name: 'Mystery Novel', description: 'Bestseller', category_id: 4, brand_id: null, base_price: 12.5, cost_price: 6, stock_quantity: 30, min_stock: 5, unit: 'pcs', prices: {}, price_template_id: null, discount_policy: {}, merch: {}, barcodes: ['0123456789014'], is_active: true }
			]);
  			const categoryMap = computed(() => {
				const map = new Map<number, Category>();
				categories.value.forEach(c => map.set(c.category_id, c));
				return map;
			});
  


			const productSearch = ref('');
			const onlyActive = ref(true);
			const categoryFilterId = ref<number | null>(null);
  
			const productColumns = computed<QTableProps['columns']>(() => [
				{ name: 'sku', label: t('pm.product.sku'), field: 'sku', align: 'left', sortable: true },
				{ name: 'name', label: t('pm.product.name'), field: 'name', align: 'left', sortable: true },
				{ name: 'category', label: t('pm.product.category'), field: 'category_id', align: 'left', sortable: false },
				{ name: 'base_price', label: t('pm.product.basePrice'), field: 'base_price', align: 'right', sortable: true },
				{ name: 'cost_price', label: t('pm.product.costPrice'), field: 'cost_price', align: 'right', sortable: true },
				{ name: 'stock_quantity', label: t('pm.product.stock'), field: 'stock_quantity', align: 'right', sortable: true },
				{ name: 'is_active', label: t('pm.product.active'), field: 'is_active', align: 'center' },
				{ name: 'actions', label: '', field: 'product_id', align: 'right' }
			]);
  
			const pagination = ref({ page: 1, rowsPerPage: 10 });
  
			const categorySelectOptions = computed(() =>
				categories.value
					.slice()
					.sort((a, b) => a.name.localeCompare(b.name))
					.map(c => ({ label: c.name, value: c.category_id }))
			);
  
			const displayedProducts = computed(() => {
				const term = productSearch.value.trim().toLowerCase();
				const catFilter = categoryFilterId.value;
  
				return products.value.filter(p => {
					if (onlyActive.value && !p.is_active) return false;
					if (catFilter && p.category_id !== catFilter) return false;
					if (!term) return true;
					const hay = `${p.sku} ${p.name} ${p.description || ''}`.toLowerCase();
					return hay.includes(term);
				});
			});
  
			const categoryDialog = ref<{
				open: boolean;
				model: Nullable<Category>;
				mode: 'create' | 'edit';
			}>({ open: false, model: null, mode: 'create' });
  
			function openCategoryDialog(seed?: Partial<Category>) {
				const model: Category = {
					category_id: seed?.category_id ?? 0,
					name: seed?.name ?? ''
				};
				categoryDialog.value = {
					open: true,
					model,
					mode: seed && seed.category_id ? 'edit' : 'create'
				};
			}
  
			function saveCategory() {
				const model = categoryDialog.value.model!;
				if (!model.name.trim()) return;
  
						if (categoryDialog.value.mode === 'edit' && model.category_id) {
							const idx = categories.value.findIndex(c => c.category_id === model.category_id);
										if (idx !== -1) {
											const existing = categories.value[idx];
											if (existing) {
												// mutate in-place to keep required properties intact
												existing.name = model.name;
												// no parent field to update
											}
										}
						} else {
					const newId = Math.max(0, ...categories.value.map(c => c.category_id)) + 1;
					categories.value.push({ ...model, category_id: newId });
				}
				categoryDialog.value.open = false;
			}
  
			function confirmDeleteCategory(cat: Category) {
				const message = t('pm.confirmDeleteCategory.simple');
				confirm.value = {
					open: true,
					message,
					onConfirm: () => {
			categories.value = categories.value.filter(c => c.category_id !== cat.category_id);
			// Unassign products for this category
						products.value = products.value.map(p => (p.category_id === cat.category_id ? { ...p, category_id: null } : p));
						if (categoryFilterId.value === cat.category_id) categoryFilterId.value = null;
					}
				};
			}

			const productDialog = ref<{
				open: boolean;
				model: Nullable<Product>;
				mode: 'create' | 'edit';
			}>({ open: false, model: null, mode: 'create' });
  
			const barcodesEdit = ref<string[]>([]);
			const pricesJson = ref<string>('{}');
			const discountPolicyJson = ref<string>('{}');
			const merchJson = ref<string>('{}');
  
			function openProductDialog(seed?: Partial<Product>) {
				const model: Product = {
					product_id: seed?.product_id ?? 0,
					sku: seed?.sku ?? '',
					name: seed?.name ?? '',
					description: seed?.description ?? '',
					category_id: seed?.category_id ?? categoryFilterId.value ?? null,
					brand_id: seed?.brand_id ?? null,
					base_price: seed?.base_price ?? 0,
					cost_price: seed?.cost_price ?? 0,
					stock_quantity: seed?.stock_quantity ?? 0,
					min_stock: seed?.min_stock ?? 0,
					unit: seed?.unit ?? 'pcs',
					prices: seed?.prices ?? {},
					price_template_id: seed?.price_template_id ?? null,
					discount_policy: seed?.discount_policy ?? {},
					merch: seed?.merch ?? {},
					barcodes: seed?.barcodes ?? [],
					is_active: seed?.is_active ?? true
				};
				productDialog.value = {
					open: true,
					model,
					mode: seed && seed.product_id ? 'edit' : 'create'
				};
  
				barcodesEdit.value = [...model.barcodes];
				pricesJson.value = safeJSONStringify(model.prices);
				discountPolicyJson.value = safeJSONStringify(model.discount_policy);
				merchJson.value = safeJSONStringify(model.merch);
			}
  
			function parseJSONOrDefault<T>(text: string, fallback: T): T {
				try {
					const v = JSON.parse(text);
					return v as T;
				} catch {
					return fallback;
				}
			}
			function safeJSONStringify(v: unknown): string {
				try { return JSON.stringify(v ?? {}, null, 2); } catch { return '{}'; }
			}
			function saveProduct() {
				const model = productDialog.value.model!;
				// merge editors back
				const merged: Product = {
					...model,
					barcodes: [...barcodesEdit.value],
					prices: parseJSONOrDefault(pricesJson.value, {}),
					discount_policy: parseJSONOrDefault(discountPolicyJson.value, {}),
					merch: parseJSONOrDefault(merchJson.value, {})
				};
				if (productDialog.value.mode === 'edit' && merged.product_id) {
					const idx = products.value.findIndex(p => p.product_id === merged.product_id);
					if (idx !== -1) products.value[idx] = { ...merged };
				} else {
					const newId = Math.max(0, ...products.value.map(p => p.product_id)) + 1;
					products.value.push({ ...merged, product_id: newId });
				}
				productDialog.value.open = false;
			}
  
			function duplicateProduct(row: Product) {
				const newId = Math.max(0, ...products.value.map(p => p.product_id)) + 1;
				const copy: Product = {
					...row,
					product_id: newId,
					sku: `${row.sku}-COPY`,
				  name: `${row.name} (Copy)`,
					sku_generated: null
				};
				products.value.push(copy);
			}
  
	    function confirmDeleteProduct(row: Product) {
				confirm.value = {
					open: true,
					message: t('pm.confirmDeleteProduct', { name: row.name }),
					onConfirm: () => {
						products.value = products.value.filter(p => p.product_id !== row.product_id);
					}
				};
			}
  
			function onToggleActive(row: Product) {
				// Hook for API update
				void row;
			}
  
			function refreshAll() {
				getCategories();
			}

			// Category option actions from the filter combobox
			function onEditCategoryOption(id: number) {
				const cat = categoryMap.value.get(id);
				if (cat) openCategoryDialog(cat);
			}
			function onDeleteCategoryOption(id: number) {
				const cat = categoryMap.value.get(id);
				if (cat) confirmDeleteCategory(cat);
			}
  
			// Helpers
			function categoryLabelById(id: number | null): string | null {
				if (id == null) return null;
				const c = categoryMap.value.get(id);
				return c ? c.name : null;
			}
			function formatCurrency(n: number): string {
				return `$${(n ?? 0).toFixed(2)}`;
			}
  
			// Confirm dialog
			const confirm = ref<{ open: boolean; message: string; onConfirm?: () => void }>({
				open: false,
				message: ''
			});
			function getCategories() {
				fetch('http://192.168.1.2:8080/api/categories', { method: 'GET' })
					.then(res => {
						if (!res.ok) throw new Error(`HTTP ${res.status}`);
						return res.json();
					})
					.then((data: Array<{ id: number; name: string; created_at?: string }>) => {
						// Map backend shape (id,name,created_at) -> frontend Category
						categories.value = data.map(item => ({
							category_id: item.id,
							name: item.name,
							...(item.created_at ? { created_at: item.created_at } : {})
						}));
						// Reset filter if it no longer exists
						if (categoryFilterId.value && !categories.value.some(c => c.category_id === categoryFilterId.value)) {
							categoryFilterId.value = null;
						}
					})
					.catch(error => {
						console.error('Error fetching categories:', error);
					});
			}
            onMounted(() => {
                getCategories();
            });
  
			return {
					t,
  
				categories,
				products,
  
				// categories
				categorySelectOptions,
                getCategories,
				onEditCategoryOption,
				onDeleteCategoryOption,
				// products table
				productSearch,
				categoryFilterId,
				onlyActive,
				displayedProducts,
				productColumns,
				pagination,
  
				// dialogs and actions
				categoryDialog,
				productDialog,
				openCategoryDialog,
				saveCategory,
				confirmDeleteCategory,
  
				openProductDialog,
				saveProduct,
				duplicateProduct,
				confirmDeleteProduct,
				onToggleActive,
				refreshAll,
  
				// product editors
				barcodesEdit,
				pricesJson,
				discountPolicyJson,
				merchJson,
  
				// helpers
				categoryLabelById,
				formatCurrency,
  
				// confirm
				confirm
			};
		}
	});
	</script>
  
	<style scoped>
	.full-height {
		height: 100%;
	}
  
	.ellipsis-2-lines {
		display: -webkit-box;
		line-clamp: 2;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
	</style>
