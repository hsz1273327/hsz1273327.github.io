import Koa from 'koa'
import logger from 'koa-pino-logger'
import router from 'koa-route'
const app = new Koa()
app.use(logger())
app.use(router.get('/test', async ctx => {
    ctx.response.type = 'application/json'
    ctx.body = JSON.stringify({
        'message': 'Hello World'
    })
}))

app.listen(4000)