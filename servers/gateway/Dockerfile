# Created image based on alpine
# add ca-cerificates 
# for Go API gateway server
FROM alpine
RUN apk update && apk add --no-cache ca-certificates
COPY gateway /gateway
EXPOSE 443
ENTRYPOINT [ "/gateway" ]