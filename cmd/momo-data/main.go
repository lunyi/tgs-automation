package main

import (
	"bytes"
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/postgresql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/alexmullins/zip"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/tealeg/xlsx"
)

func main() {
	config := util.GetConfig()
	app := postgresql.NewMomoDataInterface(config.Postgresql)
	defer app.Close()

	today := time.Now().Format("2006-01-02") // Today's date in "YYYY-MM-DD" format
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	data := []map[string]string{
		{"Brand": "MOVN2", "Field": "pdp.first_deposit_on", "File": "deposit"},
		{"Brand": "MOVN2", "Field": "p.registered_on", "File": "register"},
		{"Brand": "MOPH", "Field": "pdp.first_deposit_on", "File": "deposit"},
		{"Brand": "MOPH", "Field": "p.registered_on", "File": "register"},
	}

	session := createSession(config.AwsS3)
	// Print the slice of maps
	for _, item := range data {
		fmt.Printf("Code: %s, Field: %s\n", item["Brand"], item["Field"])

		players, err := app.GetMomodata(item["Brand"], item["Field"], yesterday, today, "+08:00")
		if err != nil {
			log.LogFatal(err.Error())
		}

		filename, err := CreateExcel(players, item["Brand"], item["File"])
		if err != nil {
			log.LogFatal(err.Error())
		}
		password := fmt.Sprintf("PG%s", time.Now().Format("20060102"))
		zipFile, err := Zipfile(filename, password)
		if err != nil {
			log.LogFatal(err.Error())
		}
		uploadFileToS3(session, config.AwsS3.Bucket, zipFile, fmt.Sprintf("./%s", zipFile))
	}
}

func CreateExcel(players []postgresql.PlayerInfo, brand string, dateField string) (string, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("PlayerInfo")
	if err != nil {
		log.LogFatal(fmt.Sprintf("AddSheet failed: %s", err))
		return "", err
	}

	headerRow := sheet.AddRow()
	headerTitles := []string{"Agent", "Host", "PlayerName", "DailyDepositAmount", "DailyDepositCount", dateField}
	for _, title := range headerTitles {
		cell := headerRow.AddCell()
		cell.Value = title
	}

	// Populating data
	for _, player := range players {
		row := sheet.AddRow()
		row.AddCell().Value = player.Agent.String
		row.AddCell().Value = player.Host
		row.AddCell().Value = player.PlayerName
		row.AddCell().SetFloat(player.DailyDepositAmount)
		row.AddCell().SetInt(player.DailyDepositCount)
		row.AddCell().Value = player.FirstDepositOn.Format(time.RFC3339)
	}

	filename := fmt.Sprintf("PG%s-%s-%s.xlsx", time.Now().Format("20060102"), brand, dateField)

	// Save the file to the disk
	err = file.Save(filename)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Save failed:: %s", err))
		return "", err
	}

	log.LogInfo("File saved successfully.")
	return filename, nil
}

func Zipfile(fileToZip string, password string) (string, error) {

	zipFileName := strings.Replace(fileToZip, ".xlsx", ".zip", -1)

	log.LogInfo(fmt.Sprintf("Creating ZIP file: %s", zipFileName))

	outFile, err := os.Create(zipFileName)
	if err != nil {
		log.LogFatal(err.Error())
		return "", err
	}
	defer outFile.Close()

	// Create a new ZIP writer
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Open the file to be zipped
	inFile, err := os.Open(fileToZip)
	if err != nil {
		log.LogFatal(err.Error())
		return "", err
	}
	defer inFile.Close()

	// Get the file info to replicate it in the ZIP
	info, err := inFile.Stat()
	if err != nil {
		log.LogFatal(err.Error())
		return "", err
	}

	// Create a header based on the file info
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		log.LogFatal(err.Error())
		return "", err
	}
	header.Name = fileToZip // Ensure filename is correct

	// Encrypt and write the file into the ZIP
	zipFileWriter, err := zipWriter.Encrypt(zipFileName, password)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error creating encrypted zip entry: %s", err))
		return "", err
	}

	_, err = io.Copy(zipFileWriter, inFile)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error writing file to zip: %s", err))
		return "", err
	}

	if err := zipWriter.Close(); err != nil {
		log.LogFatal(fmt.Sprintf("Error closing zip writer: %s", err))
		return "", err
	}

	log.LogInfo("ZIP file created successfully.")
	return zipFileName, nil
}

func createSession(config util.AwsS3Config) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKey, config.AccessSecret, ""),
	})
	if err != nil {
		log.LogFatal(fmt.Sprintf("Failed to create AWS session: %s", err))
	}

	return sess
}

func uploadFileToS3(sess *session.Session, bucketName, fileKey, filePath string) {
	// Read the file
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Unable to read file: %s", err))
	}

	// Create an uploader with the session and default options
	uploader := s3.New(sess)

	// Upload the file's bytes to S3
	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
		Body:   bytes.NewReader(fileBytes),
	})
	if err != nil {
		log.LogFatal(fmt.Sprintf("Unable to upload file to S3: %s", err))
	}

	log.LogInfo("Successfully uploaded file to S3")
}
