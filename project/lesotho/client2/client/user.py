import state
import service
from colors import bcolors


def _log_out():
    state.clear_session_user()


def _edit_document_by_id(doc_id: int):
    if not service.check_access(state.STATE['id'], doc_id, 'editor').body['authorized']:
        print("You are not authorized to edit this document.")
        return
    print("Editing...")


def _read_document_by_id(doc_id: int):
    if not service.check_access(state.STATE['id'], doc_id, 'viewer').body['authorized']:
        print("You are not authorized to view this document.")
        return
    print("Viewing...")



def _select_document():
    all_docs = service.get_all_docs().body

    while True:
        for i, doc in enumerate(all_docs):
            print(f"[{i + 1}.]\t{bcolors.BOLD} {doc['name']} {bcolors.ENDC}")
        inp = input("Enter ordinal number of document:").strip()
        if len(inp) == 0:
            continue
        if inp[-1] == '.':
            inp = inp[:-1]

        try:
            k = int(inp) - 1
            return all_docs[k]
        except:
            continue


def _edit_document():
    doc = _select_document()
    _edit_document_by_id(doc['id'])


def _read_document():
    doc = _select_document()
    _read_document_by_id(doc['id'])


def _new_document():
    doc_name = input("Document name:")
    doc = service.new_doc(state.STATE['id'], doc_name).body
    _edit_document_by_id(doc['id'])


def _prompt_dashboard_dict():
    return {
        'doc-new': { 'desc': "Create a new document", 'func': _new_document },
        'doc-edit': { 'desc': "Edit a document", 'func': _edit_document },
        'doc-read': { 'desc': "Read a document", 'func': _read_document },
        'log-out': { 'desc': "Sign out", 'func': _log_out },
    }


def prompt_user_dict():
    return _prompt_dashboard_dict()