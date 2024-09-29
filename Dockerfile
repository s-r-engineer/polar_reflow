FROM golang AS builder
WORKDIR /build
COPY ./go.mod /build/go.mod
# COPY ./go.mod /build/go.sum
RUN go mod download
COPY . /build/
RUN go build -o ./polar_reflow main.go

FROM ubuntu
COPY --from=builder /build/polar_reflow /polar_reflow
# COPY  polar_reflow /polar_reflow
CMD [ "/polar_reflow" ]