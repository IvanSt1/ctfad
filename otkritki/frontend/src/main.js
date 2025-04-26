import { createApp } from 'vue';

import PrimeVue from 'primevue/config'
import "primevue/resources/themes/aura-dark-green/theme.css"
import "primevue/resources/primevue.min.css"
import "primeicons/primeicons.css"
import "primeflex/primeflex.css";

import router from './router/index.js'
import App from './App.vue'
import ToastService from 'primevue/toastservice';

export const app = createApp(App)

app.use(router)
app.use(PrimeVue)
app.use(ToastService);
app.mount('#app')

