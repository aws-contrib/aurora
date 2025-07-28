package ent_test

import (
	"fmt"

	"github.com/aws-contrib/aurora/internal/database/ent"
	"github.com/jackc/pgx/v5/pgxpool"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gateway", func() {
	var gateway ent.Gateway

	BeforeEach(func() {
		var err error
		gateway, err = NewGateway()
		Expect(err).NotTo(HaveOccurred())
		Expect(gateway).NotTo(BeNil())
	})

	AfterEach(func() {
		gateway.Close()
	})

	Describe("Open", func() {
		When("the uri is invalid", func() {
			It("returns an error", func(ctx SpecContext) {
				gateway, err := ent.Open(ctx, "fake://localhost:5432/aurora")
				Expect(gateway).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		When("the options are invalid", func() {
			It("returns an error", func(ctx SpecContext) {
				option := ent.GatewayOptionFunc(
					func(_ *pgxpool.Config) error {
						return fmt.Errorf("oh no")
					},
				)

				gateway, err := ent.Open(ctx, "postgres://localhost:5432/aurora", option)
				Expect(gateway).To(BeNil())
				Expect(err).To(MatchError("oh no"))
			})
		})

		When("the connection cannot be established", func() {
			It("returns an error", func(ctx SpecContext) {
				option := ent.GatewayOptionFunc(
					func(config *pgxpool.Config) error {
						config.MaxConns = -1
						return nil
					},
				)

				gateway, err := ent.Open(ctx, "postgres://localhost:5432/aurora", option)
				Expect(gateway).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("RunInTx", func() {
		It("runs in transaction", func(ctx SpecContext) {
			done := make(chan any)

			pipeline := ent.NewQueryPipeline(
				func(x ent.Querier) error {
					defer close(done)
					querier, ok := x.(*ent.Queries)
					Expect(ok).To(BeTrue())
					Expect(querier).NotTo(BeNil())
					return nil
				},
			)

			Expect(gateway.RunInTx(ctx, pipeline)).To(Succeed())
			Expect(done).To(BeClosed())
		})

		When("the connection is closed", func() {
			BeforeEach(func() {
				gateway.Close()
			})

			It("returns an error", func(ctx SpecContext) {
				done := make(chan any)

				pipeline := ent.NewQueryPipeline(
					func(x ent.Querier) error {
						defer close(done)
						return nil
					},
				)

				Expect(gateway.RunInTx(ctx, pipeline)).To(MatchError("closed pool"))
				Expect(done).NotTo(BeClosed())
			})
		})

		When("the callback returns an error", func() {
			It("returns an error", func(ctx SpecContext) {
				done := make(chan any)

				pipeline := ent.NewQueryPipeline(
					func(x ent.Querier) error {
						defer close(done)
						return fmt.Errorf("oh no")
					},
				)

				Expect(gateway.RunInTx(ctx, pipeline)).To(MatchError("oh no"))
				Expect(done).To(BeClosed())
			})
		})
	})

	Describe("Ping", func() {
		It("pings the gateway", func(ctx SpecContext) {
			Expect(gateway.Ping(ctx)).To(Succeed())
		})
	})

	Describe("Close", func() {
		It("closes the gateway", func() {
			Expect(gateway.Close).NotTo(Panic())
		})
	})
})
