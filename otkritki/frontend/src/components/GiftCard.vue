<script setup>
    import Card from 'primevue/card';
</script>

<template>
    <Card style="width: 30rem; overflow: hidden" class="relative mx-auto">
        <template #header>
            <img alt="card gif" :src="`/gifs/${card.imageType}.gif`" />
        </template>
        <template #title>From: {{ card.from }}</template>
        <template v-if="card.text" #content>
            <p class="m-0">
                {{ card.text }}
            </p>
        </template>
    </Card>
</template>

<script>
export default {
    props: ['card']
}
</script>

<style scoped>
img {
    height: 30rem;
    width: 25rem;
    position: relative;
    left: 8%;
}
</style>
