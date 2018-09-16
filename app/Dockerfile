FROM alpine
# Update
RUN apk update && apk upgrade
# Create /app
RUN mkdir /app && mkdir /app/static &&  mkdir /app/templates
WORKDIR /app
# Copy (folders templates and static must be copied into build folder)
COPY tags-drive .
COPY static ./static
COPY templates ./templates
ENTRYPOINT [ "./tags-drive" ]