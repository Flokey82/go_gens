# genthesaurus

A simple thesaurus generator for Go using the ea-thesaurus json found here:

https://github.com/dariusk/ea-thesaurus

## TODO

- [ ] Find words by tag
- [ ] Traverse associations (and filter them by tag)
- [ ] Improve UX

## Usage

```go
package main

import (
    "fmt"

    "github.com/Flokey82/go_gens/genthesaurus"
)

func main() {
    // Load thesaurus
    thesaurus := genthesaurus.New()
    
    // Add a new word association.
    thesaurus.AddAssociation("good", "evil", 1, "antonym")

    // Add some tags.
    thesaurus.Add("good", "positive")
    thesaurus.Add("evil", "negative")

    thesaurus.Log()
}
```
