# Docker File
# RIGHT NOW IS BROKEN
FROM openjdk:17

WORKDIR /app
COPY out/paper-autoupdater-linux.bin /app/paper-autoupdater
COPY start.bash /app/start.bash
