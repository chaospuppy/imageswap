FROM golang:1.17 as build
WORKDIR /app
COPY . .
RUN go mod download && \
  CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o imageswap main.go

FROM gcr.io/distroless/base
COPY --from=build /app/imageswap /
EXPOSE 8443

CMD ["/imageswap"]
