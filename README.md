# golang-assignment
This assignment is based on a viral 
[1brc problem](https://github.com/gunnarmorling/1brc/tree/main)

But for smaller machines like one of ours, i am using 100M lines of data

# Problem statement
So, the problem statement effectively is

- we have a txt file generated using a script which has 100M or 0.1B lines in which each line has a station and temperature data

- Here is a sample

```txt
Rankweil;-74.5
Vārānasi;29.6
San Felipe Orizatlán;0.3
Kudowa-Zdrój;-69.0
Chicago Ridge;-93.0
Asfarvarīn;13.6
Schwyz;19.2
Natal;81.5
Amānganj;49.3
Cedar Park;59.6
```

- So we have to calculate min, mean, max temperatures for each station, given that maximum number of stations generated is 10,000

# Generate weather data

- For generating txt file of 100M lines, run the following command (or I will upload this `1.5GB file measurements-100M.txt` in another branch `main_with_data`)


```sh
python3 create-measurements.py 100_000_000
```
- Note: create-measurements.py, weather_stations.csv are from 1brc repository
- Note: Can use 10M (10_000_000) for faster program, measurements.txt uploaded in this branch has 1M records
# Run weather calculation
- Once txt file is generated, run main.go
```sh
go run main.go
```

# Plan
- We will read txt file in chunks and write all lines to a channel which will be shared among a worker pool in a `fan-out` design after which all partial results coming from each worker will be collected by main goroutine in a `fan-in` design

# Results
- This program will approximately take about 20 seconds using 3 cores cpu to execute for 100M records