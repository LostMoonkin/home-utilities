PROJECT = "Air780e_Uart_SMS_Forward"
VERSION = "1.0.0"

log.setLevel("DEBUG")
log.info("main", PROJECT, VERSION)

sys = require("sys")
config = require("config")
util = require("util")
processor = require("processor")

-- 添加硬狗防止程序卡死
wdt.init(9000)
sys.timerLoopStart(wdt.feed, 3000)

-- SIM 自动恢复, 周期性获取小区信息, 网络遇到严重故障时尝试自动恢复等功能
mobile.setAuto(10000, 30000, 8, true, 60000)

-- 串口初始化
uart.setup(1, 115200, 8, 1, uart.NONE)

-- fskv 初始化
if fskv.init() then
    log.info("fskv", "fskv init success.")
end

-- 模块初始化
sys.taskInit(function()
    -- 等待网络环境准备就绪
    sys.waitUntil("IP_READY", 1000 * 60 * 5)
    local imei = mobile.imei()
    local number = mobile.number()
    local status = mobile.status()
    log.info("main", "device startup: ", string.format("Device startup, imei=%s, number=%s, status=%d", imei, number, status))
    util.uart_send("", "MESSAGE", {imei = imei, number = number, status = status})

    sys.wait(60000);
    -- EC618配置小区重选信号差值门限，不能大于15dbm，必须在飞行模式下才能用
    mobile.flymode(0, true)
    mobile.config(mobile.CONF_RESELTOWEAKNCELL, 10)
    mobile.config(mobile.CONF_STATICCONFIG, 1) -- 开启网络静态优化
    mobile.flymode(0, false)
end)

-- 定时开关飞行模式
if type(config.FLYMODE_INTERVAL) == "number" and config.FLYMODE_INTERVAL >= 1000 * 60 then
    sys.timerLoopStart(function()
        sys.taskInit(function()
            log.info("main", "change flymode.")
            mobile.reset()
            sys.wait(1000)
            mobile.flymode(0, true)
            mobile.flymode(0, false)
        end)
    end, config.FLYMODE_INTERVAL)
end

-- 新短信处理
sms.setNewSmsCb(function(num, txt, metas)
    -- num 手机号码
    -- txt 文本内容
    -- metas 短信的元数据,例如发送的时间,长短信编号
    -- 注意, 长短信会自动合并成一条txt
    log.info("main", "process message callback: ", num, txt, metas and json.encode(metas) or "")
    processor.pushMessage("", {
        number = num,
        txt = txt,
        metas = metas,
    })
end)

-- 后端消息处理
uart.on(1, "receive", function(id, len)
    local data = uart.read(id, len)
    log.debug("main", "uart receive: ", id, len, data)
    -- 短信ack
    local dataId, dataType, body = string:match("([^,]+):([^,]+):(.*)")
    if dataType == "ACK" then
        log.debug("main", "ACK message: ", dataId)
        processor.ackMessage(dataId)
    end
end)

-- 心跳包
if config.ENABLE_HEART_BEAT then
    sys.timerLoopStart(function()
        local imei = mobile.imei()
        local number = mobile.number()
        local status = mobile.status()
        log.info("main", "heart beat: ", string.format("Device startup, imei=%s, number=%s, status=%d", imei, number, status))
        util.uart_send("", "MESSAGE", {imei = imei, number = number, status = status})
    end, config.HEART_BEAT_INTERVAL)
end


sys.run()

