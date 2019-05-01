# -*- coding: utf-8 -*-

import argparse
import datetime
import hashlib
import json
import os
import sys
import time

import bs4
import requests
from selenium import common, webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.chrome.options import Options

import config as CONFIG
import mysql_article

QIITA_URL = 'https://qiita.com/'
GIGAZINE_URL = 'https://gigazine.net/'
TECHABLE_URL = 'https://techable.jp/ranking'
GOOGLE_URL = 'https://news.google.com/foryou?hl=ja&gl=JP&ceid=JP%3Aja'

FAVICON = 'http://www.google.com/s2/favicons?domain='

LOCAL = 'local'
PRODUCT = 'product'
DB_CONFIG = None

def calc_hash(URL):
    return hashlib.sha1(URL.encode()).hexdigest()

def insert_articles(articles, now):
    for title, url, image_url in articles:
        hash_ = calc_hash(url)
        with mysql_article.Database_article(DB_CONFIG) as db:
            article_id = db.insert_article(title, url, image_url, hash_, now)


def get_qiita_daily_trend():
    response = requests.get(QIITA_URL)
    response_html = bs4.BeautifulSoup(response.text, "html.parser")

    article_block = response_html.find('div', class_='p-home_main mb-3 mr-0@s').find('div')
    articles = article_block['data-hyperapp-props']
    article_json = json.loads(articles)
    articles = []
    for aj in article_json['trend']['edges']:
        title = aj['node']['title']
        url = QIITA_URL + aj['node']['author']['urlName'] + '/items/' + aj['node']['uuid']
        image_url = FAVICON + url.split('/')[2]
        articles.append([title, url, image_url])
    return articles

def get_gigazine():
    print('get_gigazine')
    response = requests.get(GIGAZINE_URL)
    response_html = bs4.BeautifulSoup(response.text, "html.parser")

    article_blocks = response_html.find('div', class_='content').find_all('div', class_='card')
    articles = []
    for article_block in article_blocks:
        title = article_block.find('h2').find('a')['title']
        url = article_block.find('h2').find('a')['href']
        try:
            image_url = article_block.find('img')['data-src']
        except KeyError:
            try:
                image_url = article_block.find('img')['src']
            except KeyError:
                image_url = '#'
        articles.append([title, url, image_url])
    return articles

def get_techable_ranking():
    print('get_techable_ranking')
    response = requests.get(TECHABLE_URL)
    response_html = bs4.BeautifulSoup(response.text, "html.parser")

    article_blocks = response_html.find('div', class_='te-article-tab-panel__panel row').find_all('a', class_='te-articles__list__item__hit')
    articles = []
    for article_block in article_blocks:
        title = article_block.find('h3', class_='te-articles__list__item__content__title').text
        url = article_block['href']
        image_url = '#'
        for img in article_block.find('div', class_='te-articles__list__item__thumb__img')['style'].split():
            if 'http' in img:
                image_url = img
        articles.append([title, url, image_url])
    return articles

def get_google_recommend():
    print('get_google_recommend')
    chrome_options = Options()
    #chrome_options.add_argument('--headless')
    chrome_options.add_argument('--no-sandbox')
    chrome_options.add_argument('--disable-gpu')
    driver = webdriver.Chrome(os.environ['CHOROMEDRIVER_PATH'], options=chrome_options)
    #driver = webdriver.Chrome(os.environ['CHOROMEDRIVER_PATH'])
    driver.set_page_load_timeout(5)
    driver.implicitly_wait(5)
    driver.get(GOOGLE_URL)

    # login
    element = driver.find_element_by_id('identifierId')
    element.send_keys(CONFIG.EMAIL, Keys.ENTER)
    element = driver.find_element_by_id('password').find_element_by_tag_name('input')
    element.send_keys(CONFIG.PASSWORD, Keys.ENTER)

    article_blocks = driver.find_elements_by_class_name('NiLAwe')
    articles = []
    for article_block in article_blocks:
        article_part = article_block.find_element_by_tag_name('h3')
        title = article_part.find_element_by_tag_name('a').text
        url = article_part.find_element_by_tag_name('a').get_attribute('href')
        try:
            image_url = article_block.find_element_by_tag_name('img').get_attribute('src')
        except common.exceptions.NoSuchElementException:
            image_url = '#'
        print(title)
        articles.append([title, url, image_url])
    driver.close()
    driver.quit()
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
    articles = get_gigazine()
    insert_articles(articles, now)

    # techable
    articles = get_techable_ranking()
    insert_articles(articles, now)

    # google おすすめ
    # articles = get_google_recommend()
    # insert_articles(articles, now)
