package rates

import (
	"fmt"
	"log"
	"testing"
)

func TestFetch(t *testing.T) {
	c, err := Country("NL")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c.Name)
	fmt.Println(c.Periods[0].Rates["standard"])
}
