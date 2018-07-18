package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

func SendRandomDataToInflux(addr, dbname string) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: addr,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbname,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create 2h or data in the past
	walk := 1000.0
	for t := time.Now().Add(time.Duration(-2) * time.Hour); t.Before(time.Now()); t = t.Add(time.Second) {
		// Create a point and add to batch
		walk += rand.Float64() - 0.5
		tags := map[string]string{}
		fields := map[string]interface{}{
			"value": rand.Float64() * 10,
			"walk":  walk,
		}

		pt, err := client.NewPoint("random", tags, fields, t)
		fmt.Println(pt)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	for {
		newbp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  dbname,
			Precision: "s",
		})
		if err != nil {
			log.Fatal(err)
		}
		walk += rand.Float64() - 0.5
		tags := map[string]string{}
		fields := map[string]interface{}{
			"value": rand.Float64() * 10,
			"walk":  walk,
		}

		pt, err := client.NewPoint("random", tags, fields, time.Now())
		fmt.Println(pt)
		if err != nil {
			log.Fatal(err)
		}
		newbp.AddPoint(pt)
		// Write the batch
		if err := c.Write(newbp); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}

	// Close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}
