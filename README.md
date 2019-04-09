# errors

A drop in replacement for Go's stdlib errors package with support for previous errors.

This is a pretty simple way of errors stacking. It also provides a way of printing/marshaling the stack.

It would be relatively straightforward to add a Kind or Op field like [Error handling in Upspin](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html) and [upspin.io/errors](https://godoc.org/upspin.io/errors) as Upspin was the motivation.
