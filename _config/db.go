package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout           = 30
	connectionStringTemplate = "mongodb://%s:%s@%s"
)

type Resource struct {
	DB    *mongo.Database
	DBLog *mongo.Database
}

func CreateResource() (*Resource, error) {
	_ = godotenv.Load()
	var err error
	var client *mongo.Client
	var dbName string
	var connectionURI string
	dbName = os.Getenv("MONGODB_DB_NAME")
	connectionURI = os.Getenv("MONGODB_ENDPOINT")
	client, err = mongo.NewClient(
		options.Client().ApplyURI(connectionURI),
		options.Client().SetMinPoolSize(3),
		options.Client().SetMaxPoolSize(10),
		options.Client().SetMaxConnIdleTime(5*time.Minute),
	)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	_ = client.Connect(ctx)
	err = client.Ping(ctx, nil)
	if err != nil {
		lineNotifyAlert(err)
		return nil, err
	}
	color.Green("Connect database successfully")
	color.Green(connectionURI)

	return &Resource{DB: client.Database(dbName)}, nil
}

func (r *Resource) Close() {
	ctx, cancel := InitContext()
	defer cancel()

	if err := r.DB.Client().Disconnect(ctx); err != nil {
		color.Red("Close connection falure, Something wrong...")
		return
	}
	if r.DBLog != nil {
		if err := r.DBLog.Client().Disconnect(ctx); err != nil {
			color.Red("Close connection falure, Something wrong...")
			return
		}
	}

	color.Cyan("Close connection successfully")
}

func (r *Resource) CloseLog() {
	ctx, cancel := InitContext()
	defer cancel()

	if err := r.DB.Client().Disconnect(ctx); err != nil {
		return
	}
}

func InitContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	return ctx, cancel
}

func lineNotifyAlert(msg error) error {
	type reslinenoti struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	var response reslinenoti
	accesstoken := "pKahNdB8I8ifXmaeekmD66EbvXSgKMF5hBYUt5bgiNa"
	payload := fmt.Sprintf(
		"message= \n  API-SERVICE %s Error Connect DB \n  Error : %s \n", os.Getenv("SERVICE"), msg)

	URL := "https://notify-api.line.me/api/notify"
	Header := map[string][]string{
		"Authorization": {"Bearer " + accesstoken},
		"Content-Type":  {"application/x-www-form-urlencoded"},
	}
	if err := externalCALL(URL, "POST", Header, payload, &response); err != nil {
		fmt.Println("err_call_line:", err)
		return nil
	}
	return nil
}

func externalCALL(URL string, method string, headers map[string][]string, bodyPayload string, obj interface{}) error {
	payload := strings.NewReader(bodyPayload)
	client := &http.Client{}
	req, err := http.NewRequest(method, URL, payload)
	if err != nil {
		return err
	}
	req.Header = headers
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, obj)
	if res.StatusCode != 200 {
		return errors.New("CANNOT CALL API")
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, obj); err != nil {
		return err
	}
	return nil
}
