<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>

export default {
  name: 'app',
  created() {
    // 初始化store中的jwttoken
    this.$store.dispatch("initJWTToken")
  }
}
</script>

<style>
</style>
