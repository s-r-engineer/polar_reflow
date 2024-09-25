FROM golang AS builder
ENV GO111MODULE="on"
ENV GOOS="linux"
ENV GOARCH="amd64"
COPY . /build/
WORKDIR /build
RUN go mod vendor
RUN go build -o ./polar_reflow main.go

FROM ubuntu
COPY --from=builder /build/polar_reflow /polar_reflow
CMD [ "/polar_reflow" ]