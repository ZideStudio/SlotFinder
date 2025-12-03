# Stage 1: Extract Volta versions from package.json
FROM node:alpine AS volta-reader

WORKDIR /app

COPY package.json ./

# Extract and validate Volta configuration
RUN node << 'EOF' > /tmp/volta-info.sh
const pkg = require('./package.json');
if (!pkg.volta || !pkg.volta.node || !pkg.volta.npm) {
    console.error('Error: Volta configuration not found in package.json');
    process.exit(1);
}
console.log('export NODE_VERSION="' + pkg.volta.node + '"');
console.log('export NPM_VERSION="' + pkg.volta.npm + '"');
console.log('export NODE_MAJOR="' + pkg.volta.node.split('.')[0] + '"');
EOF

# Stage 2: Builder with Volta-specified Node.js version
# NOTE: Update the version tag below if Volta's Node.js major version changes in package.json
# The validation step below will warn if there's a mismatch
FROM node:24-alpine AS builder

WORKDIR /app

# Copy the Volta version info
COPY --from=volta-reader /tmp/volta-info.sh /tmp/

# Validate Node.js version and install the correct npm version from Volta
RUN . /tmp/volta-info.sh && \
    ACTUAL_NODE=$(node --version | sed 's/v//') && \
    echo "Volta specifies Node.js ${NODE_VERSION}, using ${ACTUAL_NODE}" && \
    if [ "${ACTUAL_NODE%%.*}" != "${NODE_MAJOR}" ]; then \
        echo "WARNING: Node.js major version mismatch. Update FROM node:${NODE_MAJOR}-alpine in Dockerfile"; \
    fi && \
    echo "Installing npm@${NPM_VERSION} from Volta configuration" && \
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
