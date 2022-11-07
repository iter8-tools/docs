
FROM node:18
# RUN mkdir -p /usr/app
WORKDIR /usr/app
COPY frontend/node/package*.json .

RUN npm install

COPY frontend/node .

CMD ["node", "main.js"]