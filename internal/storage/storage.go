package storage

import (
	types "github.com/araza95/learn0x-GO/internal/models"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
}
