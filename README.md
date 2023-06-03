# chatchat

此项目为前后端分离开发||[后端接口文档链接](https://console-docs.apipost.cn/preview/56cc7639487d33cb/20154c443f1b0168)

[![My Skills](https://skillicons.dev/icons?i=golang,docker,mysql,redis,postman,git)](https://skillicons.dev)



## 💡  简介

基于 [gin ](https://github.com/gin-gonic/gin)框架实现的一个简单聊天室后端，模拟了第三方登录



## 🚀 功能

### 认证系统

- 第三方登录颁发 `oauth2 token`  (包含 `access token` 和 `refresh token`)
- 登录颁发 `token`
- 对令牌进行解析获取认证信息与用户信息

### 用户系统

- 用户的登录与邮箱验证码注册
- 第三方登录&&绑定邮箱
- 用户改名,改头像，改简介，改密码
- 获取用户信息

### 好友系统

- 添加好友
- 删除好友

### 群组系统

- 退出群组
- 删除群组（仅群主可操作）
- 踢人（仅群主可操作）
- 关键词模糊查找群组
- 创建群组
- 加入群组

### 聊天系统

- 实时发送消息，实时接收

- 离线状态下也可接收消息，上线后可获取

- 未读消息提示

- 获取历史聊天记录（群组和好友）

  

  

## 🌟 亮点

### 其他功能

- 验证码注册

  ![image-20230602113259130](./images/image-20230602113259130.png)

- 云储存图片，节省空间

  ![image-20230602113403012](./images/image-20230602113403012.png)
  
- 基于Oauth2的自定义第三方登录，实现了两种登录方式

  1，用code换取token，安全性更高
  
  ![img](https://upload-images.jianshu.io/upload_images/9173561-cf2e57e89aa63a95.png?imageMogr2/auto-orient/strip|imageView2/2/w/1010/format/webp)
  
  2，用pwd直接换取token，可用于受信任的项目
  
  
  
  

> 详情请见源码和接口文档



### 技术栈

- [gin](https://gin-gonic.com/zh-cn/docs/introduction/)

> Gin 是一个用 Go (Golang) 编写的 Web 框架。
>

```
本项目功能不是很繁杂，所以用不上微服务，所以采用了写单体比较合适的gin框架
```

- [mysql](https://www.mysql.com/)

> 一个关系型数据库管理系统，由瑞典MySQL AB 公司开发，属于 Oracle 旗下产品。MySQL 是最流行的关系型数据库管理系统关系型数据库管理系统之一，在 WEB 应用方面，MySQL是最好的 RDBMS (Relational Database Management System，关系数据库管理系统) 应用软件之一

```
遇事不决还得是mysql，以后重构可以考虑mongodb
```

- [redis](https://redis.io/)

> 一个开源的、使用C语言编写的、支持网络交互的、可基于内存也可持久化的Key-Value数据库

```
缓存存储还是选型最普遍的redis
```

- [sentinel](https://sentinelguard.io/zh-cn/)

> 一个面向分布式、多语言异构化服务架构的流量治理组件

```
我用来做限流
```

- [jaeger](https://www.jaegertracing.io/)

> 由Uber开源的分布式追踪系统

```
日志管理之余，使用jaeger做追踪系统
```

- [docker](https://www.docker.com/)

> Google 公司推出的 Go 语言 进行开发实现，基于 Linux 内核的 cgroup，namespace，以及 AUFS 类的 Union FS 等技术的一个容器服务

```
容器用docker-compose部署
```

## 📂 存储设计

### 表设计

#### 用户系统

##### `users`

##### 模拟第三方用户账号

![](./images/image-20230530171635238.png)

![image-20230530171919804](./images/image-20230530171919804.png)

##### `user_bases`

聊天室项目所用用户账号数据库

第三方登录绑定后将在此数据库中添加数据

![image-20230530172157094](./images/image-20230530172157094.png)

![image-20230530172053843](./images/image-20230530172053843.png)

`user_auths`

第三方用户与聊天室用户绑定

![image-20230530172239074](./images/image-20230530172239074.png)

![image-20230530172313806](./images/image-20230530172313806.png)

#### 聊天系统

##### `friend`

![image-20230530172413540](./images/image-20230530172413540.png)

![image-20230530172427961](./images/image-20230530172427961.png)



##### `groups`

![image-20230530172454833](./images/image-20230530172454833.png)

![image-20230530172508923](./images/image-20230530172508923.png)



##### `group_members`

![image-20230530172538903](./images/image-20230530172538903.png)

![image-20230530172548966](./images/image-20230530172548966.png)







##### `message`

![image-20230530172708114](./images/image-20230530172708114.png)

![image-20230530172732979](./images/image-20230530172732979.png)

### 缓存设计

#### 表缓存

#### 用户信息缓存

使用哈希表存储已经注册的用户信息

key 名称：`user:[username]`

![image-20230530173836184](./images/image-20230530173836184.png)

#### 用户账号id缓存

使用 有序集合 存储用户 id 与用户名（QQ邮箱），方便通过 唯一id 获取用户名

key 格式：`userID`

#### ![image-20230530173956438](./images/image-20230530173956438.png)

#### 第三方账号绑定缓存

key 格式：`Oauth2User`

![image-20230530174354581](./images/image-20230530174354581.png)

#### 消息缓存

好友消息 key 格式：`friend:[id1]to[id2]`

![image-20230530174504096](./images/image-20230530174504096.png)

群组消息 key 格式: `group:[group_id]`



![image-20230530174558672](./images/image-20230530174558672.png)

离线消息缓存，在被获取之后自动删除

key 格式: `id`

![image-20230530174644328](./images/image-20230530174644328.png)

#### 邮箱验证码缓存

key 格式：`mail:[username]`,120s后过期删除

![image-20230530174759479](./images/image-20230530174759479.png)

## 📖 API文档

[接口文档](https://console-docs.apipost.cn/preview/56cc7639487d33cb/20154c443f1b0168)

## 👁️可观测性

日志记录

![image-20230530180305078](./images/image-20230530180305078.png)

jaeger 链路追踪，可以更便捷地进行debug

![image-20230530180353380](./images/image-20230530180353380.png)

![image-20230530180414994](./images/image-20230530180414994.png)



