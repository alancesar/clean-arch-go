package user

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (UseCase) GetAllUsers() ([]User, error) {
	return []User{
		{
			ID:    1,
			Email: "someuser@mail.com",
		},
	}, nil
}
