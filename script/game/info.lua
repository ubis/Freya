local function getRuntimeOS()
    local OS = os.getenv("OS")

    if OS == "Windows_NT" then
        return "Windows"
    elseif os.execute("uname -s >/dev/null") then
        local uname = io.popen("uname -s"):read("*l")

        if uname == "Linux" then
            return "Linux"
        elseif uname == "Darwin" then
            return "macOS"
        end
    end

    return "Unknown"
end

local function sendServerInfo(session)
    local go = getGoVersion()
    local os = getRuntimeOS()
    local build = getBuildInfo()

    sendClientMessage(session, 'Welcome to Freya - CABAL Server Emulator!')
    sendClientMessage(session,
        'Running on '..os..' OS with ' ..go..' and '.._VERSION)
    sendClientMessage(session, 'Build: #'..build)
    sendClientMessage(session, '')
	sendClientMessage(session, 'Type #help to see available commands.')
end

addEventHandler('onPlayerJoin', function(session)
    sendServerInfo(session)
end)

addCommandHandler('help', function(session)
    sendClientMessage(session, 'Available commands:')
    sendClientMessage(session, ' #reload - reload scripts')
end)

addCommandHandler('reload', function(session)
    sendClientMessage(session, 'Reloading scripts...')
    reloadScripts()
end)