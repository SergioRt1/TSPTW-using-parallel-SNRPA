# TSPTW using Parallel Stable NRPA

## Description

Implementation of a parallel Stable NRPA to solve the TSPTW as describe by Tristan Cazenave
in [Stabilized Nested Rollout Policy Adaptation](https://arxiv.org/pdf/2101.03563.pdf).

The implementation uses an actors model to handle the concurrent computation of the NRPA leaves (playout) and the
concurrent runs of different trees.

## Prerequisites

* [Golang](https://golang.org) 1.14+

You need to download the repository and place it in the src folder inside your GOPATH folder.

## Build

```
go build main.go
```

## Run

```
./main -cpuprof=nrpa.prof -runs=5 -all (args...)
```

### Program arguments

| Command | Type   | Description                                                | Default        |
|---------|--------|------------------------------------------------------------|----------------|
| cpuprof | string | Filename of the generated report                           |                |
| file    | string | File name of a test case                                   | rc_201.1.txt   |
| time    | int    | Execution timeout in sec                                   | 10000          |
| factor  | int    | Tree stabilization factor                                  | 10             |
| iter    | int    | Iterations in the next level (number of children per node) | 10             |
| levels  | int    | Tree levels                                                | 5              |
| runs    | int    | Number of trees to run                                     | 4              |
| all     | bool   | read all case files                                        | false          |
| nactors | int    | number of actors that computes the NRPA tree               | same as runs   |
| pactors | int    | number of actors that computes the leaves.                 | same as factor |

## Author

* **[Sergio Rodríguez](https://github.com/SergioRt1)**

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](LICENSE) file for details