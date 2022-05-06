# toip
批量转换 URL/IP/DOMAIN 成为IP形式，懒癌福音 all to ip

# 用法
处理不规则输入
`echo google.com/sssssss|toip`

`echo google.com:9999/sssssss|toip`

`echo  http://google.com:9999/sssssss|toip    -p`

域名解析

`echo  http://sogou.com:3333/sogogo|toip    -p  -n`


从指定文件读取，保留port，域名解析，输出到指定文件
`cat url.txt|toip  -p  -n   -o result.txt`

![image.png](https://note.youdao.com/yws/res/34761/WEBRESOURCEc18097a1f1f4ecc75931db66e7105b86)


静默输出
`cat url.txt|toip  -p  -n   -o result.txt   -q`

![image.png](https://note.youdao.com/yws/res/34766/WEBRESOURCEb7ce4b1434a2af96200d462adb87a651)

# 参数
Usage of toip:
  -l string
        url文件路径 (default "url.txt")
  -n    对域名进行DNS解析
  -o string
        保存到文件
  -p    保留PORT
  -q    直接输出IP/域名，其他内容不输出

