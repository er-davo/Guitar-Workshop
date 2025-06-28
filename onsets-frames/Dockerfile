FROM node:22-slim

WORKDIR /app

COPY app/package.json /app/

RUN npm install

COPY app/ /app/

CMD [ "node", "server.js" ]