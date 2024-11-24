<template>
  <header>
    <Menubar :model="items"/>
  </header>
  <main><slot/></main>
</template>

<script lang="ts" setup>
import Menubar from 'primevue/menubar';
const items = ref([
    {
        label: 'Home',
        icon: 'pi pi-home'
    },
    {
        label: 'Features',
        icon: 'pi'}])
</script>

<style>
main{

  display: flex;
  flex-direction: column;
  align-items: center;
  height: calc(100vh - 70px);
}
</style>