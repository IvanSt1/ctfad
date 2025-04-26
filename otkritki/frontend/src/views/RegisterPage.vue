<script setup>
    import InputText from 'primevue/inputtext';
    import Button from 'primevue/button';
    import InputSwitch from 'primevue/inputswitch';
    import Buttom from 'primevue/button';

    import { ref } from 'vue';
    const checked = ref(false);

    import { register } from '@/service/auth';
    import { authorize} from '@/service/authorize';
</script>

<script>
export default {
    name: "Register",
    data() {
        return {
            message: "",
            messageColor: "red"
        }
    },
    methods: {
        async handleRegistration() {
            let userGender = "female"
            if (document.getElementById('maleSwitch').checked === true ) { // )))))))))))))0
                userGender = "male"
            }
            try {
                await register(this.username, this.password, userGender);
                authorize();
                this.$router.push({path: '/user'});
            } catch(error) {
                this.message = error.message;
            }
        }
    }
}
</script>

<template>
    <div class="center-container">
    <div class="container mt-5">
            <h1 class="text-center">Register</h1>
            <form>
                <div class="mb-3">
                    <label for="username" @keyup="(event) => this.message = ''" class="form-label">Username:</label>
                    <InputText class="form-control login-input" id="username" v-model="username" placeholder="10 symbols+"/>
                </div>

                <div class="mb-3">
                    <label for="password" class="form-label">Password:</label>
                    <InputText v-model="password" type="password" class="form-control login-input" id="password" placeholder="10 symbols+" />
                </div>

                <div class="m-10 card flex justify-content-center">
                    <label for="maleSwitch" class="form-label flex justify-content-center">I'm a male</label>
                    <InputSwitch inputId="maleSwitch" class="ml-3" v-model="checked" />
                </div>

            <Button label="Submit" @click.prevent="handleRegistration" class="btn btn-primary"/>
                <p :style="{ color: messageColor }" class="mt-3">{{ message }}</p>
            </form> 
    </div>
    </div>
</template>

<style lang='css' src='../assets/login_forms.css' scoped></style>
