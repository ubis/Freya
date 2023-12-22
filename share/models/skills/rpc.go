package skills

type QuickLinkRequest struct {
	Server  byte
	Id      int32
	Command string
	OldLink *Link
	NewLink *Link
}

type QuickLinkResponse struct {
	Result bool
}
