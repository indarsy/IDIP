<template>
  <div class="video-player">
    <video ref="videoPlayer" class="video-js vjs-default-skin" controls preload="auto" width="100%" height="auto">
      <source :src="src" :type="videoType">
      您的浏览器不支持视频播放
    </video>
  </div>
</template>

<script>
import videojs from 'video.js'
import 'video.js/dist/video-js.css'

export default {
  name: 'VideoPlayer',

  props: {
    src: {
      type: String,
      required: true
    },
    options: {
      type: Object,
      default: () => ({})
    }
  },

  data() {
    return {
      player: null,
      defaultOptions: {
        autoplay: false,
        controls: true,
        fluid: true,
        playbackRates: [0.5, 1, 1.5, 2],
        controlBar: {
          children: [
            'playToggle',
            'volumePanel',
            'currentTimeDisplay',
            'timeDivider',
            'durationDisplay',
            'progressControl',
            'playbackRateMenuButton',
            'fullscreenToggle'
          ]
        }
      }
    }
  },

  computed: {
    videoType() {
      const ext = this.src.split('.').pop().toLowerCase()
      switch (ext) {
        case 'mp4':
          return 'video/mp4'
        case 'webm':
          return 'video/webm'
        case 'm3u8':
          return 'application/x-mpegURL'
        default:
          return 'video/mp4'
      }
    }
  },

  mounted() {
    this.initPlayer()
  },

  beforeUnmount() {
    this.disposePlayer()
  },

  watch: {
    src: {
      handler(newVal) {
        if (this.player && newVal) {
          this.player.src({ src: newVal, type: this.videoType })
        }
      },
      immediate: true
    },
    options: {
      handler(newVal) {
        if (this.player) {
          this.player.options(newVal)
        }
      },
      deep: true
    }
  },

  methods: {
    initPlayer() {
      const options = {
        ...this.defaultOptions,
        ...this.options
      }

      this.player = videojs(this.$refs.videoPlayer, options, () => {
        this.player.src({
          src: this.src,
          type: this.videoType
        })

        // 添加错误处理
        this.player.on('error', () => {
          const error = this.player.error()
          console.error('视频播放错误:', error)
          this.$emit('error', error)
        })

        // 添加加载事件处理
        this.player.on('loadstart', () => {
          console.log('开始加载视频:', this.src)
        })
      })
    },

    handleError() {
      const error = this.player.error()
      console.error('视频错误:', error)
      this.$emit('error', error)
    },

    disposePlayer() {
      if (this.player) {
        this.player.dispose()
        this.player = null
      }
    },


    play() {
      if (this.player) {
        this.player.play()
      }
    },

    pause() {
      if (this.player) {
        this.player.pause()
      }
    },

    reset() {
      if (this.player) {
        this.player.currentTime(0)
        this.player.pause()
      }
    }
  }
}
</script>

<style scoped>
.video-player {
  width: 100%;
  background: #000;
}

:deep(.video-js) {
  width: 100%;
  height: auto;
  aspect-ratio: 16/9;
}

:deep(.vjs-fullscreen) {
  aspect-ratio: unset;
}
</style>