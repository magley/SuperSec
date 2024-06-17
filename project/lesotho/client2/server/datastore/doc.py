import json
import bcrypt

class DocRepo:
    def __init__(self):
        self.salt = bcrypt.gensalt()
        self.data = []
        self.id_counter = 0
        self.fname = './data/doc.json'
        self.load()

    def load(self):
        try:
            with open(self.fname) as f:
                data_json = f.read()
                self.data = json.loads(data_json)
                for e in self.data:
                    if e['id'] > self.id_counter:
                        self.id_counter = e['id'] 
        except FileNotFoundError:
            self.flush()

    def flush(self):
        with open(self.fname, 'w') as f:
            f.write(json.dumps(self.data))

    def create(self, owner_id: int, name: str):
        self.id_counter += 1
        u = {
            'id': self.id_counter,
            'owner_id': owner_id,
            'name': name,
        }
        self.data.append(u)
        self.flush()
        return u
    
    def get_all(self):
        return self.data[:]

    def find_by_id(self, id: int):
        for o in self.data:
            if o['id'] == id:
                return o
        return None
    
    def find_by_owner_id(self, owner_id: int):
        docs = []
        for o in self.data:
            if o['owner_id'] == owner_id:
                docs.append(o)
        return docs