FROM golang:1.20 as build
WORKDIR /build
COPY . .
RUN go build -o hello-server

FROM busybox
COPY --from=build /build/hello-server /bin
EXPOSE 8080
CMD [ "hello-server" ]