<template>
  <div class="goods">
    <van-swipe class="goods-swipe" :autoplay="3000">
      <van-swipe-item v-for="thumb in product.Images" :key="thumb.ID">
        <img :src="staticBase + thumb.ImageSmall" />
      </van-swipe-item>
    </van-swipe>

    <van-cell-group>
      <van-cell>
        <div class="goods-title">{{ product.Name }}</div>
        <div class="goods-price">￥ {{ product.Price }}</div>
      </van-cell>
      <van-cell class="goods-express">
        <van-col span="14">剩余：{{ product.quantity || 0 }}</van-col>
      </van-cell>
    </van-cell-group>

    <van-cell-group class="goods-cell-group">
      <van-cell value="进入店铺" icon="shop-o" is-link @click="sorry">
        <template slot="title">
          <span class="van-cell-text">官方店铺</span>
        </template>
      </van-cell>
      <van-cell title="线下门店" icon="location-o" is-link @click="sorry" />
    </van-cell-group>

    <van-cell-group class="goods-cell-group">
      <van-cell title="查看商品详情" is-link @click="sorry" />
    </van-cell-group>

    <van-cell-group class="goods-cell-group" v-if="product.GroupID">
      <van-cell title="型号" :label="product.ModelInfo" is-link @click="selectModel" />
    </van-cell-group>

    <van-popup v-model="selectModelShow" position="bottom" v-if="product.GroupID">
      <van-list finished-text="没有可选型号" @load="onLoad">
        <van-cell v-for="p in product.Group.Products" :key="p.ID" :title="p.ModelInfo" />
      </van-list>
    </van-popup>

    <van-goods-action>
      <van-goods-action-icon icon="chat-o" @click="sorry">客服</van-goods-action-icon>
      <van-goods-action-icon icon="cart-o" @click="onClickCart">购物车</van-goods-action-icon>
      <van-goods-action-button type="warning" @click="sorry">加入购物车</van-goods-action-button>
      <van-goods-action-button type="danger" @click="sorry">立即购买</van-goods-action-button>
    </van-goods-action>
  </div>
</template>

<script>
import {
  Tag,
  Col,
  Icon,
  Cell,
  CellGroup,
  Swipe,
  Toast,
  SwipeItem,
  GoodsAction,
  GoodsActionIcon,
  GoodsActionButton,
  Popup,
  List
} from "vant";
import base, { staticBase } from "../plugin/api";

export default {
  components: {
    [Tag.name]: Tag,
    [Col.name]: Col,
    [Icon.name]: Icon,
    [Cell.name]: Cell,
    [CellGroup.name]: CellGroup,
    [Swipe.name]: Swipe,
    [SwipeItem.name]: SwipeItem,
    [GoodsAction.name]: GoodsAction,
    [GoodsActionIcon.name]: GoodsActionIcon,
    [GoodsActionButton.name]: GoodsActionButton,
    [Popup.name]: Popup,
    [List.name]: List
  },
  data() {
    return {
      staticBase,
      goods: {
        title: "美国伽力果（约680g/3个）",
        price: 2680,
        express: "免运费",
        remain: 19,
        thumb: [
          "https://img.yzcdn.cn/public_files/2017/10/24/e5a5a02309a41f9f5def56684808d9ae.jpeg",
          "https://img.yzcdn.cn/public_files/2017/10/24/1791ba14088f9c2be8c610d0a6cc0f93.jpeg"
        ]
      },
      product: {},
      selectModelShow: false
    };
  },
  mounted() {
    this.refreshProduct(this.$route.query.ID);
  },
  methods: {
    refreshProduct(ID) {
      this.axios
        .get(base + "product-info", {
          params: {
            ID
          }
        })
        .then(resp => {
          if (resp.data.error) {
            this.product = {};
            return;
          }
          this.product = resp.data.data;
        });
    },

    selectModel() {
      this.selectModelShow = true;
    },

    formatPrice() {
      return "¥" + (this.goods.price / 100).toFixed(2);
    },
    onClickCart() {
      this.$router.push("cart");
    },
    sorry() {
      Toast("暂无后续逻辑~");
    }
  }
};
</script>

<style lang="less">
.goods {
  padding-bottom: 50px;
  &-swipe {
    img {
      width: 100%;
      display: block;
    }
  }
  &-title {
    font-size: 16px;
  }
  &-price {
    color: #f44;
  }
  &-express {
    color: #999;
    font-size: 12px;
    padding: 5px 15px;
  }
  &-cell-group {
    margin: 15px 0;
    .van-cell__value {
      color: #999;
    }
  }
  &-tag {
    margin-left: 5px;
  }
}
</style>