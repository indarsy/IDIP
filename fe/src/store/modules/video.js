import {
    getVideoList,
    getVideoDetail,
    deleteVideo,
    batchDeleteVideos,
    getVideoStats
} from '@/api/video'

const state = {
    videoList: [],
    total: 0,
    currentVideo: null,
    loading: false,
    stats: {
        totalCount: 0,
        totalSize: 0,
        todayCount: 0
    }
}

const mutations = {
    SET_VIDEO_LIST(state, { list, total }) {
        state.videoList = list
        state.total = total
    },

    SET_CURRENT_VIDEO(state, video) {
        state.currentVideo = video
    },

    SET_LOADING(state, status) {
        state.loading = status
    },

    SET_STATS(state, stats) {
        state.stats = stats
    },

    REMOVE_VIDEO(state, id) {
        state.videoList = state.videoList.filter(video => video.id !== id)
        state.total--
    },

    REMOVE_VIDEOS(state, ids) {
        state.videoList = state.videoList.filter(video => !ids.includes(video.id))
        state.total -= ids.length
    }
}

const actions = {
    // 获取视频列表
    async getVideoList({ commit }, params) {
        commit('SET_LOADING', true)
        try {
            const { data } = await getVideoList(params)
            commit('SET_VIDEO_LIST', {
                list: data.list,
                total: data.total
            })
        } finally {
            commit('SET_LOADING', false)
        }
    },

    // 获取视频详情
    async getVideoDetail({ commit }, id) {
        commit('SET_LOADING', true)
        try {
            const { data } = await getVideoDetail(id)
            commit('SET_CURRENT_VIDEO', data)
            return data
        } finally {
            commit('SET_LOADING', false)
        }
    },

    // 删除视频
    async deleteVideo({ commit }, id) {
        try {
            await deleteVideo(id)
            commit('REMOVE_VIDEO', id)
        } catch (error) {
            // 做一些错误处理，比如记录日志
            console.error('删除视频失败:', error)
            // 或者转换错误信息
            throw new Error(`删除视频失败: ${error.message}`)
        }
    },

    // 批量删除视频
    async batchDeleteVideos({ commit }, ids) {
        try {
            await batchDeleteVideos(ids)
            commit('REMOVE_VIDEOS', ids)
        } catch (error) {
            console.error('删除视频失败:', error)
            // 或者转换错误信息
            throw new Error(`删除视频失败: ${error.message}`)
        }
    },

    // 获取统计信息
    async getVideoStats({ commit }) {
        try {
            const { data } = await getVideoStats()
            commit('SET_STATS', data)
        } catch (error) {
            console.error('Failed to get video stats:', error)
        }
    }
}

const getters = {
    // 获取视频总数
    videoCount: state => state.total,

    // 获取当前视频
    currentVideo: state => state.currentVideo,

    // 获取加载状态
    isLoading: state => state.loading,

    // 获取统计信息
    videoStats: state => state.stats,

    // 格式化存储大小
    formattedTotalSize: state => {
        const size = state.stats.totalSize
        const units = ['B', 'KB', 'MB', 'GB', 'TB']
        let index = 0
        let formattedSize = size

        while (formattedSize >= 1024 && index < units.length - 1) {
            formattedSize /= 1024
            index++
        }

        return `${formattedSize.toFixed(2)} ${units[index]}`
    }
}

export default {
    namespaced: true,
    state,
    mutations,
    actions,
    getters
} 