import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    {
        path: '/',
        redirect: '/video'
    },
    {
        path: '/video',
        name: 'Video',
        component: () => import('@/views/video/VideoManager.vue')
    },
    {
        path: '/workshop',
        name: 'Workshop',
        component: () => import('@/views/workshop/WorkShopManager.vue')
    },
    {
        path: '/video/capture',
        name: 'VideoCapture',
        component: () => import('@/views/video/VideoCapture.vue')
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router
