FROM golang:1.17 as build
WORKDIR /app
COPY go.mod go.sum .

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o imageswap main.go

FROM gcr.io/distroless/base
COPY --from=build /app/imageswap /
EXPOSE 8443 8080

CMD ["/imageswap"]
