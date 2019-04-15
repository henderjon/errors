# errors

A drop in replacement for Go's stdlib errors package with support for previous errors.

This module is largely copied-n-pasted from [upspin.io/errors](https://godoc.org/upspin.io/errors). It removes of the various `Kind`s offered in upspin and adds `Location` to report the file and line of the invocation of `Here()`.

After grokking upspin.io/errors, the serialization format is simply `|---int64(len(v))---|---v---|` and repeat. This requires the sender to include all values (even zeroes) and receiving end know the order. To add labels/types would be trivial but seems unnecessary as the library also includes the `Unserialize` func.

To the extent that the errors should be logged or stored in such a way that `Unserialize` isn't used, Error() allows for the string representation of the error stack (default format is `message[[[ = kind] @ location]\n\t]`). The "\n\t" can be changed via `Sep`. There is always `json.Marshal()` which will do exactly as expected. I've included `Encode()` which returns a delimiter separated value (DSV) formated string. Think CSV but using unit and field separators. It's esoteric, but that is the purpose for which they were created. Wanna read it? ` ... | tr "\036\037" "\n,"`. Again, this most likely won't be used much, but it was fun.

This library allows for arbitrary `Kind`s. Where upspin defines the error `Kind`s used within, this library is a drop-in for the stdlib's errors. Thus, the importing context needs to define the `Kind`s and the string/descr (if any). This lib simply allows for the storage/transfer of the uint8(Kind).

