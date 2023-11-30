package main

import (
	"IB3/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Создание экземпляра веб-сервиса
	serv := service.New()

	// Конфигурация HTTP-сервера
	s := &http.Server{
		Addr:    ":13999",          // Адрес и порт для прослушивания (можно использовать конфигурацию)
		Handler: serv.GetHandler(), // Обработчик запросов
	}
	s.SetKeepAlivesEnabled(true)

	// Создание контекста и функции отмены для управления жизненным циклом сервера
	ctx, cancel := context.WithCancel(context.Background())

	// Запуск HTTP-сервера в горутине
	go func() {
		log.Printf("starting http server at %d", 13999)
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Обработка грациозного завершения
	gracefullyShutdown(ctx, cancel, s)
}

// gracefullyShutdown обеспечивает грациозное завершение работы сервера при получении сигналов от системы
func gracefullyShutdown(ctx context.Context, cancel context.CancelFunc, server *http.Server) {
	// Создание канала для получения сигналов
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	// Ожидание сигнала от системы
	<-ch

	// Попытка грациозного завершения работы сервера
	if err := server.Shutdown(ctx); err != nil {
		log.Print(err)
	}

	// Вызов функции отмены для завершения работы контекста
	cancel()
}
