FROM node:lts-alpine as builder

ARG NODE_ENV=production
ENV NODE_ENV $NODE_ENV

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

RUN npm run build

CMD ["npx", "serve", "-s", "/app/dist", "-l", "31338"]
