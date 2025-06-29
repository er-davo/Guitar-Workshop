FROM node:22-slim

WORKDIR /app

COPY app/package.json /app/
COPY app/package-lock.json /app/

RUN npm install --force --unsafe-perm

COPY app/ /app/

# CMD [ "node", "server.js" ]

CMD [ "tail", "-f", "/dev/null" ]