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

	BeforeEach(func() {
		var err error
		gateway, err = NewGateway()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		gateway.Close()
	})

	Describe("Revision", func() {
		var entity *ent.Revision

		BeforeAll(func() {
			entity = NewFakeRevision()
		})

		Describe("CreateTableRevisions", func() {
			It("creates the aurora_schema_revisions table", func(ctx SpecContext) {
				Expect(gateway.CreateTableRevisions(ctx)).To(Succeed())
			})
		})

		Describe("InsertRevision", func() {
			var params *ent.InsertRevisionParams

			BeforeEach(func() {
				params = &ent.InsertRevisionParams{}
				params.SetRevision(entity)
			})

			It("inserts a revision", func(ctx SpecContext) {
				revision, err := gateway.InsertRevision(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(revision).To(BeComparableTo(entity))
			})
		})

		Describe("UpsertRevision", func() {
			var params *ent.UpsertRevisionParams

			BeforeEach(func() {
				params = &ent.UpsertRevisionParams{}
				params.SetRevision(entity)
			})

			It("inserts a revision", func(ctx SpecContext) {
				revision, err := gateway.UpsertRevision(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(revision).To(BeComparableTo(entity))
			})
		})

		Describe("UpdateRevision", func() {
			var params *ent.UpdateRevisionParams

			BeforeEach(func() {
				params = &ent.UpdateRevisionParams{}
				params.UpdateMask = []string{"total"}
				params.SetRevision(entity)
			})

			It("updates a revision", func(ctx SpecContext) {
				revision, err := gateway.UpdateRevision(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(revision).To(BeComparableTo(entity))
			})
		})

		Describe("GetRevision", func() {
			var params *ent.GetRevisionParams

			BeforeEach(func() {
				params = &ent.GetRevisionParams{}
				params.SetRevision(entity)
			})

			It("returns a revision", func(ctx SpecContext) {
				revision, err := gateway.GetRevision(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(revision).To(BeComparableTo(entity))
			})
		})

		Describe("DeleteRevision", func() {
			var params *ent.DeleteRevisionParams

			BeforeEach(func() {
				params = &ent.DeleteRevisionParams{}
				params.SetRevision(entity)
			})

			It("deletes a revision", func(ctx SpecContext) {
				revision, err := gateway.DeleteRevision(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(revision).To(BeComparableTo(entity))
			})
		})

		Describe("ExecInsertRevision", func() {
			var params *ent.ExecInsertRevisionParams

			BeforeEach(func() {
				params = &ent.ExecInsertRevisionParams{}
				params.SetRevision(entity)
			})

			It("inserts a revision", func(ctx SpecContext) {
				Expect(gateway.ExecInsertRevision(ctx, params)).To(Succeed())
			})
		})

		Describe("ExecUpsertRevision", func() {
			var params *ent.ExecUpsertRevisionParams

			BeforeEach(func() {
				params = &ent.ExecUpsertRevisionParams{}
				params.SetRevision(entity)
			})

			It("inserts a revision", func(ctx SpecContext) {
				Expect(gateway.ExecUpsertRevision(ctx, params)).To(Succeed())
			})
		})

		Describe("ExecUpdateRevision", func() {
			var params *ent.ExecUpdateRevisionParams

			BeforeEach(func() {
				params = &ent.ExecUpdateRevisionParams{}
				params.UpdateMask = []string{"total"}
				params.SetRevision(entity)
			})

			It("updates a revision", func(ctx SpecContext) {
				Expect(gateway.ExecUpdateRevision(ctx, params)).To(Succeed())
			})
		})

		Describe("ExecDeleteRevision", func() {
			var params *ent.ExecDeleteRevisionParams

			BeforeEach(func() {
				params = &ent.ExecDeleteRevisionParams{}
				params.SetRevision(entity)
			})

			It("inserts a revision", func(ctx SpecContext) {
				Expect(gateway.ExecDeleteRevision(ctx, params)).To(Succeed())
			})
		})

		Describe("ListRevisions", func() {
			var params *ent.ListRevisionsParams

			BeforeEach(func() {
				params = &ent.ListRevisionsParams{}
			})

			It("lists all revisions", func(ctx SpecContext) {
				revisions, err := gateway.ListRevisions(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(revisions).To(BeEmpty())
			})

			When("the gateway returns an error", func() {
				var gateway ent.Gateway

				BeforeEach(func() {
					db := &FakeDBTX{}
					db.QueryReturns(nil, fmt.Errorf("oh no"))
					gateway = ent.New(db)
				})

				It("returns an error", func(ctx SpecContext) {
					revisions, err := gateway.ListRevisions(ctx, params)
					Expect(err).To(MatchError("oh no"))
					Expect(revisions).To(BeEmpty())
				})
			})

			When("the gateway rows scan return an error", func() {
				var gateway ent.Gateway

				BeforeEach(func() {
					rows := &FakeRows{}
					rows.NextReturns(true)
					rows.ScanReturns(fmt.Errorf("oh no"))

					db := &FakeDBTX{}
					db.QueryReturns(rows, nil)
					gateway = ent.New(db)
				})

				It("returns an error", func(ctx SpecContext) {
					revisions, err := gateway.ListRevisions(ctx, params)
					Expect(err).To(MatchError("oh no"))
					Expect(revisions).To(BeEmpty())
				})
			})

			When("the gateway rows have an error", func() {
				var gateway ent.Gateway

				BeforeEach(func() {
					rows := &FakeRows{}
					rows.ErrReturns(fmt.Errorf("oh no"))

					db := &FakeDBTX{}
					db.QueryReturns(rows, nil)
					gateway = ent.New(db)
				})

				It("returns an error", func(ctx SpecContext) {
					revisions, err := gateway.ListRevisions(ctx, params)
					Expect(err).To(MatchError("oh no"))
					Expect(revisions).To(BeEmpty())
				})
			})
		})
	})
})
