package integration_test

import (
	"os/exec"

	. "github.com/NathMcBride/digest-authentication/integration/support"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var path string
var app *gexec.Session

const healthCheck = "http://localhost:8080/health"

var _ = BeforeSuite(func() {
	var err error
	path, err = gexec.Build("github.com/NathMcBride/digest-authentication/src")
	Expect(err).NotTo(HaveOccurred())

	DeferCleanup(gexec.CleanupBuildArtifacts)
})

var _ = BeforeEach(func() {
	cmd := exec.Command(path)
	var err error
	app, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(WaitForReady(2, healthCheck)).Should(Succeed(),
		CSprintf(Orange, "%s Not found", healthCheck))

	GinkgoWriter.Println(CSprintf(Green, "Started ") + CSprintf(Cyan, healthCheck))
})

var _ = AfterEach(func() {
	timeOut := 5

	app.Terminate()
	EventuallyWithOffset(1, app, timeOut).Should(gexec.Exit(), CSprintf(Orange, "Failed to shutdown"))
	GinkgoWriter.Println(CSprintf(Green, "Stopped ") + CSprintf(Cyan, healthCheck))
})

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
