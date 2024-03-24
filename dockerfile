FROM golang:1.19-alpine as builder
WORKDIR /app
ADD . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./auth-service-build *.go

#Run stage
FROM alpine:3.13
WORKDIR /app
COPY --from=builder ./app/ .

EXPOSE 3000

CMD [ "./auth-service-build" ]