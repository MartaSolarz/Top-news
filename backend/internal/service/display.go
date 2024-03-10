package service

type DisplayNewsService struct {
	newsDB NewsDatabase
}

func NewDisplayNewsService(newsDB NewsDatabase) *DisplayNewsService {
	return &DisplayNewsService{
		newsDB: newsDB,
	}
}
