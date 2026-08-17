import copy

class Collection:
    def __init__(self, data=None):
        self._store = {}
        if data:
            self._store.update(data)
        self._store['secret'] = ""

    def MarshalJSON(self):
        result = dict(self._store)
        if 'secret' in result:
            result['secret'] = result['secret']
        return result

    def with_secret(self, value):
        self._store['secret'] = value
        return self

    def get_secret(self):
        return self._store.get('secret', '')

    def __eq__(self, other):
        return self._store == other._store if isinstance(other, Collection) else self._store == other

    def to_dict(self):
        return self.MarshalJSON()