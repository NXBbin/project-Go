<template>
  <div>
    <van-nav-bar title="请选择收货地址" left-arrow>
      <van-button plain type="info" slot="right" size="small">编辑</van-button>
    </van-nav-bar>
    <van-address-list v-model="AddressID" :list="list" @add="onAdd" @edit="onEdit" />
  </div>
</template>

<script>
import { AddressList, NavBar,Button  } from "vant";
import base, { staticBase } from "../plugin/api";

export default {
  components: {
    [AddressList.name]: AddressList,
    [Button .name]: Button ,

    [NavBar.name]: NavBar
  },
  data() {
    return {
      AddressID: null,
      list: [
      ]
    };
  },
  mounted() {
    this.refreshAddressList();
  },
  methods: {
    refreshAddressList() {
      let token = this.$store.getters.JWTToken;
      this.axios
        .get(base + "member-address-list", {
          headers: {
            Authorization: "Bearer " + token
          }
        })
        .then(resp => {
          console.log(resp);
          this.list = resp.data.data.map(element=>{
            return {
              id: element.ID,
              name: element.Name,
              tel: element.Tel,
              address: element.Province + element.City + element.County + element.addressDetail,
              isDefault: element.IsDefault,
            }
          });
          let index = this.list.findIndex(element => element.isDefault);
          if (index != -1) {
            this.AddressID = this.list[index].id;
          }
        })
        .catch(error => {
          console.error(error);
          this.$router.push({
            path: "/user/login",
            query: { redirect: "/user/address-list" }
          });
        });
    },
    onEdit() {
      alert("edit");
    },
    onAdd() {
      this.$router.push("/user/address-add")
    }
  }
};
</script>

<style lang="less">
</style>