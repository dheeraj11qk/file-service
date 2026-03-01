package main

import (
	"context"
	"io"
	"log"
	"os"

	"file-service/proto"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	client := proto.NewFileServiceClient(conn)

	stream, err := client.UploadFile(context.Background())
	if err != nil {
		log.Fatal("Error creating stream:", err)
	}

	// Open file to upload
	file, err := os.Open("cmd/client/test.txt")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading file:", err)
		}

		err = stream.Send(&proto.UploadFileRequest{
			Filename:  "test.txt",
			ChunkData: buffer[:n],
		})
		if err != nil {
			log.Fatal("Error sending chunk:", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Error receiving response:", err)
	}

	log.Println("Upload Response:", res.FileMessage)
	log.Println("File ID:", res.FileId)
}
