package ent_test

import (
	"log/slog"
	"testing"

	. "github.com/golang-cz/devslog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEntity(t *testing.T) {
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
	RunSpecs(t, "Database Entity Suite")
}
