package open

import (
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/utils"
)

func AddApp(apps *models.Apps) (int, error) {
	apps.ClientId = utils.RandomString(6)
	apps.ClientSecret = utils.RandomString(12)
	apps.Tokens = []string{}
	id, err := models.AddApp(apps)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateApp(apps *models.Apps) error {
	app, err := models.GetApp(apps.Id)
	if err != nil {
		return err
	}
	app.Name = apps.Name
	app.Scopes = apps.Scopes
	err = models.UpdateApp(app)
	if err != nil {
		return err
	}
	return nil
}

func ResetApp(apps *models.Apps) error {
	apps.ClientSecret = utils.RandomString(12)
	apps.Tokens = []string{}
	err := models.UpdateApp(apps)
	return err
}

func DeleteApp(ids []int) error {
	for _, id := range ids {
		err := models.DeleteApp(id)
		if err != nil {
			return err
		}
	}
	return nil
}
