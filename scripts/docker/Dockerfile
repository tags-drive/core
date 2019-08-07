FROM golang:1.12-alpine as builder

ENV CGO_ENABLED=0

RUN mkdir /build

# Copy code to /build ("docker build" must be ran in root folder)
COPY . /build

# Build binary
RUN cd /build && \
	go build -o tags-drive --mod=vendor main.go


FROM alpine

RUN mkdir /app
WORKDIR /app

COPY --from=builder /build/tags-drive .
COPY ./scripts/docker/run.sh .

ENTRYPOINT [ "./run.sh" ]
