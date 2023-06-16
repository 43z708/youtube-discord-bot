FROM mariadb:10.3

RUN mv /etc/localtime /etc/localtime.org && \
    ln -s /usr/share/zoneinfo/Japan /etc/localtime

COPY ./conf.d /etc/mysql/conf.d