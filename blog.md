

## Graceful shutdown

In order for our main process to shutdown gracefully, we need to be able to pass the information that it has been killed to the child 
goroutines. In Go, we have a very powerful tool: `Context`.

There are several articles describing how to use `Context`for cancellation (https://eli.thegreenplace.net/2020/graceful-shutdown-of-a-tcp-server-in-go/, https://www.sohamkamani.com/blog/golang/2018-06-17-golang-using-context-cancellation/), but 
this article https://medium.com/@matryer/make-ctrl-c-cancel-the-context-context-bd006a8ad6ff shows an elegant way of trapping 
OS signals and propagating the cancellation using a context.

We start by adding this code of block in our main function:

```Go
ctx, cancel := context.WithCancel(context.Background())
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt)

defer func() {
    signal.Stop(c)
    cancel()
}()

go func() {
    select {
    case <-c:
        cancel()
    case <-ctx.Done():
    }
}()
```

We use a channel that will catch OS signals, because of `signal.Notify(c, os.Interrupt)`, and then we will read 
elements from this channel. When we receive an OS signal, the `cancel()` function of the context will be called.
If the `main()` method returns (for example because we have a way of stopping our process other than OS signals), 
it will also call the `cancel()` function and close the channel.



## Improvements

One important improvement would be to clean-up the `delivered` slice in `Tobcast`, this will avoid longer GC pauses and 
innecessary memory usage.
It must also be noted that in order for the tcp connections to be correctly opened, all instances in the cluster 
must be running. This can be improved by a listener and an auto-discovery system, which is out of the scope of this 
article.