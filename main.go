package main

import (
	"fmt"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	taskCreturer := func(a chan Ttype) {
		go func() {
			var i int = 1
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = fmt.Sprintf("Number of error: %v; Some error occured", i)
					i++
				}

				a <- Ttype{cT: ft, id: int(time.Now().Nanosecond())} // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	task_worker := func(a Ttype) Ttype {
		tt, err := time.Parse(time.RFC3339, a.cT)
		if err != nil {
			fmt.Printf("There is error: %v\n", err)
		}
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			t = task_worker(t)
			// go tasksorter(t)
			tasksorter(t)
		}
		close(superChan)
	}()

	result := map[int]Ttype{}
	err := []error{}
	go func() {
		for r := range doneTasks {
			result[r.id] = r
			// go func() {
			// 	result[r.id] = r
			// }()
		}
		// close(doneTasks)
	}()
	go func() {
		for r := range undoneTasks {
			err = append(err, r)
			// go func() {
			// 	err = append(err, r)
			// }()
		}
		// close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	fmt.Println("Errors:")
	for _, st := range err {
		fmt.Printf("Error: %#v\n", st)
	}

	fmt.Println("Done tasks:")
	for _, st := range result {
		fmt.Printf("Task: id: %v; Result: %v\n", st.id, string(st.taskRESULT))
	}
}