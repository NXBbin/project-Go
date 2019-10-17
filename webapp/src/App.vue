<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: "app",
  mounted() {
    // 初始化store：包括：购物车
    this.$store.dispatch('initToken')
    this.$store.dispatch('initCart')
  },
};
</script>

<style>
</style>
