import Vue from 'vue'
import Vuex from 'vuex'
Vue.use(Vuex)

import cart from './cart/store.js'
import member from './member/store.js'

const store = new Vuex.Store({
    modules: {
        cart,
        member,
    }
  })

export default store