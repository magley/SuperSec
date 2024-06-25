import requests

class LesothoClient():
    def __init__(self, lesotho_url: str, lesotho_api_key_client_name: str, lesotho_api_key: str) -> None:
        self._lesotho_url = lesotho_url
        self._lesotho_api_key_client_name = lesotho_api_key_client_name
        self._lesotho_api_key = lesotho_api_key

    def acl_update(self, namespace: str, obj: str, relation: str, user: str) -> requests.Response:
        acl_directive = {
            'object': f'{namespace}:{obj}',
            'relation': relation,
            'user': f'{user}',
        }

        headers = {
            "Authorization": f"{self._lesotho_api_key_client_name} {self._lesotho_api_key}"
        }

        return requests.post(f'{self._lesotho_url}/acl', json=acl_directive, headers=headers)

    def acl_query(self, namespace: str, obj: str, relation: str, user: str) -> bool:
        acl_directive = {
            'object': f'{namespace}:{obj}',
            'relation': relation,
            'user': f'{user}',
        }

        headers = {
            "Authorization": f"{self._lesotho_api_key_client_name} {self._lesotho_api_key}"
        }

        resp = requests.get(f'{self._lesotho_url}/acl/check', params=acl_directive, headers=headers)
        if not resp.ok:
            return False
        else:
            return resp.json()['authorized']
        
    def namespace_update(self, namespace_dict: dict) -> requests.Response:
        headers = {
            "Authorization": f"{self._lesotho_api_key_client_name} {self._lesotho_api_key}"
        }

        return requests.post(f'{self._lesotho_url}/namespace', json=namespace_dict, headers=headers)