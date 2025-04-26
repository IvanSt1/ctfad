<script setup>
    import Menubar from 'primevue/menubar';
    import AddCardMenu from '@/components/AddCardMenu.vue';
    import { logout } from '@/service/auth';
    import { isMale } from '@/service/genderLib';
    import { toastSuccess, toastError } from '@/service/toastService';

</script>


<template>
    <Menubar :model="items">
        <template v-if="isMale()" #end><AddCardMenu/></template>
    </Menubar>
</template>
  
  <script>
  export default {
    name: 'Navbar',
      data () {
          return {
              isMale: true,
              items: [
                  {
                      label: 'Home',
                      icon: 'pi pi-home',
                      url: '/user'
                  },
                  {
                      label: 'Cards',
                      icon: 'pi pi-gift',
                      url: '/cards'
                  },
                  {
                      label: 'Logout',
                      icon: 'pi pi-arrow-down-right',
                      command:  () => {this.handleLogout()}
                  },
              ],
          }
      },
      methods: {
          async handleLogout() {
              try {
                  await logout();
                  toastSuccess("Logged you out!");
              } catch (error) {
                  toastError(error.message);
              }
          }
      }
  }
      </script>
  
  
