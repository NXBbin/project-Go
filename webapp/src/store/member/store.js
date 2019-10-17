
const tokenKey = 'MEMBER_JWTTOKEN'

export default {
    // 存储定义数据
    state: {
        JWTToken:'',
    },
    // 获取数据
    getters: {
        JWTToken(state) {
            return state.JWTToken
        },
        // JWTToken: state => state.JWTToken
    },
    // 更新数据
    mutations: {
       setJWTToken(state, payload) {
            state.JWTToken = payload.token
       },
       clearJWTToken(state, payload) {
            state.JWTToken = ''
       }
    },
    // 触发更新操作
    actions: {
        setJWTToken(context, payload) {
            context.commit('setJWTToken', payload)
            // 存储到localStorage中
            window.localStorage.setItem(tokenKey, payload.token)
        },
        // 初始化
        initToken(context, payload) {
            let token = window.localStorage.getItem(tokenKey) || ''
            context.commit('setJWTToken', {token})
        },
        clearToken(context, payload) {
            context.commit('clearJWTToken')
            window.localStorage.removeItem(tokenKey)
        }
    },
  }