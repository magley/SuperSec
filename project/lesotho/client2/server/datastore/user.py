import json
import bcrypt

class UserRepo:
    def __init__(self):
        self.salt = bcrypt.gensalt()
        self.data = []
        self.id_counter = 0
        self.fname = './data/user.json'
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

    def save(self, email: str, password_raw: str):
        self.id_counter += 1
        u = {
            'id': self.id_counter,
            'email': email,
            'password': bcrypt.hashpw(password_raw.encode(), self.salt).decode("utf-8")
        }
        self.data.append(u)
        self.flush()

    def find_by_id(self, id: int):
        for o in self.data:
            if o['id'] == id:
                return o
        return None
    
    def find_by_email(self, email: str):
        for o in self.data:
            if o['email'] == email:
                return o
        return None
    
    def find_by_email_password(self, email: str, password_raw: str):
        for o in self.data:
            if o['email'] == email and bcrypt.checkpw(password_raw.encode(), o['password'].encode()):
                return o
        return None