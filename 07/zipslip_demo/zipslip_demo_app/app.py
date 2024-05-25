import tarfile
import tempfile
from flask import Flask, abort, request
import os

app = Flask(__name__)


@app.route("/")
def root():
    return '''
<form method="post" enctype="multipart/form-data" action="http://localhost:5000/upload">
  <input type="file" name="file">
  <input type="submit">
</form>
'''


@app.route('/upload', methods=['POST'])
def cache():
    if 'file' not in request.files:
        return abort(400)
    extraction = extract_from_archive(request.files['file'])
    if extraction:
        return {"list": extraction}, 200
    return '', 204


def extract_from_archive(file):
    tmp = tempfile.gettempdir()
    path = os.path.join(tmp, file.filename)
    file.save(path)

    if tarfile.is_tarfile(path):
        tar = tarfile.open(path, 'r:gz')
        tar.extractall(tmp)

        extractdir = 'uploads'
        os.makedirs(extractdir, exist_ok=True)

        extracted_filenames = []

        for tarinfo in tar:
            name = tarinfo.name
            if tarinfo.isreg():
                filename = f'{extractdir}/{name}'
                os.replace(os.path.join(tmp, name), filename)
                extracted_filenames.append(filename)
                continue

            os.makedirs(f'{extractdir}/{name}', exist_ok=True)

        tar.close()
        return extracted_filenames

    return False
