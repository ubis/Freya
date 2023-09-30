package inventory

type ItemRequest struct {
	Server  byte
	Id      int32
	Command string
	Item    Item
	NewItem *Item
}

type ItemResponse struct {
	Result bool
	Item   *Item
}
