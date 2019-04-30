# -*- coding: utf-8 -*-

import argparse
import hashlib
import datetime
import json
import sys

import bs4
import requests

import config as CONFIG
import mysql_article

QIITA_URL = 'https://qiita.com/'

LOCAL = 'local'
PRODUCT = 'product'
DB_CONFIG = None

def calc_hash(URL):
    return hashlib.sha1(URL.encode()).hexdigest()

def insert_articles(articles, now):
    for title, url in articles:
        hash_ = calc_hash(url)
        with mysql_article.Database_article(DB_CONFIG) as db:
            article_id = db.insert_article(title, url, hash_, now)


def get_qiita_daily_trend():
    response = requests.get(QIITA_URL)
    response_html = bs4.BeautifulSoup(response.text, "html.parser")

    article_block = response_html.find('div', class_='p-home_main mb-3 mr-0@s').find('div')
    articles = article_block['data-hyperapp-props']
    article_json = json.loads(articles)
    articles = []
    for aj in article_json['trend']['edges']:
        title = aj['node']['title']
        URL = QIITA_URL + aj['node']['author']['urlName'] + '/items/' + aj['node']['uuid']
        articles.append([title, URL])
    return articles


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='')
    parser.add_argument('-e', '--env')
    args = parser.parse_args()
    print('env:', args.env)
    if args.env == LOCAL:
        DB_CONFIG = CONFIG.DB_CONFIG_LOCAL
    elif args.env == PRODUCT:
        DB_CONFIG = CONFIG.DB_CONFIG
    else:
        print('need "--env" option')
        sys.exit()

    now = int(datetime.datetime.now().timestamp())
    
    # qiita
    articles = get_qiita_daily_trend()
    insert_articles(articles, now)

    # gigazine

    # techable

    # google おすすめ