# Instana Memory Calculator

**Note:** This fork adapts the excellent [`Java Buildpack Memory Calculator`](https://github.com/cloudfoundry/java-buildpack-memory-calculator) to the somewhat unusual requirements of the Instana Agent in terms of direct memory usage, as well as building and shipping the memory calculator on architectures not supported by Cloud Foundry.

The Instana Memory Calculator calculates a holistic JVM memory configuration with the goal of ensuring that the Instana Agent performs well while not exceeding a container's (or cgroup's) memory limit and being recycled.

In order to perform this calculation, the Memory Calculator requires the following input:

* `--total-memory`: total memory available to the application, typically expressed with size classification (`B`, `K`, `M`, `G`, `T`)
* `--loaded-class-count`: the number of classes that will be loaded when the application is running
* `--thread-count`: the number of user threads
* `--jvm-options`: JVM Options, typically `JAVA_OPTS`
* `--head-room`: percentage of total memory available which will be left unallocated to cover JVM overhead
* `--direct-memory-to-heap-ratio`: how to split the available memory between heap and direct memory, when neither is hard-coded to a particular value; default is `.1`, that is, 10% of available allocated to direct memory
* `--heap-young-generation-ratio`: how much of the heap to reserve for young generation; default is `.3`, that is, 30% of the heap reserved for the young generation

The Memory Calculator prints the calculated JVM configuration flags (_excluding_ any that the user has specified in `--jvm-options`).  If a valid configuration cannot be calculated (e.g. more memory must be allocated than is available), an error is printed and a non-zero exit code is returned.  In order to **override** a calculated value, users should pass any of the standard JVM configuration flags into `--jvm-options`.  The calculation will take these as fixed values and adjust the non-fixed values accordingly.

## Install  

```sh
go get -v github.com/instana/java-buildpack-memory-calculator
```

## Algorithm

The following algorithm is used to generate the holistic JVM memory configuration:

1. Headroom is calculated as `total memory * (head room / 100)`
1. If `-XX:MaxMetaspaceSize` is configured it is used for the amount of metaspace.  If not configured then the value is calculated as `(5800B * loaded class count) + 14000000b`.
1. If `-XX:ReservedCodeCacheSize` is configured it is used for the amount of reserved code cache.  If not configured `240M` (the JVM default) is used.
1. If `-Xss` is configured it is used for the size of each thread stack.  If not configured `1M` (the JVM default) is used.
1. If `-XX:MaxDirectMemorySize` is configured it is used for the amount of direct memory.  If not configured, and "direct-memory-to-heap" ratio is not configured, `10M` (in the absence of any reasonable heuristic) is used; if "direct-memory-to-heap" ratio is configured, the available space after all the previous calculations will be divided between direct memory and heap according to the specified ratio.
1. If `-Xmx` is configured it is used for the size of the heap.  If not configured and no "direct-memory-to-heap" ratio is configured, then the value is calculated as `total memory - (headroom + direct memory + metaspace + reserved code cache + (thread stack * thread count))`.

Broadly, this means that for a constant application (same number of classes), the overhead outside of heap and direct memory is a fixed value.  Any changes to the total memory will be directly reflected in the size of the heap amd direct memory.  Adjustments to memory configurations (e.g. stack size, reserved code cache) _can_ result in larger heap sizes but can also have negative runtime side-effects that must be taken into account.

The document [Java Buildpack Memory Calculator v3][v3] provides some rationale for the memory calculator externals.

[v3]: https://docs.google.com/document/d/1vlXBiwRIjwiVcbvUGYMrxx2Aw1RVAtxq3iuZ3UK2vXA/edit?usp=sharing

### Compressed class space size

According to the [HotSpot GC Tuning Guide][h]:

> The MaxMetaspaceSize applies to the sum of the committed compressed class space and the space for the other class metadata.

Therefore the memory calculator does not set the compressed class space size (`-XX:CompressedClassSpaceSize`) since the memory for the compressed class space is bounded by the maximum metaspace size (`-XX:MaxMetaspaceSize`).

[h]: https://docs.oracle.com/javase/8/docs/technotes/guides/vm/gctuning/considerations.html

## License

The Instana Memory Calculator is Open Source software released under the [Apache 2.0 license][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0.html
