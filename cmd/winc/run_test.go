package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/Microsoft/hcsshim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var _ = Describe("Run", func() {
	var (
		containerId string
		bundleSpec  specs.Spec
	)

	BeforeEach(func() {
		containerId = filepath.Base(bundlePath)
		bundleSpec = runtimeSpecGenerator(createSandbox(imageStore, rootfsPath, containerId))
	})

	AfterEach(func() {
		execute(exec.Command(wincBin, "delete", containerId))
		_, _, err := execute(exec.Command(wincImageBin, "--store", imageStore, "delete", containerId))
		Expect(err).NotTo(HaveOccurred())
	})

	It("creates a container and runs the init process", func() {
		writeSpec(bundlePath, bundleSpec)
		_, _, err := execute(exec.Command(wincBin, "run", "-b", bundlePath, "--detach", containerId))
		Expect(err).ToNot(HaveOccurred())

		Expect(containerExists(containerId)).To(BeTrue())

		pl := containerProcesses(containerId, "powershell.exe")
		Expect(len(pl)).To(Equal(1))

		containerPid := getContainerState(containerId).Pid
		Expect(isParentOf(containerPid, int(pl[0].ProcessId))).To(BeTrue())
	})

	Context("when the --detach flag is passed", func() {
		BeforeEach(func() {
			bundleSpec.Process.Args = []string{"cmd.exe", "/C", "waitfor fivesec /T 5 >NULL & exit /B 0"}
		})

		It("the process runs in the container and returns immediately", func() {
			writeSpec(bundlePath, bundleSpec)
			_, _, err := execute(exec.Command(wincBin, "run", "-b", bundlePath, "--detach", containerId))
			Expect(err).ToNot(HaveOccurred())

			pl := containerProcesses(containerId, "cmd.exe")
			Expect(len(pl)).To(Equal(1))

			containerPid := getContainerState(containerId).Pid
			Expect(isParentOf(containerPid, int(pl[0].ProcessId))).To(BeTrue())

			Eventually(func() []hcsshim.ProcessListItem {
				return containerProcesses(containerId, "cmd.exe")
			}, "10s").Should(BeEmpty())
		})
	})

	Context("when the --detach flag is not passed", func() {
		It("the process runs in the container and returns the exit code when the process finishes", func() {
			bundleSpec.Process.Args = []string{"cmd.exe", "/C", "exit /B 5"}
			writeSpec(bundlePath, bundleSpec)
			cmd := exec.Command(wincBin, "run", "-b", bundlePath, containerId)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session).Should(gexec.Exit(5))

			pl := containerProcesses(containerId, "cmd.exe")
			Expect(len(pl)).To(Equal(0))
		})

		It("passes stdin through to the process", func() {
			bundleSpec.Process.Args = []string{"C:\\temp\\read.exe"}
			bundleSpec.Mounts = []specs.Mount{
				{
					Source:      filepath.Dir(readBin),
					Destination: "C:\\temp",
				},
			}
			writeSpec(bundlePath, bundleSpec)
			cmd := exec.Command(wincBin, "run", "-b", bundlePath, containerId)
			cmd.Stdin = strings.NewReader("hey-winc\n")
			stdOut, _, err := execute(cmd)
			Expect(err).NotTo(HaveOccurred())
			Expect(stdOut.String()).To(ContainSubstring("hey-winc"))
		})

		It("captures the stdout", func() {
			bundleSpec.Process.Args = []string{"cmd.exe", "/C", "echo hey-winc"}
			writeSpec(bundlePath, bundleSpec)
			stdOut, _, err := execute(exec.Command(wincBin, "run", "-b", bundlePath, containerId))
			Expect(err).NotTo(HaveOccurred())
			Expect(stdOut.String()).To(ContainSubstring("hey-winc"))
		})

		It("captures the stderr", func() {
			bundleSpec.Process.Args = []string{"cmd.exe", "/C", "echo hey-winc 1>&2"}
			writeSpec(bundlePath, bundleSpec)
			_, stdErr, err := execute(exec.Command(wincBin, "run", "-b", bundlePath, containerId))
			Expect(err).ToNot(HaveOccurred())
			Expect(stdErr.String()).To(ContainSubstring("hey-winc"))
		})

		It("captures the CTRL+C", func() {
			bundleSpec.Process.Args = []string{"cmd.exe", "/C", "echo hey-winc & waitfor ever /T 9999"}
			writeSpec(bundlePath, bundleSpec)
			cmd := exec.Command(wincBin, "run", "-b", bundlePath, containerId)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
			}
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Consistently(session).ShouldNot(gexec.Exit(0))
			Eventually(session.Out).Should(gbytes.Say("hey-winc"))
			pl := containerProcesses(containerId, "cmd.exe")
			Expect(len(pl)).To(Equal(1))

			sendCtrlBreak(session)
			Eventually(session).Should(gexec.Exit(1067))
			pl = containerProcesses(containerId, "cmd.exe")
			Expect(len(pl)).To(Equal(0))
		})
	})

	Context("when the '--pid-file' flag is provided", func() {
		var pidFile string

		BeforeEach(func() {
			f, err := ioutil.TempFile("", "pid")
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Close()).To(Succeed())
			pidFile = f.Name()
		})

		AfterEach(func() {
			Expect(os.RemoveAll(pidFile)).To(Succeed())
		})

		It("places the container pid in the specified file", func() {
			bundleSpec.Process.Args = []string{"cmd.exe", "/C", "waitfor ever /T 9999"}
			writeSpec(bundlePath, bundleSpec)
			_, _, err := execute(exec.Command(wincBin, "run", "-b", bundlePath, "--detach", "--pid-file", pidFile, containerId))
			Expect(err).ToNot(HaveOccurred())

			containerPid := getContainerState(containerId).Pid

			pidBytes, err := ioutil.ReadFile(pidFile)
			Expect(err).ToNot(HaveOccurred())
			pid, err := strconv.ParseInt(string(pidBytes), 10, 64)
			Expect(err).ToNot(HaveOccurred())
			Expect(int(pid)).To(Equal(containerPid))
		})
	})

	Context("when the '--no-new-keyring' flag is provided", func() {
		It("ignores it and creates and starts a container", func() {
			writeSpec(bundlePath, bundleSpec)
			_, _, err := execute(exec.Command(wincBin, "run", "-b", bundlePath, "--detach", "--no-new-keyring", containerId))
			Expect(err).ToNot(HaveOccurred())
			Expect(containerExists(containerId)).To(BeTrue())
		})
	})

	Context("when the container exists", func() {
		BeforeEach(func() {
			writeSpec(bundlePath, bundleSpec)
			_, _, err := execute(exec.Command(wincBin, "create", "-b", bundlePath, containerId))
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_, _, err := execute(exec.Command(wincBin, "delete", containerId))
			Expect(err).ToNot(HaveOccurred())
		})

		It("errors", func() {
			_, stdErr, err := execute(exec.Command(wincBin, "run", "-b", bundlePath, "--detach", containerId))
			Expect(err).To(HaveOccurred())
			expectedErrorMsg := fmt.Sprintf("container with id already exists: %s", containerId)
			Expect(stdErr.String()).To(ContainSubstring(expectedErrorMsg))
		})
	})

	Context("when the bundlePath is not specified", func() {
		It("uses the current directory as the bundlePath", func() {
			writeSpec(bundlePath, bundleSpec)
			createCmd := exec.Command(wincBin, "run", "--detach", containerId)
			createCmd.Dir = bundlePath
			_, _, err := execute(createCmd)
			Expect(err).ToNot(HaveOccurred())
			Expect(containerExists(containerId)).To(BeTrue())
		})
	})
})

func writeSpec(path string, spec specs.Spec) {
	config, err := json.Marshal(&spec)
	Expect(err).NotTo(HaveOccurred())
	Expect(ioutil.WriteFile(filepath.Join(path, "config.json"), config, 0666)).To(Succeed())
}
