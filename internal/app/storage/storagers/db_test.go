package storagers

import (
	"context"
	"testing"

	"github.com/AlexTerra21/shortener/internal/app/utils"
)

func BenchmarkInsert(b *testing.B) {
	connString := "host=localhost user=shortner password=userpassword dbname=short_urls sslmode=disable"
	db := DB{}
	err := db.New(connString)
	if err != nil {
		b.Log(err.Error())
	}

	countTimes := b.N
	urls := make([]ShortenedURL, countTimes)

	for i := 0; i < countTimes; i++ {
		url := ShortenedURL{
			UUID:        1144,
			IdxShortURL: utils.RandSeq(10),
			OriginalURL: utils.RandSeq(15),
			DeletedFlag: false,
		}
		urls = append(urls, url)
	}

	ctx := context.Background()

	b.ResetTimer()

	b.Run("insert", func(b *testing.B) {
		for _, url := range urls {
			db.insertURL(ctx, url)
		}
	})

	b.Run("bulk_insert", func(b *testing.B) {
		db.insertURLs(ctx, &urls)
	})

	b.Run("delete", func(b *testing.B) {
		b.StopTimer()
		users := make([]UsersURL, countTimes)
		for _, url := range urls {
			user := UsersURL{
				UserID: url.UUID,
				URLID:  url.IdxShortURL,
			}
			users = append(users, user)
		}
		b.StartTimer()
		db.Delete(ctx, users)
	})

}
