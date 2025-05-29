1.前端使用Vue3和Typescript，用axios进行http请求， 使用echarts绘制图表，Websocket进行实时数据传输<br>
2.后端使用Gin，ORM使用Gorm，数据库使用Postgresql，使用/gorilla/websocket进行Websocket通信<br>
3.进入./go目录，执行下面命令运行软件，端口占用5432和3000，数据库初始化在./go/init.sql，运行docker-compose后自动初始化数据库<br>
```bash
docker-compose up --build
```
4.技术结构<br>
启动页面会生成一个UUID，每个页面可以投一次票，刷新页面重新生成UUID，可以重新投票<br>
浏览器加载页面时先通过http获取投票信息，再建立/ws/poll WebSocket链接<br>
每当有新的投票时，处理投票信息插入后，使用broadcastPollUpdate广播通知所有ws链接，浏览器更新数据<br>
```askii
客户端            HTTP后端             WebSocket
|                 |                        |
|---- POST /vote ---->                     |
|                 |                        |
|                 |--> 插入数据库           |
|                 |                        |
|                 |--> 执行Websocket广播--> |
|                 |                        |
|<---- HTTP 响应---|                        |
|                 |                        |
|<========== WebSocket 广播消息 ==========  |
```
5.API接口<br>
接口文档地址：https://rh2lpudi5f.apifox.cn/