/* Since we cannot use waitgroups, we have used Go channels for synchronization.
sync.mutex is used to achieve the shared buffer*/

/*LETS ANALYZE HOW THE BEHAVIOUR CHANGES WITH REGARD TO CHANGES IN M AND N*/
/* 8 readers and 2 writers are there,
      since ,
	  		no.of.reader > no.of. writer
      the writer waits for space in the buffer before they can write

	  hence ,a contention leads more frequent lock and unlocaks by mutex
*/

/* 8 reader and 8 writers are there,
	since ,
  		no.of.reader == no.of. writer
    there will be a balanced interaction

    however ,when we consider rate of writing is higher than reading, this will cause writer to wait
*/

/* 8 reader and 16 writer are there,
since ,
	2(no.of writer) == no.of.reader
the buffer may fill even more quickly

however , this may lead to frequent locks and unlocks by the mutex

*/

/* 2 reader and 8 writer are there,
since ,
	no.of.writer >>>> no.of.reader
the buffer might fill it quickly

but still writer wait for space in buffer lead to contention
*/

package main

import (
	"fmt"
	"sync"
)

const bufferSize = 10

func main() {
	var buffer []byte
	var mutex sync.Mutex
	readChannel := make(chan byte, bufferSize)
	writeChannel := make(chan byte, bufferSize)

	//  writer goroutines - N
	for i := 0; i < N; i++ {
		go func(id int) {
			for {
				value := byte(id + 'A')
				mutex.Lock()
				if len(buffer) < bufferSize {
					buffer = append(buffer, value)
					writeChannel <- value
				}
				mutex.Unlock()
			}
		}(i)
	}

	//  reader goroutines - M
	for i := 0; i < M; i++ {
		go func() {
			for {
				mutex.Lock()
				if len(buffer) > 0 {
					value := buffer[0]
					buffer = buffer[1:]
					readChannel <- value
				}
				mutex.Unlock()
			}
		}()
	}

	// Print values read by readers
	go func() {
		for {
			value := <-readChannel
			fmt.Printf("Read: %c\n", value)
		}
	}()

	// Print values written by writers
	go func() {
		for {
			value := <-writeChannel
			fmt.Printf("Write: %c\n", value)
		}
	}()

	// Keep the program running
	select {}
}

const M = 8
const N = 8
