package main

import (
	"bufio"
	"fmt"
	"mooss/harp/lex"
	"os"
)

func main() {
	fmt.Println("Harp REPL - v0.0.0")
	fmt.Println("Enter code (Ctrl+C to exit)")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		lexer := lex.NewLexer(input)

		for {
			tok, err := lexer.NextToken()
			if err != nil {
				fmt.Println(err)
				break
			}

			if tok.Type == lex.TOKEN_EOF || tok.Type == lex.TOKEN_INVALID {
				break
			}
			fmt.Printf("%+v\n", tok)
		}
	}
}
