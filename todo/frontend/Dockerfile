FROM node:16-alpine3.16 AS build-stage

WORKDIR /usr/src/app

COPY package.json yarn.lock ./

RUN yarn install

COPY . .

ENV VITE_BACKEND_URL="/api"

RUN yarn build


FROM nginx:1.23.0-alpine

COPY --from=build-stage /usr/src/app/dist /usr/share/nginx/html
