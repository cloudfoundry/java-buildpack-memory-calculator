# Java Buildpack Memory Calculator Help
```

java-buildpack-memory-calculator

Calculates JVM memory switches based on the total memory available, the number of classes the application will load, the number of threads that will be used, and any JVM options provided as input.

The output consists of any calculated memory switches.

If a calculated memory switch value is unsuitable, it can be set in the JVM options provided as input and will no longer be calculated.

Example invocation from a shell:
$ java-buildpack-memory-calculator -loadedClasses=1000 -stackThreads=10 -totMemory=1g -poolType=metaspace -vmOptions=-XX:MaxDirectMemorySize=100M\ -verbose:gc

Usage of java-buildpack-memory-calculator:
  -help
    	prints description and flag help
  -loadedClasses int
    	an estimate of the number of classes that will be loaded when the application is running
  -poolType string
    	the type of JVM pool used in the calculation. Set this to 'permgen' for Java 7 and to 'metaspace' for Java 8 and later.
  -stackThreads int
    	number of threads that will be used
  -totMemory string
    	total memory available to allocate, expressed as an integral number of bytes (B), kilobytes (K), megabytes (M) or gigabytes (G), e.g. '1G'
  -vmOptions string
    	Java VM options, typically the JAVA_OPTS specified by the user

```
