package main

import (
	"context"
	"log"
	"net"
	"sync"

	//import generate proto code via github

	pb "github.com/fahruroze/project_ta/proto/consignment"

	"google.golang.org/grpc"
)

const ( //deklarasi port yang akan digunakan
	port = ":50051"
)

type repository interface { //buat repo baru sebagai interface dari pb
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

//dummy repo
type Repository struct {
	mu           sync.RWMutex //
	consignments []*pb.Consignment
}

//Method Create new Consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	update := append(repo.consignments, consignment)
	repo.consignments = update
	repo.mu.Unlock()
	return consignment, nil
}

//Method GetAll Consignment
func (repo *Repository) GetAll() []*pb.Consignment { //deklarasi method GetAll
	return repo.consignments //return nya semua data pengiriman
}

//Service harus mengimplementasiksan semua method yang digenerate oleh protobuf

type service struct {
	repo repository
}

//CreateConsignment, method yang dibuat
// method ini bisa digunakan setelah digenerate oleh gRPC

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	//save pengiriman
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	//kirim Response sebagi nilai balik, sbg mana yg telah dididefinisikan di protobuf
	return &pb.Response{Created: true, Consignment: consignment}, nil

}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	return &pb.Response{Consignments: consignments}, nil
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

	log.Println("Running on Port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
