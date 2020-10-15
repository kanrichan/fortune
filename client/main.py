import requests
import hashlib
import time
 
# 基于验证的key
au_key = "test"

# 基于验证的时间
au_time = int(time.time())

# 将验证的key与时间合并成一个字符
au_key_time = "%s|%s"%(au_key,au_time)

# 将合并的字符进行MD5加密
m = hashlib.md5()
m.update(bytes(au_key_time,encoding='utf-8'))
authkey = m.hexdigest()

# 将生成加密的 KEY 与 时间传递至服务端
url = "http://www.kanri.ml:10086/fortune.jpg"
data = {"types":"诺亚幻想",'fromQQ':"test"}
headers = {'authkey':authkey,'autime':str(au_time)}

a = requests.post(url=url,data=data,headers=headers)
print(a.text)