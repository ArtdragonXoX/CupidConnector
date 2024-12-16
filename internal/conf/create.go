package conf

import (
	"CupidConnector/internal/utils"
	"bufio"
	"log"
	"os"
)

func CreateEnv() error {
	reader := bufio.NewReader(os.Stdin)
	for reader.Buffered() > 0 {
		_, err := reader.Discard(reader.Buffered())
		if err != nil {
			log.Println(err)
			break
		}
	}

	setDefault()

	err := saveEnv()
	if err != nil {
		return err
	}

	return nil
}

func saveEnv() error {
	log.Println("正在创建配置文件...")
	v, err := utils.IsFileExists(envPath)
	if err != nil {
		log.Println("创建配置文件失败！")
		return err
	}
	if v {
		err := os.Remove(envPath)
		if err != nil {
			log.Println("删除旧配置文件失败！")
			return err
		}
	}
	file, err := os.Create(envPath)
	if err != nil {
		log.Println("创建配置文件失败！")
		return err
	}
	defer file.Close()

	content := "STU_URL=" + defaultUrl + "\n" +
		"STU_USERNAME=" + Config.Username + "\n" +
		"STU_PASSWORD=" + Config.Password + "\n" +
		"AUTO_EXIT=" + Config.AutoExit + "\n"

	_, err = file.WriteString(content)
	if err != nil {
		log.Println("写入配置文件失败！")
		return err
	}

	log.Println("创建配置文件成功！")
	return nil
}