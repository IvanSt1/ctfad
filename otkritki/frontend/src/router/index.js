import { createRouter, createWebHistory } from "vue-router";

import Register from '@/views/RegisterPage.vue';
import Login from '@/views/LoginPage.vue';
import HomePage from '@/views/HomePage.vue';
import Cards from '@/views/Cards.vue';
import CardPage from '@/views/CardPage.vue';
import NotFound from '@/views/NotFound.vue';

import { canAccess, authorize, removeAuthorize } from '@/service/authorize';

const routes = [
    {path: '/login', name: 'login', component: Login},
    {path: '/register', name: 'register', component: Register},
    {path: '/user', name: 'home', component: HomePage},
    {path: '/cards', name: 'cards', component: Cards},
    {path: '/card/:id', component: CardPage},
    {path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound },
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

router.beforeEach(async (to, from) => {
    if (to.name !== 'login' && to.name !== 'register') {
        const allowed = await canAccess();
        if (!allowed) {
            router.push({replace: true, path:"/login"})
            removeAuthorize();
        } else {
            authorize();
        }
    }
});

export default router
