import state
import service
from colors import bcolors


def _log_out():
    state.clear_session_user()


def _edit_document_by_id(doc_id: int):
    check_access = service.check_access(state.STATE['id'], doc_id, 'viewer')
    if not check_access.ok:
        print(f"{bcolors.FAIL}{check_access.json()['error']}{bcolors.ENDC}")
        return
    else:  
        if not check_access.json()['authorized']:
            print(f"{bcolors.FAIL}You are not authorized to edit this document.{bcolors.ENDC}")
            return

    print(f"{bcolors.OKCYAN}You are in APPEND mode.\nType anything to make changes\nEnter a blank line to finish\n=========================={bcolors.ENDC}")

    ss = ""

    while True:
        i = input()
        if i == "":
            break
        ss += i + '\n'

    should_save = False
    while True:
        i = input("Save changes? [Y/n]:")
        if i in ['y', 'Y', 'yes', '1', 'true']:
            should_save = True
            break
        elif i in ['n', 'N', 'no', '0', 'false']:
            break

    if should_save:
        appenddoc = service.append_to_doc(doc_id, ss)
        if not appenddoc.ok:
            print(f"{bcolors.FAIL}{appenddoc.json()['error']}{bcolors.ENDC}")
            return

def _read_document_by_id(doc_id: int):
    check_access = service.check_access(state.STATE['id'], doc_id, 'viewer')
    if not check_access.ok:
        print(f"{bcolors.FAIL}{check_access.json()['error']}{bcolors.ENDC}")
        return
    else:  
        if not check_access.json()['authorized']:
            print(f"{bcolors.FAIL}You are not authorized to view this document.{bcolors.ENDC}")
            return

    doc = service.get_doc_by_id(doc_id)
    if doc.ok:
        doc = doc.json()
    else:
        print(f"{bcolors.FAIL}{doc.json()['error']}{bcolors.ENDC}")
        return


    print(f"{bcolors.OKCYAN}Here's the document:\n=========================={bcolors.ENDC}")
    print(doc['text'])
    print(f"{bcolors.OKCYAN}=========================={bcolors.ENDC}")


def _select_entity(entities: list, field_to_print: str, prompt: str):
    if len(entities) == 0:
        return None

    while True:
        for i, e in enumerate(entities):
            print(f"[{i + 1}.]\t{bcolors.BOLD} {e[field_to_print]} {bcolors.ENDC}")
        inp = input(prompt).strip()
        if len(inp) == 0:
            continue
        if inp[-1] == '.':
            inp = inp[:-1]

        try:
            k = int(inp) - 1
            return entities[k]
        except:
            continue

def _select_document():
    all_docs = service.get_all_docs()
    if all_docs.ok:
        all_docs = all_docs.json()
    else:
        print(f"{bcolors.FAIL}{all_docs.json()['error']}{bcolors.ENDC}")
        return None

    if len(all_docs) == 0:
        print(f"{bcolors.WARNING}There are no documents in the system{bcolors.ENDC}")
        return None
    return _select_entity(all_docs, 'name', 'Enter ordinal number of document:')


def _select_user():
    all_users = service.get_all_users()
    if all_users.ok:
        all_users = all_users.json()
        all_users = [u for u in all_users if u['email'] != state.STATE['email']]
    else:
        print(f"{bcolors.FAIL}{all_users.json()['error']}{bcolors.ENDC}")
        return None

    if len(all_users) == 0:
        print(f"{bcolors.WARNING}There are no users in the system to share with{bcolors.ENDC}")
        return None
    return _select_entity(all_users, 'email', 'Enter ordinal number of user:')


def _edit_document():
    doc = _select_document()
    if doc is None:
        return
    _edit_document_by_id(doc['id'])


def _read_document():
    doc = _select_document()
    if doc is None:
        return
    _read_document_by_id(doc['id'])


def _share_document():
    roles = [
        {'name': 'editor'},
        {'name': 'viewer'},
    ]
        
    doc = _select_document()
    if doc is None:
        return
    access_check = service.check_access(state.STATE['id'], doc['id'], 'owner')
    if access_check.ok:
        if not access_check.json()['authorized']:
            print(f"{bcolors.FAIL}You are not authorized to grant permissions to users for this document.{bcolors.ENDC}")
            return
    else:
        print(f"{bcolors.FAIL}{access_check.json()['error']}{bcolors.ENDC}")
        return
    
    user_to_share_with = _select_user()
    if user_to_share_with is None:
        return
    if user_to_share_with['id'] == state.STATE['id']:
        print(f"{bcolors.FAIL}You cannot share with yourself.{bcolors.ENDC}")
        return
    
    role = _select_entity(roles, 'name', 'Select role to grant:')
    
    sharedoc = service.share_doc(user_to_share_with['id'], doc['id'], role['name'])
    if not sharedoc.ok:
        print(f"{bcolors.FAIL}{sharedoc.json()['error']}{bcolors.ENDC}")
        return


def _new_document():
    doc_name = input("Document name:")

    doc = service.new_doc(state.STATE['id'], doc_name)
    if doc.ok:
        doc = doc.json()
    else:
        print(f"{bcolors.FAIL}{doc.json()['error']}{bcolors.ENDC}")
        return
    
    _edit_document_by_id(doc['id'])


def _prompt_dashboard_dict():
    return {
        'doc-new': { 'desc': "Create a new document", 'func': _new_document },
        'doc-edit': { 'desc': "Edit a document", 'func': _edit_document },
        'doc-read': { 'desc': "Read a document", 'func': _read_document },
        'doc-share': { 'desc': "Share a document", 'func': _share_document },
        'log-out': { 'desc': "Sign out", 'func': _log_out },
    }


def prompt_user_dict():
    return _prompt_dashboard_dict()