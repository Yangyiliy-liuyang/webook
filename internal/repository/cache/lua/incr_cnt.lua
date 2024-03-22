-- 具体业务
local key = KEYS[1]
-- 阅读数 点赞数 收藏数....
local cntKey = ARGV[1]

-- 点赞 取消点赞
local detal = tonumber(ARGV[2])

local exist = redis.call("EXISTS", key)

if exist == 1 then
    redis.call("HINCRBY", key, cntKey, detal)
    return 1
else
  return 0
end










