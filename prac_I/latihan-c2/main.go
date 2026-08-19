package main

import "fmt"

func main() {
	nama := "Syihaab"
	umur := 21
	ipk := 3.14
	isActive := true
	matkul := []string{
		"Agama Islam,",
		"Pemograman Backend,",
		"Machine Learning,",
		"Cyber Security,",
	}

	fmt.Println("Nama :", nama)
	fmt.Println("Umur :", umur)
	fmt.Println("IPK :", ipk)
	fmt.Println("Aktif :", isActive)
	fmt.Println("Mata Kuliah :", matkul)
	fmt.Println("-----------------------")

	dataIPKmhs := map[string]float64{
		"Budi": 3.5,
		"Ani":  3.7,
		"Siti": 3.8,
	}
	fmt.Println("Data IPK Awal", dataIPKmhs)

	dataIPKmhs["Rizki"] = 3.9
	fmt.Println("Data IPK Setelah Ditambah", dataIPKmhs)

	targetSearch := "Budi"
	if ipk, exists := dataIPKmhs[targetSearch]; exists {
		fmt.Println("Nama", targetSearch, "Ditemukan, dengan IPK", ipk)
	} else {
		fmt.Println("Nama", targetSearch, "Tidak Ditemukan")
	}

	delete(dataIPKmhs, targetSearch)
	fmt.Println("Data IPK Setelah Dihapus", dataIPKmhs)

	fmt.Println("Daftar Mahasiswa dan IPK :")

	for nama, ipk := range dataIPKmhs {
		fmt.Println("Nama :", nama, "IPK :", ipk)
	}
}
