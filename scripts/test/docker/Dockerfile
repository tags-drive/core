FROM golang:1.12-alpine

# Build env vars
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# Test args
ARG TEST_STORAGE_S3_ENDPOINT
ARG TEST_STORAGE_S3_ACCESS_KEY_ID
ARG TEST_STORAGE_S3_SECRET_ACCESS_KEY
ARG TEST_STORAGE_S3_SECURE

WORKDIR /build/test

# Copy source files
COPY . .

CMD [ "go", "test", "--mod=vendor", "-v", "./..." ]
