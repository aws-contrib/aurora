package ent_test

import (
	"fmt"

	"github.com/aws-contrib/aurora/internal/database/ent"

	. "github.com/aws-contrib/aurora/internal/database/ent/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RevisionRepository", func() {
	var repository *ent.RevisionRepository

	BeforeEach(func(ctx SpecContext) {
		repository = &ent.RevisionRepository{
			Gateway:    NewFakeGateway(),
			FileSystem: NewFakeFileSystem(),
		}
	})

	Describe("ApplyRevision", func() {
		var params *ent.ApplyRevisionParams

		BeforeEach(func() {
			params = &ent.ApplyRevisionParams{}
			params.Revision = NewFakeRevision()
		})

		It("applies a revision", func(ctx SpecContext) {
			Expect(repository.ApplyRevision(ctx, params)).To(Succeed())
		})

		ItReturnsError := func(msg string) {
			It("returns an error", func(ctx SpecContext) {
				Expect(repository.ApplyRevision(ctx, params)).To(MatchError(msg))
			})
		}

		When("the file system fails", func() {
			BeforeEach(func() {
				fs := repository.FileSystem.(*FakeFileSystem)
				fs.ReadFileReturns(nil, fmt.Errorf("oh no"))
			})

			ItReturnsError("oh no")
		})

		When("the gateway fails", func() {
			When("the upsert revision fails", func() {
				BeforeEach(func() {
					gateway := repository.Gateway.(*FakeGateway)
					gateway.UpsertRevisionReturns(nil, fmt.Errorf("oh no"))
				})

				ItReturnsError("oh no")
			})

			When("the execute revision fails", func() {
				BeforeEach(func() {
					row := &FakeRow{}
					row.ScanReturns(fmt.Errorf("oh no"))

					tx := repository.Gateway.(*FakeGateway).Tx().(*FakeDBTX)
					tx.QueryRowReturns(row)
				})

				It("does not return an error", func(ctx SpecContext) {
					Expect(repository.ApplyRevision(ctx, params)).To(Succeed())
					Expect(params.Revision.Error).NotTo(BeNil())
					Expect(*params.Revision.Error).To(Equal("oh no"))
				})
			})

			When("the update revision fails", func() {
				BeforeEach(func() {
					gateway := repository.Gateway.(*FakeGateway)
					gateway.ExecUpdateRevisionReturns(fmt.Errorf("oh no"))
				})

				ItReturnsError("oh no")
			})

			When("the job is not found", func() {
				BeforeEach(func() {
					gateway := repository.Gateway.(*FakeGateway)
					gateway.GetJobReturns(nil, ent.ErrNoRows)
				})

				It("applies a revision", func(ctx SpecContext) {
					Expect(repository.ApplyRevision(ctx, params)).To(Succeed())
					Expect(params.Revision.Error).To(BeNil())
				})
			})

			When("the job waiting fails", func() {
				BeforeEach(func() {
					gateway := repository.Gateway.(*FakeGateway)
					gateway.GetJobReturns(nil, fmt.Errorf("oh no"))
				})

				It("does not return an error", func(ctx SpecContext) {
					Expect(repository.ApplyRevision(ctx, params)).To(Succeed())
					Expect(params.Revision.Error).NotTo(BeNil())
					Expect(*params.Revision.Error).To(Equal("oh no"))
				})
			})

			When("the job fails", func() {
				BeforeEach(func() {
					entity := NewFakeJob()
					entity.Status = "failed"
					entity.Details = "oh no"

					gateway := repository.Gateway.(*FakeGateway)
					gateway.GetJobReturns(entity, nil)
				})

				It("does not return an error", func(ctx SpecContext) {
					Expect(repository.ApplyRevision(ctx, params)).To(Succeed())
					Expect(params.Revision.Error).NotTo(BeNil())
					Expect(*params.Revision.Error).To(Equal("oh no"))
				})
			})
		})
	})

	Describe("ListRevisions", func() {
		var params *ent.ListRevisionsParams

		BeforeEach(func() {
			params = &ent.ListRevisionsParams{}
		})

		It("list the revisions", func(ctx SpecContext) {
			revisions, err := repository.ListRevisions(ctx, params)
			Expect(err).NotTo(HaveOccurred())
			Expect(revisions).NotTo(BeEmpty())
		})

		ItReturnsError := func(msg string) {
			It("returns an error", func(ctx SpecContext) {
				revisions, err := repository.ListRevisions(ctx, params)
				Expect(err).To(MatchError(msg))
				Expect(revisions).To(BeEmpty())
			})
		}

		When("the file system fails", func() {
			BeforeEach(func() {
				fs := repository.FileSystem.(*FakeFileSystem)
				fs.GlobReturns(nil, fmt.Errorf("oh no"))
			})

			ItReturnsError("oh no")
		})

		When("the file system fails", func() {
			BeforeEach(func() {
				fs := repository.FileSystem.(*FakeFileSystem)
				fs.ReadFileReturns(nil, fmt.Errorf("oh no"))
			})

			ItReturnsError("oh no")
		})

		When("the gateway fails", func() {
			BeforeEach(func() {
				gateway := repository.Gateway.(*FakeGateway)
				gateway.GetRevisionReturns(nil, fmt.Errorf("oh no"))
			})

			ItReturnsError("oh no")
		})
	})
})
