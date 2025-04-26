<script setup>
    import Dialog from 'primevue/dialog';   
    import Button from 'primevue/button';
    import InputText from 'primevue/inputtext';
    import Textarea from 'primevue/textarea';
    import FloatLabel from 'primevue/floatlabel';
    import Dropdown from 'primevue/dropdown';
    import { toastSuccess, toastError } from '@/service/toastService';

    import { sendCard } from '@/service/cardops';

    import { ref } from 'vue';
    const visible = ref(false);
    const cards = ref([
        {name: 'Красивая открытка 1', code: "1"},
        {name: 'Красивая открытка 2', code: "2"},
        {name: 'Очень красивая открытка', code: "3"},
        {name: 'Креативная открытка', code: "4"},
        {name: 'Открытка сос мыслом', code: "5"},
        {name: 'Открытка для девочки', code: "6"},
        {name: 'Ого-го!!', code: "7"},
        {name: 'Неплохая открытка', code: "8"},
        {name: 'Тоже хорошая открытка', code: "9"},
        {name: 'Открытка', code: "10"},
    ]);
</script>

<template>
    <div>
        <Button label="Send card" @click="visible = true"/>

        <Dialog v-model:visible="visible" modal header = "Добавить открытку" :style="{ width: '30rem' }">
            <template #header>
                <div class="inline-flex align-items-center justify-content-center gap-2">
                    <span class="font-bold white-space-nowrap">Кого поздравим сегодня?</span>
                </div>
            </template>
            <div class="flex align-items-center gap-3 mb-3">
                <label for="receiver" class="font-semibold w-6rem">Имя получателя</label>
                <InputText v-model="receiver" id="receiver" class="flex-auto" autocomplete="off" />
            </div>

            <div class="flex align-items-center gap-3 mb-3">
                <Textarea v-model="text" rows="5" cols="30" placeholder="Напишите парочку приятных слов"/>
            </div>
        <Dropdown v-model="cardType" :value="cardType" :options="cards" optionLabel="name" optionValue="code" placeholder="Выберите открытку" checkmark :highlightOnSelect="false" class="w-full md:w-14rem" />

            <template #footer>
                <Button outlined type="button" label="Отправить" @click="visible=false; handleSendCard()" autofocus />
            </template>
        </Dialog>
    </div>
</template>

<script>
export default {
    name: 'AddCardMenu',
    methods: {
        async handleSendCard() {
            try {
                await sendCard(this.receiver, this.text, this.cardType);
                toastSuccess("Открытка отправлена");
            } catch (error) {
                toastError(error.message)
            }
        }
    },
    data() {
        return {
            cardType: {name: "aaaa", code: "0"},
        }
    }
}
</script>
