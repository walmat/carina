package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"github.com/radovskyb/watcher"
)

var configFile = flag.String("config", "config/watcher/bot.toml", "path to config file")

type Config struct {
	Package      string   `toml:"package"`
	BuildArgs    []string `toml:"build_args"`
	Filter       string   `toml:"filter"`
	MaxEvents    int      `toml:"max_events"`
	DebounceTime int      `toml:"debounce_time"`
	PrintWatched bool     `toml:"print_watched"`
	Directories  []string `toml:"directories"`
}

func runCommand(ctx context.Context, config Config) {
	outputFile := "nebula"
	if runtime.GOOS == "windows" {
		outputFile = outputFile + ".exe"
	}
	outPath := path.Join(os.TempDir(), uuid.NewString(), outputFile)

	config.BuildArgs = append([]string{"build", "-o", outPath}, config.BuildArgs...)
	config.BuildArgs = append(config.BuildArgs, config.Package)

	cmd := exec.CommandContext(ctx, "go", config.BuildArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("error starting build command:", err)
		return
	}

	cmd = exec.CommandContext(ctx, outPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	if err := os.RemoveAll(outPath); err != nil {
		log.Println("error cleaning up previous build directories:", err)
	}

	// check if application was closed by user, if it was stop monitoring
	select {
	case <-ctx.Done():
		break
	default:
		os.Exit(0)
	}
}

func main() {
	flag.Parse()
	configBytes, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalln("error reading watcher config:", err)
	}

	var config Config
	if err = toml.Unmarshal(configBytes, &config); err != nil {
		log.Fatalln("error parsing watcher config:", err)
	}

	w := watcher.New()
	w.SetMaxEvents(config.MaxEvents)
	w.FilterOps(watcher.Write, watcher.Move, watcher.Create, watcher.Remove, watcher.Rename)

	r := regexp.MustCompile(config.Filter)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		lastEventTime := time.Now()
		ctx, cancelFunc := context.WithCancel(context.Background())
		go runCommand(ctx, config)
		for {
			select {
			case <-w.Event:
				if time.Now().Sub(lastEventTime) < time.Duration(config.DebounceTime) {
					lastEventTime = time.Now()
					continue
				}
				fmt.Println()
				fmt.Println("Changes detected, rebuilding...")
				fmt.Println()
				cancelFunc()
				ctx, cancelFunc = context.WithCancel(context.Background())
				go runCommand(ctx, config)
				lastEventTime = time.Now()
			case err = <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				cancelFunc()
				return
			}
		}
	}()

	for _, dir := range config.Directories {
		if err := w.AddRecursive(dir); err != nil {
			log.Fatalln(err)
		}
	}

	if config.PrintWatched {
		fmt.Println("Watched Files:")
		for filePath, f := range w.WatchedFiles() {
			fmt.Printf("%s: %s\n", filePath, f.Name())
		}
		fmt.Println()
	}

	if err = w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
