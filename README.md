# webook
webook 是一个仿简书的项目
gin + gorm + cors + regex + bcrypt + sessions + redis + JWT + sms + lua脚本 + wire
## 目录结构

## 启动
docker compose up

### 登录注册 JWT会话管理 限流 使用UserAgent增强JWT安全性
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


