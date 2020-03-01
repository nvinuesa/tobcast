---
layout: post
title: Total order broadcast in go
date: 2019-11-12T16:19:00.000Z
categories: blog
---

In this article I will explore a very simple and naive implementation of a total order broadcast algorithm in 
go language. 

# Introduction

Total order (or atomic) broadcast algorithms are a family of well studied algorithms and important problem in 
distributed systems as well as fault tolerance. 

Concensus in distributed systems is a very hard problem to solve. Lots of research has been devoted to this 
series of problems and is still today an important subject of study. Although intuitively separate problems, 
ordering / causality of events and concensus are deeply linked. This is clearly explained in chapter 9 
of Martin Kleppmann's [Designing Data-Intensive Applications](https://www.oreilly.com/library/view/designing-data-intensive-applications/9781491903063/), where he describes the ideas behind total order broadcast and why order is necessary 
to achieve concensus.

On the other hand, ordering of events is itself a hard problem to solve in distributed systems, as Leslie Lamport has 
shown in his seminal paper "Time, Clocks, and the Ordering of Events in a Distributed System", where he studies 
the phylosophical links between time and ordering.

I chose to implement a variant of the algorithm caled "Communication History" that is described in the section 7.4 
of the paper "Total Order Broadcast and Multicast Algorithms: Taxonomy and Survey" by Xavier Défago et. al.

# Implementation

The algorithm is implemented as a command line application, as well as a golang library that can be added as a dependency, 
the full project can be found in https://github.com/underscorenico/tobcast.

To begin with, I followed a widelly used go project layout: https://github.com/golang-standards/project-layout. 

The `main()` function is present in `cmd/tobcast/main.go`, and it is the entry point for the binary when building it. 
It is also where the configuration is loaded, and the different components of the library are initialized.

### Configuration

The configuration is managed by [viper](https://github.com/spf13/viper). There are many other configuration 
libraries for go, but viper is well maintained and documented. 

The configuration can be passed as a yaml, json files or as environment variables. An example configuration 
is provided in the root folder of the project.

### Build

I chose to build the project using Makefile, and to follow Vincent Bernat's example Makefile (https://vincent.bernat.ch/en/blog/2019-makefile-build-golang). 

To build the project you can simply run:

```bash
make all
```

### Graceful shutdown

In order for our main process to shutdown gracefully, we need to be able to listen for events sent from the Os (signals like 
SIGINT or SIGTERM) and react properly to them.

There are several articles describing how to use `Context`for cancellation (https://eli.thegreenplace.net/2020/graceful-shutdown-of-a-tcp-server-in-go/, https://www.sohamkamani.com/blog/golang/2018-06-17-golang-using-context-cancellation/), but 
this article https://medium.com/@matryer/make-ctrl-c-cancel-the-context-context-bd006a8ad6ff shows an elegant way of trapping 
OS signals and propagating the cancellation using a context.

I will not be implementing the shutdown using the context package, instead I just listen for the signals and react 
with a function that will properly stop the different services (as we'll see later).

Simply add this to the main function:

```Go
signals := make(chan os.Signal, 1)
signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

defer func() {
    signal.Stop(signals)
}()
go func() {
    <-signals
    tobcast.Stop()
    os.Exit(1)
}()
```

We use a channel that will catch OS signals, because of `signal.Notify(c, os.Interrupt)`, and then we will read 
elements from this channel. 
The second goroutine will read messages from the signals channel and block if it is empty. When a message arrives 
(or when the channel is closed, when `main()` returns and the `defer` is executed), then the `Stop()` method 
of the `tobcast` object (the library's import) to properly close all ressources and then exit with code 1.

# Tobcast

The object `tobcast` is the entrypoint for the library when added as a dependency. You simply need to import 
it and you will have access to its builder and stop method.

A new instance of `tobcast` can be created by calling the `tobcast.New(&config)` method, and passing it the 
configuration.

The configuration object can be found in `pkg/config`, it contains the cluster's tcp ports, the listening 
tcp port of the instance, and a `keepAliveFreq` parameter that indicates the duration between two empty 
messages that are broadcasted to the cluster for liveness.

Once you have created the instance, it will start a tcp listener and a tcp writer that will broadcast 
the messages every time you call the `Multicast(message)` method. 

## Improvements

One important improvement would be to clean-up the `delivered` slice in `Tobcast`, this will avoid longer GC pauses and 
innecessary memory usage.
It must also be noted that in order for the tcp connections to be correctly opened, all instances in the cluster 
must be running. This can be improved by a listener and an auto-discovery system, which is out of the scope of this 
article.