FROM alpine:3.13
RUN apk update && apk --no-cache add curl && apk --no-cache add tzdata
ADD main /main
ENTRYPOINT [ "/main" ]
