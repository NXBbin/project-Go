
import actions from './userActoins.js'

export default {
    // 存储定义数据
    state: {
      jwtToken: ""
    },
    // 获取数据
    getters: {
        JWTToken: state => {
            return state.jwtToken
        }
    },
    // 更新数据
    mutations: {
      updateJWTToken (state, token) {
        state.jwtToken = token
      }
    },
    // 触发更新操作
    actions,
  }