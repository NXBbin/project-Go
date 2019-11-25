<template>
  <div>
    <el-row class="login-row">
      <el-col :span="8" :offset="8">
        <el-card shadow="always">
          <div slot="header" class="clearfix">
            <span>用户登录</span>
          </div>

          <el-form
            :model="loginForm"
            status-icon
            :rules="loginRules"
            ref="loginForm"
            label-width="100px"
          >
            <el-form-item label="用户名" prop="User">
              <el-input v-model="loginForm.User" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="密码" prop="Password">
              <el-input type="password" v-model="loginForm.Password" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="验证码" prop>
              <el-input v-model="loginForm.Code"></el-input>
            </el-form-item>
            <el-form-item>
              <el-image :src="codeSrc">
                <div slot="error" class="image-slot">
                  <i class="el-icon-picture-outline"></i>
                </div>
              </el-image>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitLoginForm">提交</el-button>
              <el-button @click="resetLoginForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import base from "../../api/uri";
export default {
  name: "Login",
  data() {
    return {
      loginForm: {},
      loginRules: {
        User: [{ required: true, message: "请输入用户名", trigger: "blur" }],
        Password: [{ required: true, message: "请输入密码", trigger: "blur" }]
      },
      codeSrc: base + 'check-code',
    };
  },
  mounted() {
   
  },
  methods: {
  
    submitLoginForm() {
      this.$refs["loginForm"].validate(valid => {
        if (valid) {
          let params = new URLSearchParams();
          params.append("User", this.loginForm.User);
          params.append("Password", this.loginForm.Password);
          params.append("Code", this.loginForm.Code);
          this.axios.post(base + "user/auth", params).then(resp => {
            if (resp.data.error) {
              this.$message.error(resp.data.error);
              return;
            }
            // 将token存储到localStorage中
            // window.localStorage.setItem("jwt-token", resp.data.token)
            // 使用store分发action的方式存储：
            this.$store.dispatch("setJWTToken", resp.data.token);
            // 跳转到来源
            let redirect = this.$route.query.redirect
              ? this.$route.query.redirect
              : "/";
            this.$router.push(redirect);
          });
        } else {
          return false;
        }
      });
    },
    resetLoginForm() {
      this.$refs["loginForm"].resetFields();
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.login-row {
  margin-top: 120px;
}
</style>
