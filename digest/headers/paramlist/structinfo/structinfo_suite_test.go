package structinfo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStructinfo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Structinfo Suite")
}
