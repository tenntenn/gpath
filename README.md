# gpath

```
type Foo struct {
    N []int
}

f := &Foo{N: []int{100}}

v, err := gpath.At(f, `N[0]`)
if err != nil {
    // handle error
}
fmt.Println(v) // 100
```
