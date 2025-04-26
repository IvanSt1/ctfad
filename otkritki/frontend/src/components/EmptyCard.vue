<script setup>
    import GiftCard from '@/components/GiftCard.vue';
</script>

<template>
    <div>
        <GiftCard :card="emptyCard" />
    </div>
</template>


<script>
export default {
    data() {
        return {
            emptyCard: {
                imageType: "no_cards",
                from: "No one",
                text: "Никто не прислал вам открыток",
            }
        }
    }
}
</script>
