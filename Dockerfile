FROM ccr.ccs.tencentyun.com/dhub.wallstcn.com/alpine:3.5

#ENV CONFIGOR_ENV ivktest
#ADD resources/* /
ADD conf/ /conf
ADD public/ public
ADD template /template
ADD server /
ENTRYPOINT [ "/server" ]
