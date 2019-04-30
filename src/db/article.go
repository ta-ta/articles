package db

type Dashboard struct {
	Date    string `db:"date_"`
	Created int64  `db:"created"`
	Read    int64  `db:"read_"`
}

type Article struct {
	ArticleID int64  `db:"article_id"`
	Title     string `db:"title"`
	URL       string `db:"url"`
	Domain    string
	Kind      int64  `db:"kind"`
	Created   string `db:"created"`
	Priority  int64  `db:"priority"`
}

// FetchDashboards Dashboardを取ってくる
func (db Database) FetchDashboards() ([]*Dashboard, error) {
	var dashboards []*Dashboard
	queryDate := "SELECT CURRENT_DATE() - INTERVAL num DAY AS date_ " +
		"FROM (SELECT @num := 0 AS num " +
		"UNION " +
		"SELECT @num := @num + 1 AS num " +
		"FROM information_schema.COLUMNS " +
		"LIMIT 14" +
		") t1"
	queryCreated := "SELECT DATE(FROM_UNIXTIME(`created`)) as date_, count(*) as created " +
		"FROM `article` " +
		"WHERE DATE(FROM_UNIXTIME(`created`)) between DATE_SUB(CURRENT_DATE(), INTERVAL 14 DAY) and DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY) " +
		"GROUP BY date_ "
	queryRead := "SELECT DATE(FROM_UNIXTIME(`read_mark`)) as date_, count(*) as read_ " +
		"FROM `article` " +
		"WHERE DATE(FROM_UNIXTIME(`read_mark`)) between DATE_SUB(CURRENT_DATE(), INTERVAL 14 DAY) and DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY) " +
		"GROUP BY date_ "
	query := "SELECT date_, COALESCE(created, 0) as created, COALESCE(read_, 0) as read_ " +
		"FROM (" + queryDate + ") as query_date " +
		"left join (" + queryCreated + ") as query_created " +
		"using (date_) " +
		"left join (" + queryRead + ") as query_read " +
		"using (date_) " +
		"ORDER BY date_ "
	if err := db.database.Select(&dashboards, query); err != nil {
		return nil, err
	}
	return dashboards, nil
}

// FetchUnread unreadを取ってくる
func (db Database) FetchUnread() ([]*int64, error) {
	var unreads []*int64
	query := "SELECT count(*) as unread " +
		"FROM `article` " +
		"WHERE kind = 0 " +
		"limit 1"
	if err := db.database.Select(&unreads, query); err != nil {
		return nil, err
	}
	return unreads, nil
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
func (db Database) UpdateRead(articleID, read, readMark int64) error {
	query := "UPDATE article SET kind = ?, read_mark = ? " +
		"WHERE id = ?"
	_, err := db.database.Exec(query, read, readMark, articleID)
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
