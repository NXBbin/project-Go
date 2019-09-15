<template>
  <div class>
    <el-row class="main-header">
      <el-col :span="12">
        <el-page-header content="产品列表"></el-page-header>
      </el-col>

      <el-col :span="12">
        <el-button type="primary" class="float-right" @click="setDialogVisible = true">添加</el-button>
        <el-dialog title="设置" :visible.sync="setDialogVisible">
          <el-form :model="item" ref="itemSetForm" :rules="itemSetRules" label-width="140px">
            <el-tabs v-model="setDialogActiveName">
              <el-tab-pane label="基本信息" name="general"></el-tab-pane>

              <el-tab-pane label="SEO 信息" name="seo"></el-tab-pane>
            </el-tabs>
          </el-form>

          <div slot="footer" class="dialog-footer">
            <el-button @click="setDialogVisible = false">取 消</el-button>
            <el-button type="primary" @click="submitItemSetForm">确 定</el-button>
          </div>
        </el-dialog>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="24">
        <el-table ref="itemsTable" :data="items" tooltip-effect="dark" style="width: 100%">
          <el-table-column type="selection" width="55"></el-table-column>
          <el-table-column prop="Name" label="产品"></el-table-column>
          <el-table-column prop="Price" label="价格"></el-table-column>
          <el-table-column prop="Category.Name" label="所属分类"></el-table-column>
          <el-table-column fixed="right" label="操作" width="120">
            <template slot-scope>
              <el-button type="text" size="small">移除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import base from "../../api/uri.js";
export default {
  name: "ProductList",
  data() {
    return {
      items: [],

      item: {},
      itemSetRules: {},
      setDialogVisible: false,
      setDialogActiveName: "general"
    };
  },
  mounted() {
    this.refreshItems();
  },
  methods: {
    refreshItems() {
      this.axios.get(base + "products").then(resp => {
        if (resp.data.error == "") {
            this.items = resp.data.data
            console.log(this.items)
        } else {
            this.items = []
        }
      });
    },
    submitItemSetForm() {}
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
