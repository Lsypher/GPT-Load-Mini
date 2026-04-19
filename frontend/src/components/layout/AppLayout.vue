<template>
  <div class="app-layout">
    <!-- Sidebar -->
    <aside class="sidebar" :class="{ 'sidebar--open': sidebarOpen }">
      <div class="sidebar-header">
        <div class="logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="3" y="3" width="18" height="18" rx="3" stroke="currentColor" stroke-width="2"/>
            <path d="M8 12h8M12 8v8" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <span class="logo-text">GPT-Load Mini</span>
        </div>
      </div>

      <nav class="sidebar-nav">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ 'nav-item--active': isActive(item.path) }"
          @click="closeMobileSidebar"
        >
          <component :is="item.icon" class="nav-icon" />
          <span class="nav-label">{{ item.label }}</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="theme-indicator">
          <span class="theme-dot"></span>
          <span class="theme-text">{{ theme === 'dark' ? '深色模式' : '浅色模式' }}</span>
        </div>
      </div>
    </aside>

    <!-- Overlay for mobile -->
    <div
      v-if="sidebarOpen"
      class="sidebar-overlay"
      @click="closeMobileSidebar"
    ></div>

    <!-- Main content -->
    <div class="main-wrapper">
      <!-- Header -->
      <header class="header">
        <button class="menu-toggle show-mobile-only" @click="sidebarOpen = !sidebarOpen">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 12h18M3 6h18M3 18h18" stroke-linecap="round"/>
          </svg>
        </button>

        <h1 class="page-title">{{ pageTitle }}</h1>

        <div class="header-actions">
          <button class="btn btn-ghost btn-icon" @click="toggleTheme" :title="theme === 'dark' ? '切换到浅色模式' : '切换到深色模式'">
            <svg v-if="theme === 'dark'" class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5"/>
              <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
            </svg>
            <svg v-else class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
            </svg>
          </button>
        </div>
      </header>

      <!-- Content -->
      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  DataAnalysis,
  FolderOpened,
  Key,
  Document,
  Setting,
  Connection
} from '@element-plus/icons-vue'

const route = useRoute()
const sidebarOpen = ref(false)
const theme = ref<'dark' | 'light'>('dark')

const navItems = [
  { path: '/dashboard', label: '仪表盘', icon: DataAnalysis },
  { path: '/groups', label: '分组', icon: FolderOpened },
  { path: '/keys', label: 'API 密钥', icon: Key },
  { path: '/logs', label: '日志', icon: Document },
  { path: '/proxy-test', label: '代理测试', icon: Connection },
  { path: '/settings', label: '设置', icon: Setting },
]

const pageTitle = computed(() => {
  const map: Record<string, string> = {
    '/dashboard': '仪表盘',
    '/groups': '分组管理',
    '/keys': 'API 密钥管理',
    '/logs': '请求日志',
    '/proxy-test': '代理测试',
    '/settings': '设置',
  }
  if (route.path.startsWith('/groups/new')) return '创建分组'
  if (route.path.includes('/edit')) return '编辑分组'
  return map[route.path] || '仪表盘'
})

function isActive(path: string) {
  return route.path === path || (path !== '/dashboard' && route.path.startsWith(path))
}

function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  localStorage.setItem('theme', theme.value)
  window.__setTheme?.(theme.value)
}

function closeMobileSidebar() {
  sidebarOpen.value = false
}

function handleResize() {
  if (window.innerWidth > 768) {
    sidebarOpen.value = false
  }
}

onMounted(() => {
  const saved = localStorage.getItem('theme') as 'dark' | 'light' | null
  if (saved) {
    theme.value = saved
  } else if (window.matchMedia('(prefers-color-scheme: light)').matches) {
    theme.value = 'light'
  }
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  background: var(--bg-primary);
}

/* ==================== Sidebar ==================== */
.sidebar {
  width: var(--sidebar-width);
  background: var(--bg-secondary);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 100;
  transition: transform var(--transition-base);
}

.sidebar-header {
  padding: var(--space-4) var(--space-5);
  border-bottom: 1px solid var(--border-subtle);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.logo-icon {
  width: 28px;
  height: 28px;
  color: var(--accent);
}

.logo-text {
  font-size: 1.0625rem;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}

.sidebar-nav {
  flex: 1;
  padding: var(--space-4) var(--space-3);
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9375rem;
  font-weight: 500;
  transition: all var(--transition-fast);
}

.nav-item:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.nav-item--active {
  background: var(--accent-muted);
  color: var(--accent);
}

.nav-item--active:hover {
  background: var(--accent-muted);
  color: var(--accent);
}

.nav-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.nav-label {
  white-space: nowrap;
}

.sidebar-footer {
  padding: var(--space-4) var(--space-5);
  border-top: 1px solid var(--border-subtle);
}

.theme-indicator {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: 0.8125rem;
  color: var(--text-tertiary);
}

.theme-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
}

/* ==================== Overlay ==================== */
.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 99;
  display: none;
}

/* ==================== Main Wrapper ==================== */
.main-wrapper {
  flex: 1;
  margin-left: var(--sidebar-width);
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  transition: margin-left var(--transition-base);
}

/* ==================== Header ==================== */
.header {
  height: var(--header-height);
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  padding: 0 var(--space-6);
  gap: var(--space-4);
  position: sticky;
  top: 0;
  z-index: 50;
}

.menu-toggle {
  display: none;
  width: 36px;
  height: 36px;
  padding: var(--space-2);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  cursor: pointer;
}

.menu-toggle:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.menu-toggle svg {
  width: 100%;
  height: 100%;
}

.page-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.header-actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.icon {
  width: 20px;
  height: 20px;
}

/* ==================== Content ==================== */
.content {
  flex: 1;
  padding: var(--space-6);
  background: var(--bg-primary);
}

/* ==================== Responsive ==================== */
@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar--open {
    transform: translateX(0);
  }

  .sidebar-overlay {
    display: block;
  }

  .main-wrapper {
    margin-left: 0;
  }

  .menu-toggle {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .content {
    padding: var(--space-4);
  }

  .header {
    padding: 0 var(--space-4);
  }
}
</style>
