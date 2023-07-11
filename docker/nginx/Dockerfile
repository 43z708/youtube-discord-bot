FROM nginx
COPY default.conf /etc/nginx/conf.d/
ARG USERNAME
ARG PASSWORD
RUN echo "${USERNAME}:$(openssl passwd -apr1 ${PASSWORD})" > /etc/nginx/.htpasswd
