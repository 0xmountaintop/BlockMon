// curl -X POST http://101.37.164.153:9888/get-block -d '{"block_height": 26}'
package main

import(
    "log"
    // "fmt"
    "strconv"
    "time"
    "encoding/json"
    "io/ioutil"
    "net/http"
    
    "github.com/parnurzeal/gorequest"
    // "github.com/bytom/consensus/difficulty"
)


type t_data struct {
    BlockCount  uint64  `json:"block_count, omitempty"`
    Hash        string  `json:"hash, omitempty"`
    PreBlckHsh  string  `json:"previous_block_hash, omitempty"`
    Size        uint64  `json:"size, omitempty"`
    Version     uint8   `json:"version, omitempty"`
    Height      uint64  `json:"height, omitempty"`
    Timestamp   int64  `json:"timestamp, omitempty"`
    Nonce       uint64  `json:"nonce, omitempty"`
    Bits        uint64  `json:"bits, omitempty"`
    Diffi        string  `json:"difficulty, omitempty"`
}

type t_resp struct {
    Status      string  `json:"status"`
    Data        t_data  `json:"data"`
}


const (
    walletAddr = "http://101.37.164.153:9888/"
    retargetSeconds = 60
)


func main() {
    go http.ListenAndServe(":80", http.FileServer(http.Dir(".")))

dododo:
    request := gorequest.New()
    var resp t_resp

    _, body, _ := request.Post(walletAddr + "get-block-count").
        End()

    json.Unmarshal([]byte(body), &resp)
    if resp.Status != "success" {
        log.Fatalln("Request fail!")    
    }
    log.Printf("Block Count: %d\n\n", resp.Data.BlockCount)

    var dataStr string
    var diffiStr string
    var jsonDiffiStr string
    var elapsed time.Duration
    var diffi_elapsed time.Duration
    last_blck_timestamp := int64(0)
    last_diffi_blck_height := uint64(1)
    last_diffi := ""
    last_diffi_timestamp := int64(0)


    jsonDiffiStr = `
                {
                    "data": {
                        "lines": [
                `


    for i := uint64(1); i <= resp.Data.BlockCount; i++ {
        _, body, _ = request.Post(walletAddr + "get-block").
            Send(`{
                    "block_height": `+ strconv.FormatUint(i, 10) + `
                    }`).
            End()
        json.Unmarshal([]byte(body), &resp)
        // log.Printf("Block %d of %d:\n\tTimestamp: %d, %v\n\tNonce: %d\n\tBits: %d, e.g., %d\n\tDiff: %s",
        // log.Printf("Block %d of %d:\n\tTimestamp: %d, %v\n\tNonce: %d\n\tBits: %d\n\tDiff: %s",

        if i == 1 {
            log.Printf("Block %d of %d:\n\tTimestamp: %d, %v\n\tDiffi: %s",
                resp.Data.Height, resp.Data.BlockCount,
                resp.Data.Timestamp, time.Unix(resp.Data.Timestamp,0),
                // resp.Data.Nonce,
                // resp.Data.Bits, difficulty.CompactToBig(resp.Data.Bits),
                // resp.Data.Bits,
                resp.Data.Diffi)
        } else {
            elapsed = time.Unix(resp.Data.Timestamp,0).Sub(time.Unix(last_blck_timestamp,0))
            if elapsed.Seconds() >= retargetSeconds {
                log.Printf("Block %d of %d:\n\tTimestamp: %d, %v, elapsed: %v\tToo long!!!\n\tDiffi: %s ",
                    resp.Data.Height, resp.Data.BlockCount,
                    resp.Data.Timestamp, time.Unix(resp.Data.Timestamp,0), elapsed.String(),
                    // resp.Data.Nonce,
                    // resp.Data.Bits, difficulty.CompactToBig(resp.Data.Bits),
                    // resp.Data.Bits,
                    resp.Data.Diffi)
            } else{
                log.Printf("Block %d of %d:\n\tTimestamp: %d, %v, elapsed: %v\n\tDiffi: %s",
                    resp.Data.Height, resp.Data.BlockCount,
                    resp.Data.Timestamp, time.Unix(resp.Data.Timestamp,0), elapsed.String(),
                    // resp.Data.Nonce,
                    // resp.Data.Bits, difficulty.CompactToBig(resp.Data.Bits),
                    // resp.Data.Bits,
                    resp.Data.Diffi)
            }
        }

        dataStr += strconv.FormatUint(resp.Data.Height, 10)
        dataStr += "\t"
        dataStr += strconv.FormatInt(resp.Data.Timestamp, 10)
        dataStr += "\t"
        dataStr += resp.Data.Diffi
        dataStr += "\t"
        dataStr += time.Unix(resp.Data.Timestamp,0).String()
        dataStr += "\t"
        if i > 1 {
            dataStr += "elapsed:"
            dataStr += elapsed.String()
            dataStr += "\t"
            if elapsed.Seconds() >= retargetSeconds {
                dataStr += "Too long!!!"
            }
        }
        dataStr += "\n"

        if resp.Data.Diffi != last_diffi {
            diffi_elapsed = time.Unix(resp.Data.Timestamp,0).Sub(time.Unix(last_diffi_timestamp,0))
            if i > 1 {
                log.Printf("Diffi changes!!!\n\t%v at height %v\n\tto\n\t%v at height %v\n\tBlock interval: %v\n\tTaking: %v\n",
                    last_diffi, last_diffi_blck_height, 
                    resp.Data.Diffi, resp.Data.Height,
                    resp.Data.Height - last_diffi_blck_height,
                    diffi_elapsed.String())
            }
            diffiStr += strconv.FormatUint(resp.Data.Height, 10)
            diffiStr += "\t"
            diffiStr += strconv.FormatInt(resp.Data.Timestamp, 10)
            diffiStr += "\t"
            diffiStr += resp.Data.Diffi
            diffiStr += "\t"
            diffiStr += time.Unix(resp.Data.Timestamp,0).String()
            diffiStr += "\t"
            if i > 1 {
                diffiStr += "blocks interval:"
                diffiStr += strconv.FormatUint(resp.Data.Height - last_diffi_blck_height, 10)
                diffiStr += "\t"
                diffiStr += diffi_elapsed.String()
                diffiStr += "\t"
            }
            diffiStr += "\n"

            jsonDiffiStr += `   [`
            millisec := time.Unix(resp.Data.Timestamp,0).UnixNano() / int64(time.Millisecond)
            jsonDiffiStr += strconv.FormatInt(millisec, 10)
            jsonDiffiStr += `               ,`
            jsonDiffiStr += resp.Data.Diffi
            jsonDiffiStr += `               ,`
            jsonDiffiStr += resp.Data.Diffi
            jsonDiffiStr += `               ,`
            jsonDiffiStr += resp.Data.Diffi
            jsonDiffiStr += `               ,`
            jsonDiffiStr += resp.Data.Diffi
            jsonDiffiStr += `               ,
                                    0.0
                                ],`

            last_diffi_blck_height = resp.Data.Height
            last_diffi = resp.Data.Diffi
            last_diffi_timestamp = resp.Data.Timestamp
        }

        last_blck_timestamp = resp.Data.Timestamp
    }

    jsonDiffiStr = jsonDiffiStr[0:len([]rune(jsonDiffiStr))-1]
    jsonDiffiStr += `
                                    ]
                                },
                        "success": true
                    }
                `


    err := ioutil.WriteFile("./all-blocks.csv", []byte(dataStr), 0644)
    check(err)
    err = ioutil.WriteFile("./diffi-changes.csv", []byte(diffiStr), 0644)
    check(err)
    err = ioutil.WriteFile("./data/diffiData.json", []byte(jsonDiffiStr), 0644)
    check(err)

goto dododo
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}
