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
	"path/filepath"
	"time"

	"github.com/alexmullins/zip"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tealeg/xlsx"
)

func getData() map[string][]map[string]string {
	dataMOVN2 := []map[string]string{
		{"Brand": "MOVN2", "Field": "pdp.first_deposit_on", "Column": "first_deposit_on", "File": "deposit"},
		{"Brand": "MOVN2", "Field": "p.registered_on", "Column": "registered_on", "File": "registered"},
	}

	dataMOPH := []map[string]string{
		{"Brand": "MOPH", "Field": "pdp.first_deposit_on", "Column": "first_deposit_on", "File": "deposit"},
		{"Brand": "MOPH", "Field": "p.registered_on", "Column": "registered_on", "File": "registered"},
	}

	return map[string][]map[string]string{
		"MOVN2": dataMOVN2,
		"MOPH":  dataMOPH,
	}
}

func main() {
	config := util.GetConfig()
	app := postgresql.NewMomoDataInterface(config.Postgresql)
	defer app.Close()

	today := time.Now().Format("2006-01-02") // Today's date in "YYYY-MM-DD" format
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	prefilename := time.Now().AddDate(0, 0, -1).Format("0102")
	session := createSession(config.AwsS3)
	password := fmt.Sprintf("PG%s", time.Now().Format("20060102"))
	data := getData()

	for brand, records := range data {
		fmt.Printf("Records for %s:\n", brand)
		filenames := []string{}
		for _, record := range records {
			filenames = getMomoData(record, app, brand, yesterday, today, filenames)
		}
		zipAndUpoload(prefilename, brand, filenames, password, session, config)

		fmt.Println("-----")
	}
}

func getMomoData(
	record map[string]string,
	app postgresql.GetMomoDataInterface,
	brand string,
	yesterday string,
	today string,
	filenames []string) []string {

	fmt.Printf("Field: %s, File: %s\n", record["Field"], record["File"])
	players, err := app.GetMomodata(brand, record["Field"], yesterday, today, "+08:00")
	if err != nil {
		log.LogFatal(err.Error())
	}

	filename, err := CreateExcel(players, brand, record["Column"], record["File"])
	if err != nil {
		log.LogFatal(err.Error())
	}
	filenames = append(filenames, filename)
	return filenames
}

func zipAndUpoload(
	prefilename string,
	brand string,
	filenames []string,
	password string,
	session *session.Session,
	config util.TgsConfig) {

	zipfilename := fmt.Sprintf("%s_%s.zip", prefilename, brand)
	zipFile, err := Zipfiles(zipfilename, filenames, password)
	if err != nil {
		log.LogFatal(err.Error())
	}
	filePath := fmt.Sprintf("./%s", zipFile)
	uploadFileToS3(session, config.AwsS3.Bucket, zipFile, filePath)
	TelegramNotify(config.MomoTelegram, filePath, fmt.Sprintf("%s Data", brand))

	deleteFiles()
}

func CreateExcel(players []postgresql.PlayerInfo, brand string, column string, fileField string) (string, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("PlayerInfo")
	if err != nil {
		log.LogFatal(fmt.Sprintf("AddSheet failed: %s", err))
		return "", err
	}

	boldStyle := xlsx.NewStyle()
	boldStyle.Font.Bold = true

	headerRow := sheet.AddRow()
	headerTitles := []string{"Agent", "Host", "PlayerName", "DailyDepositAmount", "DailyDepositCount", column}
	for _, title := range headerTitles {
		cell := headerRow.AddCell()
		cell.Value = title
		cell.SetStyle(boldStyle)
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

	filename := fmt.Sprintf("%s_%s_%s.xlsx", time.Now().AddDate(0, 0, -1).Format("0102"), brand, fileField)

	// Save the file to the disk
	err = file.Save(filename)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Save failed:: %s", err))
		return "", err
	}

	log.LogInfo("File saved successfully.")
	return filename, nil
}

func Zipfiles(zipFileName string, fileToZip []string, password string) (string, error) {

	newZipFile, err := os.Create(zipFileName)
	if err != nil {
		log.LogFatal(err.Error())
		panic(err)
	}
	defer newZipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, filename := range fileToZip {
		addFileToZipWithPassword(zipWriter, filename, password)
	}
	return zipFileName, nil
}

func addFileToZipWithPassword(zipWriter *zip.Writer, filename string, password string) {
	// Open the file to be added to the zip file
	fileToZip, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fileToZip.Close()

	// Create a writer for the file entry in the zip file, with encryption
	writer, err := zipWriter.Encrypt(filename, password)
	if err != nil {
		panic(err)
	}

	// Copy the file content to the zip file
	if _, err = io.Copy(writer, fileToZip); err != nil {
		panic(err)
	}
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

func TelegramNotify(config util.MomoTelegramConfig, file string, message string) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Failed to create Telegram bot: %s", err))
		return err
	}

	bot.Debug = true
	chatID := config.ChatId
	msg := tgbotapi.NewDocumentUpload(chatID, file)
	msg.Caption = message

	if _, err := bot.Send(msg); err != nil {
		log.LogFatal(fmt.Sprintf("Failed to send message: %s", err))
		return err
	}
	return nil
}

func TelegramNotify2(config util.MomoTelegramConfig, file string, message string) error {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Failed to create Telegram bot: %s", err))
		return err
	}

	bot.Debug = true
	chatID := config.ChatId
	msg := tgbotapi.NewDocumentUpload(chatID, file)
	msg.Caption = message

	if _, err := bot.Send(msg); err != nil {
		log.LogFatal(fmt.Sprintf("Failed to send message: %s", err))
		return err
	}
	return nil
}

func deleteFiles() {
	patterns := []string{
		"./*.xlsx",
		"./*.zip",
	}

	for _, pattern := range patterns {
		// Use filepath.Glob to find all files that match the pattern
		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.LogFatal(err.Error())
		}

		// Loop through the matching files and delete them
		for _, match := range matches {
			err := os.Remove(match)
			if err != nil {
				log.LogInfo(fmt.Sprintf("Failed to delete %s: %s", match, err))
			} else {
				log.LogInfo(fmt.Sprintf("Deleted %s", match))
			}
		}
	}
}
