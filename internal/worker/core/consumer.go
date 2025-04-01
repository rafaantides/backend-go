package core

import (
	"backend-go/internal/worker/services/queue"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type ProcessMessageFunc func([]byte) error

type Consumer struct {
	queueService   *queue.RabbitMQService
	prefetchCount  int
	stopChan       chan struct{} // Canal para sinalizar parada segura
	mu             sync.Mutex
	processMessage ProcessMessageFunc
}

func NewConsumer(queueService *queue.RabbitMQService, processMessage ProcessMessageFunc, prefetchCount int) *Consumer {
	return &Consumer{
		queueService:   queueService,
		prefetchCount:  prefetchCount,
		stopChan:       make(chan struct{}), // Inicializa o canal de parada
		processMessage: processMessage,
	}
}
func (c *Consumer) Start() {
	log.Println("Iniciando Worker")

	msgs, err := c.queueService.ConsumeMessages()
	if err != nil {
		log.Printf("Erro ao iniciar consumo da fila %s: %v", c.queueService.GetQueueName(), err)
		return
	}

	log.Printf("Worker iniciado na fila '%s' processando até %d mensagens simultaneamente...", c.queueService.GetQueueName(), c.prefetchCount)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, c.prefetchCount) // Controla workers simultâneos

	for {
		select {
		case <-c.stopChan:
			log.Println("Sinal de parada recebido. Encerrando worker...")
			goto cleanup

		case msg, ok := <-msgs:
			if !ok {
				log.Println("Canal de mensagens fechado. Encerrando worker...")
				goto cleanup
			}

			wg.Add(1)
			semaphore <- struct{}{} // Bloqueia se já houver prefetchCount workers em execução

			go func(msgBody []byte, msg amqp.Delivery) {
				defer wg.Done()
				defer func() { <-semaphore }() // Libera um slot ao final

				if err := c.processMessage(msgBody); err != nil {
					log.Printf("Erro ao processar mensagem: %v", err)
					// TODO: integrar com o discord
					c.Stop() // Para o worker
				}

				c.queueService.AckMessage(msg)
			}(msg.Body, msg)
		}
	}

cleanup:
	wg.Wait() // Aguarda todas as goroutines antes de encerrar
	close(semaphore)
	log.Println("Worker finalizado com sucesso.")
}

// Método para parar o worker com segurança
func (c *Consumer) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.stopChan:
		// Se já foi fechado, não faz nada
	default:
		close(c.stopChan)
	}
}
