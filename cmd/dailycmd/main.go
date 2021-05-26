package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kyoukaya/genshindaily/internal/genshindaily"
)

const usage = "usage: ./daily < conf.json"

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	check(err)
	conf := &genshindaily.Config{}
	check(json.Unmarshal(b, conf))
	s, err := genshindaily.HandleMessage(context.Background(), conf)
	check(err)
	fmt.Println(s)
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, strings.TrimSpace(err.Error()))
		fmt.Println(usage)
		os.Exit(1)
	}
}
