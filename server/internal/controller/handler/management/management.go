// Package management предоставляет реализацию методов для работы различными типами приватных данных
package management

import (
	pb "github.com/GophKeeper/api/proto"
	"github.com/GophKeeper/server/internal/app/bankcard"
	"github.com/GophKeeper/server/internal/app/filedata"
	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/GophKeeper/server/internal/app/textdata"
)

// ManagementServer представляет обработчик для управления приватными данными.
type ManagementServer struct {
	pb.UnimplementedManagementServiceServer
	LoginDataApp    *logindata.LoginDatas
	TextDataApp     *textdata.TextDatas
	FileDataApp     *filedata.FileDatas
	BankCardDataApp *bankcard.BankCardDatas
}

// NewManagementHandler создает новый экземпляр ManagementServer.
func NewManagementHandler(loginData *logindata.LoginDatas, textData *textdata.TextDatas, fileData *filedata.FileDatas, bankCard *bankcard.BankCardDatas) *ManagementServer {
	return &ManagementServer{
		LoginDataApp:    loginData,
		TextDataApp:     textData,
		FileDataApp:     fileData,
		BankCardDataApp: bankCard,
	}
}
