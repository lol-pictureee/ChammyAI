package telegramBot

import (
	"bot_tg/api/fusebrain"
	"bot_tg/api/nicefish"
	"bot_tg/api/quotes"
	"bot_tg/envhandler"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

func Run_bot() {
	runtime.GOMAXPROCS(4)
	cfg, errCfg := envhandler.LoadConfig()
	if errCfg != nil {
		log.Fatal(errCfg)
	}

	bot, errToken := tgbotapi.NewBotAPI(cfg.BotToken)
	if errToken != nil {
		log.Panic(errToken)
	}

	bot.Debug = true
	slog.Info(
		"Авторизован",
		"username", bot.Self.UserName,
	)
	//log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go func(update tgbotapi.Update) {
			if update.Message.Chat != nil {

				userID := strconv.FormatInt(update.Message.From.ID, 10)

				dir := filepath.Join("users", userID)
				err := os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					log.Println("Error creating directory:", err)
					return
				}

				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, это Чамик AI. Отправь мне промт (текст) или фото кем я стану :))")
						bot.Send(msg)
					default:
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Такой команды нет, чо по парам?")
						bot.Send(msg)
					}
					return
				}

				if update.Message.Text != "" {
					userPrompt := update.Message.Text
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Генерация поставлена в очередь.")

					_, errSend := bot.Send(msg)
					if errSend != nil {
						log.Printf("Error sending message: %v", errSend)
						return
					}
					imagePath, errRequest := fusebrain.Request(userPrompt, userID)

					if errRequest != nil {
						log.Printf("Ошибка генерации изображения для пользователя %s: %v", userID, errRequest)

						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка при генерации изображения. Попробуйте еще раз.")
						_, errRequest = bot.Send(errMsg)
						if errRequest != nil {
							log.Printf("Ошибка отправки сообщения об ошибке пользователю %s: %v", userID, errRequest)
							return
						}
						return
					}

					if imagePath != "" {

						errWrite := WriteMessageToFile(update, dir)
						if errWrite != nil {
							log.Println("Error writing to file:", errWrite)
						}

						//urlString := urlupload.RunUrlUpload(imagePath)

						//fmt.Println(urlString) // сразу в SwapService

						requestID, errRequestID := nicefish.RequestID(imagePath, "api/nicefish/image2.jpg", userID)
						if errRequestID != nil {
							errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сгенерировать изображение.")
							_, errSend = bot.Send(errMsg)
							log.Printf("Ошибка при сохранении файла: %v", errRequestID)
							return
						}

						fmt.Println(requestID, userID)

						time.Sleep(10 * time.Second)
						downloadUrl, errDownload := nicefish.DownloadPhoto(requestID)
						if errDownload != nil {
							errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сгенерировать изображение.")
							_, errSend = bot.Send(errMsg)
							log.Printf("Ошибка при сохранении файла requestID: %v", errDownload)
							return

						}

						fmt.Println(downloadUrl)
						resp, err := http.Get(downloadUrl)
						if err != nil {
							errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сгенерировать изображение.")
							_, errSend = bot.Send(errMsg)
							log.Printf("Ошибка при сохранении файла HTTP.get: %v", err)
							return
						}
						defer resp.Body.Close()

						if resp.StatusCode != http.StatusOK {
							log.Printf("Статус не ОК")
							return
						}

						filename := uuid.NewString() + ".jpg"

						file, err := os.Create(filename)
						if err != nil {
							log.Printf("Ошибка при создании файла: %v", err)
							return
						}
						defer file.Close()

						_, err = io.Copy(file, resp.Body)
						if err != nil {
							log.Printf("Ошибка при сохранении файла: %v", err)
							return
						}

						fmt.Println("Файл успешно сохранен:", filename)

						//segmind.RequestSwapService("https://sun9-33.userapi.com/impg/Lq4jz5yFOPA5R7hR7v-R0ft67Ev36ydvrhjJEg/Y7ojl0HLDDc.jpg?size=1024x1024&quality=95&sign=7bffdefd5e0cf653db1b99ed651bd315&type=album", urlString) // КОНЧИЛИСЬ КРЕДИТЫ

						//photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(imagePath))
						//resultImagePath := "output_image.jpeg" // ---> в папку того чат ID, от которого запрос
						photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(filename))
						quote, errGetQuote := quotes.GetQuote()
						if errGetQuote != nil {
							fmt.Println(errGetQuote)
							return
						}
						photo.Caption = quote

						_, errSend = bot.Send(photo)
						if errSend != nil {
							log.Printf("Error sending photo: %v", errSend)
						}

					} else {

						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сгенерировать изображение.")

						_, errSend = bot.Send(errMsg)
						if errSend != nil {
							log.Printf("Error sending message: %v", errSend)
							return
						}
					}
				}

				if update.Message.Photo != nil {

					userID := fmt.Sprintf("%d", update.Message.From.ID)

					photo := update.Message.Photo[len(update.Message.Photo)-1]

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Обрабатываю изображение...")
					_, errSend := bot.Send(msg)
					if errSend != nil {
						log.Printf("Ошибка отправки сообщения в бота: %v", errSend)
						return
					}

					userDir := filepath.Join("users", userID)

					if _, err := os.Stat(userDir); os.IsNotExist(err) {
						if err := os.MkdirAll(userDir, 0755); err != nil {
							log.Printf("Ошибка создания директории пользователя: %v", err)
							errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обработать изображение (ошибка при создании директории).")
							bot.Send(errMsg)
							return
						}
					}

					uuidGen := uuid.NewString()
					imageName := fmt.Sprintf("%s.jpg", uuidGen)
					outputPath := filepath.Join(userDir, imageName)

					fileConfig := tgbotapi.FileConfig{FileID: photo.FileID}
					file, err := bot.GetFile(fileConfig)
					if err != nil {
						log.Printf("Ошибка получения информации о файле: %v", err)
						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обработать изображение (ошибка получения файла).")
						bot.Send(errMsg)
						return
					}

					fileURL := file.Link(bot.Token)
					fmt.Println(fileURL)
					fmt.Println(outputPath)
					client := &http.Client{}
					resp, err := client.Get(fileURL)
					if err != nil {
						log.Printf("Ошибка скачивания файла: %v", err)
						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обработать изображение (ошибка скачивания).")
						bot.Send(errMsg)
						return
					}
					defer resp.Body.Close()

					outFile, err := os.Create(outputPath)
					if err != nil {
						log.Printf("Ошибка создания файла: %v", err)
						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обработать изображение (ошибка создания файла для сохранения).")
						bot.Send(errMsg)
						return
					}
					defer outFile.Close()

					_, err = io.Copy(outFile, resp.Body)
					if err != nil {
						log.Printf("Ошибка копирования файла: %v", err)
						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обработать изображение (ошибка сохранения файла).")
						bot.Send(errMsg)
						return
					}

					msgSuccess := tgbotapi.NewMessage(update.Message.Chat.ID, "Фото скачано и сохранено.")
					bot.Send(msgSuccess)

					requestID, errRequestID := nicefish.RequestID(outputPath, "api/nicefish/image2.jpg", userID)

					if errRequestID != nil {
						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сгенерировать изображение.  Попробуйте позже.") // Более понятное сообщение
						_, errSend = bot.Send(errMsg)
						log.Printf("Ошибка при обработке файла nicefish.RequestID: %v", errRequestID)
						return
					}

					fmt.Println(requestID)

					fmt.Println(requestID, userID)

					time.Sleep(10 * time.Second)
					downloadUrl, errDownload := nicefish.DownloadPhoto(requestID)
					if errDownload != nil {
						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сгенерировать изображение.")
						_, errSend = bot.Send(errMsg)
						log.Printf("Ошибка при сохранении файла requestID: %v", errDownload)
						return

					}

					fmt.Println(downloadUrl)
					respUserPhoto, err := http.Get(downloadUrl)
					if err != nil {
						errMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сгенерировать изображение.")
						_, errSend = bot.Send(errMsg)
						log.Printf("Ошибка при сохранении файла HTTP.get: %v", err)
						return
					}
					defer respUserPhoto.Body.Close()

					if respUserPhoto.StatusCode != http.StatusOK {
						log.Printf("Статус не ОК")
						return
					}

					filename := uuid.NewString() + ".jpg"

					file2, err := os.Create(filename)
					if err != nil {
						log.Printf("Ошибка при создании файла: %v", err)
						return
					}
					defer file2.Close()

					_, err = io.Copy(file2, respUserPhoto.Body)
					if err != nil {
						log.Printf("Ошибка при сохранении файла: %v", err)
						return
					}

					fmt.Println("Файл успешно сохранен:", filename)
					photo2 := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(filename))
					_, errSend = bot.Send(photo2)
					if errSend != nil {
						log.Printf("Error sending photo: %v", errSend)
					}

				} else {

					if err != nil {
						log.Printf("Error sending message: %v", err)
						return

					}
					return
				}
			}
		}(update)
	}

}
