<script setup>
import { Head, Link } from '@inertiajs/vue3';

defineProps({
  text: String,
});

</script>

<template>
    <div>Hello, {{ text }}!</div>
</template>
