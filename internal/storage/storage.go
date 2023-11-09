package storage

var urldb map[string]string = make(map[string]string)

func GetUrldb() *map[string]string {

	return &urldb
}

// func SetUrldb(shortUrl, longUrl string) {
// 	urldb[shortUrl] = longUrl
// }

// func (sd *DataStorage) GetUrldb() map[string]string {

// 	// if sd.getUrldb() == nil {
// 	// 	sd.Urldb = map[string]string{}
// 	// }
// 	return sd.Urldb
// }
