from flask import Flask, send_from_directory
import os

root_dir = os.getenv("MU_ROOT_DIR")
lib_dir = os.path.join(os.path.dirname(os.path.realpath(__file__)),"mangojs/")

app = Flask(__name__)

@app.route('/mangojs/<path:path>')
def do_mangojs(path):
    return send_from_directory(lib_dir, path)

@app.route('/<path:path>')
def others(path):
    return send_from_directory(root_dir, path)

app.run(port=int(os.getenv("MU_HTTP_PORT")),debug=True)
