FROM debian:buster-slim

LABEL maintainer="tangxianmeng@agora.io"

COPY bin/rtm-sdk /

COPY internal/rtmlib/libagora_rtm_sdk.so /usr/lib

ENTRYPOINT ["/rtm-sdk"]