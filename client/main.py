import requests
from user import User
from session import Session
from models import Post
from random import randrange

localhost = 'http://127.0.0.1:8080/'

def test_new_post(c):
    c.new_post({'String': 'YirenExclusive', 'Valid': True},
        'Test #' + str(randrange(65535)),  
        {'String': 'www.zhouyiren.com', 'Valid': True}, 
        {'String': 'section', 'Valid': False}, 
        "Foobar#" + str(randrange(65535)), 
        {'String': '', 'Valid': False},
        ['tag_1', 'tag_2'],
        ['test1']
    )

def test_log_in(user_id, password=None):
    return User.log_in(localhost, user_id, password)

def test_log_out(c):
    c.log_out()

def test_sign_up(user_id, password=None):
    User.sign_up(localhost, user_id, 'Larry', 'Zhou', password)

def main():
    session = Session()
    session.start()

if __name__ == '__main__':
    main()
    