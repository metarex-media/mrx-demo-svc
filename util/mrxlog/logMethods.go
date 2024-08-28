// Package mrxlog logs the mrx path chain through the register.
package mrxlog

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// ColourConsole wraps github.com/lmittmann/tint to log the output
// straight to the console with colour.
func ColourConsole(opts *slog.HandlerOptions) {

	// set global logger with custom options
	// assign colours based on Operating system
	colourStart(opts, false)

}

// ColourConsole wraps github.com/lmittmann/tint to log the output
// straight to the console *without* colour.
func Console(opts *slog.HandlerOptions) {
	colourStart(opts, true)
}

func colourStart(opts *slog.HandlerOptions, noColour bool) {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:       opts.Level,
			TimeFormat:  time.RFC3339,
			ReplaceAttr: opts.ReplaceAttr,
			AddSource:   opts.AddSource,
			NoColor:     noColour,
		}),
	))
}
