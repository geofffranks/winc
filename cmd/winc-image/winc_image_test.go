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

	"github.com/Microsoft/hcsshim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var _ = Describe("WincImage", func() {
	var (
		storePath   string
		containerId string
	)

	BeforeEach(func() {
		var err error
		containerId = randomContainerId()
		storePath, err = ioutil.TempDir("", "container-store")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(storePath)).To(Succeed())
	})

	It("creates and deletes a sandbox", func() {
		stdout, _, err := execute(wincImageBin, "--store", storePath, "create", rootfsPath, containerId)
		Expect(err).NotTo(HaveOccurred())

		var desiredSpec specs.Spec
		Expect(json.Unmarshal(stdout.Bytes(), &desiredSpec)).To(Succeed())
		volumeGuid := getVolumeGuid(storePath, containerId)
		Expect(desiredSpec.Version).To(Equal(specs.Version))
		Expect(desiredSpec.Root.Path).To(Equal(volumeGuid))
		Expect(desiredSpec.Windows.LayerFolders).ToNot(BeEmpty())
		Expect(desiredSpec.Windows.LayerFolders[0]).To(Equal(rootfsPath))
		for _, layer := range desiredSpec.Windows.LayerFolders {
			Expect(layer).To(BeADirectory())
		}

		driverInfo := hcsshim.DriverInfo{HomeDir: storePath, Flavour: 1}
		Expect(hcsshim.LayerExists(driverInfo, containerId)).To(BeTrue())

		_, _, err = execute(wincImageBin, "--store", storePath, "delete", containerId)
		Expect(err).To(Succeed())

		Expect(hcsshim.LayerExists(driverInfo, containerId)).To(BeFalse())
		Expect(filepath.Join(driverInfo.HomeDir, containerId)).NotTo(BeADirectory())
	})

	Context("when provided --log <log-file>", func() {
		var (
			logFile string
			tempDir string
		)

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "log-dir")
			Expect(err).NotTo(HaveOccurred())

			logFile = filepath.Join(tempDir, "winc-image.log")
		})

		AfterEach(func() {
			_, _, err := execute(wincImageBin, "--store", storePath, "delete", containerId)
			Expect(err).To(Succeed())
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		})

		Context("when the provided log file path does not exist", func() {
			BeforeEach(func() {
				logFile = filepath.Join(tempDir, "some-dir", "winc-image.log")
			})

			It("creates the full path", func() {
				_, _, err := execute(wincImageBin, "--log", logFile, "--store", storePath, "create", rootfsPath, containerId)
				Expect(err).NotTo(HaveOccurred())

				Expect(logFile).To(BeAnExistingFile())
			})
		})

		Context("when it runs successfully", func() {
			It("does not log to the specified file", func() {
				_, _, err := execute(wincImageBin, "--log", logFile, "--store", storePath, "create", rootfsPath, containerId)
				Expect(err).NotTo(HaveOccurred())

				contents, err := ioutil.ReadFile(logFile)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(contents)).To(BeEmpty())
			})

			Context("when provided --debug", func() {
				It("outputs debug level logs", func() {
					_, _, err := execute(wincImageBin, "--log", logFile, "--debug", "--store", storePath, "create", rootfsPath, containerId)
					Expect(err).NotTo(HaveOccurred())

					contents, err := ioutil.ReadFile(logFile)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(contents)).NotTo(BeEmpty())
				})
			})
		})

		Context("when it errors", func() {
			It("logs errors to the specified file", func() {
				execute(wincImageBin, "--log", logFile, "--store", storePath, "create", "garbage-something", containerId)

				contents, err := ioutil.ReadFile(logFile)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(contents)).NotTo(BeEmpty())
			})
		})
	})

	Context("when using unix style rootfsPath", func() {
		var (
			tempRootfs string
			tempdir    string
			err        error
		)

		BeforeEach(func() {
			tempdir, err = ioutil.TempDir("", "rootfs")
			Expect(err).ToNot(HaveOccurred())
			err := exec.Command("cmd.exe", "/c", fmt.Sprintf("mklink /D %s %s", filepath.Join(tempdir, "rootfs"), rootfsPath)).Run()
			Expect(err).ToNot(HaveOccurred())

			tempRootfs = strings.Replace(tempdir, "C:", "", -1) + "/rootfs"
		})

		AfterEach(func() {
			// remove symlink so we don't clobber rootfs dir
			err := exec.Command("cmd.exe", "/c", fmt.Sprintf("rmdir %s", filepath.Join(tempdir, "rootfs"))).Run()
			Expect(err).ToNot(HaveOccurred())
			Expect(os.RemoveAll(tempdir)).To(Succeed())
		})

		destToWindowsPath := func(input string) string {
			vol := filepath.VolumeName(input)
			if vol == "" {
				input = filepath.Join("C:", input)
			}
			return filepath.Clean(input)
		}

		It("creates and deletes a sandbox with unix rootfsPath", func() {
			stdout, _, err := execute(wincImageBin, "--store", storePath, "create", tempRootfs, containerId)
			Expect(err).NotTo(HaveOccurred())

			var desiredSpec specs.Spec
			Expect(json.Unmarshal(stdout.Bytes(), &desiredSpec)).To(Succeed())
			volumeGuid := getVolumeGuid(storePath, containerId)
			Expect(desiredSpec.Version).To(Equal(specs.Version))
			Expect(desiredSpec.Root.Path).To(Equal(volumeGuid))
			Expect(desiredSpec.Windows.LayerFolders).ToNot(BeEmpty())
			Expect(desiredSpec.Windows.LayerFolders[0]).To(Equal(destToWindowsPath(tempRootfs)))
			for _, layer := range desiredSpec.Windows.LayerFolders {
				Expect(layer).To(BeADirectory())
			}

			driverInfo := hcsshim.DriverInfo{HomeDir: storePath, Flavour: 1}
			Expect(hcsshim.LayerExists(driverInfo, containerId)).To(BeTrue())

			_, _, err = execute(wincImageBin, "--store", storePath, "delete", containerId)
			Expect(err).To(Succeed())

			Expect(hcsshim.LayerExists(driverInfo, containerId)).To(BeFalse())
		})
	})

	Context("when provided a disk limit", func() {
		Context("when the disk limit is valid", func() {
			var (
				mountPath          string
				diskLimitSizeBytes int
				volumeGuid         string
			)

			BeforeEach(func() {
				diskLimitSizeBytes = 50 * 1024 * 1024
			})

			JustBeforeEach(func() {
				_, _, err := execute(wincImageBin, "--store", storePath, "create", "--disk-limit-size-bytes", strconv.Itoa(diskLimitSizeBytes), rootfsPath, containerId)
				Expect(err).NotTo(HaveOccurred())

				mountPath, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())

				volumeGuid = getVolumeGuid(storePath, containerId)
				Expect(exec.Command("mountvol", mountPath, volumeGuid).Run()).To(Succeed())
			})

			AfterEach(func() {
				Expect(exec.Command("mountvol", mountPath, "/D").Run()).To(Succeed())
				Expect(os.RemoveAll(mountPath)).To(Succeed())
				_, _, err := execute(wincImageBin, "--store", storePath, "delete", containerId)
				Expect(err).To(Succeed())
			})

			It("doesn't allow files large than the limit to be created", func() {
				largeFilePath := filepath.Join(mountPath, "file.txt")
				Expect(exec.Command("fsutil", "file", "createnew", largeFilePath, strconv.Itoa(diskLimitSizeBytes+1)).Run()).ToNot(Succeed())
				Expect(largeFilePath).ToNot(BeAnExistingFile())
			})

			It("allows files at the limit to be created", func() {
				largeFilePath := filepath.Join(mountPath, "file.txt")
				Expect(exec.Command("fsutil", "file", "createnew", largeFilePath, strconv.Itoa(diskLimitSizeBytes)).Run()).To(Succeed())
				Expect(largeFilePath).To(BeAnExistingFile())
			})

			Context("when the provided disk limit is 0", func() {
				BeforeEach(func() {
					diskLimitSizeBytes = 0
				})

				It("does not set a limit", func() {
					output, err := exec.Command("dirquota", "quota", "list", fmt.Sprintf("/Path:%s", mountPath)).CombinedOutput()
					Expect(err).To(HaveOccurred())
					Expect(string(output)).To(ContainSubstring("The requested object was not found"))
				})
			})
		})

		Context("when the provided disk limit is below 0", func() {
			It("errors", func() {
				_, _, err := execute(wincImageBin, "--store", storePath, "create", "--disk-limit-size-bytes", "-5", rootfsPath, containerId)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("when creating the sandbox layer fails", func() {
		It("errors", func() {
			_, stderr, err := execute(wincImageBin, "create", "some-bad-rootfs", "")
			Expect(err).To(HaveOccurred())
			Expect(stderr.String()).To(ContainSubstring("rootfs layer does not exist"))
		})
	})

	Context("deleting when provided a nonexistent containerId", func() {
		var logFile string

		BeforeEach(func() {
			logF, err := ioutil.TempFile("", "winc-image.log")
			Expect(err).NotTo(HaveOccurred())
			logFile = logF.Name()
			logF.Close()
		})

		AfterEach(func() {
			Expect(os.Remove(logFile)).To(Succeed())
		})

		It("logs a warning", func() {
			_, _, err := execute(wincImageBin, "--log", logFile, "delete", "some-bad-container-id")
			Expect(err).ToNot(HaveOccurred())

			contents, err := ioutil.ReadFile(logFile)
			Expect(string(contents)).To(ContainSubstring("Layer `some-bad-container-id` not found. Skipping delete."))
		})
	})

	Context("when create is called with the wrong number of args", func() {
		It("prints the usage", func() {
			stdOut, _, err := execute(wincImageBin, "create")
			Expect(err).To(HaveOccurred())
			Expect(stdOut.String()).To(ContainSubstring("Incorrect Usage"))
		})
	})

	Context("when delete is called with the wrong number of args", func() {
		It("prints the usage", func() {
			stdOut, _, err := execute(wincImageBin, "delete")
			Expect(err).To(HaveOccurred())
			Expect(stdOut.String()).To(ContainSubstring("Incorrect Usage"))
		})
	})

	Context("deleting after failed attempts", func() {
		var (
			driverInfo    hcsshim.DriverInfo
			sandboxLayers []string
		)

		BeforeEach(func() {
			driverInfo = hcsshim.DriverInfo{HomeDir: storePath, Flavour: 1}

			parentLayerChain, err := ioutil.ReadFile(filepath.Join(rootfsPath, "layerchain.json"))
			Expect(err).NotTo(HaveOccurred())
			parentLayers := []string{}
			Expect(json.Unmarshal(parentLayerChain, &parentLayers)).To(Succeed())

			sandboxLayers = append([]string{rootfsPath}, parentLayers...)
		})

		Context("when a layer has been created but is not activated", func() {
			It("destroys the layer", func() {
				Expect(hcsshim.CreateSandboxLayer(driverInfo, containerId, rootfsPath, sandboxLayers)).To(Succeed())

				_, _, err := execute(wincImageBin, "--store", storePath, "delete", containerId)
				Expect(err).NotTo(HaveOccurred())
				Expect(hcsshim.LayerExists(driverInfo, containerId)).To(BeFalse())
			})
		})

		Context("when a layer has been created and activated but is not prepared", func() {
			It("destroys the layer", func() {
				Expect(hcsshim.CreateSandboxLayer(driverInfo, containerId, rootfsPath, sandboxLayers)).To(Succeed())
				Expect(hcsshim.ActivateLayer(driverInfo, containerId)).To(Succeed())

				_, _, err := execute(wincImageBin, "--store", storePath, "delete", containerId)
				Expect(err).NotTo(HaveOccurred())
				Expect(hcsshim.LayerExists(driverInfo, containerId)).To(BeFalse())
			})
		})
	})
})
