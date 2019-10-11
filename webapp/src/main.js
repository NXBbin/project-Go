import Vue from 'vue'
import App from './App.vue'

// 扩展
import router from './plugin/router.js' 
import './plugin/axios.js'
import './plugin/vant.js'

// vuex
// store
import store from './store/'

// 配置
Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app')
