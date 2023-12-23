package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	wg              sync.WaitGroup
	mutex                 = sync.Mutex{}
	widgetInventory int32 = 1000 //Package-level variable will avoid all the pointers
	newPurchase           = sync.NewCond(&mutex)
)

func main() {
	fmt.Println("Starting inventory count = ", widgetInventory)
	wg.Add(2)
	go makeSales()
	go newPurchases()
	wg.Wait()
	fmt.Println("Ending inventory count = ", widgetInventory)
}

func makeSales() { // 1000000 widgets sold
	for i := 0; i < 30000; i++ {
		mutex.Lock()
		if widgetInventory-100 < 0 {
			newPurchase.Wait()
		}
		widgetInventory -= 100
		if widgetInventory < 0 {
			fmt.Println(widgetInventory)
		}
		mutex.Unlock()
	}

	wg.Done()
}

func newPurchases() { // 1000000 widgets purchased
	time.Sleep(time.Microsecond * 50) // delay to add more inventory to allow makeSales consumes inventory in advance

	for i := 0; i < 30000; i++ {
		mutex.Lock()
		widgetInventory += 100
		if widgetInventory < 0 {
			fmt.Println(widgetInventory)
		}
		newPurchase.Signal()
		mutex.Unlock()
	}
	wg.Done()
}
