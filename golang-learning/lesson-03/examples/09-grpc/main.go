package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	mode := flag.String("mode", "server", "运行模式: server 或 client")
	addr := flag.String("addr", ":50051", "服务器地址 (server模式) 或连接地址 (client模式)")
	flag.Parse()
	
	switch *mode {
	case "server":
		log.Println("Starting gRPC server...")
		if err := startServer(*addr); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	case "client":
		log.Println("Starting gRPC client...")
		if *addr == ":50051" {
			*addr = "localhost:50051"
		}
		runClientDemo(*addr)
	default:
		log.Printf("Unknown mode: %s. Use 'server' or 'client'\n", *mode)
		os.Exit(1)
	}
}

