# counter
A GoLang implementation of counters.

It provides an object: `counters.CounterBox` and `Counter`, `Min` and `Max` implementation.

The example usages:

    import "github.com/orian/counters"
    
    //... some code
    
    cb := counters.NewCounterBox()
    c := cb.GetCounter("ex")
    c.Increment()
    c.IncrementBy(2)
    c.Value()  // returns 7
    cb.GetCounter("ex").Value()  // Returns 7

For convenience there is also an `http.HandleFunc` provided which prints values of counters.

One may use subpackage `globals` if want to use a global counters.

    import "github.com/orian/counters/global"
    
    //... some code
    
    c := global.GetCounter("ex")
    c.Increment()
    c.IncrementBy(2)
    c.Value()  // returns 7
    global.GetCounter("ex").Value()  // Returns 7

The library and all objects are thread safe.
