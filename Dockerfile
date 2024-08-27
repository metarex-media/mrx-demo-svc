# syntax=docker/dockerfile:1

FROM golang:1.22

# Set destination for
WORKDIR /workspace/mrx-demo-svc

# Copy all files within the repo
COPY . ./

# Build the main server
RUN CGO_ENABLED=0 GOOS=linux go build -o mrx-demo-svc
# Enable things for the user to use
RUN chmod u+x ./clogrc/run.sh
RUN go mod tidy

# build the services API
WORKDIR /workspace/mrx-demo-svc/api
RUN go build

# Get back to the start
WORKDIR /workspace/mrx-demo-svc

# Expose the main server port
# Not the api server
EXPOSE 8080



# add user for security
RUN useradd -m appuser
USER appuser

# health check to ensure both servers are running etc
HEALTHCHECK --interval=180s --timeout=90s \ 
    CMD curl --fail http://localhost:8080/test || exit 1

# Run both servers using the run script
CMD ["bash", "-c", "echo $PATH && ./clogrc/run.sh"]
