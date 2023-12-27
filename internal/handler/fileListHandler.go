package handler

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pierre0210/discord-drive/internal/storage"
)

func renameFileInDB(oldFileName string, newFileName string) error {
    log.Println("Renaming file:", oldFileName, "to", newFileName)

    sqlStatement := `UPDATE discssd SET files = files - $1 || jsonb_build_object($2::text, files-> $1) WHERE files ? $1`

    result, err := storage.GetDB().Exec(sqlStatement, oldFileName, newFileName)
    if err != nil {
        log.Println("Error executing UPDATE statement:", err)
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Println("Error getting rows affected:", err)
        return err
    }

    if rowsAffected == 0 {
        log.Println("No rows affected. File may not exist in the database.")
    }

    return nil
}

func RenameFileHandler(ctx *gin.Context) {
	log.Println("RenameFileHandler is hit")
    oldName := ctx.Param("oldName")
    newName := ctx.Param("newName")

    if oldName == "" || newName == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
        return
    }

    err := renameFileInDB(oldName, newName)
    if err != nil {
        log.Fatal(err)
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "File renamed successfully"})
}


func deleteFileFromDB(fileName string) error {

    _, err := storage.GetDB().Exec("DELETE FROM discssd WHERE files ? $1", fileName)
    return err
}

func DeleteFileHandler(ctx *gin.Context) {
    fileName := ctx.Param("name")

    if fileName == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
        return
    }

    err := deleteFileFromDB(fileName)
    if err != nil {
        log.Fatal(err)
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}


func GetFileList(ctx *gin.Context) {

    rows, err := storage.GetDB().Query("SELECT jsonb_object_keys(files) FROM discssd")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var files []string
    for rows.Next() {
        var file string
        if err := rows.Scan(&file); err != nil {
            log.Fatal(err)
        }
        files = append(files, file)
    }

    ctx.JSON(http.StatusOK, gin.H{
        "files": files,
    })
}