package cmd_test

import (
	"log/slog"
	"os"
	"testing"

	. "github.com/golang-cz/devslog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	os.Setenv("EXAMPLE_API_DATABASE_HOST", "example.com")
	os.Setenv("EXAMPLE_API_AWS_REGION", "us-east-1")

	slog.SetDefault(
		slog.New(
			NewHandler(GinkgoWriter,
				&Options{
					HandlerOptions: &slog.HandlerOptions{
						AddSource: true,
						Level:     slog.LevelDebug,
					},
					NewLineAfterLog:   true,
					StringIndentation: true,
				},
			),
		),
	)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Schema Suite")
}
