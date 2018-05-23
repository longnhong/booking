package mongodb

type mongoDBConfig struct {
	MaxPool int
	Path    string
	DBname  string
	DBUser  string
	DBPass  string
}

var MongoConfig = mongoDBConfig{}

func CheckAndInitServiceConnection() {
	if service.baseSession == nil {
		service.URL = MongoConfig.Path
		service.DbUser = MongoConfig.DBUser
		service.DbPass = MongoConfig.DBPass
		err := service.New()
		if err != nil {
			logDB.Errorf("disconnected from %s", MongoConfig.Path)
			panic(err)
		}
	}
}
