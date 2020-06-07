import csv
import os
import sys

source_file = './../data/techcrunch_posts.csv'
tag_file = '../data/tags.csv'
topic_file = '../data/topics.csv'
author_file = '../data/authors.csv'

files = [
    tag_file, 
    topic_file,
    author_file
]

def prepare_file(file_path):
    """
    Returns a brand-new file object in "append mode".
    """
    if os.path.exists(file_path):
        os.remove(file_path)
    new_file = open(file_path, 'a+')
    return new_file

def prepare_load(sql_file, csv_file_path, table_name):
    """
    Write preparation statements for load data local infile.
    Returns:
    - .sql file object in append mode. 
    """
    # prepare for loading
    sql_file.write("SET foreign_key_checks = 0;\n")
    sql_file.write("SET GLOBAL local_infile = `On`;\n")
    sql_file.write("load data local infile '{0}'\n".format(csv_file_path))
    sql_file.write('into table `{0}`\n'.format(table_name))
    sql_file.write("fields terminated by ','\n")
    sql_file.write("enclosed by '{0}'\n".format('"'))
    sql_file.write("lines terminated by '{0}'\n".format('\\n'))
    sql_file.write('ignore 1 rows\n')
    return sql_file

def generate_sql_for_post_author_table(csv_file_path):
    """
    Generate assign statements for Moves table.
    """
    sql_file = prepare_file('../data/load_authors.sql')

    # prepare table
    sql_file.write('drop table if exists PostAuthor;\n')
    sql_file.write('create table PostAuthor(postID INT, authorID VARCHAR(100), primary key(postID, authorID));\n')

    # load from local infile
    sql_file = prepare_load(sql_file, csv_file_path, 'PostAuthor')
    sql_file.write("(@postID, @authorID)\n")
    sql_file.write("set\n   postID=@postID,\n   authorID=@authorID;\n")
    sql_file.close()

def generate_sql_for_post_tag_table(csv_file_path):
    """
    Generate assign statements for Moves table.
    """
    sql_file = prepare_file('../data/load_tags.sql')

    # prepare table
    sql_file.write('drop table if exists PostTag;\n')
    sql_file.write('create table PostTag(postID INT, tag VARCHAR(100), primary key(postID, tag));\n')

    # load from local infile
    sql_file = prepare_load(sql_file, csv_file_path, 'PostTag')
    sql_file.write("(@postID, @tag)\n")
    sql_file.write("set\n   postID=@postID,\n   tag=@tag;\n")
    sql_file.close()

def generate_sql_for_post_topic_table(csv_file_path):
    """
    Generate assign statements for Moves table.
    """
    sql_file = prepare_file('../data/load_topics.sql')

    # prepare table
    sql_file.write('drop table if exists PostTopic;\n')
    sql_file.write('create table PostTopic(postID INT, topic VARCHAR(100), primary key(postID, topic));\n')

    # load from local infile
    sql_file = prepare_load(sql_file, csv_file_path, 'PostTopic')
    sql_file.write("(@postID, @topic)\n")
    sql_file.write("set\n   postID=@postID,\n   topic=@topic;\n")
    sql_file.close()


def check_source_file(path):
    if not os.path.exists(path):
        raise OSError("OSError: File [%s] does not exist!\n", path)

def clear_file(paths):
    for path in paths:
        if os.path.exists(path):
            os.remove(path)

def create_csv_source_categories(dictionary):
    csv_file_categories = open('../data/categories.csv', 'a+')
    csv_file_categories.write('postID,category\n')
    for postId in dictionary:
        if 'category' in dictionary[postId]:
            csv_file_categories.write(postId + ',' + dictionary[postId]['category'] + '\n')
    csv_file_categories.close()

def create_csv_source_authors(dictionary):
    csv_file_authors = open('../data/authors.csv', 'a+')
    csv_file_authors.write('postID,authorID\n')
    for postId in dictionary:
        for author in dictionary[postId]['authors']:
            csv_file_authors.write(postId + ',' + '.'.join(author.split()) + '\n')
    csv_file_authors.close()

def create_csv_source_tags(dictionary):
    csv_file_tags = open('../data/tags.csv', 'a+')
    csv_file_tags.write('postID,tag\n')
    for postId in dictionary:
        for tag in dictionary[postId]['tags']:
            if len(tag) > 0:
                csv_file_tags.write(postId + ',' + tag + '\n')    
    csv_file_tags.close()

def create_csv_source_topics(dictionary):
    csv_file_topics = open('../data/topics.csv', 'a+')
    csv_file_topics.write('id,topic\n')

    for postId in dictionary: 
        for topic in dictionary[postId]['topics']:
            if len(topic) > 0:
                csv_file_topics.write(postId + ',' + topic + '\n')
    csv_file_topics.close()

def extract_post_info(source_file):
    csv.field_size_limit(sys.maxsize)
    post_info = {}
    try:
        check_source_file(source_file)
    except Exception as e:
        print("Error [%s] occurred when reading [%s]\n", str(e), source_file)

    with open(source_file) as master_file:
        csv_reader = csv.reader(master_file)
        next(csv_reader)
        for row in csv_reader:
            authors, category, content, date, postId, img_src, section, tags, title, topics, url = row
            tag_list = tags.split(',')
            author_list = authors.split(',')
            topic_list = topics.split(',')
            # add tags, authors, topics to post
            if postId not in post_info:
                post_info[postId] = {}
                post_info[postId]['authors'] = set()
                post_info[postId]['tags'] = set()
                post_info[postId]['topics'] = set()
            for tag in tag_list:
                if len(tag) > 0 and tag not in post_info[postId]['tags']:
                    post_info[postId]['tags'].add(tag)
            for author in author_list:
                if len(author) > 0 and author not in post_info[postId]['authors']:
                    post_info[postId]['authors'].add(author)
            for topic in topic_list:
                if len(topic) > 0 and topic not in post_info[postId]['topics']:
                    post_info[postId]['topics'].add(topic)
            if len(category) > 0:
                post_info[postId]['category'] = category
    create_csv_source_categories(post_info)
        
def main():
    # clear_file(files)
    extract_post_info(source_file)
    # generate_sql_for_post_author_table(author_file)
    # generate_sql_for_post_tag_table(tag_file)
    # generate_sql_for_post_topic_table(topic_file)

if __name__ == "__main__":
    main() 
