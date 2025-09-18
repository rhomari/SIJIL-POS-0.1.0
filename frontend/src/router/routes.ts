import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
  { path: '', name: 'pos', component: () => import('pages/SalesTerminal.vue') },
  { path: 'products', name: 'products', component: () => import('pages/ProductsManagement.vue') },
  { path: 'customers', name: 'customers', component: () => import('pages/CustomersPage.vue') }
    ],
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
