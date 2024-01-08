# BUILD
FROM docker.io/library/golang:1.21.5 as build

WORKDIR /build

COPY . ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o ./drone-crowdin-v2 ./cmd/drone-crowdin-v2/*.go


# RUN
FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/drone-crowdin-v2 /
COPY --from=build /tmp /tmp

CMD [ "/drone-crowdin-v2" ]
