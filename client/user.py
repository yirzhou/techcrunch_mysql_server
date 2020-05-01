from models import Post, SQLNullString
from time import sleep
from utils import Utils
import requests
import sys

default_server_url='http://127.0.0.1:8080/'

class User:
    def __init__(self, user_id, password, server_url=default_server_url):
        self.user_id = user_id
        self.password = password
        self.server_url = server_url

    def __get_category(self):
        Utils.clear()
        print("\nWe have the following categories for you:\n")
        categories = requests.get(self.get_route('categories')).json()[5:]
        return Utils.print(categories, 'categories')[0].sql_null_string()
    
    def __get_topics(self):
        user_topics = []
        print("\nWe have the following topics available for you:\n")
        topics = requests.get(self.get_route('topics')).json()[5:]
        topics_sql = Utils.print(topics)
        for topic_sql in topics_sql:
            user_topics.append(topic_sql.string)
        return user_topics
    
    def __get_tags(self):
        Utils.clear()
        user_tags = []
        tag = input('Please enter a tag for your post (type [:q] to exit): ')
        while tag != ':q':
            if tag != ':q':
                user_tags.append(tag.strip())
            else:
                break
            tag = input('\nPlease enter another tag for your post (type [:q] to exit): ')
        return user_tags
    
    def __get_content(self):
        Utils.clear()
        print('\nPlease enter the content of your post. Type [ctrl + d] to save and end your post:\n')
        content = sys.stdin.readlines()
        return ''.join(content)

    def __new_post(self, category, content, image_src=SQLNullString('').sql_null_string(), section=SQLNullString('').sql_null_string(), title="", url=SQLNullString('').sql_null_string(), tags=[], topics=[]):
        post_obj = Post(self.user_id, category, content, image_src, section, title, url, tags, topics)
        post_json = {
            'post_id': post_obj.post_id,
            'authors': [post_obj.author],
            'category': post_obj.category,
            'content': post_obj.content.strip(),
            'image_src': post_obj.image_src,
            'section': post_obj.section,
            'title': post_obj.title.strip(),
            'url': post_obj.url,
            'tags': post_obj.tags,
            'topics': post_obj.topics    
        }
        print(post_json)

        r = requests.post(url=self.get_route('posts/new'), params={'user_id':self.user_id}, json=post_json)
        return r.status_code

    def get_route(self, route):
        return self.server_url + route

    def get_posts(self):
        posts_dictionary = requests.get(url=self.get_route('posts')).json()
        Utils.read_posts(self, posts_dictionary)

    def get_topics(self):
        topics = requests.get(url=self.get_route('topics')).json()[5:]
        Utils.print(topics)

    def follow_topic(self, topic):
        return requests.put(self.get_route('users/{0}/topics/add'.format(self.user_id)), params={'topic':topic}).status_code

    def react_to_post(self, switch, post_id):
        action = ''
        if switch == 1:
            action = 'thumb_up'
        elif switch == -1:
            action = 'thumb_down'
        status_code = requests.post(self.get_route('posts/{0}'.format(str(post_id))), params={'user_id':self.user_id, 'action':action}).status_code
        Utils.handle_status_code(status_code)

    def get_authors(self):
        authors = requests.get(url=self.get_route('authors')).json()[5:]
        for author in authors:
            print("postID: [%d], authorID: [%s]\n" % (author['post_id'], author['author_id']))
        return authors

    def get_new_posts_since_last_login(self):
        Utils.clear()
        posts = requests.get(url=self.get_route('users/{0}/new_posts'.format(self.user_id))).json()
        if len(posts) == 0:
            print('Nothing new to show you! Maybe check again later.')
            sleep(2)
            return 
        Utils.read_posts(self, posts)
        
    def log_out(self):
        r = requests.post(url=self.get_route('users/logout'), data={'user_id':self.user_id, 'password':self.password})
        if r.status_code == 202:
            Utils.clear()
            print("See you next time, {0}!\n".format(self.user_id))
        else:
            Utils.clear()
            print("Oops! Something is wrong.")\

    def create_post(self):
        user_category = self.__get_category()
        user_topics = self.__get_topics()
        user_content = self.__get_content()
        user_img_src = SQLNullString(input('\nPlease enter the URL of your image:')).sql_null_string()
        user_title = input('\nPlease enter the title of your post:')

        user_tags = self.__get_tags()
        Utils.handle_status_code(self.__new_post(user_category, user_content, user_img_src, SQLNullString('').sql_null_string(), user_title, SQLNullString('').sql_null_string(), user_tags, user_topics))
        
    def join_group(self):
        groups_dict = requests.get(url=self.get_route('groups')).json()
        group_id = Utils.print_group(groups_dict)
        status_code = requests.put(self.get_route('groups/{0}/add'.format(str(group_id))), params={'user_id': self.user_id}).status_code
        Utils.handle_status_code(status_code)
    
    def create_group(self):
        status_code = requests.post(self.get_route('groups/new'), params={'user_id':self.user_id}).status_code
        if status_code >= 400:
            print('failed to create a new group... please try again later.')
            sleep(2)
        else:
            print('group created!')
            sleep(2)

    def give_thumb(self, post_id, thumb):
        if thumb == 'y':
            return requests.post(url=self.get_route('posts/{0}'.format(str(post_id))), params={'user_id': self.user_id, 'action':'thumbup'}).status_code
        elif thumb == 'n':
            return requests.post(url=self.get_route('posts/{0}'.format(str(post_id))), params={'user_id': self.user_id, 'action':'thumbdown'}).status_code
        return 404

    def get_thumbs(self, post_id):
        return requests.get(url=self.get_route('posts/{0}/thumbs'.format(str(post_id)))).json()['thumbs']
    
    def run(self):
        try:
            while True:
                Utils.clear()
                options = '''
                
                Please choose one of the following:\n
                1. Check new posts of the topics you like.    
                2. Join a group.
                3. Create a group.
                4. Just show me all of the posts!
                5. Write a new post!
                6. Log out

                Your choice: '''
                opt = input(options)

                if opt == '1':
                    self.get_new_posts_since_last_login()
                elif opt == '2':
                    self.join_group()
                elif opt == '3':
                    self.create_group()
                elif opt == '4':
                    self.get_posts()
                elif opt == '5':
                    self.create_post()
                elif opt == '6':
                    self.log_out()
                    break

        except KeyboardInterrupt:
            self.log_out()
            return

    @staticmethod
    def get_guest_route(root_url, route):
        return root_url + route

    @staticmethod
    def sign_up(root_url, user_id, first_name, last_name, password=None):
        return requests.post(url=User.get_guest_route(root_url, 'users/new'), data={'user_id':user_id, 'first_name':first_name, 'last_name':last_name, 'password':password}).status_code

    @staticmethod
    def log_in(root_url, user_id, password):
        try:
            r = requests.post(url=User.get_guest_route(root_url, 'users/login'), data={'user_id':user_id, 'password':password})
        
            if 200 <= r.status_code < 300:
                print("Welcome, {0}!".format(user_id))
                return User(user_id, password, root_url)
            else:
                Utils.clear()
                print('wrong combination of username and password.')
                return None
        except Exception:
            raise Exception('The server might be down. Maybe try again later!\n')
        sys.exit(1)
