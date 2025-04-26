<script setup>
    import GiftCard from '@/components/GiftCard.vue'
    import { getCard } from '@/service/cardops';
    import { toastError } from '@/service/toastService';
</script>

<template>
    <div>
        <GiftCrd :card="card"/>
    </div>
</template>

<script>
export default {
    data() {
        return {
            card: null
        }
    },
    methods: {
        getCardData: async function() {
            try {
                this.card = await getCard(this.$route.params.id);
            } catch(error) {
                toastError(error.message);
            }
        }
    },
    beforeMount() {
        this.getCardData();
    },
}
</script>
