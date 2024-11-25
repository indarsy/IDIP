<template>
  <el-config-provider>
    <div class="app-container">
      <!-- 顶部导航 -->
      <el-header class="app-header">
        <div class="header-content">
          <div class="logo">
            <h2>视频管理系统</h2>
          </div>
          <div class="header-right">
            <el-space>
              <el-button type="primary" @click="refreshPage">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </el-space>
          </div>
        </div>
      </el-header>

      <!-- 主要内容区 -->
      <el-container class="main-container">
        <!-- 侧边栏 -->
        <el-aside width="200px" class="app-aside">
          <el-menu
            :default-active="activeMenu"
            class="app-menu"
            router
          >
            <!-- 视频管理子菜单 -->
            <el-sub-menu index="/video">
              <template #title>
                <el-icon><VideoCamera /></el-icon>
                <span>视频管理</span>
              </template>
              
              <el-menu-item index="/">
                <el-icon><List /></el-icon>
                <span>视频列表</span>
              </el-menu-item>
              
              <el-menu-item index="/video/capture">
                <el-icon><Camera /></el-icon>
                <span>视频采集</span>
              </el-menu-item>
            </el-sub-menu>
            <el-menu-item index="/workshop">
              <el-icon><House /></el-icon>
              <span>车间管理</span>
            </el-menu-item>
          </el-menu>
        </el-aside>

        <!-- 内容区 -->
        <el-main class="app-main">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </el-main>
      </el-container>
    </div>
  </el-config-provider>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
//import { ElMessage } from 'element-plus'
import { Refresh, VideoCamera, House } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

// 当前激活的菜单项
const activeMenu = computed(() => route.path)

// 刷新页面
const refreshPage = () => {
  const { fullPath } = route
  router.replace(fullPath)
}
</script>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0;
}

.header-content {
  height: 60px;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo h2 {
  margin: 0;
  font-size: 20px;
  color: #409EFF;
}

.main-container {
  flex: 1;
  overflow: hidden;
}

.app-aside {
  background-color: #fff;
  border-right: 1px solid #dcdfe6;
}

.app-menu {
  height: 100%;
  border-right: none;
}

.app-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

/* 路由过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式布局 */
@media (max-width: 768px) {
  .app-aside {
    width: 64px !important;
  }
  
  .el-menu {
    width: 64px;
  }
  
  .el-menu-item {
    padding: 0 !important;
    text-align: center;
  }
  
  .el-menu-item span {
    display: none;
  }
}
</style>