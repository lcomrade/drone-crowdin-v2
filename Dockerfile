# BUILD
FROM golang:1.18.8-alpine as build

WORKDIR /build

RUN apk update && apk upgrade && apk add --no-cache make git

COPY . ./

RUN make


# RUN
FROM alpine as run

WORKDIR /

COPY --from=build /build/dist/* /usr/local/bin/

CMD [ "/usr/local/bin/drone-crowdin-v2" ]
