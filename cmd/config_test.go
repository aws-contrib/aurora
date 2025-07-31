package cmd_test

import (
	"os"

	"github.com/aws-contrib/aurora/cmd"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var config *cmd.Config

	BeforeEach(func() {
		config = &cmd.Config{}
	})

	Describe("UnmarshalText", func() {
		It("unmarshals valid HCL text", func() {
			data, err := os.ReadFile("aurora.hcl")
			Expect(err).ToNot(HaveOccurred())
			Expect(config.UnmarshalText(data)).To(Succeed())

			Expect(config.Environments).To(HaveLen(1))
			Expect(config.Environments[0].Name).To(Equal("aws"))
			Expect(config.Environments[0].Migration.GetDir()).To(Equal("file://database/migration"))
			Expect(config.Environments[0].GetURL()).To(Equal("postgres://example-api:DSQL_TOKEN@example.com/example-api"))
			Expect(config.Data).To(HaveLen(1))
			Expect(config.Variables).To(HaveLen(3))
		})
	})
})
