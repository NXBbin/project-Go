import Vue from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)


import Layout from '../components/partial/Layout.vue'
import CategoryTree from '../components/category/CategoryTree.vue'
// import 
// 定义路由
const routes = [
    // 非常规，不需要上菜单footer的组件
    // {path: '/login', component: Login},
    // 嵌套，常规的组件
    { 
        path: '/', component: Layout, 
        children: [
            { path: 'category-tree', component: CategoryTree, },
            { path: 'products', component: ()=>import('../components/product/ProductList.vue'), },
            { path: 'brand', component: ()=>import('../components/brand/BrandList.vue'), },
            { path: 'user', component: ()=>import('../components/user/UserList.vue'), },
        ]
    },
  ]

const router = new VueRouter({
    routes // (缩写) 相当于 routes: routes
})

export default router