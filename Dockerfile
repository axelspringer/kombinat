FROM alpine:3.6
MAINTAINER Sebastian Döll <sebastian.doell@axelspringer.com>

RUN \
    apk --update --no-cache add ca-certificates supervisor \
	&& mkdir -p /var/log/supervisor

ADD \
    /bin/kombinat_0.0.2_linux_amd64 /bin/kombinat

ADD \
    supervisord.conf /etc/supervisord.conf

RUN \
    chmod +x /bin/kombinat

EXPOSE 80
ENTRYPOINT ["/usr/bin/supervisord", "-n",  "-c", "/etc/supervisord.conf"]
