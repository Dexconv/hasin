package db

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type FileDesc struct {
	gorm.Model
	Name string         `gorm:"column:name"`
	Type string         `gorm:"column:type"`
	Ext  string         `gorm:"column:ext"`
	File string         `gorm:"column:file"`
	Tags pq.StringArray `gorm:"type:text[];column:tags"`
}

type FileDescs []FileDesc

// Creat a single file description
func (fd FileDesc) Add() error {
	return DB.Create(&fd).Error
}

// Remove a single file description
func (fd FileDesc) Remove() error {
	return DB.Unscoped().Delete(&fd).Error
}

// Remove a list of file desciptions
func (fds FileDescs) Remove() error {
	return DB.Unscoped().Delete(&fds).Error
}

// Find file desctiptions using a list of tags
func FindWithTags(tags []string) (FileDescs, error) {
	res := []FileDesc{}
	tagsStr := ""
	for i, t := range tags {
		if i != 0 {
			tagsStr += `,`
		}
		tagsStr += `"` + t + `"`
	}
	err := DB.Raw(`SELECT *
	FROM file_descs
	WHERE '{` + tagsStr + `}' <@ tags
	UNION
	SELECT * from (SELECT *
	FROM file_descs
	WHERE NOT EXISTS (SELECT 1 FROM file_descs WHERE '{` + tagsStr + `}' <@ tags) limit 1) as foo;`,
	).Scan(&res).Error
	return res, err
}

// Find file descriptions using a filename
func FindWithName(name string) (FileDescs, error) {
	res := []FileDesc{}
	err := DB.Raw(`SELECT *
	FROM file_descs
	WHERE Name = ? 
	UNION
	SELECT * from (SELECT *
	FROM file_descs
	WHERE NOT EXISTS (SELECT 1 FROM file_descs WHERE Name = ?) limit 1) as foo;`,
		name, name).Scan(&res).Error
	return res, err
}
