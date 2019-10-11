import Vue from 'vue'
import App from './App.vue'

// 扩展
import router from './plugin/router.js' 
import './plugin/axios.js'
import './plugin/vant.js'

// 配置
Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
