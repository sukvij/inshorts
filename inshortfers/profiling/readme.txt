pprof --> collect run time performance data --> cpu, memory, mutex
    cpu profiling -- > which function is consuming most cpu times.
    memory profiling --> Detects memory leaks or excessive allocations
    go routine --> find blocked go routines


CPU: curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.pprof
Memory: curl http://localhost:8080/debug/pprof/heap > heap.pprof
Goroutines: curl http://localhost:8080/debug/pprof/goroutine > goroutine.pprof
Analyze with: go tool pprof cpu.pprof (or heap.pprof, etc.).

Use go tool pprof -http=:8081 cpu.pprof for a web-based UI.


hey command install
go install github.com/rakyll/hey@latest


hey -n 1000 -c 50 http://localhost:8080/delivery?app=com.duolingo.ludokinggame&country=US&os=Android