FROM golang:1.12-alpine as builder

ENV CGO_ENABLED=0

# Copy code to /temp ("docker build" must be ran in root folder)
COPY . /temp

# Build binary
RUN cd /temp && go build -o tags-drive --mod=vendor


FROM alpine

RUN mkdir /app
WORKDIR /app

COPY --from=builder /temp/tags-drive .
COPY ./scripts/docker/run.sh .

ENTRYPOINT [ "./run.sh" ]
