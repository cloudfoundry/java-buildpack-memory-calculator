# Java Buildpack Memory Calculator

The Java Buildpack Memory Calculator calculates a holistic JVM memory configuration with the goal of ensuring that applications perform well while not exceeding a container's memory limit and being recycled.

In order to perform this calculation, the Memory Calculator requires the following input:
* `--total-memory`: total memory available to the application, typically expressed with size classification (`B`, `K`, `M`, `G`, `T`)
* `--loaded-class-count`: the number of classes that will be loaded when the application is running
* `--thread-count`: the number of user threads
* `--jvm-options`: JVM Options, typically `JAVA_OPTS`
* `--head-room`: percentage of total memory available which will be left unallocated to cover JVM overhead

The Memory Calculator prints the calculated JVM configuration flags (_excluding_ any that the user has specified in `--jvm-options`).  If a valid configuration cannot be calculated (e.g. more memory must be allocated than is available), an error is printed and a non-zero exit code is returned.  In order to **override** a calculated value, users should pass any of the standard JVM configuration flags into `--jvm-options`.  The calculation will take these as fixed values and adjust the non-fixed values accordingly.

## Install  

```sh
$ go install github.com/cloudfoundry/java-buildpack-memory-calculator@latest
```

## Algorithm

The following algorithm is used to generate the holistic JVM memory configuration:

1. Headroom is calculated as `total memory * (head room / 100)`
1. If `-XX:MaxDirectMemorySize` is configured it is used for the amount of direct memory.  If not configured `10M` (in the absence of any reasonable heuristic) is used.
1. If `-XX:MaxMetaspaceSize` is configured it is used for the amount of metaspace.  If not configured then the value is calculated as `(5800B * loaded class count) + 14000000b`.
1. If `-XX:ReservedCodeCacheSize` is configured it is used for the amount of reserved code cache.  If not configured `240M` (the JVM default) is used.
1. If `-Xss` is configured it is used for the size of each thread stack.  If not configured `1M` (the JVM default) is used.
1. If `-Xmx` is configured it is used for the size of the heap.  If not configured then the value is calculated as `total memory - (headroom + direct memory + metaspace + reserved code cache + (thread stack * thread count))`.

Broadly, this means that for a constant application (same number of classes), the non-heap overhead is a fixed value.  Any changes to the total memory will be directly reflected in the size of the heap.  Adjustments to the non-heap memory configuration (e.g. stack size, reserved code cache) _can_ result in larger heap sizes but can also have negative runtime side-effects that must be taken into account.

### Compressed class space size

According to the [HotSpot GC Tuning Guide][h]:

> The MaxMetaspaceSize applies to the sum of the committed compressed class space and the space for the other class metadata.

Therefore the memory calculator does not set the compressed class space size (`-XX:CompressedClassSpaceSize`) since the memory for the compressed class space is bounded by the maximum metaspace size (`-XX:MaxMetaspaceSize`).

[h]: https://docs.oracle.com/javase/8/docs/technotes/guides/vm/gctuning/considerations.html

## License
The Java Buildpack Memory Calculator is Open Source software released under the [Apache 2.0 license][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0.html
