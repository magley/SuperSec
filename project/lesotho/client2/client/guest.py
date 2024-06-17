import getpass
import service
import state
from colors import bcolors


def log_in():
    email = input("Enter email:")
    password = getpass.getpass('Enter password:')

    resp = service.log_in(email, password)
    if resp.status_code == 200:
        state.set_session_user(resp.body['id'], resp.body['email'])
        print(f"{bcolors.OKCYAN}Logged in as {resp.body['email']}{bcolors.ENDC}")
    else:
        print(bcolors.FAIL + 'Failed to log in' + bcolors.ENDC)


def register():
    email = input("Enter email:")
    password = getpass.getpass('Enter password:')
    password_confirm = getpass.getpass('Confirm password:')

    if password != password_confirm:
        print(bcolors.FAIL + "Passwords do not match" + bcolors.ENDC)
        return

    service.register(email, password)
    print(f"{bcolors.OKCYAN}Success{bcolors.ENDC}")


def prompt_guest_dict():
    return {
        'log-in': { 'desc': "Sign in into the system", 'func': log_in },
        'register': { 'desc': "Create an account", 'func': register },
    }