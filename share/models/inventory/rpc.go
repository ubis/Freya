package inventory

type ItemRequest struct {
	Id      int32
	Command string
	Item    Item
	NewItem *Item
}

type ItemResponse struct {
	Result bool
	Item   *Item
}
