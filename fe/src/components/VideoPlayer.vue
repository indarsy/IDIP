<template>
  <div class="video-player">
    <video ref="video" autoplay playsinline controls></video>
  </div>
</template>

<script>

export default {
  name: "VideoPlayer",
  
  props: {
    form: {
      type: Object,
      required: true,
      default: () => ({
        rtsp_url: ''
      })
    }
  },

  data() {
    return {
      pc: null
    }
  },

  watch: {
    'form.rtsp_url': {
      immediate: true,
      async handler(newUrl, oldUrl) {
        if (newUrl !== oldUrl) {
          await this.cleanupConnection()
          if (newUrl) {
            await this.startStream()
          }
        }
      }
    }
  },

  beforeUnmount() {
    this.cleanupConnection()
  },
  
  methods: {
    cleanupConnection() {
      if (this.pc) {
        this.pc.getSenders().forEach(sender => {
          if (sender.track) {
            sender.track.stop()
          }
        })
        
        this.pc.getReceivers().forEach(receiver => {
          if (receiver.track) {
            receiver.track.stop()
          }
        })

        this.pc.close()
        this.pc = null
      }

      if (this.$refs.video) {
        const videoElement = this.$refs.video
        if (videoElement.srcObject) {
          const tracks = videoElement.srcObject.getTracks()
          tracks.forEach(track => track.stop())
          videoElement.srcObject = null
        }
      }
    },

    async startStream() {
      console.log('startStream', this.form.rtsp_url)
      if (!this.form.rtsp_url) {
        return
      }

      try {
        await this.cleanupConnection()
        
        await new Promise(resolve => setTimeout(resolve, 100))

        this.pc = new RTCPeerConnection()

        this.pc.ontrack = (event) => {
          const videoElement = this.$refs.video
          videoElement.srcObject = event.streams[0]
        }

        this.pc.addTransceiver('video', {
          direction: 'recvonly'
        })

        const offer = await this.pc.createOffer({
          offerToReceiveVideo: true,
          offerToReceiveAudio: false
        })
        
        await this.pc.setLocalDescription(offer)

        const response = await fetch("http://localhost:5000/offer", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            sdp: this.pc.localDescription.sdp,
            type: this.pc.localDescription.type,
            rtsp_url: this.form.rtsp_url,
          }),
        })

        if (!response.ok) {
          const errorData = await response.json()
          throw new Error(errorData.error || '服务器错误')
        }

        const data = await response.json()
        await this.pc.setRemoteDescription(data)

      } catch (error) {
        console.error("视频流启动失败:", error)
        this.$emit('error', error)
        this.cleanupConnection()
      }
    }
  }
}
</script>

<style scoped>
.video-player {
  width: 100%;
  height: 100%;
  background: #000;
  border-radius: 4px;
  overflow: hidden;
}

video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>
