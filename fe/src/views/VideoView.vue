<template>
  <div class="video-view">
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <span>视频管理系统</span>
      </div>
      <video-manager />
    </el-card>
  </div>
</template>

<script>
import VideoManager from '@/views/video/VideoManager.vue'

export default {
  name: 'VideoView',
  
  components: {
    VideoManager
  }
}
</script>

<style scoped>
.video-view {
  padding: 20px;
  background-color: #f0f2f5;
  min-height: 100vh;
}

.box-card {
  min-height: calc(100vh - 40px);
}

.clearfix {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 