package main

import (
	"context"
	"log"
	"net"
	"sync"

	//import generate proto code
	pb "project/proto/consignment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const ( //deklarasi port yang akan digunakan
	port = ":50051"
)

type repository interface { //buat repo baru sebagai interface dari pb
	Create(*pb.Consignment) (*pb.Consignment, error)
}

//dummy repo
type Repository struct {
	mu         sync.RWMutex //
	consigment []*pb.Consignment
}

//create new Consignment
func (repo *Repository) Create(consigment *pb.Consigment) (*pb.Consignment, error) {
	repo.mu.Lock()
	update := append(repo.consigments, consigment)
	repo.consigments = update
	repo.mu.Unlock()
	return consigment, nil
}

//Service harus mengimplementasiksan semua method yang digenerate oleh protobuf

type service struct {
	repo repository
}

//CreateConsignment, method yang dibuat
// method ini bisa digunakan setelah digenerate oleh gRPC

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consigment) (*pb.Response, error) {
	//save pengiriman
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	//kirim Response sebagi nilai balik, sbg mana yg telah dididefinisikan di protobuf
	return &pb.Response{Created: true, Consignment: consignment}, nil

}

func main() {
	repo := &Repository{}

	//setup our gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	//DEKLARASI service yang udah dibikin di file proto

	pb.RegisterShippingServiceServer(s, &service{repo})

	//reflection service di gRPC srver
	reflection.Register(s)

	log.Println("Running on Port:", port)
	if err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
