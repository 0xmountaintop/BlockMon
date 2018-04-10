// curl -X POST http://101.37.164.153:9888/get-block -d '{"block_height": 26}'
package main

import(
    "log"
    "fmt"
    "strconv"
    "encoding/json"
    // "encoding/hex"
    // "encoding/binary"
    
    "github.com/parnurzeal/gorequest"
    // "github.com/bytom/protocol/bc"
    // "github.com/bytom/crypto/sha3pool"
    // "github.com/bytom/consensus/difficulty"
)


type t_data struct {
    BlockCount  uint64  `json:"block_count, omitempty"`
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
        fmt.Println(body)
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