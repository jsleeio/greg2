# greg2

## what is this?

`greg2` is Yet Another `grep`. This one uses [Go
RE2](https://github.com/google/re2/wiki/Syntax) and attempts to behave
similarly to [GNU `grep`](https://www.gnu.org/software/grep/) where
equivalent functionality is supported.

## why another grep?

I needed to implement a moving-window algorithm for something else and
it was easier to test it in a smaller app first. Also a good mental
exercise!

## options

```
Usage of ./greg2:
  -after int
    	lines of following context to print for each match
  -before int
    	lines of preceding context to print for each match
  -filenames string
    	show filenames in output (valid options: no, auto, yes) (default "auto")
  -ignorecase
    	perform case-insensitive matching
  -match string
    	RE2 regular expression to match against the input files
  -quiet
    	do not output any matches
```

## ideas

### context selection by regular expression

To do this, implement `-before-from-re` and `-after-to-re`, which is
something I've always wished `grep` had. It can also be done with `sed`:

```
sed -n '/BEFOREMATCH/,/AFTERMATCH/p' inputfile
```

### parallel searching of files

This _is_ Go, after all.

## license

Copyright 2018 John Slee.  Released under the terms of the MIT license
[as included in this repository](LICENSE).
