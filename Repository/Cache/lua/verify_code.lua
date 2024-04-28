-- 存储的验证码
local key = KEYS[1]
-- 使用的次数
local cntKey = key..":cnt"

local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get",cntKey))
local code = redis.call("get",key)

--cnt == nil说明字段不存在，系统发生错误
if cnt == nil or cnt <= 0 then
--验证次数耗尽
    return -1
end

if code == expectedCode then
    redis.call("set",cntKey,-1)
    return 0
else
    redis.call("decr",cntKey)
    --不相等，用户输入错误了
    return -2
end