from datetime import date

class SQLNullString:
    def __init__(self, string):
        self.string = string.strip()
        if len(self.string) == 0:
            self.valid = False
        else:
            self.valid = True

    def sql_null_string(self):
        return {'String': self.string, 'Valid': self.valid}
    
    def string(self):
        return self.string

class Post:
    def __init__(self, author, category, content, image_src, section, title, url, tags=[], topics=[]):
        self.post_id = -1
        self.author = author
        self.category = category
        self.content = content
        self.date = str(date.today())
        self.image_src = image_src
        self.section = section
        self.title = title
        self.url = url
        self.tags = tags
        self.topics = topics
