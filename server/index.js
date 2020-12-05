import WebSocket from 'ws'

const wsserver = new WebSocket.Server({
    port: 3000
})


wsserver.on('connection', ws => {
    ws.on('message', message => {
        console.log('received: %s', message)
        switch (message) {
            case "close":
                {
                    ws.close()
                }
                break
            case "helloworld":
                {
                    ws.send('Hello World')
                }
                break
            default:
                {
                    ws.send('unkonwn command')
                }
        }
    })
})
