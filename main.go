package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/env25/mpdlrc/internal"
	"github.com/env25/mpdlrc/internal/config"

	"github.com/spf13/pflag"
)

var (
	exitCode = 0
	usage    = false
	cfg      = config.DefaultConfig()
)

func init() {
	log.SetFlags(0)
	pflag.StringVar(&cfg.MusicDir, `musicdir`, cfg.MusicDir, `override MusicDir`)
	pflag.StringVar(&cfg.LyricsDir, `lyricsdir`, cfg.LyricsDir, `override LyricsDir`)
	pflag.StringVar(&cfg.MPD.Connection, `mpd-connection`, cfg.MPD.Connection, `override MPD.Connection (possible "unix", "tcp")`)
	pflag.StringVar(&cfg.MPD.Address, `mpd-address`, cfg.MPD.Address, `override MPD.Address (use unix socket path or "host:port")`)
	pflag.StringVar(&cfg.MPD.Password, `mpd-password`, cfg.MPD.Password, `override MPD.Password`)
	pflag.BoolVar(&cfg.Debug, `debug`, cfg.Debug, `enable debug`)
	pflag.BoolVarP(&usage, `help`, `h`, usage, `show this help message`)
}

func exit() { os.Exit(exitCode) }

func main() {
	defer exit()

	for _, fpath := range config.ConfigFiles {
		err := cfg.MergeTOMLFile(fpath)
		if err != nil {
			switch err.(type) {
			case *os.PathError:
				// no-op
			default:
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}

	pflag.Parse()

	if usage {
		pflag.Usage()
		return
	}

	if cfg.Debug {
		logBuilder := new(strings.Builder)
		log.SetOutput(logBuilder)
		defer fmt.Fprint(os.Stderr, logBuilder)
	} else {
		log.SetOutput(io.Discard)
	}

	cfg.Expand()

	if err := cfg.Assert(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}

	if err := internal.NewApplication(cfg).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
}
