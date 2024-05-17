# Use a Node.js base image
FROM node:latest as builder

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY frontend/package.json .
COPY frontend/pnpm-lock.yaml .

# Install dependencies using PNPM
RUN npm install -g pnpm
RUN pnpm install

# Copy the rest of the application code
COPY frontend .

# Build the React app
RUN pnpm run build


# Use the official Go image
FROM golang:latest

# Set working directory
WORKDIR /go/src/app


# Copy the local package files to the container's workspace
COPY backend .

# Copy frontend to web
COPY --from=builder /app/dist ./web

# Install Iris
RUN go get -u github.com/kataras/iris/v12

# Build the Go application
RUN go build -o main cmds/web/main.go
RUN go build -o food-cli cmds/cli/main.go
# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
