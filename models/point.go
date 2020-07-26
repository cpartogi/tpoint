package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cpartogi/tpoint/config"
	db "github.com/cpartogi/tpoint/dbase"
	"github.com/cpartogi/tpoint/producer"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Pointlog struct {
	Id           uuid.UUID `json:"id"`
	Member_id    int       `json:"member_id"`
	Point_type   int       `json:"point_type"`
	Point_desc   string    `json:"point_desc"`
	Point_before int       `json:"point_before"`
	Point_amount int       `json:"point_amount"`
	Created_by   string    `json:"created_by"`
	Created_at   string    `json:"created_at"`
}

type Totalpoint struct {
	Member_id   int `json:"member_id"`
	Total_point int `json:"total_point"`
}

var err error
var res Response
var ctx context.Context

func getKafkaConfig(username, password string) *sarama.Config {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.WriteTimeout = 5 * time.Second
	kafkaConfig.Producer.Retry.Max = 0

	if username != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
	}
	return kafkaConfig
}

func GetPoint(member_id int) (Response, error) {
	db.Init()
	con := db.CreateCon()
	var total_point int

	sqlselect := "SELECT total_point from tb_member where member_id=?"

	row := con.QueryRow(sqlselect, member_id)

	switch err := row.Scan(&total_point); err {
	case sql.ErrNoRows:
		res.Status = http.StatusNoContent
		res.Message = "Member Not Found"
		res.Data = ""
	case nil:
		res.Status = http.StatusOK
		res.Message = "Success"
		res.Data = map[string]int{
			"member_id":   member_id,
			"total_point": total_point,
		}
	default:
		panic(err)
	}

	return res, nil
}

func AddPoint(insert_id uuid.UUID, member_id int, point_type int, point_desc string, point_amount int, created_by string, created_at string) (Response, error) {

	//kirim producer
	point_log := Pointlog{
		Id:           insert_id,
		Member_id:    member_id,
		Point_type:   point_type,
		Point_desc:   point_desc,
		Point_amount: point_amount,
		Created_by:   created_by,
		Created_at:   created_at,
	}

	var jsonData []byte
	var pesan string
	jsonData, err := json.Marshal(point_log)

	if err != nil {
		logrus.Println(err)
	}

	pesan = string(jsonData)

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	kafkaConfig := getKafkaConfig("", "")
	producers, err := sarama.NewSyncProducer([]string{"kafka:9092"}, kafkaConfig)
	if err != nil {
		logrus.Errorf("Unable to create kafka producer got error %v", err)
	}
	defer func() {
		if err := producers.Close(); err != nil {
			logrus.Errorf("Unable to stop kafka producer: %v", err)
			return
		}
	}()

	logrus.Infof("Success create kafka sync-producer")

	kafka := &producer.KafkaProducer{
		Producer: producers,
	}

	conf := config.GetConfig()

	err = kafka.SendMessage(conf.KAFKA_TOPIC, pesan)
	if err != nil {
		panic(err)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int64{
		"result_affected": 1,
	}

	return res, nil
}

func UpdatePoint(insert_id uuid.UUID, member_id int, point_type int, point_desc string, point_amount int, created_by string, created_at string) (Response, error) {
	db.Init()
	con := db.CreateCon()
	var point_before int
	var rowAffected int

	tx, err := con.Begin()

	if err != nil {
		logrus.Fatal(err)
	}

	defer tx.Rollback()

	rowAffected = 1

	sqlselect := "SELECT total_point from tb_member where member_id=? FOR UPDATE"

	row := con.QueryRow(sqlselect, member_id)

	switch err := row.Scan(&point_before); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		rowAffected = 0
	case nil:
		logrus.Println("Update point for member_id ", member_id)
	default:
		panic(err)
	}

	if rowAffected != 0 {
		sqlupdate := "UPDATE tb_member set total_point=total_point+? where member_id=?"

		stmt, err := tx.Prepare(sqlupdate)

		defer stmt.Close()

		if _, err := stmt.Exec(point_amount, member_id); err != nil {
			logrus.Fatal("gagal update", err)
		}

		sqlinsert := "INSERT INTO tb_point_log (id, member_id, point_type, point_desc, point_before, point_amount, created_by, created_at) values (?,?,?,?,?,?,?,?)"

		_, err = con.Exec(sqlinsert, insert_id, member_id, point_type, point_desc, point_before, point_amount, created_by, created_at)

		if err != nil {
			fmt.Println(err.Error())
			return Response{}, nil
		}

		if err := tx.Commit(); err != nil {
			logrus.Fatal(err)
		}
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int{
		"result_affected": rowAffected,
	}
	return res, nil
}
