<template>
  <el-menu
    :default-active="activeIndex"
    class="el-menu-demo"
    mode="horizontal"
  >
    <el-menu-item index="1">处理中心</el-menu-item>
    <el-submenu index="2">
      <template slot="title">Hank</template>
      <el-menu-item index="2-1">退出</el-menu-item>
    </el-submenu>
    <el-menu-item index="3" disabled>消息中心</el-menu-item>
    <el-menu-item index="4">
      <a href="https://www.ele.me" target="_blank">订单管理</a>
    </el-menu-item>
  </el-menu>
</template>

<script>
export default {
  name: "Header",
   data() {
      return {
          activeIndex: "1",
      } 
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>





