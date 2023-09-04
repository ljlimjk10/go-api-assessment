package models

type Student struct {
	ID          uint   `gorm:"primaryKey"`
	Email       string `gorm:"not null;unique"`
	IsSuspended bool
}

type TeacherStudentRelation struct {
	ID           uint    `gorm:"primaryKey"`
	TeacherEmail string  `gorm:"not null"`
	StudentID    uint    `gorm:"not null"`
	Student      Student `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE"`
}
