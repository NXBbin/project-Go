import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)


const routes = [
    { path: '/', component: ()=>import('../components/Index.vue') },
    { path: '/product', component: ()=>import('../components/Product.vue'), },
    { path: '/cart', component: ()=>import('../components/Cart.vue'), },
    { path: '/user/login', component: ()=>import('../components/UserLogin.vue'), },
    { path: '/user/address-list', component: ()=>import('../components/AddressList.vue'), },
    { path: '/user/address-add', component: ()=>import('../components/AddressAdd.vue'), },
    { path: '/order', component: ()=>import('../components/Order.vue'), },
  ]

const router = new VueRouter({
    routes // (缩写) 相当于 routes: routes
})

export default router