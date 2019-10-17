<template>
  <div>
    <van-nav-bar title="会员登录" left-arrow></van-nav-bar>

    <van-cell-group>
      <van-field v-model="loginForm.User" required clearable label="用户名" placeholder="请输入用户名" />

      <van-field
        v-model="loginForm.Password"
        type="password"
        label="密码"
        placeholder="请输入密码"
        required
        clearable
      />
    </van-cell-group>
    <van-button type="danger" style="width:100%;" @click="handleLogin">登录</van-button>

    <van-row>
      <van-col span="12" style=" color: rgba(0, 0, 0, 0.4);">短信验证码登录</van-col>
      <van-col span="12" style="text-align: right; color: rgba(0, 0, 0, 0.4);">手机快速注册</van-col>
    </van-row>

    <van-divider>其他方式登录</van-divider>
  </div>
</template>

<script>
import {
  Field,
  CellGroup,
  Cell,
  Icon,
  Row,
  Col,
  Notify,
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
    [Field.name]: Field,
    [CellGroup.name]: CellGroup,
    [Icon.name]: Icon,
    [Cell.name]: Cell,
    [Row.name]: Row,
    [Col.name]: Col,
    [Notify.name]: Notify,
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
      loginForm: {
        User: "",
        Password: ""
      }
    };
  },
  mounted() {},
  methods: {
    handleLogin() {
      let params = new URLSearchParams();
      params.append("User", this.loginForm.User);
      params.append("Password", this.loginForm.Password);
      this.axios.post(base + "member-login", params).then(resp => {
        if (resp.data.error) {
          // 错误
          Notify({ type: "primary", message: "用户或密码错误" });
          return;
        }
        // // 认证通过
        this.$store.dispatch("setJWTToken", {token: resp.data.token});

        // 同步购物车
        this.$store.dispatch("cartSync")

        // 跳转到来源
        let redirect = this.$route.query.redirect
          ? this.$route.query.redirect
          : "/";
        this.$router.push(redirect);
      });
    }
  }
};
</script>

<style lang="less">
</style>