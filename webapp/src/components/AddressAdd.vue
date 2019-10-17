<template>
  <div>
    <van-nav-bar title="新增地址" left-arrow></van-nav-bar>
    <van-address-edit
      :area-list="areaList"
      show-postal
      show-set-default
      show-search-result
      @save="onSave"
    />
  </div>
</template>

<script>
import { AddressEdit, NavBar, Notify } from "vant";
import base, { staticBase } from "../plugin/api";

import areaList from "../plugin/area.js";

export default {
  components: {
    [AddressEdit.name]: AddressEdit,
    [Notify.name]: Notify,
    [NavBar.name]: NavBar
  },
  data() {
    return {
      areaList: areaList
    };
  },
  mounted() {},
  methods: {
    onSave(content) {
      this.axios
        .post(base + "member-address-add", content, {
          headers: {
            Authorization: "Bearer " + this.$store.getters.JWTToken
          }
        })
        .then(resp => {
            Notify({ type: 'success', message: '添加成功' });
          this.$router.push({
            path: "/user/address-list",
          });
        })
        .catch(error => {
          this.$router.push({
            path: "/user/login",
            query: { redirect: "/user/address-list" }
          });
        });
    }
  }
};
</script>

<style lang="less">
</style>