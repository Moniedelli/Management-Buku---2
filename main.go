package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/jung-kurt/gofpdf"
)

type Book struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	PageCount int    `json:"page_count"`
	Year      int    `json:"year"`
}

func main() {
	var library []Book
	var mutex sync.Mutex

	for {
		fmt.Println("\nMenu Perpustakaan:")
		fmt.Println("1. Tambah Buku")
		fmt.Println("2. Tampilkan Daftar Buku")
		fmt.Println("3. Hapus Buku")
		fmt.Println("4. Edit Buku")
		fmt.Println("5. Print Buku")
		fmt.Println("6. Keluar")
		fmt.Print("Pilih menu (1-6): ")

		var choice int
		fmt.Scanf("%d\n", &choice)

		switch choice {
		case 1:
			addBook(&library, &mutex)
		case 2:
			printBooks(&library)
		case 3:
			deleteBook(&library, &mutex)
		case 4:
			editBook(&library, &mutex)
		case 5:
			fmt.Println("Pilihan:")
			fmt.Println("1. Cetak semua buku")
			fmt.Println("2. Cetak berdasarkan ID")
			fmt.Print("Pilih cara pencetakan (1/2): ")
			var printChoice int
			fmt.Scanf("%d\n", &printChoice)
			switch printChoice {
			case 1:
				printAllBooks(&library)
			case 2:
				var id string
				fmt.Print("Masukkan ID Buku: ")
				fmt.Scanf("%s\n", &id)
				printBookByID(&library, id)
			default:
				fmt.Println("Input tidak valid")
			}
		case 6:
			os.Exit(0)
		default:
			fmt.Println("Input tidak valid")
		}
	}
}

func addBook(library *[]Book, mutex *sync.Mutex) {
	var book Book
	var id string
	for {
		fmt.Print("Masukkan ID Buku: ")
		fmt.Scanf("%s\n", &id)
		if isBookIDUnique(library, id) {
			break
		} else {
			fmt.Println("ID Buku sudah digunakan")
		}
	}
	book.ID = id
	fmt.Print("Masukkan Judul Buku: ")
	fmt.Scanf("%s\n", &book.Title)
	fmt.Print("Masukkan Nama Penulis: ")
	fmt.Scanf("%s\n", &book.Author)
	fmt.Print("Masukkan Nama Penerbit: ")
	fmt.Scanf("%s\n", &book.Publisher)
	var pageCountStr string
	fmt.Print("Masukkan Jumlah Halaman: ")
	fmt.Scanf("%s\n", &pageCountStr)
	pageCount, err := strconv.Atoi(pageCountStr)
	if err == nil {
		book.PageCount = pageCount
	} else {
		fmt.Println("Input tidak valid")
		return
	}
	var yearStr string
	fmt.Print("Masukkan Tahun Terbit: ")
	fmt.Scanf("%s\n", &yearStr)
	year, err := strconv.Atoi(yearStr)
	if err == nil {
		book.Year = year
	} else {
		fmt.Println("Input tidak valid")
		return
	}
	mutex.Lock()
	*library = append(*library, book)
	mutex.Unlock()
	saveBooks(library)
	fmt.Println("Buku berhasil ditambahkan")
}

func isBookIDUnique(library *[]Book, id string) bool {
	for _, book := range *library {
		if book.ID == id {
			return false
		}
	}
	return true
}

func printAllBooks(library *[]Book) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Daftar Buku")
	pdf.Ln(10)

	for _, book := range *library {
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(0, 10, fmt.Sprintf("Judul: %s", book.Title))
		pdf.Ln(8)
		pdf.Cell(0, 10, fmt.Sprintf("Penulis: %s", book.Author))
		pdf.Ln(8)
		pdf.Cell(0, 10, fmt.Sprintf("Penerbit: %s", book.Publisher))
		pdf.Ln(8)
		pdf.Cell(0, 10, fmt.Sprintf("Jumlah Halaman: %d", book.PageCount))
		pdf.Ln(8)
		pdf.Cell(0, 10, fmt.Sprintf("Tahun Terbit: %d", book.Year))
		pdf.Ln(12)
	}

	err := pdf.OutputFileAndClose("pdf/all_books.pdf")
	if err != nil {
		fmt.Println("Terjadi kesalahan saat menyimpan file PDF")
		return
	}

	fmt.Println("Data seluruh buku berhasil di-print ke PDF (all_books.pdf)")
}

func printBookByID(library *[]Book, id string) {
	var bookToPrint *Book
	for i, book := range *library {
		if book.ID == id {
			bookToPrint = &((*library)[i])
			break
		}
	}
	if bookToPrint == nil {
		fmt.Println("Buku tidak ditemukan")
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Detail Buku")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, "Judul: "+bookToPrint.Title)
	pdf.Ln(8)
	pdf.Cell(0, 10, "Penulis: "+bookToPrint.Author)
	pdf.Ln(8)
	pdf.Cell(0, 10, "Penerbit: "+bookToPrint.Publisher)
	pdf.Ln(8)
	pdf.Cell(0, 10, "Jumlah Halaman: "+strconv.Itoa(bookToPrint.PageCount))
	pdf.Ln(8)
	pdf.Cell(0, 10, "Tahun Terbit: "+strconv.Itoa(bookToPrint.Year))

	err := pdf.OutputFileAndClose(fmt.Sprintf("pdf/book_%s.pdf", id))
	if err != nil {
		fmt.Println("Terjadi kesalahan saat menyimpan file PDF")
		return
	}

	fmt.Printf("Buku dengan ID %s berhasil di-print ke PDF (book_%s.pdf)\n", id, id)
}

func printBooks(library *[]Book) {
	files, err := ioutil.ReadDir("books")
	if err != nil {
		fmt.Println("Terjadi kesalahan saat membaca direktori books")
		return
	}
	var books []Book
	for _, file := range files {
		if !file.IsDir() {
			data, err := ioutil.ReadFile(filepath.Join("books", file.Name()))
			if err != nil {
				fmt.Println("Terjadi kesalahan saat membaca file buku")
				return
			}
			var book Book
			err = json.Unmarshal(data, &book)
			if err != nil {
				fmt.Println("Terjadi kesalahan saat memproses file buku")
				return
			}
			books = append(books, book)
		}
	}
	fmt.Println("\nDaftar Buku:")
	for i, book := range books {
		fmt.Printf("%d. %s - %s (%d halaman, %d)\n", i+1, book.Title, book.Author, book.PageCount, book.Year)
	}
}

func deleteBook(library *[]Book, mutex *sync.Mutex) {
	var id string
	fmt.Print("Masukkan ID Buku: ")
	fmt.Scanf("%s\n", &id)

	mutex.Lock()
	defer mutex.Unlock()

	var newLibrary []Book
	for _, book := range *library {
		if book.ID != id {
			newLibrary = append(newLibrary, book)
		} else {
			err := os.Remove(filepath.Join("books", book.ID+".json"))
			if err != nil {
				fmt.Println("Error deleting book JSON:", err)
			}
		}
	}
	if len(newLibrary) == len(*library) {
		fmt.Println("Buku tidak ditemukan")
		return
	}
	*library = newLibrary
	saveBooks(library)
	fmt.Println("Buku berhasil dihapus")
}

func editBook(library *[]Book, mutex *sync.Mutex) {
	var id string
	fmt.Print("Masukkan ID Buku: ")
	fmt.Scanf("%s\n", &id)
	var bookToEdit *Book
	for i, book := range *library {
		if book.ID == id {
			bookToEdit = &((*library)[i])
			break
		}
	}
	if bookToEdit == nil {
		fmt.Println("Buku tidak ditemukan")
		return
	}
	var title string
	fmt.Print("Masukkan Judul Buku (kosongkan untuk tidak mengubah): ")
	fmt.Scanf("%s\n", &title)
	if title != "" {
		bookToEdit.Title = title
	}
	var author string
	fmt.Print("Masukkan Nama Penulis (kosongkan untuk tidak mengubah): ")
	fmt.Scanf("%s\n", &author)
	if author != "" {
		bookToEdit.Author = author
	}
	var publisher string
	fmt.Print("Masukkan Nama Penerbit (kosongkan untuk tidak mengubah): ")
	fmt.Scanf("%s\n", &publisher)
	if publisher != "" {
		bookToEdit.Publisher = publisher
	}
	var pageCountStr string
	fmt.Print("Masukkan Jumlah Halaman (kosongkan untuk tidak mengubah): ")
	fmt.Scanf("%s\n", &pageCountStr)
	if pageCountStr != "" {
		pageCount, err := strconv.Atoi(pageCountStr)
		if err == nil {
			bookToEdit.PageCount = pageCount
		} else {
			fmt.Println("Input tidak valid")
			return
		}
	}
	var yearStr string
	fmt.Print("Masukkan Tahun Terbit (kosongkan untuk tidak mengubah): ")
	fmt.Scanf("%s\n", &yearStr)
	if yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			bookToEdit.Year = year
		} else {
			fmt.Println("Input tidak valid")
			return
		}
	}

	mutex.Lock()
	saveBooks(library)
	mutex.Unlock()
	fmt.Println("Buku berhasil diubah")
}

func saveBooks(library *[]Book) {
	err := os.MkdirAll("books", 0755)
	if err != nil {
		fmt.Println("Terjadi kesalahan saat membuat direktori books")
		return
	}
	for _, book := range *library {
		data, err := json.Marshal(book)
		if err != nil {
			fmt.Println("Terjadi kesalahan saat memproses data buku")
			return
		}
		err = ioutil.WriteFile(filepath.Join("books", "book-"+book.ID+".json"), data, 0644)
		if err != nil {
			fmt.Println("Terjadi kesalahan saat menyimpan data buku")
			return
		}
	}
}
