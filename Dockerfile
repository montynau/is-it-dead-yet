# Use the official Go runtime as a base image
FROM golang:1.21-alpine
# Set the working directory
WORKDIR /app
# Copy all project files
COPY . .
# Compile the program
RUN go build -o is-it-dead-yet .
# Set the container to listen on port 8080
EXPOSE 8080
# Run the program
CMD ["./is-it-dead-yet"]