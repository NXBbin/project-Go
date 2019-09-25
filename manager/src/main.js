import Vue from 'vue'
import App from './App.vue'

// elementUI
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
Vue.use(ElementUI);

// vue-router
import router from './router/router.js' 

// axios vue-axios
import './plugin/axios.js'

// vuex
// store
import store from './store/store.js'

// css
import './assets/manager.css'

Vue.config.productionTip = false

new Vue({
  store,
  router,
  render: h => h(App),
}).$mount('#app')
