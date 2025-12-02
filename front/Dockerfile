# Build arguments for Volta versions (will be set dynamically)
ARG NODE_VERSION=24.11.1
ARG NPM_VERSION=11.6.4

FROM node:${NODE_VERSION}-alpine AS builder

WORKDIR /app

# Install the specific npm version provided via build arg (extracted from Volta config by Makefile)
RUN npm install -g npm@${NPM_VERSION}

COPY package*.json ./

RUN npm ci

COPY . ./

RUN npm run build

FROM nginx:1.29-alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]
