addCommandHandler('help', function(session)
    sendClientMessage(session, 'Available commands:')
    sendClientMessage(session, ' #reload - reload scripts')
    sendClientMessage(session, ' #getlevel - get current level')
    sendClientMessage(session, ' #setlevel <new_level> - set new level')
end)

addCommandHandler('reload', function(session)
    sendClientMessage(session, 'Reloading scripts...')
    reloadScripts()
end)

addCommandHandler('getlevel', function(session)
    local charLevel = getPlayerLevel(session)
    sendClientMessage(session, 'Your level is: '..charLevel)
end)

addCommandHandler('setlevel', function(session, arg)
    local level = tonumber(arg)
    if not level then
        sendClientMessage(session, 'Invalid command usage: #setlevel <new_level>')
        return
    end

    setPlayerLevel(session, level)
    sendClientMessage(session, 'Your level is set to '..level)
end)

addCommandHandler('drop', function(session, kind, opt)
    local kind_id = tonumber(kind)
    local opt_id = tonumber(opt)

    if not kind_id or not opt_id then
        sendClientMessage(session, 'Invalid command usage: #drop <kind> <opt>')
        return
    end

    local x, y = getPlayerPosition(session)
    local opcode = 204 -- NewItemList packet opcode
    local item_id = math.random(1, 100) -- use random

    -- some functions to fill buffer
    local function band(a, b)
        local result = 0
        local bitval = 1
        while a > 0 and b > 0 do
            if a % 2 == 1 and b % 2 == 1 then
                result = result + bitval
            end
            a = math.floor(a / 2)
            b = math.floor(b / 2)
            bitval = bitval * 2
        end
        return result
    end

    local function rshift(a, n)
        return math.floor(a / 2^n)
    end

    local function append(tbl, n, count)
        count = count or 4
        for i = 1, count do
            table.insert(tbl, band(rshift(n, (i - 1) * 8), 0xFF))
        end
    end

    local bytes = {}

    table.insert(bytes, 0x01) -- count
    append(bytes, item_id)    -- item idx
    append(bytes, opt)        -- item opt
    append(bytes, 0x0A000602) -- from idx
    append(bytes, kind)       -- kind Idx
    table.insert(bytes, x)    -- pos x
    table.insert(bytes, 0x00) --
    table.insert(bytes, y)    -- pos y
    table.insert(bytes, 0x00) --
    append(bytes, 0x330d, 2)  -- uniq key
    table.insert(bytes, 0x01) -- type
    table.insert(bytes, 0x04) -- unk

    sendClientPacket(session, opcode, bytes)
end)