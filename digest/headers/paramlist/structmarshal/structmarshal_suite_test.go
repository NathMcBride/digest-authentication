package structmarshal_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStructmarshal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Structmarshal Suite")
}
