package models

type Apps struct {
	Id           int      `xorm:"pk autoincr INTEGER" json:"id,omitempty"`
	Name         string   `xorm:"TEXT unique" json:"name,omitempty"`
	Scopes       []string `xorm:"JSON" json:"scopes,omitempty"`
	ClientId     string   `xorm:"TEXT" json:"client_id,omitempty"`
	ClientSecret string   `xorm:"TEXT" json:"client_secret,omitempty"`
	Tokens       []string `xorm:"JSON" json:"tokens,omitempty"`
	Createdat    string   `xorm:"not null DATETIME created" json:"createdat,omitempty"`
	Updatedat    string   `xorm:"not null DATETIME updated" json:"updatedat,omitempty"`
}

func QueryApp() ([]*Apps, error) {
	apps := make([]*Apps, 0)
	session := engine.Table(new(Apps))
	err := session.Find(&apps)
	if err != nil {
		return nil, err
	}
	return apps, err
}

func AddApp(app *Apps) (int, error) {
	_, err := engine.Table(app).Insert(app)
	if err != nil {
		return 0, err
	}
	_, _ = engine.Where("name=?", app.Name).Get(app)
	return app.Id, err
}

func GetApp(id int) (*Apps, error) {
	app := new(Apps)
	_, err := engine.ID(id).Get(app)
	return app, err
}

func GetAppById(clientId string) (*Apps, error) {
	app := new(Apps)
	_, err := engine.Table(app).Where("client_id=?", clientId).Get(app)
	return app, err
}

func UpdateApp(app *Apps) error {
	_, err := engine.Table(app).ID(app.Id).AllCols().Update(app)
	return err
}

func DeleteApp(id int) error {
	_, err := engine.Table(new(Apps)).Delete(&Apps{Id: id})
	return err
}
