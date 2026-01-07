package main

import (
  "context"
  "fmt"
  "log"
  "net"
  "os"
  "strconv"
  "time"

  "google.golang.org/grpc"
  msg "imesh.ai/grpc-test/messaging"
)

const (
  ENV_HOST   = "HOST"
  ENV_PORT   = "PORT"
  ENV_REPLY  = "REPLY"
  ENV_DELAY  = "DELAY"  
)

var (
  host   string = "localhost"
  port   int    = 8080
  reply  string = "hello from server"
  delay  int    = 0
)

func parseEnv() error {
  envHost := os.Getenv(ENV_HOST)
  envPort := os.Getenv(ENV_PORT)
  envReply := os.Getenv(ENV_REPLY)
  envDelay := os.Getenv(ENV_DELAY)

  if envHost != "" {
    host = envHost
  }

  if envPort != "" {
    envPortInt, err := strconv.Atoi(envPort)
    if err != nil {
      return nil
    }
    port = envPortInt
  }

  if envReply != "" {
    reply = envReply
  }

  if envDelay != "" {
    envDelayInt, err := strconv.Atoi(envDelay)
    if err != nil {
      return fmt.Errorf("invalid DELAY value: %v", err)
    }
    delay = envDelayInt
  }

  return nil
}

type server struct {
  msg.UnimplementedMessagingServer
}

func (s *server) BasicRequestReply(_ context.Context, in *msg.BasicMessage) (*msg.BasicMessage, error) {
  log.Printf("Received: %v", in.GetMessage())

  
  if delay > 0 {
    log.Printf("Sleeping for %d seconds", delay)
    time.Sleep(time.Duration(delay) * time.Second)
  }

  return &msg.BasicMessage{Message: reply}, nil
}

func main() {
  err := parseEnv()
  if err != nil {
    log.Fatalf("failed parsing environment variables: %v", err)
  }

  lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  s := grpc.NewServer()
  msg.RegisterMessagingServer(s, &server{})
  log.Printf("started listening at %v (delay: %ds)", lis.Addr(), delay)

  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}
 
