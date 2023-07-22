package main

import (
	"fmt"
	"go.starlark.net/starlark"
	"image"
	"tidbyt.dev/pixlet/encode"
	"tidbyt.dev/pixlet/globals"
	"tidbyt.dev/pixlet/runtime"
)

type RenderOpts struct {
	Cache         runtime.Cache
	Magnify       int
	RenderGif     bool
	MaxDuration   int
	SilenceOutput bool
	Width         int
	Height        int
}

func render(script []byte, opts RenderOpts, vars map[string]string) ([]byte, error) {
	globals.Width = opts.Width
	globals.Height = opts.Height

	// Remove the print function from the starlark thread if the silent flag is passed.
	initializers := []runtime.ThreadInitializer{}
	if opts.SilenceOutput {
		initializers = append(initializers, func(thread *starlark.Thread) *starlark.Thread {
			thread.Print = func(thread *starlark.Thread, msg string) {}
			return thread
		})
	}

	// XXX re-use cache
	runtime.InitHTTP(opts.Cache)
	runtime.InitCache(opts.Cache)

	applet := runtime.Applet{}
	err := applet.LoadWithInitializers("memory file", script, nil, initializers...)
	if err != nil {
		return nil, fmt.Errorf("failed to load applet: %w", err)
	}

	roots, err := applet.Run(vars, initializers...)
	if err != nil {
		return nil, fmt.Errorf("error running script: %w", err)
	}
	screens := encode.ScreensFromRoots(roots)

	filter := func(input image.Image) (image.Image, error) {
		if opts.Magnify <= 1 {
			return input, nil
		}
		in, ok := input.(*image.RGBA)
		if !ok {
			return nil, fmt.Errorf("image not RGBA, very weird")
		}

		out := image.NewRGBA(
			image.Rect(
				0, 0,
				in.Bounds().Dx()*opts.Magnify,
				in.Bounds().Dy()*opts.Magnify),
		)
		for x := 0; x < in.Bounds().Dx(); x++ {
			for y := 0; y < in.Bounds().Dy(); y++ {
				for xx := 0; xx < opts.Magnify; xx++ {
					for yy := 0; yy < opts.Magnify; yy++ {
						out.SetRGBA(
							x*opts.Magnify+xx,
							y*opts.Magnify+yy,
							in.RGBAAt(x, y),
						)
					}
				}
			}
		}

		return out, nil
	}

	var buf []byte

	if screens.ShowFullAnimation {
		opts.MaxDuration = 0
	}

	if opts.RenderGif {
		buf, err = screens.EncodeGIF(opts.MaxDuration, filter)
	} else {
		buf, err = screens.EncodeWebP(opts.MaxDuration, filter)
	}
	if err != nil {
		return nil, fmt.Errorf("error rendering: %w", err)
	}

	return buf, nil
}
