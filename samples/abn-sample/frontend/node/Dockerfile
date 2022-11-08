
FROM node:18
# RUN mkdir -p /usr/app
WORKDIR /usr/app
COPY samples/abn-sample/frontend/node/package*.json .

RUN npm install

COPY samples/abn-sample/frontend/node .

CMD ["node", "main.js"]