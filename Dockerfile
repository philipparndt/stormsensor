FROM node:20.5-alpine
COPY app/dist /opt/app/
WORKDIR /opt/app/

CMD ["node", "index.js", "/var/lib/storm-to-mqtt-gw/config.json"]
