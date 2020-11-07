import requests
import base64
import requests
import hashlib
import time
import random

from nonebot import *

bot = get_bot()

@on_command('运势', only_to_me = False)
async def _(session: CommandSession):
    user = session.event.user_id
    group = session.event.group_id

    if group == 12345678:
        values = ["碧蓝幻想", "原神", "公主连结", "诺亚幻想", "虚拟偶像", "阴阳师", "车万", "李清歌"]
        types = random.choice(values)
        limit = "off"
    else:
        types = "诺亚幻想"
        limit = "none"

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
    
    data = {"version":3,"types":types,"fromQQ":user,"ask":"运势","limit":limit}
    headers = {'authkey':authkey,'autime':str(au_time)}

    url = "http://127.0.0.1:8000/fortune"
    fortune = requests.post(url=url,data=data,headers=headers)

    url = "http://127.0.0.1:8000/fortune.jpg"
    pic = requests.post(url=url,data=data,headers=headers)
    
    imgdata = pic.content
    base64_data = base64.b64encode(imgdata).decode("utf-8")

    message = f'[CQ:image,file=base64://{base64_data}]'
    await session.send(message)

