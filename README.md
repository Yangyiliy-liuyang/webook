# webook
webook 是一个仿简书的项目
gin + gorm + cors + regex + bcrypt + sessions 
+ redis + JWT 
+ sms + lua脚本 + wire + mock + sqlmock 
+ 
## 目录结构
## 启动
docker compose up
### 登录注册 JWT会话管理 限流 使用UserAgent增强JWT安全性
领域对象domain.User·业务概念  数据库对象dao.User·直接映射到表结构
- 时区问题int64，统一使用UTC0的时区，返回数据才处理 服务器 go应用 数据库
- 唯一索引冲突unique
- Region 使用含糊地区代表
- 联表 使用json代表Addr string或者反向持有uid
为什么使用自增主键？
- 数据库中的数据存储是一个树型结构，自增意味着树朝一个方向增长，id相邻的大概率在磁盘上也是相邻的
，充分利用操作系统预读机制。
不是自增则意味中间插入数据，页分页
### docker kubernetes部署
### 对Profile redis缓存 异步
### 短信登录 验证码发送 校验 使用Redis缓存 lua脚本
业务规则
- 一个手机号码，一分钟只能发一次，有效时间十分钟（本身系统有限流）
- 验证通过了，不可以再用
- 三次验证失败后，不可以再用
- 验证失败，只提示用户验证码不对，不提示过于频繁的失败原因
业务层上、分布式环境下并发，不是语言层面上的并发,不可以通过channel 或者 sync.lock 
检查数据-doingSth 考虑并发安全
使用分布式锁或者lua脚本
### DDD领域理论  REST风格 依赖注入
### 面向接口编程 超前设计 最小化实现
### 单元测试 集成测试 mock sqlmock
### 微信登录


