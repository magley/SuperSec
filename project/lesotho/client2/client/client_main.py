import state
from colors import bcolors
import guest
import user


def prompt(dict):
    '''
    Prompt user for a command given the list of commands.
    The user may type `help` for a list of commands

    `dict` has the following format
    ```
    {
        "some_command_name": {
            desc: "Blabla"
            func: some_function_that_takes_no_arguments
        }
    }
    ```
    '''

    def print_commands():
        print()
        for (cmd_name, cmd_data) in dict.items():
            cmd_desc = cmd_data['desc']
            cmd_func = cmd_data['func']
            print(f'{bcolors.HEADER}{cmd_name}\n\t{cmd_desc}{bcolors.ENDC}')
        print()
    
    print_commands()

    while True:
        i = input()
        if i in dict:
            dict[i]['func']()
            return
        if i == 'help':
            print_commands()
        else:
            print(bcolors.FAIL + 'Unknown command.\nType help for a list of commands' + bcolors.ENDC)


if __name__ == "__main__":
    while True:
        d = {}
        if state.get_session_is_logged_in():
            d = user.prompt_user_dict()
        else:
            d = guest.prompt_guest_dict()

        prompt(d)