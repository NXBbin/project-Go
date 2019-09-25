import Vue from 'vue'
import VueRouter from 'vue-router'
Vue.use(VueRouter)


import Layout from '../components/partial/Layout.vue'
import CategoryTree from '../components/category/CategoryTree.vue'
// import 
// 定义路由
const routes = [
    // 非常规，不需要上菜单footer的组件
    {path: '/login', component: ()=>import('../components/user/Login.vue'), },
    
    // 嵌套，常规的组件
    { 
        path: '/', 
        component: Layout, 
        meta: {
            requireAuth: true // 路由元信息
        },
        children: [
            { path: 'category-tree', component: CategoryTree, },
            { path: 'products', component: ()=>import('../components/product/ProductList.vue'), },
            { path: 'brand', component: ()=>import('../components/brand/BrandList.vue'), },
            { path: 'user', component: ()=>import('../components/user/UserList.vue'), },
            { path: 'role', component: ()=>import('../components/role/RoleList.vue'), },
            { path: 'privilege', component: ()=>import('../components/privilege/PrivilegeList.vue'), },
            { path: 'attr-type', component: ()=>import('../components/attrType/AttrTypeList.vue'), },
            { path: 'attr-group', component: ()=>import('../components/attrGroup/AttrGroupList.vue'), },
            { path: 'attr', component: ()=>import('../components/attr/AttrList.vue'), },
            { path: 'product-attr', component: ()=>import('../components/productAttr/ProductAttrList.vue'), },
        ]
    },
  ]

const router = new VueRouter({
    routes // (缩写) 相当于 routes: routes
})

import store from '../store/store'
// 路由守卫
// 前置守卫 guard
router.beforeEach((to, from, next) => {
    // 需要认证，当没有token，则跳转到login
    let token = store.getters.JWTToken
    if (to.meta.requireAuth && !token) {
    // if (to.meta.requireAuth && !window.localStorage.getItem("jwt-token")) {
        next({
            path: '/login',
            query: { redirect: to.fullPath }  
        })
        return
    }

    // 不需要认证，或者存在token，则继续即可
    next()

})

export default router