<template>
  <div>
    <!-- <el-table :data="tableData" style="width: 100%" stripe>
      <el-table-column prop="ID" label="ID" width="180"></el-table-column>
      <el-table-column prop="Name" label="分类名" width="180"></el-table-column>
      <el-table-column prop="ParentId" label="上级分类 ID"></el-table-column>
    </el-table>-->

    <el-row>
      <el-col :span="12">
        <el-page-header content="分类列表"></el-page-header>
      </el-col>

      <el-col :span="12">
        <el-button type="primary" class="float-right" @click="addDialogVisible = true">添加</el-button>
      </el-col>
      <el-dialog title="分类添加" :visible.sync="addDialogVisible">
        <el-form
          :model="categoryAdd"
          ref="categoryAddForm"
          :rules="categoryAddRules"
          label-width="140px"
        >
          <el-tabs v-model="addDialogActiveName">
            <el-tab-pane label="基本信息" name="general">
              <el-form-item label="分类名称" prop="Name">
                <el-input v-model="categoryAdd.Name" autocomplete="off"></el-input>
              </el-form-item>

              <el-form-item label="上级分类">
                <el-select
                  v-model="categoryAdd.ParentId"
                  placeholder="请选择上级分类"
                  @change="categoryOptionChangeHandle"
                >
                  <el-option
                    v-for="item in categoryOptions"
                    :key="item.ID"
                    :value="item.ID"
                    :label="item.Name"
                    :disabled="item.ID == 1"
                  >
                    <span v-bind:style="{ paddingLeft: item.Deep*16 + 'px' }">{{ item.Name }}</span>
                  </el-option>
                </el-select>
                <el-checkbox
                  v-model="categoryAdd.IsTop"
                  style="margin-left: 24px;"
                  @change="IstopHandle"
                >是顶级分类</el-checkbox>
              </el-form-item>

              <el-form-item label="分类 Logo">
                <el-upload action>
                  <el-button size="small" type="primary">点击上传</el-button>
                  <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过500kb</div>
                </el-upload>
              </el-form-item>

              <el-form-item label="描述">
                <el-input
                  type="textarea"
                  :rows="3"
                  placeholder="请输入内容"
                  v-model="categoryAdd.Description"
                ></el-input>
              </el-form-item>

              <el-form-item label="排序" prop="SortOrder">
                <el-input v-model.number="categoryAdd.SortOrder"></el-input>
              </el-form-item>
            </el-tab-pane>

            <el-tab-pane label="SEO 信息" name="seo">
              <el-form-item label="Meta Title">
                <el-input v-model="categoryAdd.MetaTitle" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="Meta Keywords">
                <el-input v-model="categoryAdd.MetaKeywords" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="Meta Description">
                <el-input type="textarea" v-model="categoryAdd.MetaDescription" autocomplete="off"></el-input>
              </el-form-item>
            </el-tab-pane>
          </el-tabs>
        </el-form>

        <div slot="footer" class="dialog-footer">
          <el-button @click="addDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="submitCategoryAddForm">确 定</el-button>
        </div>
      </el-dialog>
    </el-row>

    <el-row>
        
      <el-col :offset="2" :span="20">
        <el-tree :data="treeData" show-checkbox node-key="ID" :props="treeDataProps">
          <span slot-scope="{ node, data }" style="flex: 1;">
            <span>{{ node.label }}</span>
            <span>
              <el-button class="float-right" type="text" size="mini" @click.prevent="() => update(node, data)">更新</el-button>
              <el-button class="float-right" type="text" size="mini" @click="() => remove(node, data)">删除</el-button>
            </span>
          </span>
        </el-tree>
      </el-col>
    </el-row>

    <el-dialog title="分类更新" :visible.sync="updateDialogVisible">
        <el-form
          :model="categoryUpdate"
          ref="categoryUpdateForm"
          :rules="categoryUpdateRules"
          label-width="140px"
        >
          <el-tabs v-model="updateDialogActiveName">
            <el-tab-pane label="基本信息" name="general">
              <el-form-item label="分类名称" prop="Name">
                <el-input v-model="categoryUpdate.Name" autocomplete="off"></el-input>
              </el-form-item>

              <el-form-item label="上级分类">
                <el-select
                  v-model="categoryUpdate.ParentId"
                  placeholder="请选择上级分类"
                  @change="categoryOptionChangeHandle"
                >
                  <el-option
                    v-for="item in categoryOptions"
                    :key="item.ID"
                    :value="item.ID"
                    :label="item.Name"
                    :disabled="item.ID == 1"
                  >
                    <span v-bind:style="{ paddingLeft: item.Deep*16 + 'px' }">{{ item.Name }}</span>
                  </el-option>
                </el-select>
                <el-checkbox
                  v-model="categoryUpdate.IsTop"
                  style="margin-left: 24px;"
                  @change="IstopHandle"
                >是顶级分类</el-checkbox>
              </el-form-item>

              <el-form-item label="分类 Logo">
                <el-upload action>
                  <el-button size="small" type="primary">点击上传</el-button>
                  <div slot="tip" class="el-upload__tip">只能上传jpg/png文件，且不超过500kb</div>
                </el-upload>
              </el-form-item>

              <el-form-item label="描述">
                <el-input
                  type="textarea"
                  :rows="3"
                  placeholder="请输入内容"
                  v-model="categoryUpdate.Description"
                ></el-input>
              </el-form-item>

              <el-form-item label="排序" prop="SortOrder">
                <el-input v-model.number="categoryUpdate.SortOrder"></el-input>
              </el-form-item>
            </el-tab-pane>

            <el-tab-pane label="SEO 信息" name="seo">
              <el-form-item label="Meta Title">
                <el-input v-model="categoryUpdate.MetaTitle" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="Meta Keywords">
                <el-input v-model="categoryUpdate.MetaKeywords" autocomplete="off"></el-input>
              </el-form-item>
              <el-form-item label="Meta Description">
                <el-input type="textarea" v-model="categoryUpdate.MetaDescription" autocomplete="off"></el-input>
              </el-form-item>
            </el-tab-pane>
          </el-tabs>
        </el-form>

        <div slot="footer" class="dialog-footer">
          <el-button @click="updateDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="submitCategoryUpdateForm">确 定</el-button>
        </div>
      </el-dialog>
  </div>
</template>

<script>
import base from '../../api/uri.js'

export default {
  name: "CategoryTree",
  data() {
    return {
      //   tableData: []
      //   原始的服务器响应数据
      categories: [],
      // 需要的树数据
      treeData: [],
      // 展示的字段名字
      treeDataProps: {
        children: "children",
        label: "Name"
      },
      //   添加dialog是否可见
      addDialogVisible: false,
      categoryAdd: {
        IsTop: true,
        SortOrder: 0
      },
      addDialogActiveName: "general",

      categoryOptions: [],
      categoryAddRules: {
        Name: [
          {
            required: true,
            message: "请输入分类名称",
            trigger: ["blur", "change"]
          }
        ],
        SortOrder: [
          {
            required: true,
            message: "请输入排序值",
            trigger: ["blur", "change"]
          },
          {
            type: "integer",
            message: "排序值必须为整数",
            trigger: ["blur", "change"]
          }
        ]
      },
      updateDialogVisible: false,
      updateDialogActiveName: "general",
      categoryUpdate: {},
      categoryUpdateRules: {
        Name: [
          {
            required: true,
            message: "请输入分类名称",
            trigger: ["blur", "change"]
          }
        ],
        SortOrder: [
          {
            required: true,
            message: "请输入排序值",
            trigger: ["blur", "change"]
          },
          {
            type: "integer",
            message: "排序值必须为整数",
            trigger: ["blur", "change"]
          }
        ]
      },
    };
  },
  //   挂载后事件
  mounted() {
    this.getTree();
  },
  // 方法列表
  methods: {
    getTree() {
      let url = base + "category-tree"
      this.axios.get(url).then(resp => {
        if (resp.data.error == "") {
          this.categories = resp.data.data
          this.treeData = this.tree(resp.data.data, 0)
          this.indentTree(resp.data.data, 0)
        } else {
          this.treeData = []
          this.categoryOptions = []
        }
      });
    },
    // 递归的获取树状分类
    tree(all, id) {
      let items = [];
      for (let item of all) {
        if (item.ParentId == id) {
          // 存在子分类
          item.children = this.tree(all, item.ID)
          items.push(item)
        }
      }
      return items;
    },
    indentTree(all, id, deep = 0) {
      for (let item of all) {
        if (item.ParentId == id) {
          item.Deep = deep
          this.categoryOptions.push(item)
          // 存在子分类
          this.indentTree(all, item.ID, deep + 1)
        }
      }
    },
    categoryOptionChangeHandle(value) {
      this.categoryAdd.IsTop = false
    },
    IstopHandle(value) {
      if (value) {
        this.categoryAdd.ParentId = null
      }
    },
    submitCategoryAddForm() {
       // 校验
      this.$refs["categoryAddForm"].validate(valid => {
        if (valid) {
          // 校验通过
          // 表单数据提交到服务器端
          let url = base + "category";
          this.axios.post(url, this.categoryAdd).then(resp => {
            if (resp.data.error == "") {
              // 分类添加成功
              // 更新 this.treeData 和 CategoryOptions
              this.categories.push(resp.data.data)
              this.treeData = this.tree(this.categories, 0)
              this.categoryOptions = []
              this.indentTree(this.categories, 0)

              // 重置表单
              this.$refs["categoryAddForm"].resetFields()
            }
          })
          // 隐藏 dialog
          this.addDialogVisible = false
        } else {
          // 校验失败
          return false
        }
      }); 
    },
    remove(node, data) {
        let ID = data.ID
        let url = base + "category?ID=" + ID
        this.axios.delete(url).then(resp => { 
            if (resp.data.error == "") {
                // 根据当前的ID，确定原始数据中哪条记录被删了。
                let index = this.categories.findIndex(item => item.ID == ID)
                this.categories.splice(index, 1)
                this.treeData = this.tree(this.categories, 0)
                this.categoryOptions = []
                this.indentTree(this.categories, 0)
            }
        })    
    },
    update(node, data) {
      this.updateDialogVisible = true
      // 确定 更新的数据
      this.categoryUpdate = this.categories.find(item => item.ID == data.ID)
    },
     submitCategoryUpdateForm() {
      // 校验
      this.$refs["categoryUpdateForm"].validate(valid => {
        if (valid) {
          // 校验通过
          // 表单数据提交到服务器端
          let url = base + "category?ID=" + this.categoryUpdate.ID;
          this.axios.put(url, this.categoryUpdate).then(resp => {
            if (resp.data.error == "") {
              // 分类添加成功
              // 更新 this.treeData 和 CategoryOptions
              let index = this.categories.findIndex(item => item.ID == this.categoryUpdate.ID)
              this.categories[index] = resp.data.data
              this.treeData = this.tree(this.categories, 0);
              this.categoryOptions = [];
              this.indentTree(this.categories, 0);

              // 重置表单
              this.$refs["categoryUpdateForm"].resetFields();
            }
          });
          // 隐藏 dialog
          this.updateDialogVisible = false;
        } else {
          // 校验失败
          return false;
        }
      });
    },
    
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>


