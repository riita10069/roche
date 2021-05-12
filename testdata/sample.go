package testdata

type NekoServiceServer interface {
	SearchNekos(context.Context, *ListNekosRequest) (*ListNekosResponse, error)
	GetNeko(context.Context, *GetNekoRequest) (*Neko, error)
	CreateNeko(context.Context, *CreateNekoRequest) (*Neko, error)
	UpdateNeko(context.Context, *UpdateNekoRequest) (*Neko, error)
	DeleteNeko(context.Context, *DeleteNekoRequest) (*emptypb.Empty, error)
	CancelCreateNeko(context.Context, *DeleteNekoRequest) (*emptypb.Empty, error)
}
