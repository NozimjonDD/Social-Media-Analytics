package repository

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"strings"
	"tgstat/internal/entity"
	"time"
)

var words = []string{"расмий", "канали", "билан", "учун", "бўйича", "реклама", "ҳақида", "бўлган", "қандай", "видеони", "янги", "kun.uz",
	"batafsil", "mobil", "ilova", "android", "yangilangan", "obuna", "bo‘ling", "daryo", "bilan", "uchun", "bo‘yicha", "reklama", "ma’lum", "yangi",
	"расмий", "rasmiy", "sport", "English", "live", "саҳифаларимизга", "обуна", "бўлинг", "батафсил", "билан", "бўлди", "қилди", "янги", "видео",
	"Каналга", "қўшилинг", "https://t.me/xabaruzofficial", "билан", "учун", "янги", "қилинди", "бўлди", "бўйича", "маълум", "қилди", "қўлга",
	"kun.uz", "telegram", "Instagram", "facebook", "официальные", "страницы", "года", "2024", "декабря", "2023", "января",
	"@gazetauz", "правах", "рекламы", "можно", "года", "будет", "также", "области", "году", "будут", "она", "они", "здесь",
	"Уланиш", "каналга", "@hudud24official", "батафсил", "билан", "учун", "ўқиш", "қандай", "бўйича", "бўлган", "куни", "ҳақида", "hudud24.uz",
	"батафсил", "@bugunrasmiy", "билан", "учун", "видео", "бўйича", "маълум", "бўлган", "эълон", "қўлга", "куни",
	"@hook_report", "будет", "будут", "также", "узбекистана", "узбекистане", "году", "года", "тысяч", "шавкат", "которые", "граждан",
	"Билан", "учун", "бўйича", "реклама", "қилиш", "янги", "2024", "мумкин", "ҳамда", "ушбу", "орқали", "жорий",
	"Бўлинг", "каналга", "@gazetauz_uzb", "обуна", "учун", "билан", "бўйича", "қилди", "эълон", "кўра", "нисбатан",
	"instagram", "илова", "android", "facebook", "tube", "telegram", "мобил", "батафсил", "билан", "бўйича", "куни", "учун",
	"rasmiy", "telegram", "facebook", "youtube", "sahifalarimiz", "batafsil", "web-site", "uchun", "bilan", "yangi", "so‘m", "kuni", "qilindi", "haqida", "o‘zbek",
}

var st = strings.Join(words, ",")

type PostRepository struct {
	db *bun.DB
}

func NewPostRepository(db *bun.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) GetAll(ctx context.Context, filter *entity.PostFilter) ([]*entity.Post, error) {
	posts := make([]*entity.Post, 0)
	q := r.db.NewSelect().Model(&posts).Where("text != ''")
	if filter.ChanelId != "" {
		q.Where("channel_id=?", filter.ChanelId)
	}
	if filter.Order != "" {
		q.OrderExpr(filter.Order)
	}
	if len(filter.ByWord) != 0 {
		q.Where("text ilike ?", fmt.Sprintf("%%%s%%", filter.ByWord))
	}
	q.Offset(filter.Offset).Limit(filter.Limit)
	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetAllForWeek(ctx context.Context, filter *entity.PostFilter) ([]*entity.Post, error) {
	posts := make([]*entity.Post, 0)
	now := time.Now()
	q := r.db.NewSelect().Model(&posts).
		Where("channel_id=?", filter.ChanelId).
		Where("text != ''").
		Where("DATE BETWEEN ? AND ?", now.Add(-time.Hour*168).Unix(), now.Unix()).
		Limit(filter.Limit).Offset(filter.Offset)

	if filter.Order != "" {
		q.OrderExpr(filter.Order)
	}
	if len(filter.ByWord) != 0 {
		q.Where("text ilike ?", fmt.Sprintf("%%%s%%", filter.ByWord))
	}

	if err := q.Scan(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetLastUpdatedDate(ctx context.Context, channelID int) (int64, error) {
	var lastDate int64
	err := r.db.QueryRow(`SELECT max(date) FROM "posts"`).Scan(lastDate)
	if err != nil {
		return 0, err
	}
	return 0, err
}

func (r *PostRepository) Create(ctx context.Context, posts []*entity.Post) error {
	_, err := r.db.NewInsert().Model(&posts).
		On("CONFLICT (id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) GetFrequencyWords(ctx context.Context, channelID int) ([]*entity.FrequencyWords, error) {
	q := `SELECT word, COUNT(*) as "count"
FROM (
    SELECT regexp_split_to_table(lower(text), '\s') as word
    FROM posts WHERE channel_id = ?
) AS words
WHERE word <> '' AND length(word) >= 4 AND word NOT IN %s
GROUP BY word
ORDER BY "count" DESC LIMIT 50`

	var ss strings.Builder
	for i, word := range words {
		if i == 0 {
			ss.WriteString("(")
		}
		ss.WriteString(fmt.Sprintf("'%s'", word))
		if len(words) != i+1 {
			ss.WriteString(",")
			continue
		}
		ss.WriteString(")")
	}

	q = fmt.Sprintf(q, ss.String())

	rows, err := r.db.QueryContext(ctx, q, channelID)
	if err != nil {
		return nil, err
	}
	//fmt.Println(st)

	var words []*entity.FrequencyWords
	for rows.Next() {
		var word entity.FrequencyWords
		if err := rows.Scan(&word.Word, &word.Count); err != nil {
			return nil, err
		}
		words = append(words, &word)
	}

	return words, nil
}

/*

SELECT word, COUNT(*) as "count"
FROM (
    SELECT regexp_split_to_table(lower(text), '\s') as word
    FROM posts WHERE channel_id = 4598
) AS words
WHERE word <> '' AND length(word) >= 4 AND word NOT IN ('kun.uz', 'расмий')
GROUP BY word
ORDER BY "count" DESC LIMIT 50

*/
