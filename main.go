package main

func main() {
	// SendRandomDataToInflux("http://localhost:8086/", "telegraf")
	ServeRandomMetrics(":8888")
}
