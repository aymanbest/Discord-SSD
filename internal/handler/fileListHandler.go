package handler

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pierre0210/discord-drive/internal/storage"
)

func GetFileList(ctx *gin.Context) {
	// Get the list of files from the database
	rows, err := storage.GetDB().Query("SELECT files FROM discssd")
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

func DeleteFileHandler(ctx *gin.Context) {
	// Get the file name from the request parameters
	fileName := ctx.Param("name")

	// Validate the file name
	if fileName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
		return
	}

	// Delete the file entry from the database
	err := deleteFileFromDB(fileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from database"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func deleteFileFromDB(fileName string) error {
	log.Println("Deleting file:", fileName)
	// Execute the SQL DELETE statement
	result, err := storage.GetDB().Exec(`DELETE FROM discssd WHERE "files" = $1`, fileName)
	if err != nil {
		log.Println("Error executing DELETE statement:", err)
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
	oldFileName := ctx.Param("oldFileName")
	newFileName := ctx.Param("newFileName")

	err := renameFileInDB(oldFileName, newFileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rename file in database"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File renamed successfully"})
}

func renameFileInDB(oldFileName string, newFileName string) error {
	log.Println("Renaming file:", oldFileName, "to", newFileName)
	result, err := storage.GetDB().Exec(`UPDATE discssd SET "files" = $1 WHERE "files" = $2`, newFileName, oldFileName)
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