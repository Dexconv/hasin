package rest

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dexconv/hasin/store/common"
	"github.com/dexconv/hasin/store/config"
	"github.com/dexconv/hasin/store/db"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func upload(c echo.Context) error {

	filesSize, err := common.DirSize("files/")
	if err != nil {
		log.Error("could not read dir size", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "could not read dir size",
		})
	}
	fmt.Println(filesSize)
	if filesSize >= config.GLB.StorageSize*common.MB {
		log.Warn("storage limit reached")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "storage limit reached",
		})
	}

	savedAs := uuid.New().String()

	tags := c.FormValue("tags")
	tagsArr := strings.Split(tags, ",")
	file, err := c.FormFile("file")
	if err != nil {
		log.Error("file was not recieved", "error", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "file was not recieved",
		})
	}
	src, err := file.Open()
	if err != nil {
		log.Error("error openning file", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "error openning file",
		})
	}
	defer src.Close()

	srcBod, err := ioutil.ReadAll(src)
	if err != nil {
		log.Error("could not read file body", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "could not detect filetype",
		})
	}

	mtype := mimetype.Detect(srcBod)

	encoded, err := common.Encrypt(srcBod, common.Secret)
	if err != nil {
		log.Error("error encrypting file", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "error encrypting file",
		})
	}

	dst, err := os.Create("files/" + savedAs)
	if err != nil {
		log.Error("error creating file", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "error creating file",
		})
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, strings.NewReader(encoded)); err != nil {
		log.Error("error writing to file", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "error writing to file",
		})
	}

	err = db.FileDesc{
		Name: file.Filename,
		Type: mtype.String(),
		Ext:  mtype.Extension(),
		File: savedAs,
		Tags: db.ToPqTextArray(tagsArr),
	}.Add()

	if err != nil {
		log.Error("error finding files with tags", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed to add file",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "file added successfully",
	})
}

func getWithTags(c echo.Context) error {
	tags := c.Param("tags")
	mode := c.Param("mode")
	fds := db.FileDescs{}
	var err error

	if mode == "tags" {
		tagsArr := strings.Split(tags, ",")
		fds, err = db.FindWithTags(tagsArr)
		if err != nil {
			log.Error("error finding files with tags", "error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status": "failed to retrieve files",
			})
		}

	} else if mode == "name" {
		fds, err = db.FindWithName(tags)
		if err != nil {
			log.Error("error finding files with name", "error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status": "failed to retrieve files",
			})
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status": `please select 'name' or 'tags' mode`,
		})
	}

	if len(fds) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "no files could be found",
		})
	}

	archiveName := uuid.New().String() + ".zip"
	archive, err := os.Create(archiveName)
	if err != nil {
		log.Warn("could not creat archive file", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "could not creat archive file",
		})
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	filesNum := 0
	cleanup := []string{archiveName}
	for i, fd := range fds {
		f1, err := os.Open("files/" + fd.File)
		if err != nil {
			if os.IsNotExist(err) {
				log.Warn("file does not exist, deleting coresponding db record", "error", err.Error())
				err = fd.Remove()
				if err != nil {
					log.Warn("could not remove file descriptions", "error", err.Error())
					return c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"message": "could not remove file descriptions",
					})
				}
			} else {
				log.Warn("could not open file, skipping...", "error", err.Error())
			}
			continue
		}
		defer f1.Close()

		fb, err := ioutil.ReadAll(f1)
		if err != nil {
			log.Warn("could not read file body", "error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "could not read file body",
			})
		}
		fbd, err := common.Decrypt(string(fb), common.Secret)
		if err != nil {
			log.Warn("could not decrypt file body", "error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "could not decrypt file body",
			})
		}

		// add files to zip archive
		w1, err := zipWriter.Create("files/" + fmt.Sprint(i, "-") + fd.Name)
		if err != nil {
			log.Warn("could not create archive files", "error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "could not create archive files",
			})
		}
		if _, err := io.Copy(w1, strings.NewReader(fbd)); err != nil {
			log.Warn("could not copy files to archive", "error", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "could not copy files to archive",
			})
		}

		// add file names to cleanup
		filesNum += 1
		cleanup = append(cleanup, fd.File)
	}
	zipWriter.Close()

	if filesNum == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "no files could be found",
		})
	}

	if err = fds.Remove(); err != nil {
		log.Warn("could not remove file descriptions", "error", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "could not remove file descriptions",
		})
	}

	// delete files
	defer func(fileNames []string) {
		for _, fn := range fileNames {
			err := os.Remove("files/" + fn)
			if err != nil {
				log.Error("could not create archive files", "error", err.Error())
			}
		}
	}(cleanup)

	return c.File(archiveName)
}
