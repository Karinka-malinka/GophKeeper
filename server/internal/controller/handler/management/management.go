// Package management предоставляет реализацию методов для работы различными типами приватных данных
package management

import (
	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/logindata"
)

// ManagementServer представляет обработчик для управления приватными данными.
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
