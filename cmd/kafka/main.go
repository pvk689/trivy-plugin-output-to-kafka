package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"trivy-plugin-output-kafka/internal/kafka"
	"trivy-plugin-output-kafka/internal/report"
)

const writeTimeout = 30 * time.Second

var version = "dev"

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("kafka", flag.ContinueOnError)
	topic := fs.String("topic", "", "Kafka topic (required)")
	brokers := fs.String("brokers", "", "Kafka brokers, comma-separated (required)")
	key := fs.String("key", "", "Optional message key")
	showVersion := fs.Bool("version", false, "Print version and exit")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *showVersion {
		fmt.Fprintf(stdout, "kafka %s\n", version)
		return nil
	}
	if *topic == "" {
		return errors.New("topic is required")
	}
	if *brokers == "" {
		return errors.New("brokers is required")
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	w := kafka.NewWriter(strings.Split(*brokers, ","), *topic)
	defer func() {
		if err := w.Close(); err != nil {
			log.Printf("failed to close Kafka writer: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()

	if err := kafka.Write(ctx, w, *key, data); err != nil {
		return fmt.Errorf("failed to write to Kafka: %w", err)
	}

	log.Println("log has been sent to Kafka")

	if err := report.Summarize(data, stdout); err != nil {
		log.Printf("failed to summarize report: %v", err)
	}

	return nil
}
