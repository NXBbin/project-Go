<template>
  <el-container>
    <el-header>
        <Header></Header>
    </el-header>
    <el-container>
      <el-aside width="200px">
          <Menu></Menu>
      </el-aside>
      <el-container>
        <el-main>
            <router-view></router-view>
        </el-main>
        <el-footer>
            <Footer></Footer>
        </el-footer>
      </el-container>
    </el-container>
  </el-container>
</template>

<script>
import Header from './Header.vue'
import Menu from './Menu.vue'
import Footer from './Footer.vue'

export default {
  name: "Layout",
  components: {
      Header, Menu, Footer,
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>


