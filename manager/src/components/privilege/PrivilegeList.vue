<template>
  <div class>
    <el-row class="main-header">
      <el-col :span="12">
        <el-page-header content="权限列表"></el-page-header>
      </el-col>

      <el-col :span="12">
        <el-button type="primary" class="float-right" @click="handleAddItem">添加</el-button>
        <div slot="footer" class="dialog-footer">
          <el-button @click="setDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitItemSetForm">确定</el-button>
        </div>
      </el-col>
    </el-row>

    <el-dialog :title="setDialogTitle" :visible.sync="setDialogVisible">
      <el-form :model="itemSetForm" ref="itemSetForm" :rules="itemSetRules" label-width="140px">
        <el-tabs v-model="setDialogActiveName">
          <el-tab-pane label="基本信息" name="general">

						<el-form-item label="权限名称" prop="Name">
							<el-input v-model="itemSetForm.Name"></el-input>
						</el-form-item>
						<el-form-item label="权限标识" prop="Key">
							<el-input v-model="itemSetForm.Key"></el-input>
						</el-form-item>
						<el-form-item label="排序" prop="SortOrder">
							<el-input v-model="itemSetForm.SortOrder"></el-input>
						</el-form-item>
						<el-form-item label="描述" prop="Description">
							<el-input v-model="itemSetForm.Description"></el-input>
						</el-form-item>
          </el-tab-pane>
        </el-tabs>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitItemSetForm">提交</el-button>
      </div>
    </el-dialog>

    <el-row class="main-content">
      <el-col :span="18" class="main-left">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>列表</span>
          </div>

          <!-- 数据展示 -->
          <el-row class="main-data">
            <el-col :span="24">
              <el-table
                ref="itemsTable"
                :data="items"
                tooltip-effect="dark"
                style="width: 100%"
                @sort-change="handleSortChange"
              >
                <el-table-column type="index" width="50"></el-table-column>
								<el-table-column prop="Name" label="权限名称" sortable="custom"></el-table-column>
								<el-table-column prop="Key" label="权限标识" sortable="custom"></el-table-column>
								<el-table-column prop="SortOrder" label="排序" sortable="custom"></el-table-column>

                <el-table-column type="selection" width="55"></el-table-column>
                <el-table-column fixed="right" label="操作" width="120">
                  <template slot-scope="scope">
                    <el-button
                      type="text"
                      size="small"
                      @click="handleRemoveItem(scope.$index, scope.row)"
                    >移除</el-button>

                    <el-button
                      type="text"
                      size="small"
                      @click="handleEditItem(scope.$index, scope.row)"
                    >编辑</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-col>
          </el-row>
        </el-card>

        <el-row class="main-pager">
          <el-col :span="24" class="pager-col">
            <el-pagination
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="currentPage"
              :page-sizes="[5, 10, 15, 20]"
              :page-size="pageSize"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
            ></el-pagination>
          </el-col>
        </el-row>
      </el-col>

      <el-col :span="6" class="main-right">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>筛选</span>
          </div>
          <!-- 数据筛选 -->
          <el-form :label-position="'top'" label-width="80px" :model="filterForm">

						<el-form-item label="权限名称">
							<el-input v-model="filterForm.filterName"></el-input>
						</el-form-item>
						<el-form-item label="权限标识">
							<el-input v-model="filterForm.filterKey"></el-input>
						</el-form-item>            
            <el-form-item>
              <el-button type="primary" @click="submitFilterForm" class="float-right">筛选</el-button>
              <el-button @click="resetFilterForm" class="float-right">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import base from "../../api/uri.js";
import { isNull } from "util";
export default {
  name: "PrivilegeList",
  data() {
    return {
      items: [],
      currentPage: 1,
      pageSize: 10,
      total: 0,
      filterForm: {},
      sortProp: "",
      sortOrder: "",

      itemSetOperation: "", //add, edit
      itemSetForm: {},
      itemSetRules: {
      },
      setDialogVisible: false,
      setDialogActiveName: "general",
      setDialogTitle: "添加",
    };
  },
  mounted() {
    this.refreshItems();
  },
  methods: {
    refreshItems(params = {}) {
      this.axios
        .get(base + "privilege", {
          params
        })
        .then(resp => {
          if (resp.data.error == "") {
            // 当前页数据
            this.items = resp.data.data;
            // 翻页需要的数据
            this.currentPage = resp.data.pager.currentPage;
            this.pageSize = resp.data.pager.pageSize;
            this.total = resp.data.pager.total;
          } else {
            this.items = [];
          }
        });
    },
    // 翻页-size改变
    handleSizeChange(size) {
      this.refreshItems(
        Object.assign(
          {
            currentPage: 1,
            pageSize: size
          },
          this.filterForm,
          {
            sortProp: isNull(this.sortOrder) ? null : this.sortProp,
            sortOrder: this.sortOrder
          }
        )
      );
    },
    // 翻页-page改变
    handleCurrentChange(page) {
      this.refreshItems(
        Object.assign(
          {
            currentPage: page,
            pageSize: this.pageSize
          },
          this.filterForm,
          {
            sortProp: isNull(this.sortOrder) ? null : this.sortProp,
            sortOrder: this.sortOrder
          }
        )
      );
    },

    // 筛选-提交
    submitFilterForm() {
      this.refreshItems(
        Object.assign(
          {
            currentPage: 1
          },
          this.filterForm
        )
      );
    },
    // 筛选-重置
    resetFilterForm() {
      this.filterForm = {};

      this.refreshItems({
        currentPage: 1,
        pageSize: this.pageSize
      });
    },
    // 排序事件
    handleSortChange(option) {
      // 记录排序方式
      this.sortProp = option.prop;
      this.sortOrder = option.order;

      // 带有排序参数请求
      let params = Object.assign(
        {
          currentPage: 1,
          pageSize: this.pageSize
        },
        this.filterForm,
        {
          sortProp: isNull(option.order) ? null : option.prop,
          sortOrder: option.order
        }
      );
      this.refreshItems(params);
    },
    // 移除
    handleRemoveItem(index, item) {
      // 确认框
      this.$confirm("是否确认删除 " + item.Name + " ?", "确认", {
        confirmButtonText: "确认",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          // 发出删除请求
          this.axios
            .delete(base + "privilege", {
              params: {
                ID: item.ID
              }
            })
            .then(resp => {
              if ("" != resp.data.error) {
                // 给出错误消息，提示框
                this.$message.error(resp.data.error);
                return;
              }

              // 删除成，更新数据
              this.items.splice(index, 1);
            });
        })
        .catch(() => {});
    },

    handleAddItem() {
      this.itemSetOperation = "add";
      // 设置为新对象
      this.itemSetForm = {
        SortOrder: 0,
      };

      this.setDialogVisible = true;
      this.setDialogTitle = "添加";
    },
    handleEditItem(index, item) {
      this.itemSetOperation = "edit";

      // 设为当前正在编辑的对象
      this.itemSetForm = item;

      this.setDialogVisible = true;
      this.setDialogTitle = "编辑";
    },
    // 提交设置表单
    submitItemSetForm() {
      this.$refs["itemSetForm"].validate(valid => {
        if (!valid) {
          return;
        }

        // 校验通过
        // 添加
        if ("add" == this.itemSetOperation) {
          this.itemSetAdd()
        }
        // 更新
        else if ("edit" == this.itemSetOperation) {
          this.itemSetEdit()
        }
      });
    },
    // 添加
    itemSetAdd() {
      this.axios.post(base + "privilege", this.itemSetForm).then(resp => {
        if (resp.data.error != "") {
          // 失败
          this.$message.error(resp.data.error);
          return;
        }
        this.items.push(resp.data.data);
        this.$refs["itemSetForm"].resetFields();
        this.setDialogVisible = false;
      });
    },
    // 编辑
    itemSetEdit() {
      this.axios
        .put(base + "privilege", this.itemSetForm, {
          params: {
            ID: this.itemSetForm.ID
          }
        })
        .then(resp => {
          if (resp.data.error != "") {
            // 失败
            this.$message.error(resp.data.error);
            return;
          }
          let index = this.items.findIndex(item => item.ID == resp.data.data.ID)
          this.items[index] = resp.data.data
          this.$refs["itemSetForm"].resetFields();
          this.setDialogVisible = false;
        });
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
