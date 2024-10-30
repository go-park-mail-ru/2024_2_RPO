package usecase

import "RPO_back/internal/pkg/board/repository"

type BoardUsecase struct {
	boardRepository *repository.BoardRepository
}

func CreateBoardUsecase(boardRepository *repository.BoardRepository) *BoardUsecase {
	return &BoardUsecase{
		boardRepository: boardRepository,
	}
}
