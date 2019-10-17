# project-Go
前后端分离产品项目：

静态资源服务器：
-->nginx.exe

pc端：
manager -- src  -->npm run serve

app移动端：
webapp -- src  -->npm run serve

后端：
backend -- src -- go run main.go

linux启动redis：
打开VMware（CentOs6.8）--编辑--虚拟网络编辑器--先还原配置--设置VMnet8的DHCP网关与系统内设置的网关段一致
启动服务器：cd redis-5.0.5 --- src/redis-server redis.conf

登录服务器查询： redis-5.0.5/src/redis-cli  
	查询订单序号：127.0.0.1:6379> get counter+日期
	获取订单信息：127.0.0.1:6379> hget tempOrder 订单号
	遍历订单队列：127.0.0.1:6379> xrange orderQueue - +
	获取订单队列长度：127.0.0.1:6379> xlen orderQueue
