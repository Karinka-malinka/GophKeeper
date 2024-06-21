package management

import (
	pb "github.com/GophKeeper/server/cmd/proto"
	"github.com/GophKeeper/server/internal/app/logindata"
)

type ManagementServer struct {
	pb.UnimplementedManagementServiceServer
	LoginDataApp *logindata.LoginDatas
	//TODO: добавить app на все типы хранимых данных
	//TextDataApp *textdata.TextDatas
}

// NewManagementHandler создает новый экземпляр ManagementServer.
// TODO: добавить app на все типы хранимых данных
func NewManagementHandler(loginData *logindata.LoginDatas) *ManagementServer {
	return &ManagementServer{
		LoginDataApp: loginData,
	}
}
