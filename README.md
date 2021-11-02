# GOHTML
An encoder for a Go struct to HTML

Using the "reflect" package and recursion this package is able to convert a complex go struct into HTML

#### Features
 - [X] Struct
 - [X] Embedded struct
 - [X] Slice
 - [X] Table encode option
 - [X] Omit Empty
 - [ ] Map (Needs implementation and testing)
 - [ ] Custom ID (Unsure on whether this is needed)

#### Examples

##### Simple Encode
```go
package main

import (
    "fmt"
    "github.com/DrGrimshaw/gohtml"
)

type Example struct {
	Name string `html:"l=Name,e=span,c=name"`
}

func main() {
    fmt.Print(gohtml.Encode(Example{Name: "John Doe"}))
}
```

##### Output
```html
<div><span>Name</span><span class='name'>John</span></div>
```
##### Table Format
```go
package main

import (
    "fmt"
    "github.com/DrGrimshaw/gohtml"
)

type Example struct {
	Name string `html:"l=Name,e=span,c=name"`
	Age string `html:"l=Age,e=span,c=age"`
	Location string `html:"l=Location,e=span,c=location"`
}

type ExampleArr struct {
    Examples []Example `html:"row"`
}

func main() {
    tableExample := ExampleArr{Examples: []Example{{Name: "John Doe",Age: "30",Location: "Washington DC"}}}

    fmt.Print(gohtml.Encode(tableExample))
}
```

##### Output
```html
<table>
    <thead>
        <tr>
            <td>Name</td>
            <td>Age</td>
            <td>Location</td>
        </tr>
    </thead>
    <tr>
        <td><span class='name'>John Doe</span></td>
        <td><span class='age'>30</span></td>
        <td><span class='location'>Washington DC</span></td>
    </tr>
</table>
```

![Gopher mascot](mascot.png)
