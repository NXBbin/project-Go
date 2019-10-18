<template>
  <div>
    <van-nav-bar title="订单确认" left-arrow @click-left="handleReturn"></van-nav-bar>

    <van-divider>产品信息</van-divider>

    <template v-if="!empty">
      <van-card
        v-for="product in products"
        :key="product.ID"
        :price="product.Price"
        :num="buyQuantities[product.ID]"
        desc="描述信息"
        :title="product.Name"
        :thumb="staticBase + product.Images[0].ImageSmall"
      >
        <div slot="tags" v-if="product.ModelInfo">
          <van-tag plain type="success">{{product.ModelInfo}}</van-tag>
        </div>
      </van-card>
    </template>

    <van-field
      readonly
      clickable
      label="配送方式"
      :value="shipping"
      placeholder="请选择"
      @click="onClickShowShippingPicker"
    />

    <van-popup v-model="showShippingPicker" position="bottom">
      <van-picker
        show-toolbar
        :columns="shippings"
        @cancel="showShippingPicker = false"
        @confirm="onShippingConfirm"
      />
    </van-popup>

    <van-field
      readonly
      clickable
      label="配送地址"
      :value="shipping"
      placeholder="请选择"
      @click="onClickShowShippingPicker"
    />

    <van-popup v-model="showShippingPicker" position="bottom">
      <van-picker
        show-toolbar
        :columns="shippings"
        @cancel="showShippingPicker = false"
        @confirm="onShippingConfirm"
        value-key="Title"
      />
    </van-popup>

    <van-field v-model="note" type="textarea" placeholder="请输入留言" rows="2" autosize />

    <van-submit-bar :price="totalPrice" button-text="生成订单" @submit="handleOrderSubmit">
      <van-checkbox v-model="all">全选</van-checkbox>
    </van-submit-bar>

    <van-overlay :show="overlayShow" @click="overlayShow = false" style="text-align: center;">
      <van-loading slot="default" style="margin-top: 60%; " />
    </van-overlay>
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
  Stepper,
  Panel,
  Popup,
  Picker,
  Field,
  Notify,
  Overlay,
  Loading
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
    [Panel.name]: Panel,
    [Popup.name]: Popup,
    [Picker.name]: Picker,
    [Field.name]: Field,
    [Overlay.name]: Overlay,
    [Loading.name]: Loading,
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
      showShippingPicker: false,
      shippings: [],
      shipping: "",
      shippingID: null,
      note: "",
      addressID: 3,
      overlayShow: false
    };
  },
  mounted() {
    this.memberAuth();
    this.refreshProductInfo();
  },
  methods: {
    memberAuth() {
      let token = this.$store.getters.JWTToken;
      // token 不存在
      if ("" == token) {
        this.user = null;
        return;
      }

      // 设置请求头，携带token
      this.axios.defaults.headers.common["Authorization"] = "Bearer " + token;

      // 校验token的合理性
      this.axios
        .get(base + "member-auth")
        .then(resp => {
          if (resp.data.error) {
            this.user = null;
            // 清理token
            this.$store.dispatch("clearToken");
            return;
          }
          this.user = resp.data.data;
        })
        .catch(error => {
          this.user = null;
        });
    },

    refreshProductInfo() {
      // 从购物车中提取信息
      let buyProducts = this.$store.getters.products;
      console.log(buyProducts);
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
    // 刷新购物车整体信息
    refreshCartInfo() {
      // 遍历全部的产品，使用单价*购买数量
      this.totalPrice = 0;
      for (let p of this.products) {
        this.totalPrice += p.Price * 100 * this.buyQuantities[p.ID];
      }
    },
    // 获取配送方式
    refreshShipping() {
      this.axios
        .get(base + "shipping", {
          headers: {
            Authorization: "Bearer " + this.$store.getters.JWTToken
          }
        })
        .then(resp => {
          this.shippings = resp.data.data;
        })
        .catch(error => {
          this.$router.push({
            path: "/user/login",
            query: { redirect: "/order" }
          });
        });
    },
    handleOrderSubmit() {
      this.overlayShow = true;
      this.axios
        .post(
          base + "order-create",
          {
            // 订单全部数据
            ShippingID: this.ShippingID,
            AddressID: this.addressID,
            BuyProductID: [1, 2]
          },
          {
            headers: {
              Authorization: "Bearer " + this.$store.getters.JWTToken
            }
          }
        )
        .then(resp => {
          // 成功，服务器响应了。会响应订单的sn
          console.log(resp);
          Notify(
            "订单处理中，前面有 " + resp.data.waitLen + " 个订单",
            "success"
          );
          // 轮询获取结果
          let intID = setInterval(() => {
            this.axios
              .get(base + "order-result", {
                params: {
                  sn: resp.data.data
                },
                headers: {
                  Authorization: "Bearer " + this.$store.getters.JWTToken
                }
              })
              .then(resp => {
                if (resp.data.data == "error") {
                  clearInterval(intID);
                  Notify("订单生成失败", "success");
                  this.overlayShow = false
                } else if (resp.data.data == "success") {
                  clearInterval(intID);
                  Notify("订单生成成功", "success");
                  this.overlayShow = false
                } else {
                  Notify("订单处理中", "success");
                }
              });
          }, 500);
        })
        .catch(error => {
          this.$router.push({
            path: "/user/login",
            query: { redirect: "/order" }
          });
        });
    },

    handleReturn() {
      this.$router.push({ path: "/cart" });
    },
    onClickShowShippingPicker() {
      this.refreshShipping();
      this.showShippingPicker = true;
    },
    onShippingConfirm(value) {
      this.shipping = value.Title;
      this.ShippingID = value.ID;
      this.showShippingPicker = false;
    }
  }
};
</script>

<style lang="less">
</style>