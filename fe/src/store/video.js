import { defineStore } from 'pinia'
import { getVideoList, getVideoDetail } from '@/api/video'

export const useVideoStore = defineStore('video', {
    state: () => ({
        videoList: [],
        total: 0,
        currentVideo: null,
        loading: false
    }),

    actions: {
        async fetchVideoList(params) {
            this.loading = true
            try {
                const { data } = await getVideoList(params)
                this.videoList = data.list
                this.total = data.total
            } finally {
                this.loading = false
            }
        },

        async fetchVideoDetail(id) {
            this.loading = true
            try {
                const { data } = await getVideoDetail(id)
                this.currentVideo = data
                return data
            } finally {
                this.loading = false
            }
        },

        clearCurrentVideo() {
            this.currentVideo = null
        }
    }
})
