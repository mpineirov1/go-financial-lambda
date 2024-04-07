package handlers

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"github.com/mpineirov1/go-financial-lambda/models"
	"github.com/mpineirov1/go-financial-lambda/utils"
	"math"
	"strconv"
	"text/template"
	"time"

	"gopkg.in/gomail.v2"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func New() func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	}))
	
	svc := s3.New(sess)

	filename := "txns.csv"
	input := &s3.GetObjectInput{
		Bucket: aws.String("go-financial-lambda-csvs"), // Cambia al nombre de tu bucket
		Key:    aws.String(filename),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		log.Print(err)
	}
	defer result.Body.Close()

	csvReader := csv.NewReader(result.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Print(err)
	}

	transactionsByMonth := make(map[int][]models.Transaction)
	var totalBalance float64

	for _, eachrecord := range records[1:] {
		if len(eachrecord) < 3 {
			log.Print("Missing data in the csv: ", eachrecord)
			os.Exit(6) // Skip to next record
		}
		id, _ := strconv.Atoi(eachrecord[0])
		date, _ := time.Parse("1/2", eachrecord[1]) // Formato de fecha MM/DD
		transaction, _ := strconv.ParseFloat(eachrecord[2], 64)
		totalBalance = transaction + totalBalance
		month := date.Format("01")
		monthNumber, _ := strconv.Atoi(month)
		transactionModel := models.Transaction{
			ID:          id,
			Date:        date,
			Transaction: transaction,
			CreatedAt:   time.Now(),
		}
		transactionsByMonth[monthNumber] = append(transactionsByMonth[monthNumber], transactionModel)
	}

	summaryData := struct {
		TotalBalance float64
		MonthSummary map[string]models.MonthSummary
	}{
		TotalBalance: totalBalance,
		MonthSummary: make(map[string]models.MonthSummary), // Inicializar el mapa
	}

	for month, transactions := range transactionsByMonth {
		monthName := utils.GetMonthName(month)
		var monthCreditBalance, monthDebitBalance float64
		var monthCreditTotal, monthDebitTotal float64
		for _, t := range transactions {
			if t.Transaction > 0 {
				monthCreditTotal++
				monthCreditBalance += t.Transaction
			} else {
				monthDebitTotal++
				monthDebitBalance += t.Transaction
			}
		}

		monthCreditAvg := monthCreditBalance / monthCreditTotal
		monthDebitAvg := monthDebitBalance / monthDebitTotal
		if math.IsNaN(monthDebitAvg) {
			monthDebitAvg = 0
		}
		if math.IsNaN(monthCreditAvg) {
			monthCreditAvg = 0
		}
		monthSummary := models.MonthSummary{
			TransactionNumber: len(transactions),
			DebitAvg:          monthDebitAvg,
			CreditAvg:         monthCreditAvg,
		}
		summaryData.MonthSummary[monthName] = monthSummary
	}

	var body bytes.Buffer

	tmpl := template.Must(template.ParseFiles("templates/mail.html"))
	tmpl.Execute(&body, summaryData)
	

	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		req, _ := json.Marshal(request)
		log.Print(string(req))
		email, _ := request.QueryStringParameters["email"]
		sendMail(email, "Summary", body.String())
		resp := "Email sent!!"
		defer log.Print("Lambda complete")
		return events.APIGatewayProxyResponse{Body: resp, StatusCode: 200}, nil
	}
}

func sendMail(to, subject, body string) {

	from := utils.GoDotEnvVariable("MAIL_FROM_ADDRESS")
	username := utils.GoDotEnvVariable("MAIL_USERNAME")
	password := utils.GoDotEnvVariable("MAIL_PASSWORD")
	smtpHost := utils.GoDotEnvVariable("MAIL_HOST")
	smtpPort, _ := strconv.Atoi(utils.GoDotEnvVariable("MAIL_PORT"))

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, username, password)

	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}
}
