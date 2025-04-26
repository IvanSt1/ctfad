<script setup>
    import DataView from 'primevue/dataview';
    import Carousel from 'primevue/carousel';
    import GiftCard from '@/components/GiftCard.vue';
    import EmptyCard from '@/components/EmptyCard.vue';

    import { useSlots } from 'vue'
    const slots = useSlots()
</script>


<template>
    <EmptyCard v-if="latestCards===null"/>
    <Carousel v-else :value="latestCards" :numVisible="3" :numScroll="1" :responsiveOptions="responsiveOptions" circular>
    <template #item="card">
        <GiftCard :card="card.data"/>
    </template>
</Carousel>
</template>

<script>
export default {
    props: ['latestCards'],
    data() {
        return {
            layout: 'grid'
        }
    },
}
</script>
