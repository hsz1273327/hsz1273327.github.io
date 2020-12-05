import redis
r = redis.Redis(host='localhost', port=5000, db=0)
r.set('foo', 'bar')
result = r.get('foo')
print(result)
