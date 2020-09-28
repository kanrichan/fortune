from rest_framework.views import APIView
from django.http import HttpResponse, HttpRequest
import os
import json
import time
import datetime
import random
from PIL import Image, ImageDraw, ImageFont

PATH = '/root/web/web/web/plugins/fortune/data'

PATH_FORTUNE = '/root/web/web/web/plugins/fortune/data/fortune'


def drawing(type,qq):
    PATH_IMG = f'/root/web/web/web/plugins/fortune/data/img/{type}'
    print(PATH_IMG)
    PATH_OUT = f'/root/web/web/web/plugins/fortune/data/out/{type}/{qq}.jpg'
    fontPath = {
        'title': '/root/web/web/web/plugins/fortune/data/font/Mamelon.otf',
        'text': '/root/web/web/web/plugins/fortune/data/font/sakura.ttf'
    }
    imgPath = PATH_IMG + '/' + random.choice(os.listdir(PATH_IMG))
    img = Image.open(imgPath)
    # Draw title
    draw = ImageDraw.Draw(img)
    text_path = PATH_FORTUNE+'/copywriting.json'
    title_path = PATH_FORTUNE+'/goodLuck.json'
    with open(text_path, 'r', encoding='utf-8') as f:
        content = f.read()
    content = json.loads(content)
    text = random.choice(content['copywriting'])
    with open(title_path, 'r', encoding='utf-8') as f:
        content = f.read()
    content = json.loads(content)
    for i in content['types_of']:
        if i['good-luck'] == text['good-luck']:
            title = i['name']
    text = text['content']
    font_size = 45
    color = '#F5F5F5'
    image_font_center = (140, 99)
    ttfront = ImageFont.truetype(fontPath['title'], font_size)
    font_length = ttfront.getsize(title)
    draw.text((image_font_center[0]-font_length[0]/2, image_font_center[1]-font_length[1]/2),
                title, fill=color,font=ttfront)
    # Text rendering
    font_size = 25
    color = '#323232'
    image_font_center = [140, 297]
    ttfront = ImageFont.truetype(fontPath['text'], font_size)
    result = decrement(text)
    if not result[0]:
        return 
    textVertical = []
    for i in range(0, result[0]):
        font_height = len(result[i + 1]) * (font_size + 4)
        textVertical = vertical(result[i + 1])
        x = int(image_font_center[0] + (result[0] - 2) * font_size / 2 + 
                (result[0] - 1) * 4 - i * (font_size + 4))
        y = int(image_font_center[1] - font_height / 2)
        draw.text((x, y), textVertical, fill = color, font = ttfront)
    # Save
    img = img.convert("RGB")
    outPath = PATH_OUT
    img.save(outPath)
    return outPath

def decrement(text):
    length = len(text)
    result = []
    cardinality = 9
    if length > 4 * cardinality:
        return [False]
    numberOfSlices = 1
    while length > cardinality:
        numberOfSlices += 1
        length -= cardinality
    result.append(numberOfSlices)
    # Optimize for two columns
    space = ' '
    length = len(text)
    if numberOfSlices == 2:
        if length % 2 == 0:
            # even
            fillIn = space * int(9 - length / 2)
            return [numberOfSlices, text[:int(length / 2)] + fillIn, fillIn + text[int(length / 2):]]
        else:
            # odd number
            fillIn = space * int(9 - (length + 1) / 2)
            return [numberOfSlices, text[:int((length + 1) / 2)] + fillIn,
                                    fillIn + space + text[int((length + 1) / 2):]]
    for i in range(0, numberOfSlices):
        if i == numberOfSlices - 1 or numberOfSlices == 1:
            result.append(text[i * cardinality:])
        else:
            result.append(text[i * cardinality:(i + 1) * cardinality])
    return result

def vertical(str):
    list = []
    for s in str:
        list.append(s)
    return '\n'.join(list)

def is_today(path):
    filemt = time.localtime(os.path.getmtime(path))
    target_date = time.strftime("%Y-%m-%d", filemt)
    """
    2020-03-25 17:03:55
    Detects if the date is current date
    :param target_date:
    :return: Boolean
    """
    # Get the year, month and day
    c_year = datetime.datetime.now().year
    c_month = datetime.datetime.now().month
    c_day = datetime.datetime.now().day

    # Disassemble the date
    date_list = target_date.split("-")
    t_year = int(date_list[0])
    t_month = int(date_list[1])
    t_day = int(date_list[2])

    final = False
    # Compare years, months and days
    if c_year == t_year and c_month == t_month and c_day == t_day:
        final = True

    return final

class return_fortune(APIView):
    def get(self, request):
        type = request.GET.get("type")
        if type not in 'vtb|pcr|noah|ys|df':
            return HttpResponse('400')
        if type == "vtb":
            type = "df"
        qq = request.GET.get("qq")
        path = f'/root/web/web/web/plugins/fortune/data/out/{type}/{qq}.jpg'
        try:
            today = is_today(path)
        except:
            today = False
            print('[fortune]\nfile is not exist')
        if today == False:
            path = drawing(type,qq)
        elif qq == 'test':
            path = drawing(type,qq)
        else:
            path = path
        file_pic = open(path, "rb")
        return HttpResponse(file_pic.read(), content_type='image/jpg')

def return_broadcast(request):
    version = request.GET.get("version")
    type = request.GET.get("type")
    qq = request.GET.get("qq")
    ask = request.GET.get("ask")
    if type == "vtb":
        broadcast = "[木理插件-运势] 公告：\n该服务已暂停~"
    elif qq == "12345678":
        broadcast = "[木理插件-运势] 公告：\n已停止服务"
    elif ask == "kkp":
        broadcast = "[木理插件-运势] 公告：\n已停止服务"
    elif int(version) < 1:
        broadcast = "[木理插件-运势] 提醒：\n存在新的版本，请尽快更新"
    else:
        broadcast = "正常运行"
    print(broadcast)
    return HttpResponse(broadcast)
