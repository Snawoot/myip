FROM golang AS build

WORKDIR /go/src/github.com/Snawoot/myip
COPY . .
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-s -w -extldflags "-static"'

FROM scratch
COPY --from=build /go/src/github.com/Snawoot/myip/myip /
ENTRYPOINT ["/myip"]
