// curl -X POST http://101.37.164.153:9888/get-block -d '{"block_height": 26}'
package main

import(
    "log"
    // "fmt"
    "strconv"
    "encoding/json"
    "time"
    // "encoding/hex"
    // "encoding/binary"
    
    "github.com/parnurzeal/gorequest"
    // "github.com/bytom/protocol/bc"
    // "github.com/bytom/crypto/sha3pool"
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
    Diff        string  `json:"difficulty, omitempty"`
}

type t_resp struct {
    Status      string  `json:"status"`
    Data        t_data  `json:"data"`
}


const (
    walletAddr = "http://101.37.164.153:9888/"
)


func main() {
    request := gorequest.New()
    var resp t_resp

    _, body, _ := request.Post(walletAddr + "get-block-count").
        End()

    json.Unmarshal([]byte(body), &resp)
    if resp.Status != "success" {
        log.Fatalln("Request fail!")    
    }
    log.Printf("Block Count: %d\n\n", resp.Data.BlockCount)

    for i := uint64(1); i <= resp.Data.BlockCount; i++ {
        _, body, _ = request.Post(walletAddr + "get-block").
            Send(`{
                    "block_height": `+ strconv.FormatUint(i, 10) + `
                    }`).
            End()
        json.Unmarshal([]byte(body), &resp)
        log.Printf("Block %d of %d:\n\tTimestamp: %d, %v\n\tNonce: %d\n\tBits: %d\n\tDiff: %s",
            resp.Data.Height, resp.Data.BlockCount,
            resp.Data.Timestamp, time.Unix(resp.Data.Timestamp,0),
            resp.Data.Nonce,
            resp.Data.Bits,
            resp.Data.Diff)
    }
}

/*
func str2bytes(instr string, leng uint8) []byte {
    // fmt.Println([]byte(instr))
    outstr := fmt.Sprintf("%064s", instr)
    // fmt.Println(outstr)

    var b [32]byte
    hex.Decode(b[:], []byte(outstr))
    if len(instr) < 64 {
        b = litE2BigE(b)    
    }
    // fmt.Println(b)

    h := bc.NewHash(b)
    // fmt.Println(h.Bytes()[0:leng])
    return h.Bytes()[0:leng]
}

func litE2BigE(buf [32]byte) [32]byte {
    blen := len(buf)
    for i := 0; i < blen/2; i++ {
        buf[i], buf[blen-1-i] = buf[blen-1-i], buf[i]
    }
    return buf
}

func ui64To8Bytes(ui64 uint64) []byte {
    bs := make([]byte, 8)
    binary.LittleEndian.PutUint64(bs, ui64)
    // fmt.Println(bs)
    return bs
}*/