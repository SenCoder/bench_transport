/*
 * Copyright (c) 2019 PingAn. All rights reserved.
 *
 * Created by yuansheng on 6/4/19 8:25 PM.
 */

package main

import (
    "eyeSky/modules/transport/model"
    "flag"
    "log"
    "net"
    "net/rpc"
    "sync"
    "sync/atomic"
    "time"
)

var (
    tcpAddress = flag.String("transport-tcp-address", "127.0.0.1:3002", "<addr>:<port> to connect to transport")
    size       = flag.Int("size", 200, "size of messages")
    batchSize  = flag.Int("batch-size", 200, "batch size of messages")
    runfor     = flag.Duration("d", 10*time.Second, "duration of time to run")
    deadline   = flag.String("deadline", "", "deadline to start the benchmark run")
    conns      = flag.Int("c", 1, "number of tcp connections")
)

var totalMsgCount int64

func main() {
    flag.Parse()
    var wg sync.WaitGroup

    log.SetPrefix("[bench_transport] ")

    msg := make([]byte, *size)
    req := model.RpcArgs{
        Metrics: make([]model.MetricValue, *batchSize),
    }

    for i := range req.Metrics {
        req.Metrics[i].Data = msg
    }

    goChan := make(chan int)
    rdyChan := make(chan int)
    for j := 0; j < *conns; j++ {
        // log.Printf("put worker %d", j)
        wg.Add(1)
        go func() {
            pubWorker(*runfor, *tcpAddress, req, rdyChan, goChan)
            wg.Done()
        }()
        <-rdyChan
    }

    if *deadline != "" {
        t, err := time.Parse("2006-01-02 15:04:05", *deadline)
        if err != nil {
            log.Fatal(err)
        }
        d := t.Sub(time.Now())
        log.Printf("sleeping until %s (%s)", t, d)
        time.Sleep(d)
    }

    start := time.Now()
    close(goChan)
    wg.Wait()
    end := time.Now()
    duration := end.Sub(start)
    tmc := atomic.LoadInt64(&totalMsgCount)
    log.Printf("duration: %s - %.03fmb/s - %.03fops/s - %.03fus/op",
        duration,
        float64(tmc*int64(*size))/duration.Seconds()/1024/1024,
        float64(tmc)/duration.Seconds(),
        float64(duration/time.Microsecond)/float64(tmc))
}

func pubWorker(td time.Duration, addr string, batch model.RpcArgs, rdyChan chan int, goChan chan int) {

    tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
    if err != nil {
        log.Fatalf("Fatal error: %s", err.Error())
    }

    conn, err := net.DialTCP("tcp4", nil, tcpAddr)
    if err != nil {
        log.Fatalf("Fatal error: %s", err.Error())
    }

    c := rpc.NewClient(conn)

    rdyChan <- 1
    <-goChan
    var msgCount int64
    endTime := time.Now().Add(td)
    for {
        res := model.RpcReply{}
        err = c.Call("Transport.Update", batch, &res)
        if err != nil {
            log.Fatal("Update error: ", err)
        }

        if res.Invalid != 0 {
            panic("invalid data transport")
        }

        msgCount += int64(len(batch.Metrics))
        if time.Now().After(endTime) {
            break
        }
    }
    atomic.AddInt64(&totalMsgCount, msgCount)
}

