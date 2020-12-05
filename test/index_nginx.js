import WebSocket from 'ws'


const ws = new WebSocket('ws://localhost:8000')
ws.on('open', () => {
    ws.send('helloworld')
})
ws.on('message', (data) => {
    console.log(data)
    ws.close()
})
ws.on('close', () => {
    console.log('disconnected');
})