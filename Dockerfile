# syntax=docker/dockerfile:1

FROM golang:1.22

# Set destination for
WORKDIR /workspace/gl-mrx-demo-svc

# Copy all files within the repo
COPY . ./

# Build the main server
RUN CGO_ENABLED=0 GOOS=linux go build -o mrx-elt-demo
# build the services API
RUN cd api && go build


# Expose the main server port
# Not the api server
EXPOSE 8080

# Run both servers using the run script
CMD chmod u+x ./clogrc/run.sh && ./clogrc/run.sh
