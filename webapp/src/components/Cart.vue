<template>
  <div>
    <van-nav-bar title="购物车" left-arrow>
      <van-icon name="wap-nav" slot="right" />
    </van-nav-bar>

    <van-cell style="text-align: center;"  v-if="!user">
      <template slot="title">
        <span class="custom-title">登录后可同步账户购物车中的商品</span>
        <van-button
          type="danger"
          size="small"
          :to="{path: '/user/login', query:{redirect:'/cart'}}"
        >登录</van-button>
      </template>
    </van-cell>
    <template v-if="!empty">
      <van-card
        v-for="product in products"
        :key="product.ID"
        :price="product.Price"
        desc="描述信息"
        :title="product.Name"
        :thumb="staticBase + product.Images[0].ImageSmall"
      >
        <div slot="tags" v-if="product.ModelInfo">
          <van-tag plain type="success">{{product.ModelInfo}}</van-tag>
        </div>
        <van-stepper
          v-model="buyQuantities[product.ID]"
          slot="num"
          @change="handleQuantityChange(product.ID)"
        />
        <div slot="footer">
          <van-button plain size="mini" type="info">收藏</van-button>
          <van-button plain size="mini" type="danger" @click="handleRemoveProduct(product.ID)">删除</van-button>
        </div>
      </van-card>
    </template>

    <van-divider v-if="empty">购物车中还没有商品，请登录同步或者购买商品</van-divider>

    <van-submit-bar :price="totalPrice" button-text="提交订单" @submit="handleOrderSubmit">
      <van-checkbox v-model="all">全选</van-checkbox>
    </van-submit-bar>
  </div>
</template>

<script>
import {
  Icon,
  Cell,
  Button,
  NavBar,
  Card,
  Tag,
  SubmitBar,
  Checkbox,
  Divider,
  Stepper
} from "vant";
import base, { staticBase } from "../plugin/api";

export default {
  components: {
    [Icon.name]: Icon,
    [Cell.name]: Cell,
    [Button.name]: Button,
    [NavBar.name]: NavBar,
    [Card.name]: Card,
    [Tag.name]: Tag,
    [SubmitBar.name]: SubmitBar,
    [Checkbox.name]: Checkbox,
    [Divider.name]: Divider,
    [Stepper.name]: Stepper
  },
  data() {
    return {
      staticBase,
      all: true,
      empty: true,
      products: [],
      buyQuantities: {},
      totalPrice: 0,
      user: null,
    };
  },
  mounted() {
    this.memberAuth();
    this.refreshProductInfo();
  },
  methods: {
    memberAuth() {
        let token = this.$store.getters.JWTToken
        // token 不存在
        if ('' == token) {
            this.user = null
            return
        }

        // 设置请求头，携带token
        this.axios.defaults.headers.common['Authorization'] = 'Bearer ' + token

        // 校验token的合理性
        this.axios.get(base + 'member-auth').then(resp=>{
            if (resp.data.error) {
                this.user = null
                // 清理token
                this.$store.dispatch('clearToken')
                return
            }
            this.user = resp.data.data
        }).catch(error => {
            this.user = null
        })
    },

    refreshProductInfo() {
      // 从购物车中提取信息
      let buyProducts = this.$store.getters.products;
      console.log(buyProducts)
      // 提取全部的产品ID
      let ids = [];
      for (let p of buyProducts) {
        ids.push(p.productID);
        this.buyQuantities[p.productID] = p.buyQuantity;
      }

      if (ids.length == 0) {
        // 购物车中没有产品
        this.empty = true;
        return;
      }

      // 非空
      this.empty = false;

      // 再从服务器端获取详细信息
      this.axios
        .get(base + "cart-product", {
          params: {
            filterIDs: ids
          }
        })
        .then(resp => {
          // 得到产品信息后，将购买数量整合到一起。
          this.products = resp.data.data;
          this.refreshCartInfo();
        });
    },
    // 删除购物车产品
    handleRemoveProduct(id) {
      // 利用store完成
      this.$store.dispatch("removeFromCart", { productID: id });
      // 更新当前的产品列表
      this.refreshProductInfo();
    },
    // 调整数量
    handleQuantityChange(id) {
      // 利用store完成
      this.$store.dispatch("setBuyQuantity", {
        productID: id,
        buyQuantity: this.buyQuantities[id]
      });
      // 更新当前的产品列表
      this.refreshProductInfo();
    },
    // 刷新购物车整体信息
    refreshCartInfo() {
      // 遍历全部的产品，使用单价*购买数量
      this.totalPrice = 0;
      for (let p of this.products) {
        this.totalPrice += p.Price * 100 * this.buyQuantities[p.ID];
      }
    },
    handleOrderSubmit() {
      // 检测是否登录

      if (!this.user) {
        // 未登录
        this.$router.push({path:'/user/login', query: {redirect: '/cart'}})
        return
      }

      // 若已经登录，进入到送货地址的选择中
      this.$router.push({path:'/order'})
    }
  }
};
</script>

<style lang="less">
</style>