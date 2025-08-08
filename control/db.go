package control

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB 初始化数据库，创建表结构
func InitDB() (*sql.DB, error) {
	var err error
	// 使用 sync.Once 确保数据库只初始化一次
	once.Do(func() {
		db, err = sql.Open("sqlite3", "./files.db")
		if err != nil {
			log.Fatal("Failed to open database:", err)
		}

		query := `CREATE TABLE IF NOT EXISTS uploaded_files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			fileId TEXT NOT NULL,
			filename TEXT NOT NULL,
			ip TEXT NOT NULL,
			time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal("Failed to create table:", err)
		}

		// 创建短链表
		shortLinkQuery := `CREATE TABLE IF NOT EXISTS short_links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			short_code TEXT UNIQUE NOT NULL,
			file_id TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			access_count INTEGER DEFAULT 0
		);`
		_, err = db.Exec(shortLinkQuery)
		if err != nil {
			log.Fatal("Failed to create short_links table:", err)
		}

		// 创建分片记录表
		chunkQuery := `CREATE TABLE IF NOT EXISTS chunk_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			upload_id TEXT NOT NULL,
			chunk_index INTEGER NOT NULL,
			chunk_id TEXT NOT NULL,
			file_name TEXT NOT NULL,
			ip TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(upload_id, chunk_index)
		);`
		_, err = db.Exec(chunkQuery)
		if err != nil {
			log.Fatal("Failed to create chunk_records table:", err)
		}
	})

	return db, err
}

type FileRecord struct {
	FileId   string    `json:"fileId"`
	Filename string    `json:"filename"`
	Ip       string    `json:"ip"`
	Time     time.Time `json:"time"`
}

type ShortLink struct {
	ID          int       `json:"id"`
	ShortCode   string    `json:"shortCode"`
	FileId      string    `json:"fileId"`
	CreatedAt   time.Time `json:"createdAt"`
	AccessCount int       `json:"accessCount"`
}

// GetFileNameByIDOrName 查询文件名
func GetFileNameByIDOrName(idOrName string) (FileRecord, error) {
	var record FileRecord
	// 执行查询，获取对应id或name的file记录
	query := "SELECT fileId, filename, ip, time FROM uploaded_files WHERE fileId = ? OR filename = ? ORDER BY time DESC LIMIT 1"
	err := db.QueryRow(query, idOrName, idOrName).Scan(&record.FileId, &record.Filename, &record.Ip, &record.Time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return FileRecord{}, fmt.Errorf("no file found with idOrName %s", idOrName)
		}
		return FileRecord{}, err
	}

	return record, nil
}

func SaveFileRecord(fileID string, fileName string, ip string) error {
	// 插入数据到数据库
	_, err := db.Exec("INSERT INTO uploaded_files (fileId, filename, ip) VALUES (?, ?, ?)", fileID, fileName, ip)
	return err
}

func SelectAllRecord() ([]FileRecord, error) {
	// 查询所有记录
	rows, err := db.Query("SELECT fileId, filename, ip, time FROM uploaded_files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []FileRecord

	// 迭代查询结果
	for rows.Next() {
		var record FileRecord
		err := rows.Scan(&record.FileId, &record.Filename, &record.Ip, &record.Time)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	// 检查查询错误
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// CreateShortLink 创建短链
func CreateShortLink(shortCode, fileId string) error {
	_, err := db.Exec("INSERT INTO short_links (short_code, file_id) VALUES (?, ?)", shortCode, fileId)
	return err
}

// ShortCodeExists 检查短链码是否已存在
func ShortCodeExists(shortCode string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM short_links WHERE short_code = ?", shortCode).Scan(&count)
	return err == nil && count > 0
}

// GetAllShortLinks 获取所有短链
func GetAllShortLinks() ([]ShortLink, error) {
	rows, err := db.Query("SELECT id, short_code, file_id, created_at, access_count FROM short_links ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shortLinks []ShortLink
	for rows.Next() {
		var link ShortLink
		err := rows.Scan(&link.ID, &link.ShortCode, &link.FileId, &link.CreatedAt, &link.AccessCount)
		if err != nil {
			return nil, err
		}
		shortLinks = append(shortLinks, link)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shortLinks, nil
}

// GetFileIdByShortCode 通过短链码获取文件ID
func GetFileIdByShortCode(shortCode string) (string, error) {
	var fileId string
	err := db.QueryRow("SELECT file_id FROM short_links WHERE short_code = ?", shortCode).Scan(&fileId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("short link not found: %s", shortCode)
		}
		return "", err
	}

	// 增加访问计数
	_, _ = db.Exec("UPDATE short_links SET access_count = access_count + 1 WHERE short_code = ?", shortCode)

	return fileId, nil
}

// GetShortCodeByFileId 通过文件ID获取短链码
func GetShortCodeByFileId(fileId string) (string, error) {
	var shortCode string
	err := db.QueryRow("SELECT short_code FROM short_links WHERE file_id = ?", fileId).Scan(&shortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("short link not found for file: %s", fileId)
		}
		return "", err
	}
	return shortCode, nil
}

// SaveChunkRecord 保存分片记录
func SaveChunkRecord(uploadId, chunkIndex, chunkId, fileName, ip string) error {
	_, err := db.Exec("INSERT OR REPLACE INTO chunk_records (upload_id, chunk_index, chunk_id, file_name, ip) VALUES (?, ?, ?, ?, ?)",
		uploadId, chunkIndex, chunkId, fileName, ip)
	return err
}

// GetChunkRecords 获取指定上传ID的所有分片记录
func GetChunkRecords(uploadId string) ([]ChunkRecord, error) {
	rows, err := db.Query("SELECT upload_id, chunk_index, chunk_id, file_name, ip, created_at FROM chunk_records WHERE upload_id = ? ORDER BY chunk_index", uploadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ChunkRecord
	for rows.Next() {
		var record ChunkRecord
		err := rows.Scan(&record.UploadId, &record.ChunkIndex, &record.ChunkId, &record.FileName, &record.Ip, &record.CreatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// CleanupChunkRecords 清理分片记录
func CleanupChunkRecords(uploadId string) error {
	_, err := db.Exec("DELETE FROM chunk_records WHERE upload_id = ?", uploadId)
	return err
}

type ChunkRecord struct {
	UploadId   string    `json:"uploadId"`
	ChunkIndex int       `json:"chunkIndex"`
	ChunkId    string    `json:"chunkId"`
	FileName   string    `json:"fileName"`
	Ip         string    `json:"ip"`
	CreatedAt  time.Time `json:"createdAt"`
}
