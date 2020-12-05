FROM --platform=$TARGETPLATFORM  node:latest

ADD server.js /app/server.js
ADD package.json /app/package.json
WORKDIR /app
#安装依赖
RUN npm config set registry https://registry.npm.taobao.org
RUN npm install
#对外暴露的端口
EXPOSE 3000
HEALTHCHECK --interval=1m30s --timeout=10s --start-period=10s --retries=3 CMD [ "curl","http://localhost:3000/ping" ]
CMD [ "npm","start"]