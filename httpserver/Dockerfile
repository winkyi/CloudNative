FROM ubuntu:20.04
ENV VERSION=test_httpserver_env
LABEL app="httpserver" author="winkyi"
ADD deploy/amd64/httpserver /usr/bin/httpserver
EXPOSE 9999
EXPOSE 8001
ENTRYPOINT /usr/bin/httpserver