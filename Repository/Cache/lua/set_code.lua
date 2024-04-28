
-- 存储的验证码
local key = KEYS[1]
-- 使用的次数
local cntKey = key..":cnt"

local val = ARGV[1]

-- 拿到Redis中的过期时间
local ttl = tonumber(redis.call("ttl",key))
if ttl == -1 then
    -- key存在，但是没有过期时间
    return -2 --go会拿到的返回码
elseif ttl == -2 or ttl < 540 then
    -- 可以发验证码
    redis.call("set",key,val)
    redis.call("expire",key,600) -- 600秒的过期时间
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)
    return 0
else
    -- 发送太频繁
    return -1
end