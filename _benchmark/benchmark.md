# Benchmarks

To run the benchmark suite, navigate to the same folder and execute the main.go file.
This suite will test the performance of hypermatch with plain Go objects, hypermatch with JSON objects, and [quamina](https://github.com/timbray/quamina) against each other.
Simply run the command go run main.go to start the benchmarks and compare the results.

Results as of Aug, 29th, 2024 with Go 1.23.0 on MacBook Pro M1 Max, 32GB RAM with 100,000 rules:

```
---Starting with hypermatch
adding 100000 rules took 0.51862s
processed 51753 events with 517530 matches in 1.00001s -> 51752.42425 evt/s
processed 103199 events with 1031990 matches in 2.00004s -> 51598.36053 evt/s
processed 156249 events with 1562490 matches in 3.00009s -> 52081.45563 evt/s
processed 209523 events with 2095230 matches in 4.00013s -> 52379.05586 evt/s
processed 262257 events with 2622570 matches in 5.00016s -> 52449.76793 evt/s

---Starting with hypermatch-json
adding 100000 rules took 1.95150s
processed 39732 events with 397320 matches in 1.00000s -> 39731.90564 evt/s
processed 80432 events with 804320 matches in 2.00003s -> 40215.31718 evt/s
processed 121915 events with 1219150 matches in 3.00006s -> 40637.55840 evt/s
processed 164093 events with 1640930 matches in 4.00009s -> 41022.33768 evt/s
processed 206235 events with 2062350 matches in 5.00013s -> 41245.96473 evt/s

---Starting with quamina
adding 100000 rules took 4.54697s
processed 4 events with 0 matches in 1.24154s -> 3.22181 evt/s
processed 7 events with 0 matches in 2.31263s -> 3.02685 evt/s
processed 11 events with 0 matches in 3.50818s -> 3.13552 evt/s
processed 15 events with 0 matches in 4.70148s -> 3.19049 evt/s
```