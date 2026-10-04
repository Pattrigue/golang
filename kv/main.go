package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"kv/store"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	db := store.New()

	printPrompt()

	for scanner.Scan() {
		input := scanner.Text()
		tokens := strings.Fields(input)

		if len(tokens) > 0 {
			command := strings.ToUpper(tokens[0])

			switch command {
			case "SET":
				handleSet(db, tokens)
			case "GET":
				handleGet(db, tokens)
			case "DEL":
				handleDel(db, tokens)
			case "KEYS":
				handleKeys(db)
			case "TTL":
				handleTTL(db, tokens)
			case "EXIT":
				return
			default:
				fmt.Printf("ERR unknown command: %q\n", command)
			}
		}

		printPrompt()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ERR read error:", err)
		os.Exit(1)
	}
}

func printPrompt() {
	fmt.Print("> ")
}

func handleSet(db *store.Store, tokens []string) {
	if len(tokens) != 3 && len(tokens) != 5 {
		fmt.Println("ERR unrecognized format")
		return
	}

	key := tokens[1]
	value := tokens[2]

	if len(tokens) == 3 {
		db.Set(key, value)
	} else {
		subCommand := strings.ToUpper(tokens[3])
		if subCommand != "EX" {
			fmt.Printf("ERR unknown SET command: %q\n", subCommand)
			return
		}

		ttl, err := strconv.Atoi(tokens[4])
		if err != nil {
			fmt.Println("ERR EX argument must be a number")
			return
		}

		db.SetWithTTL(key, value, time.Duration(ttl)*time.Second)
	}

	fmt.Println("OK")
}

func handleGet(db *store.Store, tokens []string) {
	key := tokens[1]
	value, err := db.Get(key)
	if errors.Is(err, store.ErrNotFound) {
		fmt.Println("(nil)")
	} else {
		fmt.Println(value)
	}
}

func handleDel(db *store.Store, tokens []string) {
	key := tokens[1]
	deleted := db.Delete(key)
	if deleted {
		fmt.Println(1)
	} else {
		fmt.Println(0)
	}
}

func handleKeys(db *store.Store) {
	keys := db.Keys()
	if len(keys) > 0 {
		for _, key := range keys {
			fmt.Println(key)
		}
	} else {
		fmt.Println("")
	}
}

func handleTTL(db *store.Store, tokens []string) {
	key := tokens[1]
	expiry, err := db.TTL(key)

	if errors.Is(err, store.ErrNotFound) {
		fmt.Println(-2)
		return
	}

	if expiry == store.NoExpiry {
		fmt.Println(-1)
		return
	}

	fmt.Println(int(expiry.Seconds()))
}
