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

<img width="461" alt="Snipaste_2022-05-06_18-46-49" src="https://user-images.githubusercontent.com/48342077/167119078-a60ac41d-c246-488a-91a1-deaf8801f8fe.png">


静默输出

`cat url.txt|toip  -p  -n   -o result.txt   -q`
<img width="624" alt="Snipaste_2022-05-06_18-48-25" src="https://user-images.githubusercontent.com/48342077/167119099-20fe2c0f-3463-469e-b8dc-0ebb4b4af693.png">

# 参数
Usage of toip:
  -l string
        url文件路径 (default "url.txt")
  -n    对域名进行DNS解析
  -o string
        保存到文件
  -p    保留PORT
  -q    直接输出IP/域名，其他内容不输出

