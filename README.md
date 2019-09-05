# fibonacci-sequence
implement fibonacci sequence API using Go

## QuickStart

	go run server.go
	
## API

    prev -> error
    current -> 0
    next -> 1
    next -> 1
    next -> 2
    (server crush)
    current -> 2
    previous -> 1

## Default
data.log

    0
    1
    (empty line at last)

## Expected response

We use a hashtable for "code" and "msg" to tell frontend what's going on

    codeToMessage := map[string]string{
            "00000": "Succeed",
            "10000": "There is no previous value for 0",
            "10001": "data.log not found",
            "10002": "Can't read the first line from data.log",
            "10003": "Can't read the second line from data.log",
            "10004": "Can't update data.log",
        }

The expected response looks like this. The "data" is the fibonacci value we need

    {
        "code":"00000",
        "data":"1",
        "msg":"Succeed"
    }

## Structure

![structure](https://raw.githubusercontent.com/Windsooon/fibonacci-sequence/master/structure.png)

## Benchmark

### ab -n 5000 -c 100 -k http://127.0.0.1:8080/current/

    This is ApacheBench, Version 2.3 <$Revision: 1826891 $>
    Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
    Licensed to The Apache Software Foundation, http://www.apache.org/

    Server Software:        
    Server Hostname:        127.0.0.1
    Server Port:            8080

    Document Path:          /current/
    Document Length:        44 bytes

    Concurrency Level:      100
    Time taken for tests:   0.560 seconds
    Complete requests:      5000
    Failed requests:        0
    Keep-Alive requests:    5000
    Total transferred:      880000 bytes
    HTML transferred:       220000 bytes
    Requests per second:    8920.99 [#/sec] (mean)
    Time per request:       11.210 [ms] (mean)
    Time per request:       0.112 [ms] (mean, across all concurrent requests)
    Transfer rate:          1533.29 [Kbytes/sec] received

    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        0    0   0.5      0       5
    Processing:     0   11   4.8     10      35
    Waiting:        0   11   4.8     10      35
    Total:          0   11   4.8     10      35

    Percentage of the requests served within a certain time (ms)
      50%     10
      66%     12
      75%     14
      80%     15
      90%     18
      95%     20
      98%     23
      99%     26
     100%     35 (longest request)

### ab -n 5000 -c 100 -k http://127.0.0.1:8080/next/

    This is ApacheBench, Version 2.3 <$Revision: 1826891 $>
    Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
    Licensed to The Apache Software Foundation, http://www.apache.org/

    Server Software:        
    Server Hostname:        127.0.0.1
    Server Port:            8080

    Document Path:          /next/
    Document Length:        44 bytes

    Concurrency Level:      100
    Time taken for tests:   3.483 seconds
    Complete requests:      5000
    Failed requests:        4994
       (Connect: 0, Receive: 0, Length: 4994, Exceptions: 0)
    Keep-Alive requests:    5000
    Total transferred:      3493777 bytes
    HTML transferred:       2828622 bytes
    Requests per second:    1435.50 [#/sec] (mean)
    Time per request:       69.662 [ms] (mean)
    Time per request:       0.697 [ms] (mean, across all concurrent requests)
    Transfer rate:          979.56 [Kbytes/sec] received

    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        0    0   1.3      0      12
    Processing:     6   68  38.2     58     197
    Waiting:        3   68  38.2     58     197
    Total:          6   68  38.0     58     197

    Percentage of the requests served within a certain time (ms)
      50%     58
      66%     73
      75%     88
      80%    102
      90%    131
      95%    147
      98%    156
      99%    171
     100%    197 (longest request)
