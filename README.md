# GOHTML
An encoder for a Go struct to HTML

Using the "reflect" package and recursion this package is able to convert a complex go struct into HTML

#### Features
 - [X] Struct
 - [X] Embedded struct
 - [X] Slice
 - [ ] Map (Needs implementation and testing)
 - [ ] Custom ID (Unsure on whether this is needed)
 - [ ] Table encode option (Nice to have)

#### Example 

##### Code
```go
package main

import (
    "fmt"
    "github.com/DrGrimshaw/gohtml"
)

type Example struct {
	Name string `label:"Name" element:"span" class:"name"`
}

func main() {
    fmt.Print(gohtml.Encode(Example{Name: "John Doe"}))
}
```

##### Output
```html
<div><span>Name</span><span class='name'>John</span></div>
```

![Gopher mascot](mascot.png)
