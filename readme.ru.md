# Batch Query

Эта библиотека помогает объеденить множество мелких однотипных запросов к различным ресурсам (базам данных, кэшам, сетям)
в батчи и тем самым снизить накладные расходы на соединенях и транспорте данных.

Библиотека изначально повилась в ответ на постоянные проблемы между высоконагруженным сервисом и сервером Aerospike.
Метрики самого Aerospike упорно показывали, что он недозагружен, в то время как метрики приложения указывали на огромные
тайминги запросов к Aerospike. Также создание библиотеки было вдохновлено пакетом
[singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight) - сама идея уменьшить количество мелких запросов,
в виде логичного продолжения, пришла к мысли избавиться от них вообще.

Сама идея решения очень простая: собирать мелкие запросы до тех пор, пока не будет достигнут предел по количеству
(максимальный размер батча) или пока не истечёт время сбора батча (с момента поступления первого мелкого запроса в батч).

## Использование

Центральным компонентом библиотеки является структура `BatchQuery`, которая и выполняет всю работу. Перед использованием 
её необходимо настроить и для этих целей служит структура [`Config`](config.go). Давайте рассмотрим её поля:
* `BatchSize` - сколько мелких запросов может иметь батч. Поле необязательно, значение по умолчанию `64`.
* `CollectInterval` - максимальная продолжительность сбора батча. Начинает отсчитываться с момента поступления первого запроса в батч. Значение по умолчанию `1` секунда..
* `TimeoutInterval` - ограничение на сбор, отправку батч-запроса и пост-обработку. Должно быть больше `CollectInterval`.
* `Batcher` - абстракция для конкретного хранилища, см. описание ниже. Обязательный параметр.
* `Buffer` - хранилище для собранных батчей, готовых к отправке и обработке.
* `Workers` - количество воркеров для отправки/обработки батчей. Читают из буфера (см. `Buffer`).
* `MetricsWriter` - абстракция для конкретного TSDB решения.
* `Logger` - абстракция для логгера внутренних процессов. Полезно для отладки, не рекомендуется для продакшена.

Таким образом, пример использования выглядит следующим образом:
```go
package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/koykov/batch_query"
	promw "github.com/koykov/batch_query/metrics/prometheus"
	"github.com/koykov/batch_query/mods/aerospike"
)

func main() {
	policy := as.NewBatchPolicy()
	client, _ := as.NewClientWithPolicy(as.NewClientPolicy(), "localhost", 3000)

	// Prepare config for query.
	conf := batch_query.Config{
		BatchSize:       100,
		CollectInterval: 500 * time.Microsecond,
		TimeoutInterval: 5 * time.Millisecond,
		Workers:         10,
		Buffer:          4,
		// Declare Aerospike batcher with specific params.
		Batcher: aerospike.Batcher{
			Namespace: "my_ns",
			SetName:   "my_set",
			Bins:      []string{"bin01", "bin02", "binNN"},
			Policy:    policy,
			Client:    client,
		},
		// Declare writer to export metrics.
		MetricsWriter: promw.NewWriter("my_query", promw.WithPrecision(time.Millisecond)),
		// Declare logger for debugging purposes.
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	// Initialize the query.
	bq, _ := batch_query.New(&conf)

	// Start 10k goroutines to fetch small keys.
	for i := 0; i < 10000; i++ {
		go func() {
			for {
				resp, _ := bq.Fetch(i + rand.Intn(i))
				_ = resp.(*as.Record)
			}
		}()
	}

	c := make(chan struct{})
	<-c
}
```
В этом примере 10k горутин читает одиночные ключи, которые силами query будут объединены в батчи и обработаны пакетно.
При этом каждая горутина получит ответ именно на свой ключ, который она запрашивала или ошибку.
