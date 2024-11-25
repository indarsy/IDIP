import { getWorkshopList } from '@/api/workshop'

export default {
    namespaced: true,

    state: {
        workshopList: []
    },

    mutations: {
        SET_WORKSHOP_LIST(state, workshops) {
            state.workshopList = workshops
        }
    },

    actions: {
        async fetchWorkshopList({ commit }) {
            const data = await getWorkshopList()
            commit('SET_WORKSHOP_LIST', data)
        }
    }
}
