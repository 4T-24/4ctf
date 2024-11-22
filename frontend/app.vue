<template>
  <div>
    <NuxtRouteAnnouncer />
    <NuxtWelcome />
  </div>
</template>
<script setup lang="ts">
import config from '@/data/config.yaml'

console.debug(config.project)
</script>
