package main

import (
    "bufio"
    "fmt"
    "io"
    "mime/quotedprintable"
    "net/mail"
    "os"
    "strings"
    "time"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <inputfile>")
        os.Exit(1)
    }

    inputFile := os.Args[1]

    file, err := os.Open(inputFile)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    msg, err := mail.ReadMessage(bufio.NewReader(file))
    if err != nil {
        panic(err)
    }

    outputFile := fmt.Sprintf("mail_%s.txt", time.Now().Format("20060102150405"))

    outFile, err := os.Create(outputFile)
    if err != nil {
        panic(err)
    }
    defer outFile.Close()

    fmt.Fprintf(outFile, "FROM: %s\n", msg.Header.Get("From"))
    fmt.Fprintf(outFile, "To: %s\n", msg.Header.Get("To"))
    fmt.Fprintf(outFile, "SUBJECT: %s\n\n", msg.Header.Get("Subject"))

    reader := quotedprintable.NewReader(msg.Body)
    body, err := io.ReadAll(reader)
    if err != nil {
        panic(err)
    }

    fmt.Fprintf(outFile, "本文:\n%s", strings.Trim(string(body), "\n"))

    fmt.Printf("Output written to %s\n", outputFile)
}
