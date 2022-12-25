```mermaid
graph TD
    A["/"]
    A --> B[Valid request]
    A --> C["Request malformed (400)"]
    A --> D["Invalid API key (400)"]
    A --> E["Temporary issue (500)"]

    style A fill:blue
    style B fill:green

    classDef yellow fill:yellow,color:black
    class C,D yellow

    style E fill:red,color:black
```
