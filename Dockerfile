FROM node:18.16.0 AS node
WORKDIR /app

COPY . ./

RUN npm install
RUN npx tailwindcss -i internal/static/style/input.css -o internal/static/style/output.css 

FROM golang:1.21.3 AS go
WORKDIR /app

COPY --from=node ./app/. ./

RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app ./cmd/app/.

FROM alpine:latest
WORKDIR /app

COPY --from=go ./app/app ./
COPY --from=go ./app/internal/static ./internal/static

EXPOSE 8080
ENV PORT=8080
CMD ["./app"]
