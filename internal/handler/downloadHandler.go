package handler

import (
	"fmt"
	"net/http"
	"database/sql"
	"log"

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
	// Get the indexId from the database
	row := storage.GetDB().QueryRow("SELECT id_chain FROM discssd WHERE files = $1", fileName)
	var indexId string
	err := row.Scan(&indexId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Status(http.StatusNotFound)
		} else {
			log.Fatal(err)
		}
		return
	}

	header := ctx.Writer.Header()
	header.Set("Transfer-Encoding", "chunked")
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	ctx.Writer.WriteHeader(http.StatusOK)

	for {
		fileBytes, err := discordutil.DownloadFileFromChannel(bot, indexId)
		if err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.Writer.Write(fileBytes)
		ctx.Writer.(http.Flusher).Flush()

		
		row := storage.GetDB().QueryRow("SELECT id_chain FROM discssd WHERE files = $1", indexId)
		err = row.Scan(&indexId)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			} else {
				log.Fatal(err)
			}
		}
	}
}