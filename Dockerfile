# BUILD
FROM golang:1-alpine as build

WORKDIR /build

COPY . ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o ./drone-crowdin-v2 ./cmd/drone-crowdin-v2/*.go


# RUN
FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/drone-crowdin-v2 /

CMD [ "/drone-crowdin-v2" ]
