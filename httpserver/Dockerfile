FROM ubuntu:20.04
ENV VERSION=test_httpserver_env
LABEL app="httpserver" author="winkyi"
ADD deploy/amd64/httpserver /httpserver
EXPOSE 80
EXPOSE 9999
EXPOSE 8001
ENTRYPOINT /httpserver