package skills

type QuickLinkRequest struct {
	Id      int32
	Command string
	OldLink *Link
	NewLink *Link
}

type QuickLinkResponse struct {
	Result bool
}
