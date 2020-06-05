package main

import (
    "os"
    "flag"
    "fmt"
    "strings"
    "github.com/Snawoot/extip"
    "context"
    "time"
)

const SRV_LIST_SEP = ";"

type CLIArgs struct {
    servers string
    quorum uint
    timeout time.Duration
    ipv6 bool
}

func parseArgs() CLIArgs {
    var args CLIArgs
    flag.StringVar(&args.servers, "s", strings.Join(extip.PublicServerList, SRV_LIST_SEP),
                   "STUN server list")
    flag.UintVar(&args.quorum, "q", 2, "required number of matches for success")
    flag.DurationVar(&args.timeout, "t", 0, "hard timeout. Examples values: 1m, 3s, 1s500ms")
    flag.BoolVar(&args.ipv6, "6", false, "use IPv6")
    flag.Parse()
    return args
}

func run() int {
    args := parseArgs()

    var ctx context.Context
    if args.timeout == 0 {
        ctx = context.Background()
    } else {
        ctx, _ = context.WithTimeout(context.Background(), args.timeout)
    }

    ip, err := extip.QueryMultipleServers(ctx,
                                          strings.Split(args.servers, SRV_LIST_SEP),
                                          args.quorum,
                                          args.ipv6)
    if err != nil {
        switch res := err.(type) {
        case extip.InconclusiveResult:
            fmt.Fprintf(os.Stderr, "Inconclusive result:\n")
            fmt.Fprintf(os.Stderr, "Required quorum = %v\n", res.Quorum)
            for k, v := range res.Results {
                fmt.Fprintf(os.Stderr, "Server %s responded: %s\n", k, v)
            }
            for k, v := range res.Errors {
                fmt.Fprintf(os.Stderr, "Server %s failed: %v\n", k, v)
            }
        default:
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        }
        return 3
    }
    fmt.Println(ip)
    return 0
}

func main() {
    os.Exit(run())
}
