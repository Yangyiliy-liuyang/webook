-- 发送到的key code:业务:手机号码
local key = KEYS[1]
-- 验证次数
local cntKey = key..":cnt"
-- 准备存储的验证码
local val = ARGV[1]
local ttl = tonumber(redis.call("ttl",key))
if ttl == -1 then
    -- key存在，但是没有过期时间
    return -2
elseif ttl == -2 or ttl < 540 then
    -- 可以发验证码
    redis.call("set",key,val)
    -- 600秒
    redis.call("expire",key,600)
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)
else
    -- 发送太频繁
    return -1
end
