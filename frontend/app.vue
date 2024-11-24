<template>
  <NuxtLayout> 
    <NuxtPage />
  </NuxtLayout>
</template>
<script setup lang="ts">
import config from "@/data/config.yaml";

console.debug(config.project);
</script>
