<template>
  <div class="index">
    <van-nav-bar>
      <van-icon slot="left" name="wap-nav" />
      <van-field slot="title" v-model="keyword" placeholder="ThinkPad T490s"></van-field>
      <van-icon slot="right" name="scan" />
    </van-nav-bar>

    <van-swipe style="height: 200px;">
      <van-swipe-item v-for="(image, index) in images" :key="index">
        <img v-lazy="image" />
      </van-swipe-item>
    </van-swipe>

    <van-grid :border="false" :column-num="5">
      <van-grid-item v-for="value in 10" :key="value" icon="photo-o" text="分类" />
    </van-grid>

    <van-divider>推荐商品</van-divider>
    <van-grid :border="false" :column-num="2">
      <van-grid-item
        v-for="product in promoteProducts"
        :key="product.ID"
        :to="{ path: 'product',  query: { ID: product.ID } }"
      >
        <van-image v-if="product.Images[0]" :src="staticBase + product.Images[0].ImageSmall" />
        <van-image v-else src="https://img.yzcdn.cn/vant/apple-1.jpg" />
        <span>{{product.Name}}</span>
        <div style="color: #f00;">
          ￥
          <span>{{product.Price}}</span>
        </div>
      </van-grid-item>
    </van-grid>

    <van-tabbar v-model="active">
      <van-tabbar-item icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item icon="bar-chart-o" dot>分类</van-tabbar-item>
      <van-tabbar-item icon="shopping-cart-o" info="5">购物车</van-tabbar-item>
      <van-tabbar-item icon="user-o">用户</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script>
import base, { staticBase } from "../plugin/api.js";
import {
  Swipe,
  SwipeItem,
  Grid,
  GridItem,
  Image,
  Divider,
  NavBar,
  Field,
  Button,
  Icon,
  Tabbar,
  TabbarItem
} from "vant";

export default {
  components: {
    [Swipe.name]: Swipe,
    [SwipeItem.name]: SwipeItem,
    [Grid.name]: Grid,
    [GridItem.name]: GridItem,
    [Image.name]: Image,
    [Divider.name]: Divider,
    [NavBar.name]: NavBar,
    [Field.name]: Field,
    [Button.name]: Button,
    [Icon.name]: Icon,
    [Tabbar.name]: Tabbar,
    [TabbarItem.name]: TabbarItem
  },
  data() {
    return {
      staticBase: staticBase,
      keyword: "",
      images: [
        "https://img.yzcdn.cn/vant/apple-1.jpg",
        "https://img.yzcdn.cn/vant/apple-2.jpg"
      ],
      active: 0,
      promoteProducts: []
    };
  },
  mounted() {
    this.refreshPromote({
      pageSize: 12,
      currentPage: 1
    });
  },
  methods: {
    // 刷新推荐商品
    refreshPromote(options) {
      this.axios
        .get(base + "product-promote", {
          params: options
        })
        .then(resp => {
          if (resp.data.error) {
            this.promoteProducts = [];
            return;
          }
          this.promoteProducts = resp.data.data;
        });
    }
  }
};
</script>

<style lang="less">
</style>