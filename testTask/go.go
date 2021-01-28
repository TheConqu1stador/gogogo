package main

import (
	crypto "crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"sort"
	"time"
)

//Parameters - структура, в которой хранятся значения флагов, подаваемых на вход
//arrSize (-arr-size) - размер генерируемых массивов
//writers (-writers) - количество работников, генерирующих задачи
//iterations (-iter-count) - количество итераций для
type Parameters struct {
	arrSize    int
	writers    int
	iterations int
}

//TaskInfo - вся необходимая информация для выполенения задачи основным потоком
//id - id горутины, поставившей задачу, ar - исходный сгенерированный массив
//iteration - итерация, на которой горутина сформулировала задачу, timestamp - время добавления задачи
type TaskInfo struct {
	id        int
	iteration int
	ar        []int
	timestamp time.Time
}

//Rand - криптостойкая генерация рандомного числа в пределах 100000
func Rand() int {
	safeNum, err := crypto.Int(crypto.Reader, big.NewInt(100000))
	if err != nil {
		panic(err)
	}

	return int(safeNum.Int64())
}

//GenerateArray - генерация рандомного массива размером -arr-size
func GenerateArray(ar *[]int, size int) {
	for i := 0; i < size; i++ {
		(*ar) = append((*ar), Rand())
	}
}

//Writer - функция, генерирующая задачи для MainWorker
func Writer(id int, tasks chan<- TaskInfo, params Parameters) {
	for j := 0; j < params.iterations; j++ {
		var ar []int
		GenerateArray(&ar, params.arrSize)

		fmt.Println("Worker", id, "pushed job", j, ar)
		tasks <- TaskInfo{ar: ar, id: id, iteration: j, timestamp: time.Now()}
	}
}

//MainWorker - основная функция, обрабатывающая очередь задач
func MainWorker(tasks chan TaskInfo, res chan bool, params Parameters) {
	for j := 0; j < params.iterations*params.writers; j++ {
		tmpStruct := <-tasks
		sort.Ints(tmpStruct.ar)

		fmt.Printf("{goroutine id: %d} {iteration: %d} {queueInsertionTime: %s} {min: %d} {median: %d} {max: %d}\n",
			tmpStruct.id, tmpStruct.iteration, tmpStruct.timestamp, tmpStruct.ar[0], tmpStruct.ar[params.arrSize/2], tmpStruct.ar[params.arrSize-1])
		res <- true
	}

	close(tasks)
}

func main() {
	tasks := make(chan TaskInfo)
	res := make(chan bool)

	arrSize := flag.Int("arr-size", 0, "sizeOfArray")
	writers := flag.Int("writers", 0, "amount of writers")
	iterations := flag.Int("iter-count", 0, "amount of iterations")

	flag.Parse()
	params := Parameters{*arrSize, *writers, *iterations}

	go MainWorker(tasks, res, params)

	for w := 0; w < params.writers; w++ {
		go Writer(w, tasks, params)
	}

	for j := 0; j < params.iterations*params.writers; j++ {
		<-res
	}
}
