import time
import datetime
import os

def is_today(path):
    filemt = time.localtime(os.path.getmtime(path))
    target_date = time.strftime("%Y-%m-%d", filemt)

    c_year = datetime.datetime.now().year
    c_month = datetime.datetime.now().month
    c_day = datetime.datetime.now().day

    date_list = target_date.split("-")
    t_year = int(date_list[0])
    t_month = int(date_list[1])
    t_day = int(date_list[2])

    final = False

    if c_year == t_year and c_month == t_month and c_day == t_day:
        final = True

    return final