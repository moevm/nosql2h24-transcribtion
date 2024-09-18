#!/usr/bin/python3
import pprint
import pymongo

class Connection:
    def __init__(self, host='localhost', port=27017, username='root', password=''):
        self.host = host
        self.port = port
        self.username = username
        self.password = password

    def connect(self):
        uri = f'mongodb://{self.username}:{self.password}@{self.host}:{self.port}/'
        return pymongo.MongoClient(uri)

def main():
    database_name = "db_test"
    collection_name = "test"

    connection = Connection(password='example').connect()
    database = connection[database_name]
    collection = database[collection_name]
    collection.drop()

    while True:
        data = input('Input some shite:\n')
        if data == '':
            break
        print(collection.find_one({"_id": collection.insert_one({ 'input' : data }).inserted_id})['input'])

    print('\nThat\'s what we have here:')
    for post in collection.find():
        pprint.pprint(post)

    connection.close()


if __name__ == '__main__':
    main()
