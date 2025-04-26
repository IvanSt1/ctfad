<script setup>
    import Menubar from 'primevue/menubar'
</script>

<template>
    <div>
        <Menubar :model="items" />
    </div>
</template>

<script>
export default {
    name: 'LoginNavbar',
    data() { 
        return {
            items: [
                {
                    label: 'Login',
                    url: '/login'
                },
                {
                    label: 'Register',
                    url: '/register'
                }
            ]
        }
    }
}
</script>
