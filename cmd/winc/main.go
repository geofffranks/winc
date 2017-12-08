package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"code.cloudfoundry.org/winc/config"
	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/hcs"
	"code.cloudfoundry.org/winc/mount"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	usage = `Open Container Initiative runtime for Windows

winc is a command line client for running applications on Windows packaged
according to the Open Container Initiative (OCI) format and is a compliant
implementation of the Open Container Initiative specification.`
	exactArgs = iota
	minArgs
	maxArgs
)

// gitCommit will be the hash that the binary was built from
// and will be populated by the build flags
var gitCommit = ""

func main() {
	app := cli.NewApp()
	app.Name = "winc.exe"
	app.Usage = usage

	var v []string
	if gitCommit != "" {
		v = append(v, fmt.Sprintf("commit: %s", gitCommit))
	}
	v = append(v, fmt.Sprintf("spec: %s", specs.Version))
	app.Version = strings.Join(v, "\n")

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output for logging",
		},
		cli.StringFlag{
			Name:  "log",
			Value: os.DevNull,
			Usage: "set the log file path where internal debug information is written",
		},
		cli.StringFlag{
			Name:  "log-format",
			Value: "json",
			Usage: "set the format used by logs ('json' (default), or 'text')",
		},
		cli.StringFlag{
			Name:  "image-store",
			Value: "C:\\run\\winc",
			Usage: "directory for storage of container state",
		},
	}

	app.Commands = []cli.Command{
		createCommand,
		deleteCommand,
		runCommand,
		stateCommand,
		execCommand,
		eventsCommand,
	}

	app.Before = func(context *cli.Context) error {
		debug := context.GlobalBool("debug")
		logFile := context.GlobalString("log")
		logFormat := context.GlobalString("log-format")

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		var logWriter io.Writer
		if logFile == "" || logFile == os.DevNull {
			logWriter = ioutil.Discard
		} else {
			if err := os.MkdirAll(filepath.Dir(logFile), 0666); err != nil {
				return err
			}

			f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0666)
			if err != nil {
				return err
			}

			logWriter = f
		}
		logrus.SetOutput(logWriter)

		switch logFormat {
		case "text":
			// retain logrus's default.
		case "json":
			logrus.SetFormatter(new(logrus.JSONFormatter))
		default:
			return &InvalidLogFormatError{Format: logFormat}
		}

		return nil
	}

	cli.ErrWriter = &fatalWriter{cli.ErrWriter}
	if err := app.Run(os.Args); err != nil {
		fatal(err)
	}
}

type fatalWriter struct {
	cliErrWriter io.Writer
}

func (f *fatalWriter) Write(p []byte) (n int, err error) {
	logrus.Error(string(p))
	return f.cliErrWriter.Write(p)
}

func checkArgs(context *cli.Context, expected, checkType int) error {
	var err error
	cmdName := context.Command.Name
	switch checkType {
	case exactArgs:
		if context.NArg() != expected {
			err = fmt.Errorf("%s: %q requires exactly %d argument(s)", os.Args[0], cmdName, expected)
		}
	case minArgs:
		if context.NArg() < expected {
			err = fmt.Errorf("%s: %q requires a minimum of %d argument(s)", os.Args[0], cmdName, expected)
		}
	case maxArgs:
		if context.NArg() > expected {
			err = fmt.Errorf("%s: %q requires a maximum of %d argument(s)", os.Args[0], cmdName, expected)
		}
	}

	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		_ = cli.ShowCommandHelp(context, cmdName)
		return err
	}
	return nil
}

func fatal(err error) {
	logrus.Error(err)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func wireContainerManager(imageStore, bundlePath, containerId string) (*container.Manager, error) {
	client := hcs.Client{}

	if bundlePath == "" {
		cp, err := client.GetContainerProperties(containerId)
		if err != nil {
			return nil, err
		}
		bundlePath = cp.Name
	}

	if filepath.Base(bundlePath) != containerId {
		return nil, &container.InvalidIdError{Id: containerId}
	}

	return container.NewManager(&client, &mount.Mounter{}, imageStore, bundlePath), nil
}

func createContainer(logger *logrus.Entry, bundlePath, imageStore, containerId, pidFile string) (*specs.Spec, error) {
	if bundlePath == "" {
		var err error
		bundlePath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	bundlePath = filepath.Clean(bundlePath)

	spec, err := config.ValidateBundle(logger, bundlePath)
	if err != nil {
		return nil, err
	}

	if _, err := config.ValidateProcess(logger, "", spec.Process); err != nil {
		return nil, err
	}

	cm, err := wireContainerManager(imageStore, bundlePath, containerId)
	if err != nil {
		return nil, err
	}

	if err := cm.Create(spec); err != nil {
		return nil, err
	}

	if pidFile != "" {
		state, err := cm.State()
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(pidFile, []byte(strconv.FormatInt(int64(state.Pid), 10)), 0666); err != nil {
			return nil, err
		}
	}

	return spec, nil
}

func runProcess(containerId string, spec *specs.Process, detach bool, pidFile string) error {
	cm, err := wireContainerManager("", "", containerId)
	if err != nil {
		return err
	}

	process, err := cm.Exec(spec, !detach)
	if err != nil {
		return err
	}

	if pidFile != "" {
		if err := ioutil.WriteFile(pidFile, []byte(strconv.FormatInt(int64(process.Pid()), 10)), 0666); err != nil {
			return err
		}
	}

	if !detach {
		stdin, stdout, stderr, err := process.Stdio()
		if err != nil {
			return err
		}

		var wg sync.WaitGroup

		go func() {
			_, _ = io.Copy(stdin, os.Stdin)
			_ = stdin.Close()
		}()
		go func() {
			wg.Add(1)
			_, _ = io.Copy(os.Stdout, stdout)
			_ = stdout.Close()
			wg.Done()
		}()
		go func() {
			wg.Add(1)
			_, _ = io.Copy(os.Stderr, stderr)
			_ = stderr.Close()
			wg.Done()
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			_ = process.Kill()
		}()

		err = process.Wait()
		waitWithTimeout(&wg, 1*time.Second)
		if err != nil {
			return err
		}

		exitCode, err := process.ExitCode()
		if err != nil {
			return err
		}
		os.Exit(exitCode)
	}

	return nil
}
