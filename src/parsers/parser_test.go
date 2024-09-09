package parsers_test

import (
	. "github.com/NathMcBride/digest-authentication/src/parsers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Describe("ParsList", func() {
		It("can successfully parse a list", func() {
			parser := Parser{}

			result, err := parser.ParseList(`Digest item1="a-value", item2=a-value-2`, "Digest ")

			Expect(err).NotTo(HaveOccurred())
			expected := map[string]string{
				"item1": "a-value",
				"item2": "a-value-2",
			}
			Expect(result).To(Equal(expected))
		})

		When("an empty prefix is supplied", func() {
			It("can successfully parse a list", func() {
				parser := Parser{}

				result, err := parser.ParseList(`item1="a-value", item2=a-value-2`, "")

				Expect(err).NotTo(HaveOccurred())
				expected := map[string]string{
					"item1": "a-value",
					"item2": "a-value-2",
				}
				Expect(result).To(Equal(expected))
			})
		})

		When("a prefix is supplied, and the prefix is not present in value to parse", func() {
			It("returns an error", func() {
				parser := Parser{}

				_, err := parser.ParseList(`item1="a-value", item2=a-value-2`, "Digest ")

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
