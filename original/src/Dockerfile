# GO Repo base repo
FROM golang as builder

ENV GO111MODULE=on

# RUN apk add git

# Add Maintainer Info
LABEL maintainer="ezdeazy"

RUN apt-get -y install tzdata
RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod go.sum ./

# Download all the dependencies
RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# GO Repo base repo
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl tzdata

RUN mkdir /app

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8000
EXPOSE 8000

# Run Executable
CMD ["./main"]
