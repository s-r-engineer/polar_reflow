package linker

import (
	"fmt"
	"polar_reflow/tools"
	"sync"
	"time"
)

var chain *Link
var mutex sync.Mutex

type Link struct {
	value interface{}
	next  *Link
	prev  *Link
}

func LockMe() func() {
	mutex.Lock()
	return func() {
		mutex.Unlock()
	}
}

func isEmpty() bool {
	return chain == nil
}

func count() (counter int) {
	defer LockMe()()
	for l := chain.next; l != chain; l = l.next {
		counter += 1
	}
	return counter + 1
}

func Push(value interface{}) {
	defer LockMe()()
	link := Link{value: value}
	if isEmpty() {
		link.next = &link
		link.prev = &link
	} else {
		link.next = chain
		link.prev = chain.prev
		chain.prev = &link
		link.prev.next = &link
	}
	chain = &link
}

// will return value, delete the link from the chain and return the function to add link back if some shit will happen
func Pop() (value interface{}, f func()) {
	defer LockMe()()
	value = chain.value
	chain.next.prev = chain.prev
	chain.prev.next = chain.next
	chain = chain.next
	return value, func() {
		Push(value)
	}
}

func CreateLinker(excludeSddn, excludeRmssd bool, startTime, endTime time.Time, periods map[string][]int) {
	for method, timePeriods := range periods {
		if (excludeSddn && method == "sddn") || (excludeRmssd && method == "rmssd") {
			continue
		}
		for _, timePeriod := range timePeriods {
			clearHours := timePeriod / 60
			minutesLeft := timePeriod % 60
			clearDays := clearHours / 24
			hoursLeft := clearHours % 24
			offset := time.Duration(timePeriod) * time.Minute
			timeTagLine := fmt.Sprintf("%d%s%d%s%d%s", clearDays, "d", hoursLeft, "h", minutesLeft, "m")
			for timeCounter := startTime; timeCounter.Before(endTime); timeCounter = timeCounter.Add(offset) {
				Push([]string{
					method, timeTagLine, tools.FormatTime(timeCounter), tools.FormatTime(timeCounter.Add(offset)),
				})
			}
		}
	}
}
