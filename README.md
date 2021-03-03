# TSPTW using Parallel Stable NRPA


go build main.go
./main -cpuprof=nrpa.prof -file=rc_206.2.txt -runs=5

### Program arguments

| Command | Type   | Description                                                | Default      |
|---------|--------|------------------------------------------------------------|--------------|
| cpuprof | string | Filename of the generated report                           | N/A          |
| file    | string | File name of a test case                                   | rc_201.1.txt |
| time    | int    | Execution timeout in sec                                   | 10000        |
| factor  | int    | Tree stabilization factor                                  | 10           |
| iter    | int    | Iterations in the next level (number of children per node) | 10           |
| levels  | int    | Tree levels                                                | 5            |
| runs    | int    | Number of trees to run                                     | 4            |

