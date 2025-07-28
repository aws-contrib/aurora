package ent_test

import (
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

	Describe("Job", func() {
		var entity *ent.Job

		BeforeAll(func() {
			entity = NewFakeJob()
		})

		Describe("CreateSchemaSys", func() {
			It("creates the sys schema", func(ctx SpecContext) {
				Expect(gateway.CreateSchemaSys(ctx)).To(Succeed())
			})
		})

		Describe("CreateTableJobs", func() {
			It("creates the sys.jobs table", func(ctx SpecContext) {
				Expect(gateway.CreateTableJobs(ctx)).To(Succeed())
			})
		})

		Describe("InsertJob", func() {
			var params *ent.InsertJobParams

			BeforeEach(func() {
				params = &ent.InsertJobParams{}
				params.SetJob(entity)
			})

			It("inserts a job", func(ctx SpecContext) {
				job, err := gateway.InsertJob(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(job).To(BeComparableTo(entity))
			})
		})

		Describe("GetJob", func() {
			var params *ent.GetJobParams

			BeforeEach(func() {
				params = &ent.GetJobParams{}
				params.SetJob(entity)
			})

			It("returns a job", func(ctx SpecContext) {
				job, err := gateway.GetJob(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(job).To(BeComparableTo(entity))
			})
		})

		Describe("DeleteJob", func() {
			var params *ent.DeleteJobParams

			BeforeEach(func() {
				params = &ent.DeleteJobParams{}
				params.SetJob(entity)
			})

			It("deletes a job", func(ctx SpecContext) {
				job, err := gateway.DeleteJob(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(job).To(BeComparableTo(entity))
			})
		})

		Describe("ExecInsertJob", func() {
			var params *ent.ExecInsertJobParams

			BeforeEach(func() {
				params = &ent.ExecInsertJobParams{}
				params.SetJob(entity)
			})

			It("inserts a job", func(ctx SpecContext) {
				Expect(gateway.ExecInsertJob(ctx, params)).To(Succeed())
			})
		})

		Describe("ExecDeleteJob", func() {
			var params *ent.ExecDeleteJobParams

			BeforeEach(func() {
				params = &ent.ExecDeleteJobParams{}
				params.SetJob(entity)
			})

			It("inserts a revision", func(ctx SpecContext) {
				Expect(gateway.ExecDeleteJob(ctx, params)).To(Succeed())
			})
		})
	})
})
