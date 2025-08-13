local processor = {}

function processor.pushMessage(id, content)
    if id == nil or id == "" then
        id = util.uuid()
    end

    log.info("processor", "process message, id=", id, "content=", content)
    success = fskv.set(id, {id=id, type="MESSAGE", body=content, retry=0})
    used, total, kv_count = fskv.status()
    log.info("processor", "add message to fskv=", success, "\nfskv status=", used, total, kv_count)
end

function processor.ackMessage(id)
    if id == nil or id == "" then
        return false
    end
    fskv.del(id)
    return true
end

local function processMessage()
    local queue = {}
    local used, total, kv_count = fskv.status()
    log.info("processor", "check fskv status: ", string.format("used= %d, total=%d, kv_count=%d", used, total, kv_count))
    local it = fskv.iter()
    if it then
        while true do
            local key = fskv.next(it)
            if not key then
                break
            end
            local value = fskv.get(key)
            if type(value) == "table" and value.id ~= nil and value.content ~= nil then
                table.insert(queue, value)
            else
                log.info("processor", "invalid fskv value: ", value)
            end
        end
    end
    -- start send to uart
    for k, v in pairs(queue) do
        local message = string.toBase64(json.encode(v))
        util.uart_send(v.id, "MESSAGE", v.body)
        local retry = v.retry or 0
        if fskv.get(v.id) ~= nil then
            fskv.sett(v.id, "retry", retry + 1)
        end
        sys.wait(500)
    end
end

-- 初始化
sys.taskInit(function()
    sys.waitUntil("IP_READY")
    sys.wait(10000)

    while true do
        processMessage()
        sys.wait(config.MESSAGE_PROCESS_INTERVAL)
    end
end)



return processor
