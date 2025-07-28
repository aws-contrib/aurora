package ent_test

import (
	"fmt"

	"github.com/aws-contrib/aurora/internal/database/ent"

	. "github.com/aws-contrib/aurora/internal/database/ent/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MigrationRepository", func() {
	var repository *ent.MigrationRepository

	BeforeEach(func(ctx SpecContext) {
		repository = &ent.MigrationRepository{
			Gateway:    NewFakeGateway(),
			FileSystem: NewFakeFileSystem(),
		}
	})

	Describe("ApplyMigration", func() {
		var params *ent.ApplyMigrationParams

		BeforeEach(func() {
			params = &ent.ApplyMigrationParams{}
			params.Migration = NewFakeMigration()
		})

		It("applies a revision", func(ctx SpecContext) {
			Expect(repository.ApplyMigration(ctx, params)).To(Succeed())
		})

		ItReturnsError := func(msg string) {
			It("returns an error", func(ctx SpecContext) {
				Expect(repository.ApplyMigration(ctx, params)).To(MatchError(msg))
			})
		}

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

					tx := repository.Gateway.(*FakeGateway).Database().(*FakeDBTX)
					tx.QueryRowReturns(row)
				})

				It("does not return an error", func(ctx SpecContext) {
					Expect(repository.ApplyMigration(ctx, params)).To(Succeed())
					Expect(params.Migration.Revision.Error).NotTo(BeNil())
					Expect(*params.Migration.Revision.Error).To(Equal("oh no"))
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
					Expect(repository.ApplyMigration(ctx, params)).To(Succeed())
					Expect(params.Migration.Revision.Error).To(BeNil())
				})
			})

			When("the job waiting fails", func() {
				BeforeEach(func() {
					gateway := repository.Gateway.(*FakeGateway)
					gateway.GetJobReturns(nil, fmt.Errorf("oh no"))
				})

				It("does not return an error", func(ctx SpecContext) {
					Expect(repository.ApplyMigration(ctx, params)).To(Succeed())
					Expect(params.Migration.Revision.Error).NotTo(BeNil())
					Expect(*params.Migration.Revision.Error).To(Equal("oh no"))
				})
			})

			When("the job fails", func() {
				BeforeEach(func() {
					err := "oh no"
					entity := NewFakeJob()
					entity.Status = "failed"
					entity.Details = &err

					gateway := repository.Gateway.(*FakeGateway)
					gateway.GetJobReturns(entity, nil)
				})

				It("does not return an error", func(ctx SpecContext) {
					Expect(repository.ApplyMigration(ctx, params)).To(Succeed())
					Expect(params.Migration.Revision.Error).NotTo(BeNil())
					Expect(*params.Migration.Revision.Error).To(Equal("oh no"))
				})
			})
		})
	})

	Describe("ListMigrations", func() {
		var params *ent.ListMigrationsParams

		BeforeEach(func() {
			params = &ent.ListMigrationsParams{}
		})

		It("list the migrations", func(ctx SpecContext) {
			migrations, err := repository.ListMigrations(ctx, params)
			Expect(err).NotTo(HaveOccurred())
			Expect(migrations).NotTo(BeEmpty())
		})

		ItReturnsError := func(msg string) {
			It("returns an error", func(ctx SpecContext) {
				migrations, err := repository.ListMigrations(ctx, params)
				Expect(err).To(MatchError(msg))
				Expect(migrations).To(BeEmpty())
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
