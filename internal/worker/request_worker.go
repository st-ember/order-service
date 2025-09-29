package worker

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/st-ember/ecommerceprocessor/internal/cache"
	"github.com/st-ember/ecommerceprocessor/internal/enum"
	"github.com/st-ember/ecommerceprocessor/internal/model"
	"github.com/st-ember/ecommerceprocessor/internal/redis"
	"github.com/st-ember/ecommerceprocessor/internal/storage"
)

var reqKey string = "purchase_request"

type Worker struct {
	reqCh  chan model.Purchase
	ticker *time.Ticker
	wg     *sync.WaitGroup
	quit   chan bool
}

func NewWorker(reqCh chan model.Purchase, ticker *time.Ticker, wg *sync.WaitGroup) Worker {
	return Worker{
		reqCh:  reqCh,
		ticker: ticker,
		wg:     wg,
		quit:   make(chan bool),
	}
}

func (w *Worker) requestWorker() {
	defer w.wg.Done()

	for {
		var didTick bool
		select {
		case <-w.ticker.C:
			didTick = true
		default:
			didTick = false
		}

		batchLen := 1000

		// if channel is full, or one second has passed
		// store requests
		if len(w.reqCh) >= batchLen-1 || didTick {
			reqBatch := make([]model.Purchase, batchLen)
			for req := range w.reqCh {
				reqBatch = append(reqBatch, req)
			}
			storage.BatchStorePurchase(reqBatch)
		} else { // pop from redis and store in channel
			redisItem, _ := redis.BLPop(reqKey)
			redisStr, _ := redisItem.(string)

			var cacheItem cache.Purchase
			_ = json.Unmarshal([]byte(redisStr), &cacheItem)

			purchaseItem := &model.Purchase{
				Id:          uuid.New(),
				Product:     cacheItem.Body.Product,
				Customer:    cacheItem.Body.Customer,
				PurchasedAt: time.Unix(0, cacheItem.Timestamp),
				Status:      enum.Pending,
			}

			select {
			case w.reqCh <- *purchaseItem:
			default: // if channel is full, push item back to redis
				_ = redis.RPush(reqKey, *purchaseItem)
			}
		}

		// break loop (terminate worker) if quit channel is true
		select {
		case <-w.quit:
			return
		default:
		}
	}
}

func StartRequestWorker() {
	ticker := time.NewTicker(time.Second)
	reqChan := make(chan model.Purchase, 1000)
	var wg sync.WaitGroup
	var workerCount int

	quitChs := make(map[int]chan bool)

	go func() {
		for range ticker.C {
			// Todo: add error channel to collect errors from LLen
			queLen, _ := redis.LLen("purchase_que")
			newWorkerCount := determineWorkerCount(int(queLen))

			// Todo: manage quit channels index independent from workerCount
			diff := int(newWorkerCount) - workerCount
			if diff > 0 {
				for i := 0; i < diff; i++ {
					worker := NewWorker(reqChan, ticker, &wg)
					quitChs[workerCount+1] = worker.quit // add the quit channel initialized by the worker to quitChs map
					go worker.requestWorker()            // start the worker
				}
			} else if diff < 0 {
				for i := 0; i < -diff; i++ {
					workerQuit := quitChs[workerCount-i-1]
					workerQuit <- true               // send signal to quit worker
					delete(quitChs, workerCount-i-1) // remove quit channel
				}
			}
			workerCount = newWorkerCount // reassign worker count
		}
	}()

	wg.Wait()
}

func determineWorkerCount(queLen int) int {
	switch {
	case queLen > 100000:
		return 1000
	case queLen > 10000:
		return 100
	case queLen > 1000:
		return 10
	default:
		return 5
	}
}
