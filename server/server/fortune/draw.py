from PIL import Image, ImageDraw, ImageFont
import os
import json
import random

def drawing(types,fromQQ):
    # 插件目录 C:\Users\asus\Downloads\Programs\github\fortune\server\server\fortune
    base_path = os.path.split(os.path.realpath(__file__))[0]
    img_dir = f'{base_path}/data/img/{types}/'
    img_path = img_dir + random.choice(os.listdir(img_dir))
    out_dir = f'{base_path}/data/out/{types}/'
    out_path = out_dir + f'{fromQQ}.jpg'
    text_path = f'{base_path}/data/text/copywriting.json'
    title_path = f'{base_path}/data/text/goodLuck.json'
    fontPath = {
        'title': f"{base_path}/data/font/Mamelon.otf",
        'text': f"{base_path}/data/font/sakura.ttf"
    }

    if not os.path.exists(out_dir):
       os.makedirs(out_dir)
       print("目录创建成功！")

    img = Image.open(img_path)
    # Draw title
    draw = ImageDraw.Draw(img)
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
    img.save(out_path)
    return out_path

def vertical(str):
    list = []
    for s in str:
        list.append(s)
    return '\n'.join(list)

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

#drawing("noah","test")