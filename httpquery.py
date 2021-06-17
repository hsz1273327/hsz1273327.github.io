from requests import post, get

#url = "http://localhost:8501/v1/models/half_plus_two/metadata"
#res = get(url)
url = "http://localhost:8501/v1/models/half_plus_two:predict"
res = post(url, json={"instances": [1.0, 2.0, 5.0]})
print(res.status_code)
print(res.json())
