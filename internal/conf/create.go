package conf

import (
	"cupid-connector/internal/utils"
	"log"
	"os"
)

func CreateEnv() error {
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

	content := "BASE_URL=" + defaultBaseUrl + "\n" +
		"STU_USERNAME=" + Config.Username + "\n" +
		"STU_PASSWORD=" + Config.Password + "\n" +
		"AUTO_LOGIN=" + Config.AutoLogin + "\n" +
		"AUTO_EXIT=" + Config.AutoExit + "\n"

	_, err = file.WriteString(content)
	if err != nil {
		log.Println("写入配置文件失败！")
		return err
	}

	log.Println("创建配置文件成功！")
	return nil
}
