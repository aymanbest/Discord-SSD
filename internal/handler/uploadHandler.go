package handler

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/pierre0210/discord-drive/internal/discordutil"
	"github.com/pierre0210/discord-drive/internal/storage"
)

func PostUpload(ctx *gin.Context) {
    var table storage.FileTable
    table.Files = make(map[string]string)
    table.IdChain = make(map[string]string)
    bot, _ := ctx.MustGet("bot").(*discordgo.Session)
    form, err := ctx.MultipartForm()
    if err != nil {
        log.Println(err.Error())
    }
    files := form.File["files"]

    for _, file := range files {
        var prevId string
        chunkSize, _ := strconv.Atoi(os.Getenv("CHUNKSIZE"))
        content, _ := file.Open()

        reader := bufio.NewReaderSize(content, chunkSize)
        for i := 0; i < int(math.Ceil(float64(file.Size)/float64(chunkSize))); i++ {
            splitBuff := make([]byte, chunkSize)
            size, _ := reader.Read(splitBuff)
            chunkSum := fmt.Sprintf("%x", md5.Sum(splitBuff[:size]))
            log.Printf("%d %s", size, chunkSum)
            message := discordutil.UploadFileToChannel(bot, chunkSum, bytes.NewBuffer(splitBuff[:size]))

            if i == 0 {
                table.Files[file.Filename] = message.ID
                prevId = message.ID
            } else {
                table.IdChain[prevId] = message.ID
                prevId = message.ID
            }
        }
        content.Close()
    }

    filesJson, err := json.Marshal(table.Files)
    if err != nil {
        log.Fatal(err)
    }
    idChainJson, err := json.Marshal(table.IdChain)
    if err != nil {
        log.Fatal(err)
    }

    _, err = storage.GetDB().Exec(`INSERT INTO discssd (files, id_chain) VALUES ($1, $2)`, filesJson, idChainJson)
    if err != nil {
        log.Fatal(err)
    }
}

