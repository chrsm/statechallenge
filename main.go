package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
)

type state_t struct {
	// what you'd save
}

func main() {
	s := load()

	// it does NOT matter where/how `your_print` is defined, only that you have the ability to track the state.
	// a closure is certainly *easier*, but you could make the state global. doesn't matter.
	your_print := func(txt string) {
		fmt.Printf(txt)

		// implicit save
		wait_for_input(s)
	}

	// additional, to make it clear that it's not just about printing text - let's create some (tmp) files
	// the point is that you will have N different types of things, so make state_t usable for all of them, not just
	// text printing
	create_file := func(name, content string) {
		fp, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND, 0644) // deliberate `O_APPEND`, don't write the same thing twice, that's bad.jpeg
		if err != nil {
			panic("error creating file " + name + ": " + err.Error())
		}

		fp.WriteString(content)
		fp.Close()

		fmt.Printf("(created file %s)", name)
		// implicit save
		wait_for_input(s)
	}

	// the reentrant fn
	do_thing := func( /* could be state_t here, your_print, whatever - just needs the environment somehow */ ) {
		your_print("hello [1]")
		your_print("hello [2]")
		create_file("test.txt", "some content")
		your_print("hello [3]")
		your_print("hello [3]") // deliberately duplicated
	}

	// now run
	do_thing()

	// just an example, you can save wherever you want
	save(s)
}

func wait_for_input(s *state_t) {
	c, _ := bufio.NewReader(os.Stdin).ReadByte()
	if c == 'Q' {
		save(s)

		os.Exit(0)
	}
}

/* forgive me, i'm not going to make this super complex, we'll just use gob and a hardcoded filename */
func load() *state_t {
	s := new(state_t)

	fp, err := os.Open("state.bin")
	if err != nil {
		if os.IsNotExist(err) {
			// nothing to return
			return s
		}

		panic("error opening state.bin: " + err.Error())
	}

	if err := gob.NewDecoder(fp).Decode(s); err != nil {
		panic("error decoding state: " + err.Error())
	}

	return s
}

func save(s *state_t) {
	fp, err := os.OpenFile("state.bin", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic("error writing state.bin: " + err.Error())
	}
	defer fp.Close()

	if err := gob.NewEncoder(fp).Encode(s); err != nil {
		panic("err encoding state to gob: " + err.Error())
	}
}
