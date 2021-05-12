package testdata

import "context"

type Handler struct {}

func SearchNekos(context.Context, *ListNekosRequest) (*ListNekosResponse, error) {
	return nil, nil
}

func GetNeko(context.Context, *GetNekoRequest) (*Neko, error) {
	return nil, nil
}

func CreateNeko(context.Context, *CreateNekoRequest) (*Neko, error) {
	return nil, nil
}

func UpdateNeko(context.Context, *UpdateNekoRequest) (*Neko, error) {
	return nil, nil
}
func DeleteNeko(context.Context, *DeleteNekoRequest) (*emptypb.Empty, error) {
	return nil, nil
}