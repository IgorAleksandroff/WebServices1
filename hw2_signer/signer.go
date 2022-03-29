package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var inputData = []int{0, 1}
var result = "NOT_SET"

func main() {

	freeFlowJobs := []job{
		job(Int2Chan),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(chan2Value),
	}

	//start := time.Now()
	ExecutePipeline(freeFlowJobs...)
	//end := time.Since(start)

	//time.Sleep(time.Millisecond * 350)
	fmt.Println(result)
	//fmt.Println(end)

}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	for _, j := range jobs {
		out := make(chan interface{})
		wg.Add(1)
		go func(wg *sync.WaitGroup, jb job, i, o chan interface{}) {
			defer wg.Done()
			jb(i, o)
			close(out)
		}(wg, j, in, out)
		in = out
	}
	wg.Wait()
}

func Int2Chan(in, out chan interface{}) {
	for _, fibNum := range inputData {
		out <- fibNum
	}
}

func SingleHash(in, out chan interface{}) {
	for dataRow := range in {
		wg := &sync.WaitGroup{}
		dataInt := dataRow.(int)
		data := strconv.Itoa(dataInt)
		data1, data2 := "", ""
		wg.Add(2)
		go func() {
			defer wg.Done()
			data1 = DataSignerCrc32(data)
		}()
		go func() {
			defer wg.Done()
			dataMd5 := DataSignerMd5(data)
			data2 = DataSignerCrc32(dataMd5)
		}()
		wg.Wait()
		result = fmt.Sprintf("%s~%s", data1, data2)
		out <- result
		//fmt.Println(result)
	}
}

func MultiHash(in, out chan interface{}) {
	for dataRow := range in {
		wg := &sync.WaitGroup{}
		th := 6
		multiHash := ""
		hashMap := make(map[int]string)
		data := dataRow.(string)
		wg.Add(th)
		for i := 0; i < th; i++ {
			go func(wg *sync.WaitGroup, n int) {
				defer wg.Done()
				hashMap[n] = DataSignerCrc32(strconv.Itoa(n) + data)
			}(wg, i)
		}
		wg.Wait()
		for i := 0; i < th; i++ {
			multiHash += hashMap[i]
		}
		//fmt.Println(multiHash)
		out <- multiHash
	}
}

func MultiHash1(in, out chan interface{}) {
	for dataRow := range in {
		wg := &sync.WaitGroup{}
		mu := &sync.Mutex{}
		th := 6
		multiHash := ""
		hashMap := make(map[int]string)
		data := dataRow.(string)
		wg.Add(th)
		for i := 0; i < th; i++ {
			go func(wg *sync.WaitGroup, mu *sync.Mutex, n int) {
				defer wg.Done()
				mu.Lock()
				hashMap[n] = DataSignerCrc32(strconv.Itoa(n) + data)
				mu.Unlock()
			}(wg, mu, i)
		}
		wg.Wait()
		for i := 0; i < th; i++ {
			multiHash += hashMap[i]
		}
		fmt.Println(multiHash)
		out <- multiHash
	}
}
func CombineResults(in, out chan interface{}) {
	var combineResults []string
	for dataRow := range in {
		data := dataRow.(string)
		combineResults = append(combineResults, data)
	}
	sort.Strings(combineResults)
	combineResult := strings.Join(combineResults, "_")
	out <- combineResult
	//fmt.Println(combineResult)
}

func chan2Value(in, out chan interface{}) {
	dataRaw := <-in
	data, ok := dataRaw.(string)
	if !ok {
		fmt.Println("cant convert result data to string")
	}
	result = data
}
