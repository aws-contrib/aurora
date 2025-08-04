//go:build !goverter

package ent_test

import (
	"fmt"

	"github.com/aws-contrib/aurora/internal/database/ent"

	. "github.com/aws-contrib/aurora/internal/database/ent/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobRepository", func() {
	var repository *ent.JobRepository

	BeforeEach(func(ctx SpecContext) {
		repository = &ent.JobRepository{
			Gateway: NewFakeGateway(),
		}
	})

	Describe("WaitJob", func() {
		var params *ent.WaitJobParams

		BeforeEach(func() {
			entity := NewFakeJob()
			entity.Status = "processing"

			gateway := repository.Gateway.(*FakeGateway)
			gateway.GetJobReturnsOnCall(0, entity, nil)

			params = &ent.WaitJobParams{}
			params.ID = entity.ID
		})

		When("the job is completed", func() {
			BeforeEach(func() {
				entity := NewFakeJob()
				entity.Status = "completed"

				gateway := repository.Gateway.(*FakeGateway)
				gateway.GetJobReturnsOnCall(1, entity, nil)
			})

			It("waits for the job to complete", func(ctx SpecContext) {
				job, err := repository.WaitJob(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(job).NotTo(BeNil())
			})
		})

		When("the job is failed", func() {
			BeforeEach(func() {
				entity := NewFakeJob()
				entity.Status = "failed"

				gateway := repository.Gateway.(*FakeGateway)
				gateway.GetJobReturnsOnCall(1, entity, nil)
			})

			It("waits for the job to fail", func(ctx SpecContext) {
				job, err := repository.WaitJob(ctx, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(job).NotTo(BeNil())
			})
		})

		When("the gateway fails", func() {
			BeforeEach(func() {
				gateway := repository.Gateway.(*FakeGateway)
				gateway.GetJobReturns(nil, fmt.Errorf("oh no"))
			})

			It("returns an error", func(ctx SpecContext) {
				job, err := repository.WaitJob(ctx, params)
				Expect(err).To(MatchError("oh no"))
				Expect(job).To(BeNil())
			})
		})
	})
})
