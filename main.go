package main

import (
    "os"
    "flag"
    "fmt"
    "gortc.io/stun"
    "strings"
    "sync"
)

const SRV_LIST_SEP = ";"

type CLIArgs struct {
    servers string
    quorum uint
}

func parse_args() CLIArgs {
    var args CLIArgs
    flag.StringVar(&args.servers, "s", "stun.l.google.com:19302;" +
                   "stun.ekiga.net:3478;stun.ideasip.com:3478;" +
                   "stun.schlund.de:3478;stun.voiparound.com:3478;" +
                   "stun.voipbuster.com:3478;stun.voipstunt.com:3478",
                   "STUN server list")
    flag.UintVar(&args.quorum, "q", 2, "required number of matches for success")
    flag.Parse()
    return args
}

func getAddr(server string) (string, error) {
    c, err := stun.Dial("udp", server)
    if err != nil {
        return "", err
    }
    defer c.Close()

    message, err := stun.Build(stun.TransactionID, stun.BindingRequest)
    if err != nil {
        return "", err
    }

    var result string
    var res_err error

    if err := c.Do(message, func(res stun.Event) {
		if res.Error != nil {
			result = ""
            res_err = res.Error
            return
		}
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err == nil {
            result = xorAddr.IP.String()
            res_err = nil
		} else {
            var mappedAddr stun.MappedAddress
            if err := mappedAddr.GetFrom(res.Message); err == nil {
                result = mappedAddr.IP.String()
                res_err = nil
            } else {
                result = ""
                res_err = err
            }
        }
	}); err != nil {
		return "", err
	}
    return result, res_err
}

func worker(server string, out chan<- string, err chan<- error) {
    res_addr, res_err := getAddr(server)
    if res_err != nil {
        err <- res_err
    } else {
        out <- res_addr
    }
}

func run() int {
    args := parse_args()

    srvList := strings.Split(args.servers, SRV_LIST_SEP)
    if len(srvList) == 0 {
        fmt.Fprintf(os.Stderr, "Error: empty server list\n")
        return 2
    }
    if uint(len(srvList)) < args.quorum {
        fmt.Fprintf(os.Stderr, "Error: quorum is higher than server list length\n")
    }

    results := make(chan string)
    errors := make(chan error)

    var wg sync.WaitGroup
    for _, v := range srvList {
        wg.Add(1)
        go func() {
            defer wg.Done()
            worker(v, results, errors)
        }()
    }
    done_event := make(chan struct{}, 1)
    go func() {
        wg.Wait()
        done_event <-struct{}{}
    }()

    resultMap := make(map[string]uint)
    var errorList []error

    for {
        select {
        case err := <-errors:
            errorList = append(errorList, err)
        case res := <-results:
            resultMap[res]++
            if resultMap[res] >= args.quorum {
                fmt.Println(res)
                return 0
            }
        case <-done_event:
            fmt.Fprintln(os.Stderr, "All servers queried, but quorum wasn't reached")
            for _, v := range errorList {
                fmt.Fprintf(os.Stderr, "Error: %v\n", v)
            }
            return 3
        }
    }
    return 0
}

func main() {
    os.Exit(run())
}
