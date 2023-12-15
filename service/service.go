package service

import (
	"IB3/myDes"
	"encoding/hex"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Service - структура, представляющая веб-сервис
type Service struct{}

// New создает новый экземпляр службы
func New() *Service {
	return &Service{}
}

// GetHandler возвращает обработчик запросов для веб-сервиса
func (s *Service) GetHandler() http.Handler {
	// Создаем новый роутер с использованием gorilla/mux
	router := mux.NewRouter()

	// Определяем обработчики для различных эндпоинтов
	router.HandleFunc("/home", s.Home).Methods(http.MethodGet)
	router.HandleFunc("/about", s.About).Methods(http.MethodGet)
	router.HandleFunc("/home/shifr", s.Encode).Methods(http.MethodPost)
	router.HandleFunc("/home/unshifr", s.Decode).Methods(http.MethodPost)
	router.HandleFunc("/home/download", s.Download).Methods(http.MethodGet)

	// Возвращаем роутер в качестве обработчика запросов
	return router
}

// Home обрабатывает запрос на страницу home и отдает соответствующий HTML-файл
func (s *Service) Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/home.html")
}

// About обрабатывает запрос на страницу about и отдает соответствующий HTML-файл
func (s *Service) About(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/about.html")
}

// Decode обрабатывает запрос на дешифрацию текста или файла с использованием DES
func (s *Service) Decode(w http.ResponseWriter, r *http.Request) {

	des := myDes.NewMyDES("01234567")
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	text := des.Decode(fileBytes, "Super_Secret_key")
	processedFileName := "decode_" + handler.Filename

	// Создание пути для сохранения обработанного файла
	processedFilePath := filepath.Join(".", "processed_files", processedFileName)

	// Создание директории, если её нет
	if _, err := os.Stat("processed_files"); os.IsNotExist(err) {
		if err := os.Mkdir("processed_files", 0755); err != nil {
			log.Println(err)
			http.Error(w, "Error creating processed files directory", http.StatusInternalServerError)
			return
		}
	}

	// Запись результата в обработанный файл
	processedFile, err := os.Create(processedFilePath)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error saving processed file", http.StatusInternalServerError)
		return
	}
	//processedFile.Write([]byte(text))
	processedFile.Write([]byte(text))
	defer processedFile.Close()

	// Перенаправление на страницу скачивания
	http.Redirect(w, r, "/home/download?filename="+processedFileName, http.StatusSeeOther)
}

// Encode обрабатывает запрос на шифрацию файла с использованием DES
func (s *Service) Encode(w http.ResponseWriter, r *http.Request) {
	var processedFileName string
	var shifrText string
	des := myDes.NewMyDES("01234567")
	text := r.FormValue("text")

	// Если текст передан в запросе
	if text != "" {
		log.Println(text)
		shifrText = des.Encode(text, "Super_Secret_key")
		log.Println("Зашифрованный текст", shifrText)
		processedFileName = "encode_" + text + ".txt"

	} else {
		// Если файл передан в запросе
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error uploading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Чтение данных файла
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		text = string(fileBytes)
		shifrText = des.Encode(string(fileBytes), "key")
		processedFileName = "encode_" + handler.Filename
	}

	processedFilePath := filepath.Join(".", "processed_files", processedFileName)

	// Создание директории, если её нет
	if _, err := os.Stat("processed_files"); os.IsNotExist(err) {
		if err := os.Mkdir("processed_files", 0755); err != nil {
			log.Println(err)
			http.Error(w, "Error creating processed files directory", http.StatusInternalServerError)
			return
		}
	}

	// Запись результата в обработанный файл
	processedFile, err := os.Create(processedFilePath)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error saving processed file", http.StatusInternalServerError)
		return
	}
	hexResult := hex.EncodeToString([]byte(shifrText))
	log.Print(hexResult)
	processedFile.Write([]byte(shifrText))
	defer processedFile.Close()

	// Перенаправление на страницу скачивания
	http.Redirect(w, r, "/home/download?filename="+processedFileName, http.StatusSeeOther)
}

// Download обрабатывает запрос на скачивание файла и отдает файл с указанным именем
func (s *Service) Download(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(".", "processed_files", filename)

	// Установка заголовка для скачивания
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, filePath)
}
