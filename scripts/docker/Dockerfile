FROM alpine

RUN mkdir /app
WORKDIR /app

COPY tags-drive .
COPY run.sh .

ENTRYPOINT [ "./run.sh" ]