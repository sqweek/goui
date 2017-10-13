# goui
Goui aims to experiment with concurrent graphics/user interface designs in go.
It is in very early stages and nothing in the API should be considered stable!

The initial proof of concept in [examples/worms](example/worms) consists of a
whole bunch of goroutines (2¹³ by default) all scribbling on the same raster,
which turns out to be a pretty neat visualisation. :)
