from models import SQLNullString
from os import system, name
from time import sleep

import requests

class Utils:
    @staticmethod
    def clear(): 
        if name == 'nt': 
            _ = system('cls') 
        else: 
            _ = system('clear') 

    @staticmethod
    def beautify_metadata(post_obj):
        meta_beautified = '''
        ID: {0}
        Title: {1}
        Author: {2}
        Category: {3}
        Posted at: {4}
        Topics: {5}
        '''
        print(meta_beautified.format(str(post_obj['post_id']), post_obj['title'], ', '.join(post_obj['authors']), post_obj['category']['String'], post_obj['date'], ' '.join(post_obj['topics'])))
    
    @staticmethod
    def beautify_group(group_obj):
        group_beautified = 'Group ID: {0}\nMembers: \n{1}\n'
        members = []
        for user in group_obj['users']:
            members.append('{0} - {1}'.format(user['user_id'], user['full_name']))
        print(group_beautified.format(str(group_obj['group_id']), '\n'.join(members)))
    
    @staticmethod
    def print_group(dictionary):
        """Print an object list and returns user's selections.
        """
        Utils.clear()
        keys = list(dictionary.keys())
        if len(keys) <= 10:
            Utils.clear()
            Utils.print_group_detail(keys, dictionary)
            user_input = input('\nPlease select one by entering the group id, or quit by entering [:q]: ')
            return user_input
        else:
            page_number = 0
            user_input = None
            
            while user_input != ':q':
                l = keys[page_number*10:(page_number+1)*10]
                print('Please select one by entering the group id:')
                if 0 < page_number < int(len(keys) / 10):
                    Utils.print_group_detail(l, dictionary)
                    user_input = input('next page: [d], previous page: [a], quit: [:q]\n')
                    Utils.clear()
                    if not user_input.isdigit() and (user_input == 'a' or user_input == 'd') :
                        page_number = Utils.return_page_number(page_number, user_input)
                    elif user_input.isdigit():
                        return dictionary[user_input]['group_id']
                    else:
                        print('Please enter a valid input!\n')

                elif page_number == 0:
                    Utils.print_group_detail(l, dictionary)
                    user_input = input('next page: [d], quit: [:q]\n')
                    Utils.clear()
                    if user_input == 'd' :
                        page_number = Utils.return_page_number(page_number, user_input)
                    elif user_input.isdigit():
                        return dictionary[user_input]['group_id']
                    else:
                        print('Please enter a valid input!\n')

                elif page_number != -1:
                    Utils.print_group_detail(l, dictionary)
                    user_input = input('previous page: [a], quit: [:q]\n')
                    Utils.clear()
                    if user_input == 'a':
                        page_number = Utils.return_page_number(page_number, user_input)
                    elif user_input.isdigit():
                        return dictionary[user_input]['group_id']
                    else:
                        print('Please enter a valid input!\n')
                else:
                    break   

    @staticmethod
    def handle_status_code(status_code):
        if 200 <= status_code < 400:
            print('Operation successful.\n')
            sleep(2)
        else:
            print('Operation unsuccessful... please try again later.\n')
            sleep(2)

    @staticmethod
    def print_post_content(user_instance, id, dictionary):
        print(dictionary[id]['content'])
        print('\nThumbs: {0}\n'.format(user_instance.get_thumbs(id)))
        thumb = input('\nThumb Up? Enter [y] or [n] to thumb up or down. Enter any other key to cancel: ')
        if thumb == 'y' or thumb == 'n':
            Utils.handle_status_code(user_instance.give_thumb(id, thumb))
        for topic in dictionary[id]['topics']:
            if len(topic) > 0:
                follow_topic = input('\nFollow topic [{0}]? (y/n)'.format(topic))
                if follow_topic == 'y':
                    Utils.handle_status_code(user_instance.follow_topic(topic))
    
    @staticmethod
    def read_posts(user_instance, dictionary):
        """Print an object list and returns user's selections.
        """
        Utils.clear()
        keys = list(dictionary.keys())
        if len(keys) <= 10:
            Utils.clear()
            Utils.print_post_metadata(keys, dictionary)
            user_post_id = input('\nPlease select one by entering the post id of the post, or quit by entering [:q]: ')
            while user_post_id != ':q':
                Utils.print_post_content(user_instance, user_post_id, dictionary)
                user_post_id = input('Please select one by entering the post id of the post, or quit by entering [:q]: ')
        else:
            page_number = 0
            user_post_id = None
            
            while user_post_id != ':q':
                l = keys[page_number*10:(page_number+1)*10]
                print('Please select one by entering the post id of the post:')
                if 0 < page_number < int(len(keys) / 10):
                    Utils.print_post_metadata(l, dictionary)
                    user_post_id = input('next page: [d], previous page: [a], quit: [:q]\n')
                    Utils.clear()
                    if user_post_id.isdigit():
                        Utils.print_post_content(user_instance, user_post_id, dictionary)
                    else:
                        page_number = Utils.return_page_number(page_number, user_post_id)

                elif page_number == 0:
                    Utils.print_post_metadata(l, dictionary)
                    user_post_id = input('next page: [d], quit: [:q]\n')
                    Utils.clear()
                    if user_post_id.isdigit():
                        Utils.print_post_content(user_instance, user_post_id, dictionary)
                    elif not user_post_id.isdigit():
                        page_number = Utils.return_page_number(page_number, user_post_id)
                    else:
                        print('Please enter a valid input!\n')

                elif page_number != -1:
                    Utils.print_post_metadata(l, dictionary)
                    user_post_id = input('previous page: [a], quit: [:q]\n')
                    Utils.clear()
                    if user_post_id.isdigit():
                        Utils.print_post_content(user_instance, user_post_id, dictionary)
                    elif not user_post_id.isdigit():
                        page_number = Utils.return_page_number(page_number, user_post_id)
                    else:
                        print('Please enter a valid input!\n')
                else:
                    break    
                sleep(3)

    @staticmethod
    def print_post_metadata(keys, dictionary):
        """Ranging over all of the posts.
        """
        for index, key in enumerate(keys):
            Utils.beautify_metadata(dictionary[key])
    
    @staticmethod
    def print_group_detail(keys, dictionary):
        """Ranging over all of the groups.
        """
        for index, key in enumerate(keys):
            Utils.beautify_group(dictionary[key])


    @staticmethod
    def print(obj_list, property=''):
        """Print an object list and returns user's selections.
        """
        user_selection = set()
        if len(obj_list) <= 10:
            Utils.clear()
            Utils.print_list(obj_list)
            user_input = input('Please select one by entering the index of an item, or quit by entering [:q]: ')
            while user_input != ':q':
                user_selection.add(SQLNullString(obj_list[int(user_input)]))
                if property == 'categories' and len(user_selection) > 0:
                    break
                user_input = input('Please select one by entering the index of an item, or quit by entering [:q]: ')
        else:
            page_number = 0
            user_input = None
            
            while user_input != ':q':
                l = obj_list[page_number*10:(page_number+1)*10]
                Utils.clear()
                print('Please select one by entering the index of an item.')
                if 0 < page_number < int(len(obj_list) / 10):
                    Utils.print_list(l)
                    user_input = input('next page: [d], previous page: [a], quit: [:q]\n')
                    if user_input.isdigit() and 0 <= int(user_input) <= 9:
                        user_selection.add(SQLNullString(l[int(user_input)]))
                    else:
                        page_number = Utils.return_page_number(page_number, user_input)
                    if property == 'categories' and len(user_selection) > 0:
                        break
                elif page_number == 0:
                    Utils.print_list(l)
                    user_input = input('next page: [d], quit: [:q]\n')
                    if user_input.isdigit() and 0 <= int(user_input) <= 9:
                        user_selection.add(SQLNullString(l[int(user_input)]))
                    elif not user_input.isdigit():
                        page_number = Utils.return_page_number(page_number, user_input)
                    else:
                        print('Please enter a valid input!\n')
                    if property == 'categories' and len(user_selection) > 0:
                        break
                elif page_number != -1:
                    Utils.print_list(l)
                    user_input = input('previous page: [a], quit: [:q]\n')
                    if user_input.isdigit() and 0 <= int(user_input) <= 9:
                        user_selection.add(SQLNullString(l[int(user_input)]))
                    elif not user_input.isdigit():
                        page_number = Utils.return_page_number(page_number, user_input)
                    else:
                        print('Please enter a valid input!\n')
                    if property == 'categories' and len(user_selection) > 0:
                        break
                else:
                    break    
        return list(user_selection)

    @staticmethod
    def return_page_number(page_number, user_input):
        if user_input == 'a':
            return page_number - 1
        if user_input == 'd':
            return page_number + 1
        else:
            return -1

    @staticmethod
    def print_list(obj_list):
        for idx, obj in enumerate(obj_list):
                print('{0}. {1}\n'.format(str(idx), obj))
