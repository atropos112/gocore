# atrogolib

This is an "all-in-one" library package in Go that contains all sorts of utilities, helpers, structs that can be used in variety of places. I am using it extensively in my homelab stuff.

Should certain part of atrogolib get sizable enough, it will likely be split into its own package. The aim is to let things grow organically and see what needs splitting rather than split pre-emptively stupid number of repos unnecessarily.

# Dependencies

The hope is to use as few dependencies as possible. In particular the following are used now:

### Http web framework of choice: [net/http](https://pkg.go.dev/net/http)

The net/http is fast and powerful, its not as fast as fiber I don't think but fiber comes with bunch of compromises of its own, in particular if you need something so fast that net/http is not enough for you, you should probably roll something specific to that task anyway. net/http is a standard library.

### Logging framework of choice: [slog](https://pkg.go.dev/log/slog)

Slog is fantastic logging package that I think is either the fastest or nearly the fastest, its super easy to use, and has a lot of support. Is the be all and end all, maybe not, but I'd rather settle on something rather than have 3 different logging frameworks lying around. Slog is also coming by default as a standard package now which is very nice.

### Utilities package: [lo](https://github.com/samber/lo)

It comes with crazy number of utilities that I find myself using all the time. It's a bit of a kitchen sink but I like it. If writing a utility function check if its not in lo already. lo/parallel is also supported. This is not a standard library.
