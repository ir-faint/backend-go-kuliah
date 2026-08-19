package main

import "fmt"

type Student struct {
	Id       int
	Name     string
	Grade    float64
	IsActive bool
}

func (s Student) GetInfo() string {
	return fmt.Sprintf("ID: %d, Nama: %s, IPK: %.2f, Aktif: %t", s.Id, s.Name, s.Grade, s.IsActive)
}

func (s *Student) UpdateGrade(newGrade float64) {
	s.Grade = newGrade
}

func (s *Student) Activate() {
	s.IsActive = true
}
func (s *Student) Deactivate() {
	s.IsActive = false
}

func main() {
	mhs1 := Student{
		Id:       1,
		Name:     "Ahmad",
		Grade:    3.5,
		IsActive: true,
	}

	fmt.Println("Before Change :", mhs1.GetInfo())
	mhs1.UpdateGrade(3.9)
	fmt.Println("After Update Grade :", mhs1.GetInfo())
	mhs1.Deactivate()
	fmt.Println("After Deactivate :", mhs1.GetInfo())
	mhs1.Activate()
	fmt.Println("After Activate :", mhs1.GetInfo())
}
