package controller

import (
	pb "github.com/kujilabo/cocotola/lib/proto"
)

type adminServer struct {
	pb.UnimplementedTranslatorAdminServer
}
