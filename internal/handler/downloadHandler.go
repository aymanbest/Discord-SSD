package handler

import (
	"fmt"
	"net/http"
	"database/sql"
	"log"
    "encoding/json"
    "strings"


	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/pierre0210/discord-drive/internal/discordutil"
	"github.com/pierre0210/discord-drive/internal/storage"
)

func GetFile(ctx *gin.Context) {
    bot, _ := ctx.MustGet("bot").(*discordgo.Session)
    fileName := ctx.Query("file")
    if fileName == "" {
        ctx.Status(http.StatusBadRequest)
        return
    }

    log.Println("Fetching file ID and ID chain from database...")
    row := storage.GetDB().QueryRow("SELECT files-> $1, id_chain FROM discssd WHERE files ? $1", fileName)
    var fileId string
    var idChainJson string
    err := row.Scan(&fileId, &idChainJson)
    if err != nil {
        if err == sql.ErrNoRows {
            
            ctx.Status(http.StatusNotFound)
        } else {
            
            log.Fatal(err)
        }
        return
    }

    
    var idChain map[string]string
    err = json.Unmarshal([]byte(idChainJson), &idChain)
    if err != nil {
        log.Fatal(err)
        return
    }

    header := ctx.Writer.Header()
    header.Set("Transfer-Encoding", "chunked")
    header.Set("Content-Type", "application/octet-stream")
    header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
    ctx.Writer.WriteHeader(http.StatusOK)
    currentId := strings.Trim(fileId, "\"") 
for {
    log.Printf("Downloading chunk with ID: %s", currentId)
    fileBytes, err := discordutil.DownloadFileFromChannel(bot, currentId)
    if err != nil {
        ctx.Status(http.StatusNotFound)
        return
    }
    ctx.Writer.Write(fileBytes)
    ctx.Writer.(http.Flusher).Flush()
    nextId, exists := idChain[currentId]
    if !exists {
        break
    }
    currentId = strings.Trim(nextId, "\"")
}
}