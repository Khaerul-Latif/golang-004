package main

import (
	"fmt"
	"os"
	"strings"
)

type Biodata struct {
	Nama, Alamat, Pekerjaan, Alasan string
}

func getBiodata(absen int) Biodata {
	biodataList := map[int]Biodata{
		1: {"Khaerul", "Jalan Mawar", "Developer", "Suka bahasa pemrograman Go"},
		2: {"Latif", "Jalan Gajah", "Designer", "Suka bahasa pemrograman Go"},
	}
	biodata, found := biodataList[absen]

	if !found {
		fmt.Println("Id Tidak Ditemukan")
		os.Exit(1)
	}

	return biodata
}

func showBiodata(biodata Biodata) {
	fmt.Println("Biodata")
	fmt.Println(strings.Repeat("#", 30))
	fmt.Println("Nama:", biodata.Nama)
	fmt.Println("Alamat:", biodata.Alamat)
	fmt.Println("Pekerjaan:", biodata.Pekerjaan)
	fmt.Println("Alasan Memilih Kelas Golang:", biodata.Alasan)
	fmt.Println(strings.Repeat("#", 50))
}

func main() {
	args := os.Args

	absen := args[1]

	var absenID int
	_, err := fmt.Sscanf(absen, "%d", &absenID)
	if err != nil {
		fmt.Println("Harus berupa angka")
		os.Exit(1)
	}

	biodata := getBiodata(absenID)

	showBiodata(biodata)
}
