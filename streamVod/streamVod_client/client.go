package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/BENSARI-Fathi/v1/videoStream/streamVod/streamVodpb"
	"gocv.io/x/gocv"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client Streaming Video Microservices")

	// open webcam
	webcam, _ := gocv.VideoCaptureDevice(0)
	defer webcam.Close()
	image := gocv.NewMat()
	var buffer []byte
	// Create a grpc client
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while trying to connect %v", err)
	}
	defer cc.Close()
	c := pb.NewVideoStreamServiceClient(cc)
	stream, err := c.ImageFrame(context.Background())
	if err != nil {
		panic(err)
	}
	for {
		webcam.Read(&image)
		buffer = image.ToBytes()
		err := stream.Send(
			&pb.ImageFrameRequest{
				Rows:  int32(image.Rows()),
				Cols:  int32(image.Cols()),
				Type:  int32(image.Type()),
				Frame: buffer,
			},
		)
		if err != nil {
			if err == io.EOF {
				log.Fatal("The Server close the connexion")
			}
			log.Fatalf("Error while sending video frame %v", err)
		}
	}

}
