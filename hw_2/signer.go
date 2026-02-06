package main

import (
	"fmt"
	"strconv"
	"sync"
)

// crc32(data) + "~" + crc32(md5(data))
func SingleHash(inputChannel, outputChannel chan interface{}) {
	var workerGroup sync.WaitGroup
	md5Mutex := &sync.Mutex{}

	for inputValue := range inputChannel {
		dataString := fmt.Sprintf("%v", inputValue)
		workerGroup.Add(1)

		go func(data string) {
			defer workerGroup.Done()

			var crc32FromData string
			var crc32FromMd5 string

			var interGroup sync.WaitGroup
			interGroup.Add(2)

			go func() { // crc32(data)
				defer interGroup.Done()
				crc32FromData = DataSignerCrc32(data)
			}()

			go func() { // crc32(md5(data))
				defer interGroup.Done()

				md5Mutex.Lock()
				md5Hash := DataSignerMd5(data)
				md5Mutex.Unlock()

				crc32FromMd5 = DataSignerCrc32(md5Hash)
			}()

			interGroup.Wait()
			outputChannel <- crc32FromData + "~" + crc32FromMd5
		}(dataString)
	}

	workerGroup.Wait()
}

// 6 хешей: crc32("0"+data)...crc32("5"+data)
func MultiHash(inputChannel, outputChannel chan interface{}) {
	var workerGroup sync.WaitGroup

	for inputValue := range inputChannel {
		dataString := fmt.Sprintf("%v", inputValue)
		workerGroup.Add(1)

		go func(data string) {
			defer workerGroup.Done()

			hashResults := make([]string, 6)

			var interGroup sync.WaitGroup
			interGroup.Add(6)

			for hashIndex := 0; hashIndex < 6; hashIndex++ {
				currentIndex := hashIndex

				go func() {
					defer interGroup.Done()

					hashValue := DataSignerCrc32(
						strconv.Itoa(currentIndex) + data,
					)

					hashResults[currentIndex] = hashValue
				}()
			}

			interGroup.Wait()

			combinedHash := ""
			for index := 0; index < 6; index++ {
				combinedHash += hashResults[index]
			}

			outputChannel <- combinedHash
		}(dataString)
	}

	workerGroup.Wait()
}

// CombineResults собирает, сортирует и объединяет через "_"
func CombineResults(inputChannel, outputChannel chan interface{}) {
	collectedResults := []string{}

	for inputValue := range inputChannel {
		resultString, isString := inputValue.(string)
		if isString {
			collectedResults = append(collectedResults, resultString)
		}
	}
	//sort
	for firstIndex := 0; firstIndex < len(collectedResults)-1; firstIndex++ {
		for secondIndex := firstIndex + 1; secondIndex < len(collectedResults); secondIndex++ {
			if collectedResults[firstIndex] > collectedResults[secondIndex] {
				collectedResults[firstIndex], collectedResults[secondIndex] =
					collectedResults[secondIndex], collectedResults[firstIndex]
			}
		}
	}

	finalResult := ""
	for index, resultValue := range collectedResults {
		if index > 0 {
			finalResult += "_"
		}
		finalResult += resultValue
	}

	outputChannel <- finalResult
}

func ExecutePipeline(jobs ...job) {
	var inputChannel chan interface{}
	var pipelineGroup sync.WaitGroup

	for _, jobFunc := range jobs {
		outputChannel := make(chan interface{})
		pipelineGroup.Add(1)

		go func(currentJob job, inChannel, outChannel chan interface{}) {
			defer pipelineGroup.Done()

			currentJob(inChannel, outChannel)
			close(outChannel)

		}(jobFunc, inputChannel, outputChannel)

		inputChannel = outputChannel
	}

	//job1 out -> job2 in
	//job2 out -> job3 in
	//job3 out -> ...

	pipelineGroup.Wait()
}
