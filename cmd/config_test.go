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
			Expect(config.UnmarshalText(data)).To(HaveOccurred())
		})
	})
})
