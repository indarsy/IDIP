<template>
  <div class="rtsp-player">
    <div class="video-container" :style="containerStyle">
      <video 
        ref="videoRef"
        :src="streamUrl"
        autoplay
        @error="handleStreamError"
        :style="videoStyle"
      ></video>
      
      <div v-if="!isPlaying" class="placeholder">
        <el-icon :size="32"><VideoCamera /></el-icon>
        <p>{{ placeholderText }}</p>
      </div>
      
      <div class="status-bar">
        <el-tag :type="connectionStatus.type" size="small">
          {{ connectionStatus.text }}
        </el-tag>
      </div>
      
      <div v-if="error" class="error-message">
        <el-alert
          :title="error"
          :description="getErrorDescription(error)"
          type="error"
          show-icon
          closable
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { 
  ref, 
  reactive, 
  computed, 
  onMounted, 
  onBeforeUnmount, 
  watch,
  defineProps,
  defineEmits,
  defineExpose 
} from 'vue'
import { ElMessage } from 'element-plus'
import { VideoCamera } from '@element-plus/icons-vue'
import axios from 'axios'

const API_BASE_URL = '/rtsp-api'

const props = defineProps({
  initialUrl: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['error'])

// 状态变量
const error = ref('')
const isPlaying = ref(false)
const isConnecting = ref(false)
const streamUrl = ref('')
const videoRef = ref(null)
const containerWidth = ref(800)
const videoInfo = reactive({
  width: 0,
  height: 0,
  aspectRatio: 16/9
})

// 添加 handleResize 函数定义
const handleResize = () => {
  if (videoRef.value) {
    const container = videoRef.value.parentElement
    containerWidth.value = container.offsetWidth
  }
}

let statusCheckInterval = null

// 首先定义所有需要的函数
const stopStatusCheck = () => {
  if (statusCheckInterval) {
    clearInterval(statusCheckInterval)
    statusCheckInterval = null
  }
}

const stopPlay = async () => {
  try {
    await axios.get(`${API_BASE_URL}/stream/stop`)
    ElMessage.success('已停止播放')
  } catch (error) {
    ElMessage.error('停止失败：' + error.message)
  } finally {
    isPlaying.value = false
    streamUrl.value = ''
    stopStatusCheck()
  }
}

const fetchVideoInfo = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/stream/info`)
    if (response.data.status === 'success') {
      videoInfo.width = response.data.width
      videoInfo.height = response.data.height
      videoInfo.aspectRatio = response.data.aspect_ratio
    }
  } catch (error) {
    console.error('获取视频信息失败:', error)
  }
}

const startStatusCheck = () => {
  statusCheckInterval = setInterval(async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/stream/status`)
      if (!response.data.is_active && isPlaying.value) {
        ElMessage.error('视频流已断开')
        await stopPlay()
      }
    } catch (error) {
      console.error('状态检查失败:', error)
    }
  }, 5000)
}

const startPlay = async (url = props.initialUrl) => {
  if (!url) {
    error.value = '无可用的视频源'
    return
  }

  try {
    error.value = ''  // 清除之前的错误
    isConnecting.value = true
    const encodedUrl = encodeURIComponent(url.replace('rtsp://', ''))
    
    const response = await axios.get(`${API_BASE_URL}/stream/start/${encodedUrl}`)
    
    if (response.data.status === 'success') {
      await fetchVideoInfo()
      streamUrl.value = `${API_BASE_URL}/stream/feed`
      isPlaying.value = true
      ElMessage.success('视频预览已开始')
      startStatusCheck()
    }
  } catch (err) {
    error.value = err.response?.data?.message || '视频连接失败'
    emit('error', error.value)
  } finally {
    isConnecting.value = false
  }
}

const handleStreamError = async () => {
  if (isPlaying.value) {
    ElMessage.error('视频流加载失败')
    await stopPlay()
  }
}

// 计算属性
const connectionStatus = computed(() => {
  if (isConnecting.value) {
    return { type: 'warning', text: '连接中...' }
  }
  if (isPlaying.value) {
    return { type: 'success', text: '播放中' }
  }
  return { type: 'info', text: '未连接' }
})

const containerStyle = computed(() => {
  if (!isPlaying.value) {
    return {
      minHeight: '400px'
    }
  }
  
  const height = containerWidth.value / videoInfo.aspectRatio
  return {
    width: '100%',
    height: `${height}px`,
    maxHeight: '80vh'
  }
})

const videoStyle = computed(() => {
  if (!isPlaying.value) return {}
  
  return {
    width: '100%',
    height: '100%',
    objectFit: 'contain'
  }
})

const placeholderText = computed(() => {
  if (!props.initialUrl) {
    return '暂无可用的视频源'
  }
  return isConnecting.value ? '正在连接视频源...' : '点击播放按钮开始预览'
})

// 然后是监听器
watch(() => props.initialUrl, async (newUrl) => {
  if (newUrl && !isPlaying.value) {
    await startPlay(newUrl)
  }
}, { immediate: true })

// 生命周期钩子
onMounted(() => {
  if (videoRef.value) {
    const container = videoRef.value.parentElement
    containerWidth.value = container.offsetWidth
  }
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  stopStatusCheck()
  stopPlay()
})

// 暴露方法
defineExpose({
  startPlay,
  stopPlay
})

const getErrorDescription = (error) => {
  if (error.includes('无法连接')) {
    return '请检查摄像头是否在线，或者RTSP地址是否正确'
  }
  if (error.includes('认证失败')) {
    return '请检查用户名和密码是否正确'
  }
  return '请检查网络连接和摄像头状态'
}
</script>

<style scoped>
.rtsp-player {
  width: 100%;
  height: 100%;
  position: relative;
  background-color: #1a1a1a;
  border-radius: 4px;
  overflow: hidden;
}

.video-container {
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

video {
  width: 100%;
  height: 100%;
  object-fit: contain;
  background-color: transparent;
}

.placeholder {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  color: #909399;
}

.status-bar {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 10;
}

.error-message {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80%;
  max-width: 400px;
  z-index: 100;
}
</style> 