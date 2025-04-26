<script setup>
    import ShortCardList from '@/components/ShortCardList.vue';
    import { getLatestCards } from  '@/service/cardops';
    import { toastError } from '@/service/toastService';
</script>

<template>
    <h1>Your have lately received this cards</h1>
    <div>
        <ShortCardList :latestCards="latestCards"/>
    </div>
</template>


<script>
export default {
    data() {
        return {
            latestCards: null
        }
    },
    methods: {
        setupCarousel: async function() {
            try {
                this.latestCards= await getLatestCards(12);
            } catch(error) {
                toastError("Could not load latest cards");
            }
        }
    },
    beforeMount() {
        this.setupCarousel();
    },
}
</script>
