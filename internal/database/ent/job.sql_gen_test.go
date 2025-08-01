package ent_test

import (
	"github.com/aws-contrib/aurora/internal/database/ent"

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
		var params *ent.WaitForJobParams

		BeforeEach(func() {
			params = &ent.WaitForJobParams{}
			params.JobID = "test-job-id"
		})

		It("waits for a job to be completed", func(ctx SpecContext) {
			ok, err := gateway.WaitForJob(ctx, params)
			Expect(err).To(HaveOccurred())
			Expect(ok).To(BeFalse(), "Expected job to not be completed yet")
		})
	})
})
