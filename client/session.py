import sys

from time import sleep
from user import User
from utils import Utils

class Session:
    def __init__(self, base_url='http://127.0.0.1:8080/'):
        self.base_url = base_url

    def sign_up(self):
        first_name = input('Your first name: ')
        last_name = input('Your last name: ')
        user_id = input('Your user id (no space): ')
        pwd = ''
        pwd_again = ' '
        while pwd != pwd_again:
            pwd = input('Please enter your password: ')
            pwd_again = input('Please confirm your password: ')
            if pwd != pwd_again:
                print('Your inputs do not match! Please re-enter your password.\n')
        Utils.handle_status_code(User.sign_up(self.base_url, user_id, first_name, last_name, pwd))

    def start(self):
        Utils.clear()
        start_up_panel = '''
        Welcome to FakeBook!
        1. Log In
        2. Sign Up
        3. Exit\n'''
        while True:
            user = None
            Utils.clear()
            client_opt = input(start_up_panel)
            print(client_opt)
            if client_opt == '1':
                while not user:
                    Utils.clear()
                    user_name = input("Username:")
                    password = input("Password:")
                    user = User.log_in(self.base_url, user_name, password)
                    Utils.clear()
                    user.run()
            elif client_opt == '2':
                self.sign_up()
            elif client_opt == '3':
                Utils.clear()
                print('See you next time!\n')
                sleep(1)
                Utils.clear()
                sys.exit(0)
                