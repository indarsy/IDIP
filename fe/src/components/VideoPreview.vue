<template>
  <div class="video-preview">
    <div class="video-container">
      <video 
        ref="videoRef"
        autoplay 
        playsinline 
        muted
        controls
        style="width: 100%; height: 100%; object-fit: contain;"
      ></video>
      <div v-if="error" class="error-overlay">
        <i class="el-icon-video-camera-solid"></i>
        <span>{{ error }}</span>
        <el-button size="small" type="primary" @click="retryPlay">重试</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { startWebRTC } from '@/api/webrtc'

export default {
  name: 'VideoPreview',
  props: {
    rtspUrl: {
      type: String,
      required: true
    }
  },
  emits: ['error'],
  setup(props, { emit }) {
    const videoRef = ref(null)
    const error = ref('')
    const isPlaying = ref(false)
    let peerConnection = null
    let mediaStream = null

    const cleanup = () => {
      if (mediaStream) {
        mediaStream.getTracks().forEach(track => {
          track.stop()
          mediaStream.removeTrack(track)
        })
        mediaStream = null
      }
      
      if (peerConnection) {
        peerConnection.close()
        peerConnection = null
      }

      if (videoRef.value) {
        videoRef.value.srcObject = null
      }
      
      isPlaying.value = false
    }

    const startPlay = async (url) => {
      if (!url) return
      
      try {
        error.value = ''
        cleanup()
        
        // 创建新的 MediaStream
        mediaStream = new MediaStream()
        
        // 创建 RTCPeerConnection
        peerConnection = new RTCPeerConnection({
          iceServers: [
            { urls: 'stun:stun.l.google.com:19302' }
          ]
        })

        // 处理远程视频流
        peerConnection.ontrack = (event) => {
          const track = event.track
          if (track.kind === 'video') {
            mediaStream.addTrack(track)
            
            if (videoRef.value) {
              videoRef.value.srcObject = mediaStream
              videoRef.value.play()
                .then(() => {
                  isPlaying.value = true
                })
                .catch(err => {
                  console.error('播放失败:', err)
                  error.value = '视频播放失败'
                })
            }
          }
        }

        // 监听连接状态
        peerConnection.onconnectionstatechange = () => {
          if (peerConnection.connectionState === 'failed') {
            error.value = '视频连接失败'
            emit('error', new Error('连接失败'))
          }
        }

        // 创建 offer
        const offer = await peerConnection.createOffer({
          offerToReceiveVideo: true,
          offerToReceiveAudio: false
        })
        
        await peerConnection.setLocalDescription(offer)

        // 等待 ICE 收集完成
        await new Promise(resolve => {
          if (peerConnection.iceGatheringState === 'complete') {
            resolve()
          } else {
            peerConnection.onicegatheringstatechange = () => {
              if (peerConnection.iceGatheringState === 'complete') {
                resolve()
              }
            }
          }
        })

        // 发送 offer 到后端
        const response = await startWebRTC({
          rtspUrl: url,
          sdp: peerConnection.localDescription.sdp
        })

        if (!response.data?.sdp) {
          throw new Error('服务器返回数据格式错误')
        }

        // 设置远程描述
        await peerConnection.setRemoteDescription(
          new RTCSessionDescription({
            type: 'answer',
            sdp: response.data.sdp
          })
        )

      } catch (err) {
        console.error('播放失败:', err)
        error.value = err.message || '视频加载失败'
        emit('error', err)
        cleanup()
      }
    }

    const retryPlay = () => {
      if (props.rtspUrl) {
        startPlay(props.rtspUrl)
      }
    }

    // 监听 rtspUrl 变化
    watch(() => props.rtspUrl, (newUrl) => {
      if (newUrl) {
        startPlay(newUrl)
      } else {
        cleanup()
      }
    })

    onMounted(() => {
      if (props.rtspUrl) {
        startPlay(props.rtspUrl)
      }
    })

    onBeforeUnmount(() => {
      cleanup()
    })

    return {
      videoRef,
      error,
      isPlaying,
      retryPlay
    }
  }
}
</script>

<style scoped>
.video-preview {
  width: 100%;
  height: 100%;
  background: #000;
  position: relative;
}

.video-container {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
}

.error-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  gap: 1rem;
}

.error-overlay i {
  font-size: 2rem;
  color: #f56c6c;
}
</style>
