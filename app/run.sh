#!bin/sh
# Example of run.sh script
docker run --rm -d \
--name tags-drive \
-p 80:80 \
-v /home/username/configs:/app/configs \
-v /home/username/data:/app/data \
--env-file /home/username/tags-drive.env \
kirtis/tags-drive