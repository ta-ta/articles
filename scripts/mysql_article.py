# coding: UTF-8

from database import Database

class Database_article(Database):
    def __init__(self, dbconfig):
        super(Database_article, self).__init__(dbconfig)

    # article
    def insert_article(self, title, url, image_url, hash_, created):
        """ article追加 """
        try:
            query = """
                INSERT IGNORE INTO article (title, url, image_url, hash, created)
                VALUES (?, ?, ?, ?, ?)
            """
            self.cursor.execute(query, [title, url, image_url, hash_, created])
            self.connect.commit()
        except:
            self.connect.rollback()
            raise
        article_id = self.cursor.lastrowid
        return article_id
