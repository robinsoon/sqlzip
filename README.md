## zip压缩加密文件夹 v1.0
-----
创建带加密口令的zip文件<br/>
支持对文件夹压缩后设定密码。<br/>
支持解压到指定文件夹，解压口令。<br/>
commandline命令传入的方式调用。<br/>



运行指令
--------

加密和解密，说明如下:

*   **压缩加密**  用于文件夹压缩。
    - ``压缩加密方法 ZipEncrypt`` 先压缩+再加密文件 <br/>
    - 先压缩,一旦你加密文件，将生成一个随机数据流，这是不可压缩的。压缩过程依赖于在数据中找到可压缩模式。<br/> 
    - 压缩命令启动： 
-----
    >sqlzip.exe -z filepathname zipname mypassword
-----


​	
*   **解密解压**  
    存在更新内容。<br/>
    - 解密文件方法  ``UnZipDecrypt``
    - 解压命令启动：
-----
	>sqlzip.exe -u zipfile unzipfile mypassword
-----
	- 缺省参数使用符号 _ 替代，示例： 
	
	>sqlzip.exe -u zipfile _ mypassword
-----