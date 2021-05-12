// Code generated by roche

package usecase

type NekoUsecase struct {
	NekoRepo repository.INeko
}

func NewNekoUsecase(repo repository.INeko) repository.INeko {
	return &NekoUsecase{NekoRepo: repo}
}
func (u Neko) GetList() ([]*entity.Neko, error) {
	usecases, err := u.NekoRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return usecases, err
}
func (u Neko) GetByID(id int64) (*entity.Neko, error) {
	usecase, err := u.NekoRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return usecase, err
}
func (u Neko) Create(neko_id string, name string, size string) (*entity.Neko, error) {
	entity := &entity.Neko{
		Name:   Name,
		NekoId: NekoId,
		Size:   Size,
	}
	created, err := u.NekoRepo.Create(entity)
	if err != nil {
		return nil, err
	}
	return created, err
}
func (u Neko) Update(neko_id string, name string, size string, id int64) (*entity.Neko, error) {
	entity, err := u.NekoRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	entity = &entity.Neko{
		Name:   Name,
		NekoId: NekoId,
		Size:   Size,
	}
	updated, err := u.NekoRepo.Update(entity, id)
	if err != nil {
		return nil, err
	}
	return updated, err
}
func (u Neko) Delete(id int64) error {
	entity, err := u.NekoRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	err := u.NekoRepo.Delete(entity, id)
	if err != nil {
		return nil, err
	}
	return err
}
