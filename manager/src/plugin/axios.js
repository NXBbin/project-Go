import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

Vue.use(VueAxios, axios)

import router from '../router/router.js'

import store from '../store/store.js'

// 添加请求拦截器
Vue.axios.interceptors.request.use(function (config) {
    // 在发送请求之前做些什么
    // 当存在 TOKEN 时，将其写入请求 Header
    // let token = window.localStorage.getItem("jwt-token")
    let token = store.getters.JWTToken
    if (token) {
        config.headers.Authorization = 'Bearer ' + token
    }

    return config;
  }, function (error) {
    // 对请求错误做些什么
    return Promise.reject(error);
  });

// 添加响应拦截器
Vue.axios.interceptors.response.use(function (response) {
  // 对响应数据做点什么

  return response;
}, function (error) {
  // 对响应错误做点什么
  if (error.response) {
    switch (error.response.status) {
      case 401:
        // 返回 401 清除token信息并跳转到登录页面
        router.replace({
          path: "/login",
          query: { 
            redirect: router.currentRoute.fullPath 
          }
        });
    }
  }
  return Promise.reject(error);
});