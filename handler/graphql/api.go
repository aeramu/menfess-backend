package graphql

type Err struct {
	IsError bool
	Message string
}

func Error(err error) Err {
	return Err{
		IsError: true,
		Message: err.Error(),
	}
}

var NoError = Err{
	IsError: false,
	Message: "",
}

type AuthResponse struct {
	Token string
	Error Err
}

type BasicMutationResponse struct {
	Message string
	Error   Err
}

type PostResponse struct {
	Payload Post
	Error Err
}

type FeedResponse struct {
	Payload PostConnection
	Error Err
}

type MeResponse struct {
	Payload User
	Error Err
}

type MenfessResponse struct {
	Payload UserConnection
	Error Err
}