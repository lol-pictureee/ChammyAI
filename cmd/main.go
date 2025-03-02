package main

import (
	"bot_tg/logger"
	"bot_tg/telegramBot"
)

func main() {
	logger.LogS()
	telegramBot.Run_bot()

	//requestID, errRequestID := nicefish.RequestID("api/nicefish/image1.jpg", "api/nicefish/image2.jpg")
	//if errRequestID != nil {
	//	log.Fatalf("Ошибка при сохранении файла: %v", errRequestID)
	//	return
	//
	//}
	//
	//fmt.Println(requestID)
	//
	//time.Sleep(5 * time.Second)
	//downloadUrl, errDownload := nicefish.DownloadPhoto(requestID)
	//if errDownload != nil {
	//	log.Fatalf("Ошибка при сохранении файла: %v", errDownload)
	//	return
	//
	//}
	//
	//fmt.Println(downloadUrl)
	//resp, err := http.Get(downloadUrl)
	//if err != nil {
	//	log.Fatalf("Ошибка при сохранении файла: %v", err)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	log.Fatalf("Статус не ОК")
	//	return
	//}
	//
	//filename := uuid.NewString() + ".jpg"
	//
	//file, err := os.Create(filename)
	//if err != nil {
	//	log.Fatalf("Ошибка при создании файла: %v", err)
	//	return
	//}
	//defer file.Close()
	//
	//_, err = io.Copy(file, resp.Body)
	//if err != nil {
	//	log.Fatalf("Ошибка при сохранении файла: %v", err)
	//	return
	//}
	//
	//fmt.Println("Файл успешно сохранен:", filename)

}
