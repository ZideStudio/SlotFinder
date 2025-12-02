FROM node:24.11.1-alpine AS builder

WORKDIR /app

# Install specific npm version to match Volta configuration
RUN npm install -g npm@11.6.4

COPY package*.json ./

RUN npm ci

COPY . ./

RUN npm run build

FROM nginx:1.29-alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]
