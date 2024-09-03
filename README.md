# Go Vanilla Denial Of Service Protector

It's highly configurable dos protector for net/http package of go.

## Sample usage

### Settings

Here the definition of `Init()` function, which initializes protector:

```go

func Init(
    at uint16, // attack timespan: timespan of attacks to ban attacker.  
    ac uint16, // attack count: amount of attacks to ban attacker.
    bt uint, // ban time: how long a banned attacker remain banned, as seconds.
    ec int, // error code: which error code should be returned when ban occur.
    wl []string // whitelist: the list of exceptional ip's that never be got banned by protector.
    ) VanillaDdosProtector {}

```

### Initialization of protector

```go
    // other codes

    import (
        // other packages

        "github.com/Necoo33/gvdp"
    )

    fn main(){
        // other codes

        newProtector := gvdp.Init(30, 5, 60, 429, []string{})

        // other codes
    }

```
By that, the protector will be initialized.

### Activating the protector

Now you have to handle that banning process in whichever route that you want. For example:

```go
    // other codes

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		// start and save current state of the protector by that: 
        newProtector = newProtector.HandleBanningAndAllowing(req)

        // do whatever you want to banned users:
		if newProtector.BanOccured {
			res.Write([]byte("you're banned!"))
		} else {
			res.Write([]byte("Hello from a go server!"))
		}
	})

    // other codes
```

Here is a minimal but working sample for that protector:

```go

package main

import (
	"fmt"
	"github.com/Necoo33/gvdp"

	"net/http"
)

func main() {
	newProtector := protector.Init(30, 5, 60, 429, []string{})


	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		newProtector = newProtector.HandleBanningAndAllowing(req)

		if newProtector.BanOccured {
			res.Write([]byte("you're banned!"))
		} else {
			res.Write([]byte("Hello from a go server!"))
		}
	})

	http.ListenAndServe(":2000", nil)
}

```