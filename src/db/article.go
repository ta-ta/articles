package db

type Article struct {
	ArticleID int64  `db:"article_id"`
	Title     string `db:"title"`
	URL       string `db:"url"`
	Domain    string
	Kind      int64  `db:"kind"`
	Created   string `db:"created"`
	Priority  int64  `db:"priority"`
}

// FetchMasters masterを取ってくる
func (db Database) FetchArticles(read int64, created_order string) ([]*Article, error) {
	// 検索クエリ
	var articles []*Article
	query := "SELECT `id` as article_id, `title`, `url`, `kind`, DATE(FROM_UNIXTIME(`created`)) as created, `priority` " +
		"FROM `article` " +
		"WHERE kind = ? " +
		"ORDER BY `priority` desc, `created` " + created_order
	if err := db.database.Select(&articles, query, read); err != nil {
		return nil, err
	}
	return articles, nil
}

// UpdateRead 未読、既読を切り替える
func (db Database) UpdateRead(articleID, read int64) error {
	query := "UPDATE article SET kind = ? " +
		"WHERE id = ?"
	_, err := db.database.Exec(query, read, articleID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateProprity priorityを修正する
func (db Database) UpdateProprity(articleID, priority int64) error {
	query := "UPDATE article SET priority = ? " +
		"WHERE id = ?"
	_, err := db.database.Exec(query, priority, articleID)
	if err != nil {
		return err
	}
	return nil
}
