import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/components/layout/AppLayout.vue'),
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/DashboardView.vue') },
        { path: 'groups', name: 'GroupList', component: () => import('@/views/GroupListView.vue') },
        { path: 'groups/new', name: 'GroupCreate', component: () => import('@/views/GroupFormView.vue') },
        { path: 'groups/:id/edit', name: 'GroupEdit', component: () => import('@/views/GroupFormView.vue') },
        { path: 'keys', name: 'KeyList', component: () => import('@/views/KeyListView.vue') },
        { path: 'logs', name: 'LogList', component: () => import('@/views/LogListView.vue') },
        { path: 'proxy-test', name: 'ProxyTest', component: () => import('@/views/ProxyTestView.vue') },
        { path: 'settings', name: 'Settings', component: () => import('@/views/SettingsView.vue') },
      ],
    },
  ],
})

export default router
