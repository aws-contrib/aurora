package ent_test

import (
	"fmt"

	"github.com/aws-contrib/aurora/internal/database/ent"

	. "github.com/aws-contrib/aurora/internal/database/ent/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gateway", Ordered, func() {
	var gateway ent.Gateway

	BeforeEach(func(ctx SpecContext) {
		var err error
		gateway, err = NewGateway()
		Expect(err).NotTo(HaveOccurred())
		Expect(gateway.CreateTableRevisions(ctx)).To(Succeed())
	})

	AfterEach(func() {
		gateway.Close()
	})

	Describe("Revision", func() {
		var entity *ent.Revision

		BeforeAll(func() {
			entity = NewFakeRevision()
		})

		Describe("ApplyRevision", func() {
			var params *ent.ApplyRevisionParams

			BeforeEach(func() {
				fs := &FakeReadFileFS{}
				fs.ReadFileReturns([]byte("SELECT 1;"), nil)

				params = &ent.ApplyRevisionParams{}
				params.FileSystem = fs
				params.Revision = entity
			})

			It("applies a revision", func(ctx SpecContext) {
				Expect(gateway.ApplyRevision(ctx, params)).To(Succeed())
			})

			When("the file system fails", func() {
				BeforeEach(func() {
					fs := params.FileSystem.(*FakeReadFileFS)
					fs.ReadFileReturns(nil, fmt.Errorf("oh no"))
				})

				It("returns an error", func(ctx SpecContext) {
					Expect(gateway.ApplyRevision(ctx, params)).To(MatchError("oh no"))
				})
			})

			When("the gateway fails", func() {
				BeforeEach(func() {
					fs := params.FileSystem.(*FakeReadFileFS)
					fs.ReadFileReturns([]byte("I AM WRONG;"), nil)
				})

				It("returns an error", func(ctx SpecContext) {
					Expect(gateway.ApplyRevision(ctx, params)).NotTo(Succeed())
				})
			})
		})
	})
})
