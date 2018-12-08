Import of internal packages
===========================

Found no other way to use those than copy/pasting code from
[Google pprof repo](https://github.com/google/pprof/).
The version used here is
[3ea8567](https://github.com/google/pprof/tree/3ea8567a2e5752420fc544d2e2ad249721768934)
from Dec 6th 2018.

I had to replace `https://github.com/google/pprof/internal`
by `https://github.com/ufoot/livepprof/internal/google` and that's it.

Any other options is welcome, I hate copy/pasting code,
this is going to go out-of-sync, OTOH did not feel like
rewriting all this, and found no other public, importable binary parser,
especially written in native Go.
