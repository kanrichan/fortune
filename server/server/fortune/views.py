from rest_framework.views import APIView
from django.http import HttpResponse, HttpRequest, HttpResponseRedirect

from . import utils
from . import draw

import hashlib
import time
import os

au_list = []


def response_as_json(data, foreign_penetrate=False):
    jsonString = serializer(data=data, output_type="json", foreign=foreign_penetrate)
    response = HttpResponse(
            # json.dumps(dataa, cls=MyEncoder),
            jsonString,
            content_type="application/json",
    )
    response["Access-Control-Allow-Origin"] = "*"
    return response


def json_response(data, code=200, foreign_penetrate=False, **kwargs):
    data = {
        "code": code,
        "msg": "成功",
        "data": data,
    }
    return response_as_json(data, foreign_penetrate=foreign_penetrate)


def json_error(error_string="", code=500, **kwargs):
    data = {
        "code": code,
        "msg": error_string,
        "data": {}
    }
    data.update(kwargs)
    return response_as_json(data)

JsonResponse = json_response
JsonError = json_error


def return_fortune(request):
    # 与client端一致的验证key
    au_key = "evenfortune"
    # 从请求头中取出client端 加密前的时间
    client_au_time = request.META['HTTP_AUTIME']
 
    # 将服务端的key 与 client的时间合并成字符
    server_au_key = "%s|%s" % (au_key, client_au_time)
 
    # 然后将字符也同样进行MD5加密
    m = hashlib.md5()
    m.update(bytes(server_au_key, encoding='utf-8'))
    authkey = m.hexdigest()
 
    # 取出client端加密的key
    clint_au_key = request.META['HTTP_AUTHKEY']
 
 
    # 三重验证机制
 
    # 1.超出访问时间5s后不予验证通过。
    server_time = time.time()
    if server_time - 600 > float(client_au_time):
        return HttpResponse(status=403)
 
    # 2.服务端加密的key值 跟 client发过来的加密key比对是否一致？
    if authkey != clint_au_key:
        return HttpResponse(status=403)
 
    # 3.比对当前的key值是否是以前访问过的，访问过的也不予验证通过。
    if authkey in au_list:
        return HttpResponse(status=403)
 
    # 将成功登陆的key值保存在列表中。
    #au_list.append(authkey)

    version = request.POST.get('version')
    types = request.POST.get('types')
    fromQQ = request.POST.get('fromQQ')
    ask = request.POST.get('ask')

    base_path = os.path.split(os.path.realpath(__file__))[0]
    path = f'{base_path}/data/out/{type}/{qq}.jpg'
    try:
        today = utils.is_today(path)
    except:
        today = False
        print('[fortune]\nfile is not exist')
    if today == False:
        path = draw.drawing(type,qq)
    elif qq == 'test':
        path = draw.drawing(type,qq)
    else:
        path = path

    return JsonResponse("http://127.0.0.1/fortune.jpg")



def return_pic(request):
    # 与client端一致的验证key
    au_key = "ceshi"
    # 从请求头中取出client端 加密前的时间
    client_au_time = request.META['HTTP_AUTIME']
 
    # 将服务端的key 与 client的时间合并成字符
    server_au_key = "%s|%s" % (au_key, client_au_time)
 
    # 然后将字符也同样进行MD5加密
    m = hashlib.md5()
    m.update(bytes(server_au_key, encoding='utf-8'))
    authkey = m.hexdigest()
 
    # 取出client端加密的key
    clint_au_key = request.META['HTTP_AUTHKEY']
 
 
    # 三重验证机制
 
    # 1.超出访问时间5s后不予验证通过。
    server_time = time.time()
    if server_time - 600 > float(client_au_time):
        return HttpResponse(status=403)
 
    # 2.服务端加密的key值 跟 client发过来的加密key比对是否一致？
    if authkey != clint_au_key:
        return HttpResponse(status=403)
 
    # 3.比对当前的key值是否是以前访问过的，访问过的也不予验证通过。
    if authkey in au_list:
        return HttpResponse(status=403)
 
    # 将成功登陆的key值保存在列表中。
    #au_list.append(authkey)

    types = request.POST.get('types')
    fromQQ = request.POST.get('fromQQ')
    print(types)
    print(fromQQ)
    print("types")
    base_path = os.path.split(os.path.realpath(__file__))[0]
    path = f'{base_path}/data/out/{types}/{fromQQ}.jpg'
    file_pic = open(path, "rb")
    return HttpResponse(file_pic.read(), content_type='image/jpg')