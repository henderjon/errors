[![GoDoc](https://godoc.org/github.com/henderjon/errors?status.svg)](https://godoc.org/github.com/henderjon/errors)
[![License: BSD-3](https://img.shields.io/badge/license-BSD--3-blue.svg)](https://img.shields.io/badge/license-BSD--3-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/henderjon/errors)](https://goreportcard.com/report/github.com/henderjon/errors)
[![Build Status](https://travis-ci.org/henderjon/errors.svg?branch=dev)](https://travis-ci.org/henderjon/errors)
![tag](https://img.shields.io/github/tag/henderjon/errors.svg)
![release](https://img.shields.io/github/release/henderjon/errors.svg)

# errors

A drop in replacement for Go's stdlib errors package with support for previous errors.

This module is largely copied-n-pasted from [upspin.io/errors](https://godoc.org/upspin.io/errors). It removes of the various `Kind`s offered in upspin and adds `Location` to report the file and line of the invocation of `Here()`.

After grokking upspin.io/errors, the serialization format is simply `|---int64(len(v))---|---v---|` and repeat. This requires the sender to include all values (even zeroes) and the receiving end to know the order. To add labels/types would be trivial but seems unnecessary as the library also includes the `Unserialize` func.

To the extent that the errors should be logged or stored in such a way that `Unserialize` isn't used, `Error()` allows for the string representation of the error stack (default format is `@ location; message\n\t`). The "\n\t" can be changed via `Sep`. There is always `json.Marshal()` which will do exactly as expected. I've included `Encode()` which returns a delimiter separated value (DSV) formated string. Think CSV but using unit (31) and field (30) separators ([more info](https://www.lammertbies.nl/comm/info/ascii-characters.html)). It's esoteric, but that is the purpose for which they were created. Wanna read it? ` ... | tr "\036\037" "\n,"`. Again, this most likely won't be used much, but it was fun.

This library allows for arbitrary `Kind`s. Where upspin defines the error `Kind`s used within, this library is a drop-in for the stdlib's errors. Thus, the importing context needs to define the `Kind`s and the string/descr (if any). This lib simply allows for the storage/transfer of the uint8(Kind). Further, upspin serialized `Kind` as a scalar value with no length (it was only ever always encoded as 1 int64). For now, this library chooses to add a length param simply to keep the serialized format consistent. Admittedly, the extra bytes and flops used aren't strictly necessary.

