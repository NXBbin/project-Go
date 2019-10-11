<template>
  <div class="index">
    <van-nav-bar>
      <van-icon slot="left" name="wap-nav" />
      <van-field slot="title" v-model="keyword" placeholder="ThinkPad T490s"></van-field>
      <van-icon slot="right" name="scan" />
    </van-nav-bar>

    <router-view></router-view>

    <van-tabbar v-model="active">
      <van-tabbar-item icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item icon="bar-chart-o" dot>分类</van-tabbar-item>
      <van-tabbar-item icon="shopping-cart-o" info="5">购物车</van-tabbar-item>
      <van-tabbar-item icon="user-o">用户</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script>
import base, {staticBase} from "../plugin/api.js";
import {
  NavBar,
  Field,
  Button,
  Icon,
  Tabbar,
  TabbarItem,
  Divider
} from "vant";

export default {
  components: {
    [NavBar.name]: NavBar,
    [Field.name]: Field,
    [Button.name]: Button,
    [Icon.name]: Icon,
    [Tabbar.name]: Tabbar,
    [TabbarItem.name]: TabbarItem,
    [Divider.name]: Divider
  },
  data() {
    return {
      keyword: "",
    };
  },
  mounted() {
  },
  methods: {
    // 刷新推荐商品
  }
};
</script>

<style lang="less">
</style>