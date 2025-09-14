package authenticated

type CheckAccessTokenPayload struct {
	Token     string
	ElementID string
}
