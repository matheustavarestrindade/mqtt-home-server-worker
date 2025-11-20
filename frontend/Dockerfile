# Stage 1: Build the application
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files first to leverage Docker caching
COPY package.json package-lock.json ./

# Install dependencies. 
# We use --ignore-scripts to prevent the 'prepare' script (svelte-kit sync) 
# from failing because svelte.config.js isn't copied yet.
RUN npm ci --ignore-scripts

# Copy the rest of the application code
COPY . .

# Build the application
# This will output the server to the 'build' directory by default with adapter-node
RUN npm run build

# Stage 2: Run the application
FROM node:22-alpine AS runner

WORKDIR /app

# Set environment variables
ENV SERVICE_BASE_URL=
ENV DATABASE_URL=
ENV NODE_ENV=production
ENV HOST=0.0.0.0 
ENV PORT=3000
# You may need to set ORIGIN if you encounter CSRF errors with forms
# HOST=0.0.0.0 is crucial for Docker networking so the container can be accessed
# ENV ORIGIN=http://localhost:3000

# Copy package files to install only production dependencies
COPY package.json package-lock.json ./

# Install only production dependencies (saves space)
RUN npm ci --omit=dev --ignore-scripts

# Copy the built application from the builder stage
COPY --from=builder /app/build ./build

# Expose the port
EXPOSE 3000

# Start the server
CMD ["node", "build"]
