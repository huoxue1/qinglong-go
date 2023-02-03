package env

import (
	"github.com/huoxue1/qinglong-go/models"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"time"
)

var (
	DISABLESTATUS = 1
	ENABLESTATUS  = 0
)

func AddEnv(env *models.Envs) (int, error) {
	return models.AddEnv(env)
}

func QueryEnv(searchValue string) ([]*models.Envs, error) {
	return models.QueryEnv(searchValue)
}

func UpdateEnv(env *models.Envs) error {
	env1, _ := models.GetEnv(env.Id)
	env1.Name = env.Name
	env1.Value = env.Value
	env1.Remarks = env.Remarks
	env1.Updatedat = time.Now()
	return models.UpdateEnv(env1)
}

func GetEnv(id int) *models.Envs {
	env, err := models.GetEnv(id)
	if err != nil {
		return nil
	}
	return env
}

func DisableEnv(ids []int) error {
	for _, id := range ids {
		env, err := models.GetEnv(id)
		if err != nil {
			continue
		}
		env.Status = DISABLESTATUS
		err = models.UpdateEnv(env)
		if err != nil {
			return err
		}
	}
	return nil
}

func EnableEnv(ids []int) error {
	for _, id := range ids {
		env, err := models.GetEnv(id)
		if err != nil {
			continue
		}
		env.Status = ENABLESTATUS
		err = models.UpdateEnv(env)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteEnv(ids []int) error {
	for _, id := range ids {
		err := models.DeleteEnv(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadEnvFromDb() map[string]string {
	result := make(map[string]string, 0)
	envs, err := QueryEnv("")
	if err != nil {
		log.Errorln(err.Error())
		return result
	}
	for _, env := range envs {
		if env.Status == 1 {
			continue
		}
		if _, ok := result[env.Name]; ok {
			result[env.Name] = result[env.Name] + "&" + env.Value
		} else {
			result[env.Name] = env.Value
		}
	}
	return result
}

func LoadEnvFromFile(file string) map[string]string {
	result := make(map[string]string, 0)
	data, _ := os.ReadFile(file)
	compile := regexp.MustCompile(`export\s(.*?)="(.*?)"`)
	match := compile.FindAllStringSubmatch(string(data), -1)
	for _, i := range match {
		if _, ok := result[i[1]]; ok {
			result[i[1]] = result[i[1]] + "&" + i[2]
		} else {
			result[i[1]] = i[2]
		}
	}
	return result
}

func GetALlEnv() map[string]string {
	envFromDb := LoadEnvFromDb()
	envfromFile := LoadEnvFromFile("data/config/config.sh")
	for s, s2 := range envfromFile {
		if _, ok := envFromDb[s]; ok {
			envFromDb[s] = envFromDb[s] + "&" + s2
		} else {
			envFromDb[s] = s2
		}
	}
	return envFromDb
}
