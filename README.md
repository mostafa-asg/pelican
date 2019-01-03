# Pelican InMemory Key/Value Store
![Pelican](https://github.com/mostafa-asg/pelican/blob/master/images/pelican-design.png)

## Embeding in your Go applications
### Installation
`
go get -u https://github.com/mostafa-asg/pelican
`
### Usage
```go
package main

import (
	"time"

	"github.com/mostafa-asg/pelican/store"
)

func main() {

	kvStore := store.New(10*time.Minute, // default expiration time
		store.Sliding,  // default expire strategy
		30*time.Minute) // clean up memory every 30 minutes

	// FirstName will be expired after 10 minutes if
	// within this time no Get issued for this key
	// otherwise expiration will be reset
	// why? because expire strategy has been set to Sliding
	kvStore.Put("FirstName", "Mostafa")

	value, found := kvStore.Get("FirstName")
	if found {
		// Because value's type is interface{}
		// type assertion is needed
		println("First name is " + value.(string))
	}

	// If you know the type of the key
	// You can use the helper methods provided
	value2, found := kvStore.GetString("FirstName")
	if found {
		println("First name is " + value2)
	}

	// You can also change the settings per key
	kvStore.PutWithExpire("LastName", // key
		"Asgari",       // value
		2*time.Second,  // expire after 2 seconds
		store.Absolute) // expire strategy set to Absolute

	lastname, found := kvStore.GetString("LastName")
	if found {
		println("Last name is " + lastname)
	}

	// Wait until lastname expire
	time.Sleep(3 * time.Second)

	lastname, found = kvStore.GetString("LastName")
	if !found {
		// This line shuld be executed
		println("Last name not found")
	}

}
```
## As a standalone process
There is two ways to connect to Pelican as a standalone service:
* Via Sockets
* Via Http endpoints  

The simplest way to connect to Pelican is through available clinets:
* Java applications: [jPelican](https://github.com/mostafa-asg/jPelican)
