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