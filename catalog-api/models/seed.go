package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedDB(db *gorm.DB) {
	var count int64
	db.Model(&Category{}).Count(&count)
	if count > 0 {
		return
	}

	mathID := uuid.MustParse("e40bfc37-400a-4d90-b8ef-6fc109a716eb")
	scienceID := uuid.MustParse("6b85a8ae-72cd-4273-b8fb-d0292633344a")

	math := Category{ID: mathID, Name: "Math"}
	science := Category{ID: scienceID, Name: "Science"}

	db.Create(&[]Category{math, science})

	courses := []Course{
		{ID: uuid.New(), Title: "Intro to Algebra", CategoryID: mathID},
		{ID: uuid.New(), Title: "Physics Basics", CategoryID: scienceID},
	}

	db.Create(&courses)
}
