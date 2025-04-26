<script setup>
    import InputText from 'primevue/inputtext';
    import Button from 'primevue/button';
    import InputSwitch from 'primevue/inputswitch';
    import Buttom from 'primevue/button';

    import { login } from '@/service/auth';
    import { authorize} from '@/service/authorize';
</script>
<style lang='css' src='../assets/login_forms.css' scoped></style>

<script>

export default {
    name: "Login",
    data() {
        return {
            username: "",
            password: "",
            message: "",
            messageColor: "red"
        }
    },
    methods: {
        async handleLogin() {
            try {
                await login(this.username, this.password);
                authorize();
                this.$router.push({path: '/user'});
            } catch(error) {
                this.message = error.message;
            }
        },
    }
}
</script>

<template>
    <div class="center-container">
    <div class="container mt-5">
            <h1 class="text-center">Login</h1>
            <form>
                <div class="mb-3">
                    <label for="username" class="form-label">Username:</label>
                    <InputText @keyup="(event) => this.message = ''" class="form-control login-input" id="username" v-model="username" placeholder=""/>
                </div>

                <div class="mb-3">
                    <label for="password" class="form-label">Password:</label>
                    <InputText v-model="password" type="password" class="form-control login-input" id="password" placeholder="" />
                </div>
            <Button label="Submit" @click.prevent="handleLogin" class="btn btn-primary"/>
                <p :style="{ color: messageColor }" class="mt-3">{{ message }}</p>
            </form> 
    </div>
    </div>
</template>

