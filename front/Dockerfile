# Stage 1: Extract Volta versions from package.json
FROM node:alpine AS volta-reader

WORKDIR /app

COPY package.json ./

# Extract and validate Volta configuration, then determine the correct base image tag
RUN node << 'EOF' > /tmp/volta-info.sh
const pkg = require('./package.json');
if (!pkg.volta || !pkg.volta.node || !pkg.volta.npm) {
    console.error('Error: Volta configuration not found in package.json');
    process.exit(1);
}
console.log('export NODE_VERSION=' + pkg.volta.node);
console.log('export NPM_VERSION=' + pkg.volta.npm);
EOF

# Stage 2: Builder with extracted Volta versions  
FROM node:alpine AS builder

WORKDIR /app

# Copy the Volta version info
COPY --from=volta-reader /tmp/volta-info.sh /tmp/

# Source the versions and install the correct npm version
RUN . /tmp/volta-info.sh && \
    echo "Using npm version from Volta: ${NPM_VERSION}" && \
    npm install -g npm@${NPM_VERSION}

# Install dependencies
COPY package*.json ./
RUN npm ci

# Build application
COPY . ./
RUN npm run build

# Stage 3: Production image
FROM nginx:1.29-alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]
