FROM golang:1.18.0-alpine3.15 AS build

WORKDIR /app
RUN apk add --no-cache make
COPY . .

RUN CGO_ENABLED=0 make clean build

FROM scratch
WORKDIR /app
COPY --from=build /app/_build/ /app/
EXPOSE 8080
CMD ["/app/moekyuniversity"]