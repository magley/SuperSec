import state


def log_out():
    state.clear_session_user()


def _prompt_dashboard_dict():
    return {
        'log-out': { 'desc': "Sign out", 'func': log_out },
    }


def prompt_user_dict():
    return _prompt_dashboard_dict()