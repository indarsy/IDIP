import { getCaptureList, createCapture, cancelCapture } from '@/api/capture'

export default {
    namespaced: true,

    state: {
        captureList: [],
        currentTask: null
    },

    mutations: {
        SET_CAPTURE_LIST(state, list) {
            state.captureList = Array.isArray(list) ? list : []
        },

        ADD_CAPTURE_TASK(state, task) {
            state.captureList.unshift(task)
        },

        UPDATE_TASK_STATUS(state, { taskId, status }) {
            const task = state.captureList.find(t => t.id === taskId)
            if (task) {
                task.status = status
            }
        },

        SET_CURRENT_TASK(state, task) {
            state.currentTask = task
        }
    },

    actions: {
        // 创建采集任务
        async create({ commit }, data) {
            try {
                const response = await createCapture(data)
                commit('ADD_CAPTURE_TASK', response)
                return response
            } catch (error) {
                throw new Error(error.message || '创建采集任务失败')
            }
        },

        // 获取采集任务列表
        async fetchList({ commit }, workshopId) {
            try {
                const response = await getCaptureList(workshopId)
                const list = Array.isArray(response.data) ? response.data : []
                commit('SET_CAPTURE_LIST', list)
                return response
            } catch (error) {
                console.error('获取任务列表失败:', error)
                commit('SET_CAPTURE_LIST', [])
                throw new Error(error.message || '获取采集任务列表失败')
            }
        },

        // 取消采集任务
        async cancel({ commit }, taskId) {
            try {
                await cancelCapture(taskId)
                commit('UPDATE_TASK_STATUS', {
                    taskId,
                    status: 'cancelled'
                })
            } catch (error) {
                throw new Error(error.message || '取消任务失败')
            }
        },

        // 更新任务状态
        updateStatus({ commit }, { taskId, status }) {
            commit('UPDATE_TASK_STATUS', { taskId, status })
        },

        // 设置当前任务
        setCurrentTask({ commit }, task) {
            commit('SET_CURRENT_TASK', task)
        }
    },

    getters: {
        // 获取任务列表
        taskList: state => state.captureList,

        // 获取当前任务
        currentTask: state => state.currentTask,

        // 获取正在进行的任务
        runningTasks: state => state.captureList.filter(
            task => task.status === 'running'
        ),

        // 获取已完成的任务
        completedTasks: state => state.captureList.filter(
            task => task.status === 'completed'
        )
    }
}
