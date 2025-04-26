<script setup>
    import Navbar from '@/components/Navbar.vue'
    import LoginNavbar from '@/components/LoginNavbar.vue'
    import Toast from 'primevue/toast';

    import { isAuthorized } from '@/service/authorize';
</script>

<template>
    <Toast />
    <div v-if="isAuthorized()">
        <Navbar />
    </div>
    <div v-else>
        <LoginNavbar />
    </div>
</template>

<script>
    export default {
        name: 'Index',
    }
</script>
