FROM nginx:alpine
COPY config/nginx.conf /etc/nginx/nginx.conf
COPY config/conf.d /etc/nginx/conf.d
EXPOSE 80
STOPSIGNAL SIGTERM
CMD [ "nginx","-g", "daemon off;"]