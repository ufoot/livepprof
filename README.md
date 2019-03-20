Live pprof
==========

Live pprof is a Golang library to generate and use live profiles.

Status
------

Under heavy development, unstable, work in progress, use at your own risk.

[![Build Status](https://travis-ci.org/ufoot/livepprof.svg?branch=master)](https://travis-ci.org/ufoot/livepprof/branches)

Documentation
-------------

There are several ways to use the library, one is to directly call
the low-level API:

```go
    import "github.com/ufoot/livepprof/collector/cpu"

    // ...

    collector := cpu.New("mypackage", 10*time.Second)
    data, err := collector.Collect(nil)
    // data contains a profile, do whatever you want with it
    // in itself it's not very different from a raw Go profile,
    // the main difference is that name should be *resolved*.
```

This `"mypackage"` string is here to help you have data aggregated on
functions which belong to *your* code. This assumes your files are
in something that looks like `"github.com/me/mypackage/something"`.
Most of the time you're interested by profiling your code, if a dependency
is slow you'll still know it but aggregated from a point in code that
belongs to you.

Another way is to use a higher level profile interface which heartbeats
with profiles on a regular basis. It can then be graphed, logged,
I personally recommend using [Datadog](https://www.datadoghq.com/) to do this,
but you could technically use anything.

```go
    import (
        "log"
        "github.com/ufoot/livepprof"
    )

    // ...

    p, err := livepprof.New(livepprof.WithFilter("mypackage"))
    if err!=nil {
        // handle err
    }
    go func() { // This goroutine reports data
        for cpu := range p.CPU() {
            // This is going to be called every minute by default.
            for i, entry := range cpu.Entries {
                // Do whatever you want with entry.
                log.Printf("livepprof entry i=%d, function=%s, file=%s, stack=%s, value=%0.3f",
                    i,
                    entry.Key.Function,
                    entry.Key.File,
                    entry.Key.Stack,
                    entry.Value, // this is the actual CPU usage of that entry
                )
             }
        }
    }()

    // Your code that does things, here.

    p.Stop() // Stop the goroutine reporting data
```

Godoc links:

* [livepprof](https://godoc.org/github.com/ufoot/livepprof)
* [livepprof/objfile](https://godoc.org/github.com/ufoot/livepprof/objfile)
* [livepprof/collector](https://godoc.org/github.com/ufoot/livepprof/collector)
* [livepprof/collector/cpu](https://godoc.org/github.com/ufoot/livepprof/collector/cpu)
* [livepprof/collector/heap](https://godoc.org/github.com/ufoot/livepprof/collector/heap)

Bugs
----

Again, super experimental, among other things:

* it requires to have [GNU binutils](https://www.gnu.org/software/binutils/) installed, which is akward as Go as [builtin support](https://golang.org/pkg/debug/elf/) to analyze binaries.
* the heap profiles look wrong, generally speaking if you want outliers, you're going to get the right ones, but data is skewed, need to figure out where the problem is exactly

Authors
-------

* Christian Mauduit <ufoot@ufoot.org> : main developper, project
  maintainer.

License
-------

Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright
  notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright
  notice, this list of conditions and the following disclaimer in the
  documentation and/or other materials provided with the distribution.
* Neither the name of the copyright holder nor the
  names of its contributors may be used to endorse or promote products
  derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Live pprof homepage: https://github.com/ufoot/livepprof

Contact author: ufoot@ufoot.org
