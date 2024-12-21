package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	fmt.Print("> Aguardando Postgres aceitar conexões")
	for {
		cmd := exec.Command("docker", "exec", "postgres-dev", "pg_isready", "--host", "localhost")
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		output := out.String()

		if err == nil && strings.Contains(output, "accepting connections") {
			fmt.Println("\n> Postgres está pronto e aceitando conexões!")
			os.Exit(0) // Sucesso
		}

		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
}
