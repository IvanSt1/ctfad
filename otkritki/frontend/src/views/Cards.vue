<script setup>
    import Card from 'primevue/card';
    import DataView from 'primevue/dataview';
    import GiftCard from '@/components/GiftCard.vue';
    import EmptyCard from '@/components/EmptyCard.vue';

    import { getLatestCards } from "@/service/cardops";
    import { toastError } from '@/service/toastService';

</script>


<template>
    <div class="card">
        <DataView :value="cards" paginator :rows="6" :layout="grid">
            <template #empty>
                <EmptyCard />
            </template>
            <template #list="slotProps">
                <div class="grid grid-nogutter">
                    <div v-for="(item, index) in slotProps.items" :key="index" class="col-10 sm:col-6 md:col-4 xl:col-6 p-3">
                        <GiftCard :card="item" />
                    </div>
                </div>
            </template>
        </DataView>
    </div>
</template>

<script>
export default {
    name: 'Cards',

    data() {
        return {
            layout: "grid",
            cards: null,
        }
    },
    methods: {
        getCards: async function() {
            try {
                this.cards = await getLatestCards()
            } catch(error) {
                toastError(error.message);
            }
        }
    },
    beforeMount() {
        this.getCards()
    }
    
}
</script>
