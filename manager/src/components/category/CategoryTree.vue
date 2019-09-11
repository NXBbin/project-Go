<template>
  <div>

    <el-table :data="tableData" style="width: 100%" stripe>
      <el-table-column prop="id" label="ID" width="180"></el-table-column>
      <el-table-column prop="name" label="分类名" width="180"></el-table-column>
      <el-table-column prop="parent_id" label="上级分类ID"></el-table-column>
    </el-table>

  </div>
</template>

<script>
export default {
  name: "CategoryTree",
  data() {
      return {
        tableData: [],
      }
  },
  //   挂载后事件
  mounted() {
    this.getTree()
  },
  // 方法列表
  methods: {
    getTree() {
      let url = "http://localhost:8088/category-tree";
      this.axios.get(url).then(resp => {
        // console.log(resp)
        if (resp.data.error == "") {
            this.tableData = resp.data.data
        } else {
            this.tableData = []
        }
      });
    },
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>


