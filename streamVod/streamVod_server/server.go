package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/BENSARI-Fathi/v1/videoStream/streamVod/streamVodpb"
	"gocv.io/x/gocv"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedVideoStreamServiceServer
}

func (s *Server) ImageFrame(reqStream pb.VideoStreamService_ImageFrameServer) error {
	fmt.Println("ImageFrame Method Was Invoked")
	var (
		rows, cols, types int
		buffer            []byte
	)
	//opencv setup
	window := gocv.NewWindow("Client stream Video")
	defer window.Close()

	for {
		frameData, err := reqStream.Recv()
		if err == io.EOF {
			return reqStream.SendAndClose(&pb.ImageFrameResponse{
				StatusCode: pb.Status_GoodStream,
			})
		}
		if err != nil {
			return err
		}
		rows = int(frameData.GetRows())
		cols = int(frameData.GetCols())
		types = int(frameData.GetType())
		buffer = frameData.GetFrame()
		image, err := gocv.NewMatFromBytes(rows, cols, gocv.MatType(types), buffer)
		if err != nil {
			panic(err)
		}
		defer image.Close()
		if image.Empty() {
			continue
		} else {
			window.IMShow(image)
		}
		if window.WaitKey(1) >= 0 {
			fmt.Println("Closing the client stream !")
			return reqStream.SendAndClose(&pb.ImageFrameResponse{
				StatusCode: pb.Status_GoodStream,
			})
		}

	}

}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Can't listen %v", err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	pb.RegisterVideoStreamServiceServer(s, &Server{})
	log.Println("Grpc server Listening on 0.0.0.0:50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v\n", err)
	}
}
