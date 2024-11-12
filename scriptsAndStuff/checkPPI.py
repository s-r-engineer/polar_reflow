from dateutil import parser
# from datetime import datetime
import json
import multiprocessing
import os
from threading import Lock
import sqlite3

filellist = []

for a, b, c in os.walk('../../data'):
    for f in c:
        if "ppi_samples" in f:
            filellist.append(f'{a}/{f}')
lines1 = []
list_lock = Lock()


def runner(f):
    print(f)
    lines = []
    with open(f, 'r') as r:
        data = json.load(r)
    for i in data:
        for sample in i['devicePpiSamplesList']:
            for ppiSample in sample['ppiSamples']:
                length = ppiSample['pulseLength']
                date = parser.parse(ppiSample['sampleDateTime'])
                lines.append([date, length])
    return lines


with multiprocessing.Pool(12) as ppol:
    poolshot = ppol.map(runner, filellist)
print("done with reading")
con = sqlite3.connect("ppi.db")
cur = con.cursor()
cur.execute("CREATE TABLE if not exists ppi(time datetime unique, value integer)")
cur.execute("CREATE index if not exists ppi_ppi on ppi(time)")
for i in poolshot:
    counter3 = 0
    while True:
        d = i[counter3:counter3 + 10001]
        if len(d) == 0:
            break
        cur.executemany("INSERT or ignore INTO ppi VALUES(?, ?)", d)
        con.commit()
        counter3 += 10000