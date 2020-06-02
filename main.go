package main

import (
    "os"
    "flag"
    "fmt"
    "gortc.io/stun"
)

type CLIArgs struct {
    server string
}

func parse_args() CLIArgs {
    var args CLIArgs
    flag.StringVar(&args.server, "s", "stun.l.google.com:19302", "STUN server address")
    flag.Parse()
    return args
}

func getAddr(server string) (string, error) {
    c, err := stun.Dial("udp", server)
    if err != nil {
        return "", err
    }

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
		}
		// Decoding XOR-MAPPED-ADDRESS attribute from message.
		var addr stun.XORMappedAddress
		if err := addr.GetFrom(res.Message); err != nil {
			result = ""
            res_err = err
		}
		result = addr.IP.String()
        res_err = nil
	}); err != nil {
		return "", err
	}
    return result, res_err
}

func run() int {
    args := parse_args()

    addr, err := getAddr(args.server)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v", err)
        return 3
    }
    fmt.Println(addr)
    return 0
}

func main() {
    os.Exit(run())
}