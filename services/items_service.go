package services

type itemService struct{}
type itemServiceInterface interface {
	GetItem()
	SaveItem()
}

var ItemService itemServiceInterface = &itemService{}

func (*itemService) GetItem() {

}

func (*itemService) SaveItem() {

}
