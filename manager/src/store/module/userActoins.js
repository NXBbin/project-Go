export default {
    setJWTToken (context, token) {
        // 提交状态修改
        context.commit('updateJWTToken', token)
        // 存储到localStorage中
        window.localStorage.setItem("jwt-token", token)
    },
    initJWTToken(context) {
        // 从 localStorage 中读取token
        let token = window.localStorage.getItem("jwt-token")
        if (token) {
            context.commit('updateJWTToken', token) 
        }
    }
}